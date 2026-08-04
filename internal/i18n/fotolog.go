package i18n

// FotologText son los textos de interfaz del fotolog público (galería y
// detalle) en cada idioma. Los títulos y descripciones de cada foto NO van
// aquí: viven en la BD por idioma. Esto es solo el marco de la página.
type FotologText struct {
	Lang string
	// BasePath es el prefijo de idioma del fotolog: "/fotolog" en español,
	// "/en/fotolog" y "/ru/fotolog" en los demás. Sirve para construir
	// enlaces y URLs canónicas.
	BasePath string

	// Galería
	GalleryTitle string // <title> y og:title
	GalleryDesc  string // meta description
	GalleryH1    string
	GalleryIntro string // apostilla con guasa
	Empty        string // estado vacío
	NoTitle      string // placeholder cuando una foto no tiene título

	// Navegación
	Prev string
	Next string
	Back string // volver a la galería

	// Metadatos (etiquetas de la tabla EXIF)
	MetaHeading  string
	LabelDate    string
	LabelCamera  string
	LabelLens    string
	LabelFocal   string
	LabelApert   string
	LabelShutter string
	LabelISO     string
	LabelDims    string
	LabelGPS     string

	// Descargas
	DownloadsHeading string
	DownloadMedium   string // "Mediana"
	DownloadLarge    string // "Grande"
	DownloadOriginal string // "Original"

	// Licencia
	LicenseHeading string
	LicenseBody    string // texto humano de la licencia (con guasa)
	LicenseHowto   string // cómo atribuir
	Attribution    string // el crédito exacto a copiar
}

// LicenseURL es la licencia Creative Commons Attribution 4.0.
const LicenseURL = "https://creativecommons.org/licenses/by/4.0/"

// LicenseName es el nombre corto de la licencia.
const LicenseName = "CC BY 4.0"

// Author es el crédito de autoría para schema.org y atribución.
const Author = "Antón Zekeriev Rodin"

// AttributionText es el crédito que se pide copiar al reutilizar una foto.
const AttributionText = "Foto: Antón Zekeriev Rodin — azr.es (CC BY 4.0)"

var fotolog = map[string]FotologText{
	"es": {
		Lang:     "es",
		BasePath: "/fotolog",

		GalleryTitle: "Fotolog — Antón Zekeriev Rodin",
		GalleryDesc:  "Fotos de Antón Zekeriev Rodin. Uso libre, incluso comercial, con atribución (CC BY 4.0). Metadatos de cámara y descargas en varios tamaños.",
		GalleryH1:    "Fotolog",
		GalleryIntro: "Fotos hechas casi siempre con el móvil. Sin pretensiones y sin Photoshop. Úsalas si te sirven; solo di de quién son.",
		Empty:        "Todavía no hay fotos públicas. Paciencia, que esto va a ritmo de fin de semana.",
		NoTitle:      "Sin título",

		Prev: "← Anteriores",
		Next: "Siguientes →",
		Back: "← Volver al fotolog",

		MetaHeading:  "Datos de la foto",
		LabelDate:    "Fecha",
		LabelCamera:  "Cámara",
		LabelLens:    "Objetivo",
		LabelFocal:   "Focal",
		LabelApert:   "Apertura",
		LabelShutter: "Exposición",
		LabelISO:     "ISO",
		LabelDims:    "Dimensiones",
		LabelGPS:     "Ubicación",

		DownloadsHeading: "Descargar",
		DownloadMedium:   "Mediana",
		DownloadLarge:    "Grande",
		DownloadOriginal: "Original",

		LicenseHeading: "Licencia",
		LicenseBody:    "Esta foto está bajo Creative Commons Attribution 4.0. Tradúzcase: puedes usarla para lo que quieras, incluso para ganar dinero, sin pedirme permiso.",
		LicenseHowto:   "El único requisito es citar la autoría. Copia esto y quedamos en paz:",
		Attribution:    AttributionText,
	},
	"en": {
		Lang:     "en",
		BasePath: "/en/fotolog",

		GalleryTitle: "Photolog — Antón Zekeriev Rodin",
		GalleryDesc:  "Photos by Antón Zekeriev Rodin. Free to use, even commercially, with attribution (CC BY 4.0). Camera metadata and downloads in several sizes.",
		GalleryH1:    "Photolog",
		GalleryIntro: "Photos taken mostly with a phone. No pretensions, no Photoshop. Use them if they're useful; just say whose they are.",
		Empty:        "No public photos yet. Hang in there — this moves at weekend pace.",
		NoTitle:      "Untitled",

		Prev: "← Newer",
		Next: "Older →",
		Back: "← Back to the photolog",

		MetaHeading:  "Photo data",
		LabelDate:    "Date",
		LabelCamera:  "Camera",
		LabelLens:    "Lens",
		LabelFocal:   "Focal length",
		LabelApert:   "Aperture",
		LabelShutter: "Shutter",
		LabelISO:     "ISO",
		LabelDims:    "Dimensions",
		LabelGPS:     "Location",

		DownloadsHeading: "Download",
		DownloadMedium:   "Medium",
		DownloadLarge:    "Large",
		DownloadOriginal: "Original",

		LicenseHeading: "License",
		LicenseBody:    "This photo is under Creative Commons Attribution 4.0. In plain terms: use it for whatever you like, even to make money, no need to ask me.",
		LicenseHowto:   "The only requirement is crediting the author. Copy this and we're square:",
		Attribution:    AttributionText,
	},
	"ru": {
		Lang:     "ru",
		BasePath: "/ru/fotolog",

		GalleryTitle: "Фотолог — Антон Зекериев Родин",
		GalleryDesc:  "Фотографии Антона Зекериева Родина. Свободное использование, в том числе коммерческое, с указанием авторства (CC BY 4.0). Метаданные камеры и загрузка в нескольких размерах.",
		GalleryH1:    "Фотолог",
		GalleryIntro: "Фотографии, снятые в основном на телефон. Без претензий и без фотошопа. Пользуйтесь, если пригодятся; просто укажите, чьи они.",
		Empty:        "Публичных фотографий пока нет. Терпение — тут всё в темпе выходного дня.",
		NoTitle:      "Без названия",

		Prev: "← Новее",
		Next: "Старее →",
		Back: "← Назад в фотолог",

		MetaHeading:  "Данные снимка",
		LabelDate:    "Дата",
		LabelCamera:  "Камера",
		LabelLens:    "Объектив",
		LabelFocal:   "Фокусное расстояние",
		LabelApert:   "Диафрагма",
		LabelShutter: "Выдержка",
		LabelISO:     "ISO",
		LabelDims:    "Размер",
		LabelGPS:     "Местоположение",

		DownloadsHeading: "Скачать",
		DownloadMedium:   "Средний",
		DownloadLarge:    "Большой",
		DownloadOriginal: "Оригинал",

		LicenseHeading: "Лицензия",
		LicenseBody:    "Эта фотография под лицензией Creative Commons Attribution 4.0. Проще говоря: используйте как хотите, даже для заработка, разрешения спрашивать не нужно.",
		LicenseHowto:   "Единственное условие — указать автора. Скопируйте это, и мы в расчёте:",
		Attribution:    AttributionText,
	},
}

// GetFotolog devuelve los textos del fotolog para el idioma indicado.
func GetFotolog(lang string) (FotologText, bool) {
	t, ok := fotolog[lang]
	return t, ok
}
