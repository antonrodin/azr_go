package mysqlite

import (
	"database/sql"

	"github.com/antonrodin/azr/internal/models"
)

type PageModel struct {
	DB *sql.DB
}

func (m *PageModel) Insert(title, slug, content string) error {
	stmt := `INSERT INTO pages (title, slug, content) VALUES (?, ?, ?)`

	_, err := m.DB.Exec(stmt, title, slug, content)

	if err != nil {
		return err
	}

	return nil
}

func (m *PageModel) Update(title, slug, content string, published int) error {
	sql := `UPDATE pages 
			SET title = ?, slug = ?, content = ?, published = ?, updatedAt = CURRENT_TIMESTAMP
			WHERE slug = ?`

	_, err := m.DB.Exec(sql, title, slug, content, published, slug)

	if err != nil {
		return err
	}

	return nil
}

func (m *PageModel) GetBySlug(slug string) (*models.Page, error) {
	stmt := `
		SELECT id, title, slug, published, content, createdAt, updatedAt
		FROM pages
		WHERE slug = ?
		`

	page := &models.Page{}
	row := m.DB.QueryRow(stmt, slug)
	err := row.Scan(&page.ID, &page.Title, &page.Slug, &page.Published, &page.Content, &page.CreatedAt, &page.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return page, nil
}

func (m *PageModel) All() ([]models.Page, error) {
	sttm := `
		SELECT id, title, slug, published, content, createdAt, updatedAt 
		FROM pages 
		ORDER BY id DESC
		`

	rows, err := m.DB.Query(sttm)
	if err != nil {
		return nil, err
	}

	pages := []models.Page{}
	for rows.Next() {
		page := models.Page{}
		err := rows.Scan(&page.ID, &page.Title, &page.Slug, &page.Published, &page.Content, &page.CreatedAt, &page.UpdatedAt)
		if err != nil {
			return nil, err
		}

		pages = append(pages, page)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return pages, nil
}
