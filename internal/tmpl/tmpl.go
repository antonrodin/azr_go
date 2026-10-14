package tmpl

import (
	"html/template"
	"strings"
	"sync"
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

// Cache activa la caché de plantillas parseadas. En prod (ENV=prod) se pone
// a true desde main; en dev se reparsean en cada petición para ver los
// cambios sin reiniciar.
var Cache bool

var (
	mu    sync.RWMutex
	cache = map[string]*template.Template{}
)

// Get devuelve la plantilla path junto con los partials y el layout. Un error
// de parseo se devuelve al llamador (que responde 500) en vez de tumbar el
// servidor.
func Get(path string) (*template.Template, error) {
	if Cache {
		mu.RLock()
		t, ok := cache[path]
		mu.RUnlock()
		if ok {
			return t, nil
		}
	}

	splitedPath := strings.Split(path, "/")
	fileName := splitedPath[len(splitedPath)-1]

	// Parse all templates
	t, err := template.New(fileName).Funcs(funcMap).ParseFiles(
		"./templates/"+path,
		"./templates/partials/menu.go.tmpl",
		"./templates/partials/header.go.tmpl",
		"./templates/partials/footer.go.tmpl",
		"./templates/partials/sidebar.go.tmpl",
		"./templates/partials/brand.go.tmpl",
		"./templates/layouts/app.go.tmpl",
	)
	if err != nil {
		return nil, err
	}

	if Cache {
		mu.Lock()
		cache[path] = t
		mu.Unlock()
	}

	return t, nil
}
