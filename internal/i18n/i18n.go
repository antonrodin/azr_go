// Package i18n contiene las traducciones de la landing en los idiomas
// soportados. El español es el idioma por defecto (ruta "/").
package i18n

// Default es el idioma que se sirve en la raíz "/".
const Default = "es"

// Content agrupa todos los textos traducibles de la home.
// Los fragmentos "Aside" son las apostillas sarcásticas (en cursiva gris).
type Content struct {
	Lang string // código ISO del idioma (es, en, ru)
	Path string // ruta canónica de esta versión (/, /en, /ru)

	// Meta / SEO
	MetaTitle       string
	MetaDescription string
	OgTitle         string
	OgDescription   string

	// Claims (columna de branding)
	ClaimFather     string
	ClaimSWE        string
	ClaimAI         string
	ClaimAIAside    string
	ClaimWorks      string
	ClaimWorksAside string

	// Skills
	SkillsHeading string
	SkillsNote    string

	Skill1Title string
	Skill1Aside string
	Skill1Body  string

	BackendTitle string
	BackendBody  string

	DevOpsTitle string
	DevOpsBody  string
	DevOpsAside string

	FrontendTitle string
	FrontendBody  string

	AITitle string
	AIBody  string
	AIAside string

	OtherTitle string
	OtherBody  string
	OtherAside string

	// Colofón
	ColophonPre  string
	ColophonPost string
}

