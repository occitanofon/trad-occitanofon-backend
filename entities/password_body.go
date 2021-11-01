package entities

type PasswordResetBody struct {
	Email string `json:"email"`
}

type PasswordUpdateBody struct {
	Token                       string                      `json:"token"`
	SecretQuestionsAndResponses SecretQuestionsAndResponses `json:"secretQuestionsAndResponses"`
	Password                    string                      `json:"password"`
}
