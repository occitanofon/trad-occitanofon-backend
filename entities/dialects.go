package entities

type DialectSubdialectDocument struct {
	Dialect    string `json:"dialect"`
	Subdialect string `json:"subdialect"`
}

type Occitan struct {
	Dialect     string       `json:"dialect"`
	Subdialects []Subdialect `json:"subdialects"`
}

type Subdialect struct {
	Name                        string `json:"name"`
	TotalTranslated             int    `json:"totalTranslated"`
	TotalTranslatedByTranslator int    `json:"totalTranslatedByTranslator"`
}
