package pointer

// GetString returns the string value of the pointer.
func GetString(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

// ToStringPtr returns a pointer to the string value.
func ToStringPtr(s string) *string {
	return &s
}

// ToStringPtrIfNotEmpty returns a pointer to the string value if it's not empty, otherwise returns nil.
func ToStringPtrIfNotEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
