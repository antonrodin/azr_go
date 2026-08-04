package photos

import (
	"encoding/json"
	"fmt"
	"math"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Meta son los metadatos que extraemos de una foto con exiftool.
type Meta struct {
	TakenAt     string // ISO8601 o ""
	CameraMake  string
	CameraModel string
	Lens        string
	Focal       string // "26 mm"
	Aperture    string // "f/1.8"
	Shutter     string // "1/120"
	ISO         string
	HasGPS      bool
	GPSLat      float64
	GPSLng      float64
	Width       int
	Height      int
	ExifJSON    string // volcado completo (valores legibles) para guardar
}

// numeric es el subconjunto que pedimos a exiftool con -n (valores numéricos,
// fiables para parsear). El volcado completo se hace aparte sin -n para que
// los valores guardados sean legibles tal cual.
type numeric struct {
	DateTimeOriginal string      `json:"DateTimeOriginal"`
	CreateDate       string      `json:"CreateDate"`
	OffsetTimeOriginal string    `json:"OffsetTimeOriginal"`
	Make             string      `json:"Make"`
	Model            string      `json:"Model"`
	LensModel        string      `json:"LensModel"`
	LensID           string      `json:"LensID"`
	FocalLength      json.Number `json:"FocalLength"`
	FNumber          json.Number `json:"FNumber"`
	ExposureTime     json.Number `json:"ExposureTime"`
	ISO              json.Number `json:"ISO"`
	GPSLatitude      json.Number `json:"GPSLatitude"`
	GPSLongitude     json.Number `json:"GPSLongitude"`
	GPSLatitudeRef   string      `json:"GPSLatitudeRef"`
	GPSLongitudeRef  string      `json:"GPSLongitudeRef"`
	ImageWidth       int         `json:"ImageWidth"`
	ImageHeight      int         `json:"ImageHeight"`
	Orientation      int         `json:"Orientation"`
}

func runExiftool(args ...string) ([]byte, error) {
	out, err := exec.Command("exiftool", args...).Output()
	if err != nil {
		return nil, fmt.Errorf("exiftool %s: %w", strings.Join(args, " "), err)
	}
	return out, nil
}

// ExtractMeta lee los metadatos de la foto en path. Nunca falla por metadatos
// incompletos: campos ausentes quedan vacíos.
func ExtractMeta(path string) (*Meta, error) {
	// Volcado completo legible (se guarda tal cual en la BD).
	full, err := runExiftool("-json", "-g0", path)
	if err != nil {
		return nil, err
	}

	// Campos concretos en modo numérico para parsear con fiabilidad.
	numOut, err := runExiftool("-json", "-n",
		"-DateTimeOriginal", "-CreateDate", "-OffsetTimeOriginal",
		"-Make", "-Model", "-LensModel", "-LensID",
		"-FocalLength", "-FNumber", "-ExposureTime", "-ISO",
		"-GPSLatitude", "-GPSLongitude", "-GPSLatitudeRef", "-GPSLongitudeRef",
		"-ImageWidth", "-ImageHeight", "-Orientation", path)
	if err != nil {
		return nil, err
	}

	var nums []numeric
	if err := json.Unmarshal(numOut, &nums); err != nil || len(nums) == 0 {
		return nil, fmt.Errorf("parsear salida de exiftool: %w", err)
	}
	n := nums[0]

	m := &Meta{
		CameraMake:  strings.TrimSpace(n.Make),
		CameraModel: strings.TrimSpace(n.Model),
		Lens:        firstNonEmpty(n.LensModel, n.LensID),
		ExifJSON:    string(full),
	}

	// Fecha de captura: exiftool da "2024:06:14 20:19:20".
	m.TakenAt = parseExifDate(firstNonEmpty(n.DateTimeOriginal, n.CreateDate), n.OffsetTimeOriginal)

	if f, ok := numFloat(n.FocalLength); ok {
		m.Focal = strconv.FormatFloat(f, 'f', -1, 64) + " mm"
	}
	if f, ok := numFloat(n.FNumber); ok {
		m.Aperture = "f/" + strconv.FormatFloat(f, 'f', -1, 64)
	}
	if f, ok := numFloat(n.ExposureTime); ok && f > 0 {
		if f < 1 {
			m.Shutter = fmt.Sprintf("1/%d", int(math.Round(1/f)))
		} else {
			m.Shutter = strconv.FormatFloat(f, 'f', -1, 64) + "s"
		}
	}
	if s := n.ISO.String(); s != "" {
		m.ISO = s
	}

	// GPS: con -n exiftool suele dar el compuesto ya con signo, pero por si
	// acaso aplicamos la referencia (S/W = negativo).
	lat, okLat := numFloat(n.GPSLatitude)
	lng, okLng := numFloat(n.GPSLongitude)
	if okLat && okLng {
		if strings.EqualFold(n.GPSLatitudeRef, "S") && lat > 0 {
			lat = -lat
		}
		if strings.EqualFold(n.GPSLongitudeRef, "W") && lng > 0 {
			lng = -lng
		}
		m.HasGPS = true
		m.GPSLat = lat
		m.GPSLng = lng
	}

	// Dimensiones tal y como se ve la foto: las orientaciones EXIF 5-8
	// giran 90º, así que ancho y alto van intercambiados.
	m.Width, m.Height = n.ImageWidth, n.ImageHeight
	if n.Orientation >= 5 && n.Orientation <= 8 {
		m.Width, m.Height = m.Height, m.Width
	}

	return m, nil
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func numFloat(n json.Number) (float64, bool) {
	if n.String() == "" {
		return 0, false
	}
	f, err := n.Float64()
	if err != nil {
		return 0, false
	}
	return f, true
}

// parseExifDate convierte "2024:06:14 20:19:20" (+ offset opcional "+02:00")
// a ISO8601. Devuelve "" si no se puede parsear.
func parseExifDate(v, offset string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}

	if offset != "" {
		if t, err := time.Parse("2006:01:02 15:04:05-07:00", v+offset); err == nil {
			return t.Format(time.RFC3339)
		}
	}
	if t, err := time.Parse("2006:01:02 15:04:05", v); err == nil {
		return t.Format("2006-01-02T15:04:05")
	}
	return ""
}
