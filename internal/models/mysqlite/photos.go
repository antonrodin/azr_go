package mysqlite

import (
	"database/sql"

	"github.com/antonrodin/azr/internal/models"
)

type PhotoModel struct {
	DB *sql.DB
}

// EnsureSchema crea la tabla photos si no existe. Se llama al arrancar.
func (m *PhotoModel) EnsureSchema() error {
	stmt := `
	CREATE TABLE IF NOT EXISTS photos (
		id           INTEGER PRIMARY KEY AUTOINCREMENT,
		slug         TEXT NOT NULL UNIQUE,
		status       TEXT NOT NULL DEFAULT 'private',
		title_es     TEXT NOT NULL DEFAULT '',
		title_en     TEXT NOT NULL DEFAULT '',
		title_ru     TEXT NOT NULL DEFAULT '',
		desc_es      TEXT NOT NULL DEFAULT '',
		desc_en      TEXT NOT NULL DEFAULT '',
		desc_ru      TEXT NOT NULL DEFAULT '',
		taken_at     TEXT NOT NULL DEFAULT '',
		uploaded_at  TEXT NOT NULL DEFAULT (datetime('now')),
		camera_make  TEXT NOT NULL DEFAULT '',
		camera_model TEXT NOT NULL DEFAULT '',
		lens         TEXT NOT NULL DEFAULT '',
		focal        TEXT NOT NULL DEFAULT '',
		aperture     TEXT NOT NULL DEFAULT '',
		shutter      TEXT NOT NULL DEFAULT '',
		iso          TEXT NOT NULL DEFAULT '',
		gps_lat      REAL,
		gps_lng      REAL,
		show_gps     INTEGER NOT NULL DEFAULT 0,
		width        INTEGER NOT NULL DEFAULT 0,
		height       INTEGER NOT NULL DEFAULT 0,
		orig_ext     TEXT NOT NULL DEFAULT '',
		orig_bytes   INTEGER NOT NULL DEFAULT 0,
		exif_json    TEXT NOT NULL DEFAULT '{}'
	);
	CREATE INDEX IF NOT EXISTS idx_photos_status ON photos (status);
	`
	_, err := m.DB.Exec(stmt)
	return err
}

