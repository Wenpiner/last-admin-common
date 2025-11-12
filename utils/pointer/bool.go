package pointer

// GetBool returns the value of v, if v is nil, returns def.
func GetBool(v *bool) bool {
	if v == nil {
		return false
	}
	return *v
}

// ToBoolPtr returns a pointer to the value.
func ToBoolPtr(v bool) *bool {
	return &v
}

// ToBoolPtrIfNotFalse returns a pointer to the value if it's not false, otherwise returns nil.
func ToBoolPtrIfNotFalse(v bool) *bool {
	if !v {
		return nil
	}
	return &v
}
