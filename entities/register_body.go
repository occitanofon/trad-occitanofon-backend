package entities

import "strings"

type RegisterBody struct {
	Username                    string                      `json:"username"`
	Email                       string                      `json:"email"`
	Password                    string                      `json:"password"`
	SecretQuestionsAndResponses SecretQuestionsAndResponses `json:"secretQuestionsAndResponses"`
}

func (rb *RegisterBody) TrimFields() {
	rb.Email = strings.TrimSpace(rb.Email)
	rb.Username = strings.TrimSpace(rb.Username)
	rb.Password = strings.TrimSpace(rb.Password)
}
