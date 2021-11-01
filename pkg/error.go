package pkg

const (
	ErrDefault                   string = "ua error es susvenguda"
	ErrSecretQuestionsNotFound   string = "les questions secrêtes sont introuvables"
	ErrBadCredentials            string = "lo pseudonim o lo còdi d'activacion es incorrècte"
	ErrEmailUsed                 string = "cet email est déjà utilisé"
	ErrPseudoUsed                string = "ce pseudonyme est déjà utilisé"
	ErrRefreshTokenNotFound      string = "une erreur est survenue lors de la réauthentification de votre compte"
	ErrResetPasswordTokenExpired string = "le token permettant de procéder au changement de mot de passe a expiré"
	ErrTranslatorNotFound        string = "ce traducteur n'existe pas"
	ErrDialectNotFound           string = "ce dialecte n'existe pas"
	ErrDialectsNotFound          string = "ces dialectes n'existent pas"
)

type DBError struct {
	Code    int
	Message string
	Wrapped error
}

func (e DBError) Error() string {
	return e.Wrapped.Error()
}