var content = map[string]Content{
	"es": {
		Lang: "es", Path: "/",
		MetaTitle:       "Antón Zekeriev Rodin — Ingeniero del Software",
		MetaDescription: "Antón Zekeriev Rodin. Padre, ingeniero del software. Hago cosas que funcionan. Go, TypeScript, JavaScript, PHP y lo que haga falta.",
		OgTitle:         "Antón Zekeriev Rodin — Ingeniero del Software",
		OgDescription:   "Padre. Ingeniero del software. Hago cosas que funcionan.",
		ClaimFather:     "Padre.",
		ClaimSWE:        "Ingeniero del Software.",
		ClaimAI:         "Ingeniero de IA",
		ClaimAIAside:    "(esta web está vibecodeada, así que cuenta)",
		ClaimWorks:      "Hago cosas que funcionan",
		ClaimWorksAside: "(la mayoría de las veces)",
		SkillsHeading:   "Skills",
		SkillsNote:      "Por orden de importancia. Spoiler: la primera es la única que de verdad cuenta.",
		Skill1Title:     "Skill nº 1",
		Skill1Aside:     "— la importante",
		Skill1Body:      "Capacidad de aprender y adaptarme a cualquier cosa. Aprendo hoy lo que me hará falta mañana. Las demás de esta lista son, básicamente, esta misma aplicada.",
		BackendTitle:    "Backend",
		BackendBody:     "Go, TypeScript, JavaScript, PHP… lo que haga falta.",
		DevOpsTitle:     "DevOps",
		DevOpsBody:      "Kubernetes, Prometheus, Nginx, AWS… he usado un poco de todo",
		DevOpsAside:     "(lo justo para romper producción con conocimiento de causa)",
		FrontendTitle:   "Frontend",
		FrontendBody:    "HTML, CSS, BEM, Tailwind, Vue, React Native… etc.",
		AITitle:         "IA / Vibecoding",
		AIBody:          "Sí, soy Ingeniero de IA. ¿Pruebas? Esta misma web está vibecodeada de arriba abajo. Yo pongo el criterio; la IA, la mecanografía",
		AIAside:         "(y alguna que otra broma)",
		OtherTitle:      "Otras tonterías",
		OtherBody:       "Agile, MySQL, PostgreSQL, MongoDB, Jira, GitHub, Microservicios…",
		OtherAside:      "y los buzzwords que toquen ese trimestre",
		ColophonPre:     "Esta página está vibecodeada con",
		ColophonPost:    ". Ni una línea escrita a mano, y mira qué legible ha quedado.",
	},
	"en": {
		Lang: "en", Path: "/en",
		MetaTitle:       "Antón Zekeriev Rodin — Software Engineer",
		MetaDescription: "Antón Zekeriev Rodin. Father, software engineer. I build things that work. Go, TypeScript, JavaScript, PHP and whatever it takes.",
		OgTitle:         "Antón Zekeriev Rodin — Software Engineer",
		OgDescription:   "Father. Software engineer. I build things that work.",
		ClaimFather:     "Father.",
		ClaimSWE:        "Software Engineer.",
		ClaimAI:         "AI Engineer",
		ClaimAIAside:    "(this site is vibecoded, so it counts)",
		ClaimWorks:      "I build things that work",
		ClaimWorksAside: "(most of the time)",
		SkillsHeading:   "Skills",
		SkillsNote:      "In order of importance. Spoiler: the first one is the only one that really matters.",
		Skill1Title:     "Skill #1",
		Skill1Aside:     "— the important one",
		Skill1Body:      "The ability to learn and adapt to anything. I learn today what I'll need tomorrow. The rest of this list is basically that, applied.",
		BackendTitle:    "Backend",
		BackendBody:     "Go, TypeScript, JavaScript, PHP… whatever it takes.",
		DevOpsTitle:     "DevOps",
		DevOpsBody:      "Kubernetes, Prometheus, Nginx, AWS… a bit of everything",
		DevOpsAside:     "(just enough to break production knowingly)",
		FrontendTitle:   "Frontend",
		FrontendBody:    "HTML, CSS, BEM, Tailwind, Vue, React Native… etc.",
		AITitle:         "AI / Vibecoding",
		AIBody:          "Yes, I'm an AI Engineer. Proof? This very site is vibecoded top to bottom. I bring the judgement; the AI brings the typing",
		AIAside:         "(and the occasional joke)",
		OtherTitle:      "Other nonsense",
		OtherBody:       "Agile, MySQL, PostgreSQL, MongoDB, Jira, GitHub, Microservices…",
		OtherAside:      "and whatever buzzwords are trending this quarter",
		ColophonPre:     "This page is vibecoded with",
		ColophonPost:    ". Not a single line written by hand, and look how readable it turned out.",
	},
	"ru": {
		Lang: "ru", Path: "/ru",
		MetaTitle:       "Антон Зекериев Родин — Инженер-программист",
		MetaDescription: "Антон Зекериев Родин. Отец, инженер-программист. Делаю вещи, которые работают. Go, TypeScript, JavaScript, PHP и что угодно.",
		OgTitle:         "Антон Зекериев Родин — Инженер-программист",
		OgDescription:   "Отец. Инженер-программист. Делаю вещи, которые работают.",
		ClaimFather:     "Отец.",
		ClaimSWE:        "Инженер-программист.",
		ClaimAI:         "ИИ-инженер",
		ClaimAIAside:    "(этот сайт сделан вайбкодингом, так что засчитывается)",
		ClaimWorks:      "Делаю вещи, которые работают",
		ClaimWorksAside: "(в большинстве случаев)",
		SkillsHeading:   "Навыки",
		SkillsNote:      "В порядке важности. Спойлер: по-настоящему важен только первый.",
		Skill1Title:     "Навык №1",
		Skill1Aside:     "— самый важный",
		Skill1Body:      "Умение учиться и адаптироваться к чему угодно. Сегодня учу то, что понадобится завтра. Остальное в этом списке — по сути, то же самое на практике.",
		BackendTitle:    "Бэкенд",
		BackendBody:     "Go, TypeScript, JavaScript, PHP… что угодно.",
		DevOpsTitle:     "DevOps",
		DevOpsBody:      "Kubernetes, Prometheus, Nginx, AWS… всего понемногу",
		DevOpsAside:     "(ровно столько, чтобы ронять прод осознанно)",
		FrontendTitle:   "Фронтенд",
		FrontendBody:    "HTML, CSS, BEM, Tailwind, Vue, React Native… и т. д.",
		AITitle:         "ИИ / Вайбкодинг",
		AIBody:          "Да, я ИИ-инженер. Доказательства? Этот самый сайт сделан вайбкодингом от и до. Я задаю критерии, ИИ печатает",
		AIAside:         "(и иногда шутит)",
		OtherTitle:      "Прочая ерунда",
		OtherBody:       "Agile, MySQL, PostgreSQL, MongoDB, Jira, GitHub, микросервисы…",
		OtherAside:      "и любые баззворды этого квартала",
		ColophonPre:     "Эта страница сделана вайбкодингом с помощью",
		ColophonPost:    ". Ни строчки вручную — и посмотрите, как читабельно получилось.",
	},
}

// Get devuelve el contenido del idioma indicado y si está soportado.
func Get(lang string) (Content, bool) {
	c, ok := content[lang]
	return c, ok
}
