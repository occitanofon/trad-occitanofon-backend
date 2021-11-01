package data

type Occitan struct {
	Dialect     string   `json:"dialect"`
	Subdialects []string `json:"subdialects"`
}

var OCCITAN = []Occitan{
	{
		Dialect:     "auvernhat",
		Subdialects: []string{"estandard", "brivadés", "septentrional"},
	},
	{
		Dialect:     "gascon",
		Subdialects: []string{"estandard", "aranés", "bearnés"},
	},
	{
		Dialect:     "lengadocian",
		Subdialects: []string{"estandard", "agenés", "besierenc", "carcassés", "roergat"},
	},
	{
		Dialect:     "lemosin",
		Subdialects: []string{"estandard", "marchés", "peiregordin"},
	},
	{
		Dialect:     "provençau",
		Subdialects: []string{"estandard", "maritime", "niçard", "rodanenc"},
	},
	{
		Dialect:     "vivaroaupenc",
		Subdialects: []string{"estandard", "aupenc", "gavòt", "vivarodaufinenc"},
	},
}
