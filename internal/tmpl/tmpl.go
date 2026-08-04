package tmpl

import (
	"html/template"
	"log"
	"strings"
)

var funcMap = template.FuncMap{
	"toUpper":   toUpper,
	"unescaped": unescaped,
	"safeJS":    safeJS,
}

func toUpper(str string) string {
	return strings.ToUpper(str)
}

func unescaped(str string) template.HTML {
	return template.HTML(str)
}

// safeJS marca una cadena como JS ya seguro (p.ej. JSON-LD generado con
// json.Marshal, que escapa < > &). Sin esto, html/template la reescaparía
// como literal de cadena dentro de <script>.
func safeJS(str string) template.JS {
	return template.JS(str)
}

func Get(path string) *template.Template {
	var parsedTemplate *template.Template
	var err error

	splitedPath := strings.Split(path, "/")
	fileName := splitedPath[len(splitedPath)-1]

	parsedTemplate = template.New(fileName).
		Funcs(funcMap)

	// Parse all templates
	parsedTemplate, err = parsedTemplate.ParseFiles(
		"./templates/"+path,
		"./templates/partials/menu.go.tmpl",
		"./templates/partials/header.go.tmpl",
		"./templates/partials/footer.go.tmpl",
		"./templates/partials/sidebar.go.tmpl",
		"./templates/layouts/app.go.tmpl",
	)

	//parsedTemplate = template.Must(template.ParseGlob("./templates/partials/*.go.tmpl"))

	if err != nil {
		log.Fatal(err)
	}

	return parsedTemplate
}
