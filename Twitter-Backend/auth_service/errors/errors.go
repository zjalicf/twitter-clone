package errors

const (
	InvalidRequestFormatError = "invalid request format"
	InvalidTokenError         = "token is invalid"
	InvalidUserTokenError     = "invalid user token"
	ExpiredTokenError         = "verification token has expired"
	InvalidResendMailError    = "invalid resend mail"
	NotFoundMailError         = "mail not found in database"
)
