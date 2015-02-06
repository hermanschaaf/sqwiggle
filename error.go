package sqwiggle

// Error represents an error that could be returned by the Sqwiggle API
type Error struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	Details string    `json:"details"`
	Param   string    `json:"param"`
}

// Error is an implementation of the error interface
func (err Error) Error() string {
	return err.Message
}

// ErrorType is a type that describes the type of error returned
// by the Sqwiggle API. The constants below are defined to use this
// type.
type ErrorType string

func (e ErrorType) String() string {
	return string(e)
}

// These errors define the types of errors expected to be returned
// by the Sqwiggle API.
const (
	ErrAuthentication ErrorType = "authentication"
	ErrAuthorization            = "authorization"
	ErrInvalidParam             = "invalid_param"
	ErrUnknownParam             = "unknown_param"
	ErrLimitReached             = "limit_reached"
	ErrValidation               = "validation"
	ErrUnknown                  = "unknown"
)
