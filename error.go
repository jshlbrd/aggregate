package aggregate

// Error is used for custom Aggregate errors.
type Error string

func (e Error) Error() string { return string(e) }
