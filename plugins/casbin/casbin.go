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

// NewCasbin
func (l CasbinConf) NewCasbin(dbType, dsn string) (*casbin.Enforcer, error) {
	adapter, err := ent_adapter.NewAdapter(dbType, dsn)
	logx.Must(err)

	var text string
	if l.ModelText == "" {
		text = `
			[request_definition]
			r = sub, obj, act

			[policy_definition]
			p = sub, obj, act

			[role_definition]
			g = _, _

			[policy_effect]
			e = some(where (p.eft == allow))

			[matchers]
			m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && r.act == p.act
		`
	} else {
		text = l.ModelText
	}
	m, err := model.NewModelFromString(text)
	logx.Must(err)

	enforcer, err := casbin.NewEnforcer(m, adapter)
	logx.Must(err)

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
		rediswatcher.DefaultUpdateCallback(cbn)(data)
	})
	err := cbn.SetWatcher(w)
	logx.Must(err)
	err = cbn.SavePolicy()
	logx.Must(err)
	return cbn
}
