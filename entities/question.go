package entities

import (
	"strings"
)

type SecretQuestionAndResponse struct {
	Question string `json:"question"`
	Response string `json:"response"`
}

type SecretQuestionsAndResponses []*SecretQuestionAndResponse

func (sqar SecretQuestionsAndResponses) Trimer() {
	for i := 0; i < len(sqar); i++ {
		sqar[i].Question = strings.TrimSpace(sqar[i].Question)
		sqar[i].Response = strings.TrimSpace(sqar[i].Response)
	}
}

func (sqar SecretQuestionsAndResponses) ToLowerResponses() {
	for i := 0; i < len(sqar); i++ {
		sqar[i].Response = strings.ToLower(sqar[i].Response)
	}
}