// Insert crea la fila con los metadatos extraídos y devuelve el id.
// El slug definitivo se pone después (necesita el id para ser único).
func (m *PhotoModel) Insert(p *models.Photo) (int, error) {
	stmt := `INSERT INTO photos
		(slug, status, taken_at, camera_make, camera_model, lens, focal,
		 aperture, shutter, iso, gps_lat, gps_lng, width, height,
		 orig_ext, orig_bytes, exif_json)
		VALUES (?, 'private', ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	var lat, lng interface{}
	if p.HasGPS {
		lat, lng = p.GPSLat, p.GPSLng
	}

	res, err := m.DB.Exec(stmt, p.Slug, p.TakenAt, p.CameraMake, p.CameraModel,
		p.Lens, p.Focal, p.Aperture, p.Shutter, p.ISO, lat, lng,
		p.Width, p.Height, p.OrigExt, p.OrigBytes, p.ExifJSON)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *PhotoModel) UpdateSlug(id int, slug string) error {
	_, err := m.DB.Exec(`UPDATE photos SET slug = ? WHERE id = ?`, slug, id)
	return err
}

// UpdateTexts guarda títulos, descripciones y el toggle de GPS.
func (m *PhotoModel) UpdateTexts(id int, p *models.Photo) error {
	stmt := `UPDATE photos SET
		title_es = ?, title_en = ?, title_ru = ?,
		desc_es = ?, desc_en = ?, desc_ru = ?,
		show_gps = ?, slug = ?
		WHERE id = ?`

	showGPS := 0
	if p.ShowGPS {
		showGPS = 1
	}

	_, err := m.DB.Exec(stmt, p.TitleES, p.TitleEN, p.TitleRU,
		p.DescES, p.DescEN, p.DescRU, showGPS, p.Slug, id)
	return err
}

func (m *PhotoModel) UpdateStatus(id int, status string) error {
	_, err := m.DB.Exec(`UPDATE photos SET status = ? WHERE id = ?`, status, id)
	return err
}

func (m *PhotoModel) Delete(id int) error {
	_, err := m.DB.Exec(`DELETE FROM photos WHERE id = ?`, id)
	return err
}

const photoColumns = `id, slug, status, title_es, title_en, title_ru,
	desc_es, desc_en, desc_ru, taken_at, uploaded_at,
	camera_make, camera_model, lens, focal, aperture, shutter, iso,
	gps_lat, gps_lng, show_gps, width, height, orig_ext, orig_bytes, exif_json`

func (m *PhotoModel) scan(row interface{ Scan(...interface{}) error }) (*models.Photo, error) {
	p := &models.Photo{}
	var lat, lng sql.NullFloat64
	var showGPS int

	err := row.Scan(&p.ID, &p.Slug, &p.Status, &p.TitleES, &p.TitleEN, &p.TitleRU,
		&p.DescES, &p.DescEN, &p.DescRU, &p.TakenAt, &p.UploadedAt,
		&p.CameraMake, &p.CameraModel, &p.Lens, &p.Focal, &p.Aperture,
		&p.Shutter, &p.ISO, &lat, &lng, &showGPS,
		&p.Width, &p.Height, &p.OrigExt, &p.OrigBytes, &p.ExifJSON)
	if err != nil {
		return nil, err
	}

	if lat.Valid && lng.Valid {
		p.HasGPS = true
		p.GPSLat = lat.Float64
		p.GPSLng = lng.Float64
	}
	p.ShowGPS = showGPS == 1

	return p, nil
}

func (m *PhotoModel) Get(id int) (*models.Photo, error) {
	row := m.DB.QueryRow(`SELECT `+photoColumns+` FROM photos WHERE id = ?`, id)
	return m.scan(row)
}

func (m *PhotoModel) GetBySlug(slug string) (*models.Photo, error) {
	row := m.DB.QueryRow(`SELECT `+photoColumns+` FROM photos WHERE slug = ?`, slug)
	return m.scan(row)
}

// All devuelve todas las fotos (privadas y públicas), las más recientes primero.
func (m *PhotoModel) All() ([]*models.Photo, error) {
	rows, err := m.DB.Query(`SELECT ` + photoColumns + ` FROM photos
		ORDER BY CASE WHEN taken_at != '' THEN taken_at ELSE uploaded_at END DESC, id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []*models.Photo
	for rows.Next() {
		p, err := m.scan(rows)
		if err != nil {
			return nil, err
		}
		photos = append(photos, p)
	}

	return photos, rows.Err()
}

// SlugTakenByOther dice si otro registro distinto de id ya usa ese slug.
func (m *PhotoModel) SlugTakenByOther(slug string, id int) (bool, error) {
	var n int
	err := m.DB.QueryRow(`SELECT COUNT(*) FROM photos WHERE slug = ? AND id != ?`, slug, id).Scan(&n)
	return n > 0, err
}

// publicOrder es el criterio de la galería pública: por fecha de captura
// (o de subida si no hay EXIF) descendente.
const publicOrder = ` ORDER BY CASE WHEN taken_at != '' THEN taken_at ELSE uploaded_at END DESC, id DESC`

// CountPublic cuenta las fotos públicas (para paginar).
func (m *PhotoModel) CountPublic() (int, error) {
	var n int
	err := m.DB.QueryRow(`SELECT COUNT(*) FROM photos WHERE status = 'public'`).Scan(&n)
	return n, err
}

// PublicPage devuelve una página de fotos públicas.
func (m *PhotoModel) PublicPage(limit, offset int) ([]*models.Photo, error) {
	rows, err := m.DB.Query(`SELECT `+photoColumns+` FROM photos WHERE status = 'public'`+publicOrder+` LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Photo
	for rows.Next() {
		p, err := m.scan(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// AllPublic devuelve todas las fotos públicas (para el sitemap).
func (m *PhotoModel) AllPublic() ([]*models.Photo, error) {
	rows, err := m.DB.Query(`SELECT ` + photoColumns + ` FROM photos WHERE status = 'public'` + publicOrder)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Photo
	for rows.Next() {
		p, err := m.scan(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// GetPublicBySlug devuelve una foto pública por slug (nil si no existe o es privada).
func (m *PhotoModel) GetPublicBySlug(slug string) (*models.Photo, error) {
	row := m.DB.QueryRow(`SELECT `+photoColumns+` FROM photos WHERE slug = ? AND status = 'public'`, slug)
	return m.scan(row)
}
