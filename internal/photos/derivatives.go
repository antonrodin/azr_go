package photos

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
)

// Derivative es un fichero generado (o el original) listo para subir a S3.
type Derivative struct {
	Name        string // nombre del objeto dentro de {prefijo}/{id}/, p.ej. "1200.webp"
	Path        string // ruta local temporal
	ContentType string
}

// Los tamaños web (WebP) y de descarga (JPEG). El original se sube aparte.
var variants = []struct {
	width   int
	ext     string
	quality int
}{
	{2048, "webp", 82},
	{1200, "webp", 82},
	{400, "webp", 75},
	{2048, "jpg", 85},
	{1200, "jpg", 85},
}

// ContentTypeByExt mapea extensiones (sin punto, en minúsculas) a MIME.
func ContentTypeByExt(ext string) string {
	switch ext {
	case "jpg", "jpeg":
		return "image/jpeg"
	case "png":
		return "image/png"
	case "webp":
		return "image/webp"
	case "heic", "heif":
		return "image/heic"
	default:
		return "application/octet-stream"
	}
}

// BuildDerivatives genera con vips los derivados de src dentro de dir.
// vips thumbnail respeta la orientación EXIF (auto-rotate), no amplía
// (--size down) y con strip no arrastra los metadatos a los derivados.
func BuildDerivatives(src, dir string) ([]Derivative, error) {
	var out []Derivative

	for _, v := range variants {
		name := fmt.Sprintf("%d.%s", v.width, v.ext)
		dst := filepath.Join(dir, name)
		target := fmt.Sprintf("%s[Q=%d,strip]", dst, v.quality)

		cmd := exec.Command("vips", "thumbnail", src, target,
			strconv.Itoa(v.width),
			"--height", strconv.Itoa(v.width),
			"--size", "down")
		if b, err := cmd.CombinedOutput(); err != nil {
			return nil, fmt.Errorf("vips thumbnail %s: %w: %s", name, err, string(b))
		}

		out = append(out, Derivative{
			Name:        name,
			Path:        dst,
			ContentType: ContentTypeByExt(v.ext),
		})
	}

	return out, nil
}
