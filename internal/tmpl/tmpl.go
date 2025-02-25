package tmpl

import (
	"html/template"
	"log"
	"strings"
)

var funcMap = template.FuncMap{
	"toUpper":   toUpper,
	"unescaped": unescaped,
}

func toUpper(str string) string {
	return strings.ToUpper(str)
}

func unescaped(str string) template.HTML {
	return template.HTML(str)
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
