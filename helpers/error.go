package helpers

const (
	ErrDefault                string = "una error es subrevenguda"
	ErrOnlyOneSpaceByUsername string = "pas mai que 1 espaci"
	ErrUsernameTooShort       string = "ton pseudo est trop court"
	ErrEmailTooShort          string = "l'email est trop court"
	ErrEmailTooLong           string = "l'email est trop long"
	ErrEmailBadFormat         string = "cette email n'est pas dans un bon format"
	ErrEmailDomainNotExist    string = "ce domaine n'existe pas"
)

type HelperError struct {
	Code    int
	Message string
	Wrapped error
}

func (e HelperError) Error() string {
	return e.Wrapped.Error()
}
