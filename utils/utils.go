package utils

func NewPtr[T any](val T) *T {
	return &val
}

func GetValFromPtr[T any](val *T) T {
	var ret T
	if val == nil {
		return ret
	}
	return *val
}
