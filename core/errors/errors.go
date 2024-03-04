package errors

// Error type
type Error string

// Error returns the error message
func (e Error) Error() string {
	return string(e)
}

const (
	NA                = Error("N/A") // Special error type.
	NotYetImplemented = Error("Not yet implemented")
	Unauthenticated   = Error("Failed auth")   // Authentication - Who are you?
	Unauthorized      = Error("Access Denied") // Authorization -Are you allow to do that?
	InvalidToken      = Error("Invalid token")
	InvalidRequest    = Error("Invalid request")
	InvalidHeader     = Error("Invalid header")
	MissingHeader     = Error("Missing header")
	ServiceDown       = Error("Service unavailable")
	InvalidParameter  = Error("Invalid parameter") // Common error case about invalid parameter
	ItemNotFound      = Error("Item not found")    // Common error case item not found
)
