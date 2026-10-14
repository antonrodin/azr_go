// Package blog carga los artículos del blog desde ficheros Markdown.
//
// Cada artículo vive en content/blog/<clave>.<idioma>.md con un frontmatter
// sencillo (líneas "campo: valor" entre dos "---"). La clave une las
// traducciones de un mismo artículo (hreflang y conmutador de idioma); el
// slug es propio de cada idioma.
package blog

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

type Article struct {
	Key         string // compartida entre idiomas
	Lang        string
	Slug        string
	Title       string
	Description string
	Date        time.Time
	Image       string // ruta JPG; se espera un .webp hermano
	ImageAlt    string
	ImageCredit string
	Body        template.HTML
}

// ImageWebP es la variante WebP de la imagen de portada.
func (a *Article) ImageWebP() string {
	return strings.TrimSuffix(a.Image, filepath.Ext(a.Image)) + ".webp"
}

// Store guarda los artículos indexados por idioma.
type Store struct {
	bySlug map[string]map[string]*Article // lang → slug → artículo
	byKey  map[string]map[string]*Article // clave → lang → artículo
	list   map[string][]*Article          // lang → artículos, más reciente primero
}

var md = goldmark.New(goldmark.WithExtensions(extension.GFM, extension.Typographer))

// Load lee todos los .md de dir.
func Load(dir string) (*Store, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		return nil, err
	}

	s := &Store{
		bySlug: map[string]map[string]*Article{},
		byKey:  map[string]map[string]*Article{},
		list:   map[string][]*Article{},
	}
	for _, f := range files {
		a, err := parseFile(f)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", f, err)
		}
		if s.bySlug[a.Lang] == nil {
			s.bySlug[a.Lang] = map[string]*Article{}
		}
		if _, dup := s.bySlug[a.Lang][a.Slug]; dup {
			return nil, fmt.Errorf("%s: slug %q repetido en %s", f, a.Slug, a.Lang)
		}
		s.bySlug[a.Lang][a.Slug] = a
		if s.byKey[a.Key] == nil {
			s.byKey[a.Key] = map[string]*Article{}
		}
		s.byKey[a.Key][a.Lang] = a
		s.list[a.Lang] = append(s.list[a.Lang], a)
	}
	for _, l := range s.list {
		sort.Slice(l, func(i, j int) bool { return l[i].Date.After(l[j].Date) })
	}
	return s, nil
}

// Get devuelve el artículo de un idioma por su slug.
func (s *Store) Get(lang, slug string) (*Article, bool) {
	a, ok := s.bySlug[lang][slug]
	return a, ok
}

// Translation devuelve la versión de un artículo en otro idioma, si existe.
func (s *Store) Translation(a *Article, lang string) (*Article, bool) {
	t, ok := s.byKey[a.Key][lang]
	return t, ok
}

// List devuelve los artículos de un idioma, del más reciente al más antiguo.
func (s *Store) List(lang string) []*Article {
	return s.list[lang]
}

// parseFile lee <clave>.<idioma>.md.
func parseFile(path string) (*Article, error) {
	name := strings.TrimSuffix(filepath.Base(path), ".md")
	dot := strings.LastIndex(name, ".")
	if dot < 1 {
		return nil, fmt.Errorf("el nombre debe ser <clave>.<idioma>.md")
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	meta, body, err := splitFrontmatter(raw)
	if err != nil {
		return nil, err
	}

	a := &Article{
		Key:         name[:dot],
		Lang:        name[dot+1:],
		Slug:        meta["slug"],
		Title:       meta["title"],
		Description: meta["description"],
		Image:       meta["image"],
		ImageAlt:    meta["image_alt"],
		ImageCredit: meta["image_credit"],
	}
	if a.Slug == "" || a.Title == "" {
		return nil, fmt.Errorf("faltan slug o title en el frontmatter")
	}
	if a.Date, err = time.Parse("2006-01-02", meta["date"]); err != nil {
		return nil, fmt.Errorf("date: %w", err)
	}

	var buf bytes.Buffer
	if err := md.Convert(body, &buf); err != nil {
		return nil, err
	}
	// El Markdown lo escribimos nosotros: es contenido de confianza.
	a.Body = template.HTML(buf.String())
	return a, nil
}

func splitFrontmatter(raw []byte) (map[string]string, []byte, error) {
	const sep = "---"
	sc := bufio.NewScanner(bytes.NewReader(raw))
	if !sc.Scan() || strings.TrimSpace(sc.Text()) != sep {
		return nil, nil, fmt.Errorf("falta el frontmatter (---)")
	}
	meta := map[string]string{}
	n := len(sc.Bytes()) + 1
	for sc.Scan() {
		line := sc.Text()
		n += len(sc.Bytes()) + 1
		if strings.TrimSpace(line) == sep {
			return meta, raw[min(n, len(raw)):], nil
		}
		if k, v, ok := strings.Cut(line, ":"); ok {
			meta[strings.TrimSpace(k)] = strings.TrimSpace(v)
		}
	}
	return nil, nil, fmt.Errorf("frontmatter sin cerrar")
}
