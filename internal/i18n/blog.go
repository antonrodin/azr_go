package i18n

import (
	"fmt"
	"time"
)

// BlogText son los textos de interfaz del blog. Los artículos en sí viven en
// content/blog/*.md.
type BlogText struct {
	Lang     string
	BasePath string // "/es/blog", "/en/blog", "/ru/blog"

	IndexTitle string // <title> y og:title del listado
	IndexDesc  string // meta description del listado
	IndexH1    string
	IndexIntro string // apostilla con guasa
	Empty      string
	Back       string // volver al listado
	ReadMore   string
	ImageBy    string // prefijo del crédito de la imagen

	months [12]string
	date   string // formato con %d día, %s mes, %d año
}

// FormatDate pinta una fecha a la manera de cada idioma.
func (t BlogText) FormatDate(d time.Time) string {
	m := t.months[d.Month()-1]
	if t.Lang == "en" {
		return fmt.Sprintf(t.date, m, d.Day(), d.Year())
	}
	return fmt.Sprintf(t.date, d.Day(), m, d.Year())
}

var blog = map[string]BlogText{
	"es": {
		Lang:       "es",
		BasePath:   "/es/blog",
		IndexTitle: "Blog — Antón Zekeriev Rodin",
		IndexDesc:  "Artículos de Antón Zekeriev Rodin sobre desarrollo de software, certificaciones y lo que se tercie.",
		IndexH1:    "Blog",
		IndexIntro: "Cosas que he aprendido, casi siempre por las malas. Escritas para que tú no tengas que hacerlo.",
		Empty:      "Todavía no hay artículos. Se está cocinando.",
		Back:       "← Volver al blog",
		ReadMore:   "Leer →",
		ImageBy:    "Imagen:",
		months:     [12]string{"enero", "febrero", "marzo", "abril", "mayo", "junio", "julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre"},
		date:       "%d de %s de %d",
	},
	"en": {
		Lang:       "en",
		BasePath:   "/en/blog",
		IndexTitle: "Blog — Antón Zekeriev Rodin",
		IndexDesc:  "Articles by Antón Zekeriev Rodin about software development, certifications and whatever comes up.",
		IndexH1:    "Blog",
		IndexIntro: "Things I've learned, usually the hard way. Written down so you don't have to.",
		Empty:      "No articles yet. Something's cooking.",
		Back:       "← Back to the blog",
		ReadMore:   "Read →",
		ImageBy:    "Image:",
		months:     [12]string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"},
		date:       "%s %d, %d",
	},
	"ru": {
		Lang:       "ru",
		BasePath:   "/ru/blog",
		IndexTitle: "Блог — Антон Зекериев Родин",
		IndexDesc:  "Статьи Антона Зекериева Родина о разработке, сертификациях и обо всём, что придёт в голову.",
		IndexH1:    "Блог",
		IndexIntro: "То, чему я научился — обычно на собственных ошибках. Записано, чтобы вам не пришлось их повторять.",
		Empty:      "Статей пока нет. Что-то готовится.",
		Back:       "← Назад в блог",
		ReadMore:   "Читать →",
		ImageBy:    "Изображение:",
		months:     [12]string{"января", "февраля", "марта", "апреля", "мая", "июня", "июля", "августа", "сентября", "октября", "ноября", "декабря"},
		date:       "%d %s %d г.",
	},
}

// GetBlog devuelve los textos del blog para un idioma.
func GetBlog(lang string) (BlogText, bool) {
	t, ok := blog[lang]
	return t, ok
}
