package auth

type Service interface {
	Register(u Credentials) error
	Login(cred Credentials) (SignedUserID, error)
	Validate(sgn SignedUserID) error
}
