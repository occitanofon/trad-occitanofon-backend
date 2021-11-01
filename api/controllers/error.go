package controllers

const (
	ErrDefault                    string = "ua error es escadença"
	ErrAccountNotConfirmed        string = "ton compte n'a pas encore été activé, contacte l'administrateur"
	ErrAccountSuspended           string = "ton compte es estat suspendut, contacta l'administrator"
	ErrDialectNotProvided         string = "Aucun dialect n'a été trouvé"
	ErrTooMuchTranslationsFetched string = "As revirat 300 frasas en doas oras, tòrna mai tard"
	ErrNoMoreDataset              string = "i a pas cap de frasa de traduire"
	ErrPasswordTooShort           string = "Ton mot de passe doit contenir au moins 10 caractères"
	ErrSecretQuestions            string = "Tu n'as pas saissi les 2 secret questions"
	ErrMailerNotAllowed           string = "Tu as déjà récemment essayé de changer ton mot de passe, ressaye ultérieurement"
	ErrSecretQuestionsNoMatch     string = "Une erreur est survenue avec les questions secrètes"
	ErrBadCredentials             string = "lo pseudonim o lo còdi d'activacion es incorrècte"
	ErrNoPermDialect              string = "Aucune dialect ne t'a été attribué, contacte l'administrateur"
	ErrBadFullDialectFormat       string = "Une erreur avec le dialect saissi"
)
