package service_errors

const (
	// Token
	UnExpectedError                    = "Expected error"
	ClaimsNotFound                     = "Claims not found"
	TokenRequired                      = "token required"
	TokenExpired                       = "token expired"
	TokenInvalid                       = "token invalid"
	TokenNotFound                      = "token not found"
	TokenVerificationFailed            = "token verification failed"
	TokenNotFoundOrVerify              = "token not found or verify"
	YourEmailIsNotVerified             = "Your email is not verified"
	YourEmailIsNotVerifiedCheckedEmail = "Your email is not verified Please check your email and verify your email"

	// User
	InvalidCredentials  = "Invalid credentials"
	EmailIsNotValid     = "Email is not valid"
	EmailExists         = "Email exists"
	UsernameExists      = "Username exists"
	UserNotFound        = "User not found"
	PermissionDenied    = "Permission denied"
	EmailNotFound       = "Email not found"
	PasswordsDoNotMatch = "Passwords do not match"
	PasswordVAlidation  = "The password must contain uppercase and lowercase letters, numbers and symbols"

	// DB
	RecordNotFound      = "record not found"
	DeleteRecordError   = "delete record error"
	UnsupportedRelation = "unsupported relation"

	// email
	ErrorInSendEmail = "error in send email"
	ValueIsNotBool   = "value is not boolean"

	// subscription error
	SubscriptionRecordNotFound = "subscription record not found"
	FailedToGetMonthlyEntries  = "failed to get monthly entries"
	FailedToGetWeeklyEntries   = "failed to get weekly entries"

	AmountIsNotZero = "The amount must be greater than 2"
	// Appointment error
	TimeFormatIsNotValid       = "invalid appointment date format"
	SlotNotFound               = "appointment slot not found"
	InvalidAppointmentTime     = "invalid appointment time"
	AppointmentAlreadyReserved = "appointment slot already reserved"

	// Empty

	URlORFileEMPTY = "url is empty"
)
