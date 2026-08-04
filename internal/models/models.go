package models

import (
	"time"
)

type Page struct {
	ID        int
	Published int
	Title     string
	Slug      string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Photo es una foto del fotolog. Los títulos y descripciones van en columnas
// por idioma (es/en/ru). Status es "private" o "public" y coincide con el
// prefijo del bucket S3 donde viven sus ficheros.
type Photo struct {
	ID         int
	Slug       string
	Status     string // "private" | "public"
	TitleES    string
	TitleEN    string
	TitleRU    string
	DescES     string
	DescEN     string
	DescRU     string
	TakenAt    string // ISO8601 o "" si el EXIF no lo trae
	UploadedAt string
	CameraMake string
	CameraModel string
	Lens       string
	Focal      string
	Aperture   string
	Shutter    string
	ISO        string
	GPSLat     float64 // 0 si no hay GPS (ver HasGPS)
	GPSLng     float64
	HasGPS     bool
	ShowGPS    bool // en fotos públicas el GPS va oculto salvo que se active
	Width      int
	Height     int
	OrigExt    string // extensión del original subido (jpg, heic…)
	OrigBytes  int64
	ExifJSON   string // volcado completo de exiftool -json
}

// Prompt son las instrucciones editables que se pasan a la IA (traducción,
// dictado…). Se guardan en SQLite y se pueden editar desde /admin/prompts.
// Content son SOLO las instrucciones de estilo/reglas; el contrato de salida
// (qué JSON devolver) lo añade el código, no es editable.
type Prompt struct {
	Key         string // identificador estable (translate, dictate…)
	Label       string // nombre para la UI
	Description string // ayuda para la UI
	Content     string // instrucciones editables
	UpdatedAt   string
}
