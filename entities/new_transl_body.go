package entities

import "strings"

type TranslationsBody struct {
	Translations []*Translation `json:"translations"`
}

func (tb TranslationsBody) TrimFields() {
	for i := 0; i < len(tb.Translations); i++ {
		tb.Translations[i].Oc = strings.TrimSpace(tb.Translations[i].Oc)
		tb.Translations[i].Fr = strings.TrimSpace(tb.Translations[i].Fr)
		tb.Translations[i].En = strings.TrimSpace(tb.Translations[i].En)
		tb.Translations[i].DatasetID = strings.TrimSpace(tb.Translations[i].DatasetID)
		tb.Translations[i].Occitan = strings.TrimSpace(tb.Translations[i].Occitan)
	}
}
