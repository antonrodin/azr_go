package handlers

import (
	"net/http"
	"strings"

	"github.com/antonrodin/azr/internal/models"
	"github.com/antonrodin/azr/internal/openai"
	"github.com/go-chi/chi/v5"
)

// Claves estables de los prompts editables.
const (
	PromptTranslate = "translate"
	PromptDictate   = "dictate"
)

// PromptRegistry son los prompts que la app conoce, con sus textos por defecto.
// Sirve para (1) sembrar la BD al arrancar y (2) restaurar valores por defecto.
// Añadir un prompt nuevo = añadir aquí su entrada y usar su clave en el código.
var PromptRegistry = []models.Prompt{
	{
		Key:         PromptTranslate,
		Label:       "Traducción (ES → EN/RU)",
		Description:  "Cómo se traducen título y descripción al inglés y al ruso cuando guardas una foto.",
		Content:     openai.DefaultTranslateInstructions,
	},
	{
		Key:         PromptDictate,
		Label:       "Dictado (voz → título + descripción)",
		Description:  "Cómo se convierte lo que dictas por voz en un título y una descripción en español.",
		Content:     openai.DefaultComposeInstructions,
	},
}

// promptOr devuelve el contenido del prompt key, o def si no hay BD, falla o
// está vacío. Así la IA nunca se queda sin instrucciones.
func (app *AppRepository) promptOr(key, def string) string {
	if app.Prompts == nil {
		return def
	}
	c, err := app.Prompts.Get(key)
	if err != nil || strings.TrimSpace(c) == "" {
		return def
	}
	return c
}

type promptsPage struct {
	adminPage
	Prompts []models.Prompt
}

func (app *AppRepository) PromptsList(w http.ResponseWriter, r *http.Request) {
	list, err := app.Prompts.All()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	app.render(w, "admin/prompts.go.tmpl", promptsPage{
		adminPage: adminPage{Lang: "es", Flash: app.Session.PopString(r.Context(), "flash")},
		Prompts:   list,
	})
}

func (app *AppRepository) PromptUpdate(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	content := strings.TrimSpace(r.PostFormValue("content"))

	if err := app.Prompts.Update(key, content); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	app.Session.Put(r.Context(), "flash", "Prompt guardado.")
	http.Redirect(w, r, "/admin/prompts", http.StatusSeeOther)
}

func (app *AppRepository) PromptReset(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	for _, d := range PromptRegistry {
		if d.Key == key {
			if err := app.Prompts.Update(key, d.Content); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			break
		}
	}
	app.Session.Put(r.Context(), "flash", "Prompt restaurado al valor por defecto.")
	http.Redirect(w, r, "/admin/prompts", http.StatusSeeOther)
}
