package pointer

// GetUint32 returns the value of the pointer if it's not nil, otherwise returns 0.
func GetUint32(v *uint32) uint32 {
	if v == nil {
		return 0
	}

	return *v
}

// GetInt32 returns the value of the pointer if it's not nil, otherwise returns 0.
func GetInt32(v *int32) int32 {
	if v == nil {
		return 0
	}

	return *v
}

// GetInt64 returns the value of the pointer if it's not nil, otherwise returns 0.
func GetInt64(v *int64) int64 {
	if v == nil {
		return 0
	}

	return *v
}

// ToUint32Ptr returns a pointer to the value.
func ToUint32Ptr(v uint32) *uint32 {
	return &v
}

// ToUint64Ptr returns a pointer to the value.
func ToUint64Ptr(v uint64) *uint64 {
	return &v
}

// ToInt64Ptr returns a pointer to the value.
func ToInt64Ptr(v int64) *int64 {
	return &v
}

// ToInt32Ptr returns a pointer to the value.
func ToInt32Ptr(v int32) *int32 {
	return &v
}

// ToUint32PtrIfNotZero returns a pointer to the uint32 value if it's not zero, otherwise returns nil.
func ToUint32PtrIfNotZero(v uint32) *uint32 {
	if v == 0 {
		return nil
	}
	return &v
}

// ToInt64PtrIfNotNil returns a pointer to the int64 value if it's not nil, otherwise returns nil.
func ToInt64PtrIfNotNil(v int64) *int64 {
	if v == 0 {
		return nil
	}
	return &v
}

func ToInt32PtrIfNotZero(v int32) *int32 {
	if v == 0 {
		return nil
	}
	return &v
}
