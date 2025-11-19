package casbin

import (
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	ent_adapter "github.com/casbin/ent-adapter"
	rediswatcher "github.com/casbin/redis-watcher/v2"
	redis2 "github.com/redis/go-redis/v9"
	"github.com/wenpiner/last-admin-common/config"
	last_redis "github.com/wenpiner/last-admin-common/last-redis"
	"github.com/zeromicro/go-zero/core/logx"
)

type CasbinConf struct {
	ModelText string `json:"ModelText,optional,env=CASBIN_MODEL_TEXT"`
}

// permissionMatch 自定义权限匹配函数
// 处理权限继承关系：
// - write 权限包含 read 权限（用于 configuration 域）
// - admin 权限包含所有权限
// 参数类型为 interface{} 以兼容 Casbin 的 AddFunction 要求
func permissionMatch(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return false, nil
	}

	grantedAction, ok1 := args[0].(string)
	requiredAction, ok2 := args[1].(string)

	if !ok1 || !ok2 {
		return false, nil
	}

	// 完全匹配
	if grantedAction == requiredAction {
		return true, nil
	}

	// 权限继承规则
	// 1. write 包含 read（用于 configuration 域）
	if requiredAction == "read" && grantedAction == "write" {
		return true, nil
	}

	// 2. admin 包含所有权限
	if grantedAction == "admin" {
		return true, nil
	}

	return false, nil
}

// NewCasbin
func (l CasbinConf) NewCasbin(dbType, dsn string) (*casbin.Enforcer, error) {
	adapter, err := ent_adapter.NewAdapter(dbType, dsn)
	logx.Must(err)

	var text string
	if l.ModelText == "" {
		text = `
				[request_definition]
				# r = sub, dom, obj, act
				# r.sub: 主体 (用户)
				# r.dom: 域 (在这里就是您的“资源类型”，如 'api', 'config', 'menu')
				# r.obj: 客体 (资源路径)
				# r.act: 动作 (HTTP方法 或 'read'/'write'/'view')
				r = sub, dom, obj, act

				[policy_definition]
				# p = sub, dom, obj, act
				# p.sub: 主体 (角色)
				# p.dom: 域 (资源类型)
				# p.obj: 客体 (资源路径)
				# p.act: 动作
				p = sub, dom, obj, act

				[role_definition]
				# g = _, _
				# 我们保持角色是全局的，不分域
				g = _, _

				[policy_effect]
				e = some(where (p.eft == allow))

				[matchers]
				# 使用自定义 permissionMatch 函数处理权限继承
				# 1. 检查用户(r.sub)是否属于角色(p.sub)
				# 2. 检查请求的域(r.dom)是否与策略的域(p.dom)完全匹配
				# 3. 检查资源路径(r.obj)是否匹配 (继续使用 keyMatch2)
				# 4. 检查动作(r.act)是否匹配，支持权限继承
				m = g(r.sub, p.sub) && r.dom == p.dom && keyMatch2(r.obj, p.obj) && permissionMatch(p.act, r.act)
		`
	} else {
		text = l.ModelText
	}
	m, err := model.NewModelFromString(text)
	logx.Must(err)

	enforcer, err := casbin.NewEnforcer(m, adapter)
	logx.Must(err)

	// 注册自定义权限匹配函数
	enforcer.AddFunction("permissionMatch", permissionMatch)

	err = enforcer.LoadPolicy()
	logx.Must(err)

	return enforcer, nil
}

// MustNewCasbin
func (l CasbinConf) MustNewCasbin(dbType, dsn string) *casbin.Enforcer {
	csb, err := l.NewCasbin(dbType, dsn)
	if err != nil {
		logx.Errorw("initialize Casbin failed", logx.Field("detail", err.Error()))
		log.Fatalf("initialize Casbin failed, error: %s", err.Error())
		return nil
	}

	return csb
}

func (l CasbinConf) MustNewRedisWatcher(c config.RedisConfig, f func(string2 string)) persist.Watcher {
	w, err := rediswatcher.NewWatcher(c.Host, rediswatcher.WatcherOptions{
		Options: redis2.Options{
			Network:  "tcp",
			Password: c.Password,
		},
		Channel:    string(last_redis.CasbinChannel),
		IgnoreSelf: false,
	})
	logx.Must(err)

	err = w.SetUpdateCallback(f)
	logx.Must(err)

	return w
}

// MustNewCasbinWithRedisWatcher returns Casbin Enforcer with Redis watcher.
func (l CasbinConf) MustNewCasbinWithRedisWatcher(dbType, dsn string, c config.RedisConfig) *casbin.Enforcer {
	cbn := l.MustNewCasbin(dbType, dsn)
	w := l.MustNewRedisWatcher(c, func(data string) {
		// 在 watcher 回调中重新注册自定义函数，确保策略更新后函数仍然可用
		cbn.AddFunction("permissionMatch", permissionMatch)
		rediswatcher.DefaultUpdateCallback(cbn)(data)
	})
	err := cbn.SetWatcher(w)
	logx.Must(err)
	err = cbn.SavePolicy()
	logx.Must(err)
	return cbn
}

// Check
func Check(cbn *casbin.Enforcer, domain string, rolesIds []string, obj, act string) bool {
	var reqs [][]any
	for _, v := range rolesIds {
		reqs = append(reqs, []any{v, domain, obj, act})
	}

	res, err := cbn.BatchEnforce(reqs)
	if err != nil {
		logx.Errorw("验证 Casbin 异常", logx.Field("error", err))
		return false
	}

	for _, v := range res {
		if v {
			return true
		}
	}

	return false
}
