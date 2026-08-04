package mysqlite

import (
	"database/sql"

	"github.com/antonrodin/azr/internal/models"
)

type PromptModel struct {
	DB *sql.DB
}

func (m *PromptModel) EnsureSchema() error {
	_, err := m.DB.Exec(`
	CREATE TABLE IF NOT EXISTS prompts (
		key         TEXT PRIMARY KEY,
		label       TEXT NOT NULL DEFAULT '',
		description TEXT NOT NULL DEFAULT '',
		content     TEXT NOT NULL DEFAULT '',
		updated_at  TEXT NOT NULL DEFAULT (datetime('now'))
	);`)
	return err
}

// Seed inserta los prompts que aún no existan; NO pisa el contenido de los que
// ya estén (para no borrar ediciones de Anton). Sí refresca label/description,
// que son metadatos de la interfaz.
func (m *PromptModel) Seed(defaults []models.Prompt) error {
	for _, p := range defaults {
		if _, err := m.DB.Exec(
			`INSERT OR IGNORE INTO prompts (key, label, description, content) VALUES (?, ?, ?, ?)`,
			p.Key, p.Label, p.Description, p.Content); err != nil {
			return err
		}
		if _, err := m.DB.Exec(
			`UPDATE prompts SET label = ?, description = ? WHERE key = ?`,
			p.Label, p.Description, p.Key); err != nil {
			return err
		}
	}
	return nil
}

// Get devuelve el contenido (instrucciones) de un prompt.
func (m *PromptModel) Get(key string) (string, error) {
	var content string
	err := m.DB.QueryRow(`SELECT content FROM prompts WHERE key = ?`, key).Scan(&content)
	return content, err
}

func (m *PromptModel) All() ([]models.Prompt, error) {
	rows, err := m.DB.Query(`SELECT key, label, description, content, updated_at FROM prompts ORDER BY key`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Prompt
	for rows.Next() {
		var p models.Prompt
		if err := rows.Scan(&p.Key, &p.Label, &p.Description, &p.Content, &p.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (m *PromptModel) Update(key, content string) error {
	_, err := m.DB.Exec(`UPDATE prompts SET content = ?, updated_at = datetime('now') WHERE key = ?`, content, key)
	return err
}
