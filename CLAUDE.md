# CLAUDE.md — azr_go

Web personal de **azr.es** (sitio de Anton Zekeriev Rodin). Aplicación web en **Go** con renderizado server-side de HTML.

> **Fotolog (en vivo, fases 0-4 hechas):** galería pública CC BY en 3 idiomas (`/fotolog`,
> `/{lang}/fotolog`, paginación `/fotolog/pagina/{n}`, detalle `/fotolog/{slug}`) + zona
> privada `/admin` (login contra el servicio auth, subida desde móvil, pipeline
> exiftool→vips→S3, traducción ES→EN/RU y dictado por voz con OpenAI, prompts editables en
> `/admin/prompts`). Plan, decisiones y **avisos de infra** (bucket dev vs prod, nginx
> `client_max_body_size`) en **`FOTOLOG-PLAN.md`** — léelo antes de tocar el fotolog.

> **Blog (en vivo desde 2026-10-14):** `/{lang}/blog` y `/{lang}/blog/{slug}` (idioma siempre
> en la ruta; `/blog` y `/blog/{slug}` → 301 a `/es/…`). Sin BD ni admin: cada artículo es
> `content/blog/<clave>.<idioma>.md` con frontmatter (`title`, `slug`, `description`, `date`
> AAAA-MM-DD, `image` JPG con `.webp` hermano en `public/img/blog/`, `image_alt`,
> `image_credit`). La `<clave>` une las traducciones (hreflang, conmutador); el slug es
> propio de cada idioma. Markdown → HTML con goldmark (`internal/blog`). Se carga al
> arrancar (en prod hay que reiniciar para ver un artículo nuevo; en dev se recarga en cada
> petición). Textos de interfaz en `internal/i18n/blog.go`. Entra solo en el sitemap.
> Primer artículo: `ckad` (experiencia CKAD).

> **Portada:** landing de una sola página en **tres idiomas**: español (por defecto, ruta `/`), inglés (`/en`) y ruso (`/ru`). Diseño minimalista **negro sobre blanco**, dos columnas: izquierda branding (monograma AZR, foto, nombre, claims, redes — GitHub, YouTube, LinkedIn, Credly, blog, email — y selector de idioma), derecha "Skills" con tono **sarcástico**.
>
> **i18n:** los textos viven en `internal/i18n/i18n.go` (mapa `Content` por idioma; para tocar copy se edita ahí, no la plantilla). El handler `Home` resuelve `{lang}` del path y renderiza `home.go.tmpl` con el `Content` del idioma; idioma no soportado → 404. La plantilla emite `canonical` + `hreflang` (x-default → español). Los enlaces sociales están hardcodeados en la plantilla, compartidos por los tres idiomas. La portada lleva JSON-LD `Person` (`PersonName`/`JobTitle` en i18n) y `rel="me"` en los enlaces sociales. La portada no lee la BD; el fotolog sí (`photos`, `prompts`, `sessions`). El CRUD de `pages` existe en el código pero no está cableado a rutas.
>
> **Roadmap (deseado por Anton):** añadir Currículum, Proyectos, Experiencia. Mantener el tono sarcástico y el estilo negro sobre blanco, muy legible. Anton quiere potenciar sus certificaciones (perfil Credly enlazado desde la home).

> El `README.md` está copiado de otro proyecto ("Circuitos / Nürburgring") y se usa como plantilla de despliegue; ignora sus textos, no son de este sitio. El stack y los comandos sí aplican.

## Stack

- **Go 1.23.4**
- **Router:** `go-chi/chi v5` sobre `net/http`
- **Base de datos:** **SQLite** vía `mattn/go-sqlite3` (requiere **CGO**), fichero en `DB_FILE` (`./database/app.db`)
- **Sesiones:** `alexedwards/scs/v2` con store **SQLite** (`scs/sqlite3store`, tabla `sessions` en la misma BD, creada al arrancar) → sobreviven a los reinicios. Vida 30 días, cookie `Secure` solo con `ENV=prod`
- **Auth:** login vía `internal/authsvc` contra `AUTH_SERVICE` (`https://auth.microservices.azr.es`, servicio `azr`); el JWT se guarda en sesión y el middleware `Auth` lo valida con `golang-jwt/jwt/v5` y `JWT_SECRET` (leído de `.env`)
- **Fotos:** AWS SDK v2 (S3), `vips` y `exiftool` instalados en el sistema; **OpenAI** opcional (`internal/openai`) para traducir y transcribir
- **Utilidades:** `gosimple/slug` (URLs), `joho/godotenv` (`.env`)

## Estructura

- `cmd/web/` — entrypoint
  - `main.go` — `run()`: carga `.env`, abre SQLite, configura sesiones, monta `http.Server` en `:PORT`
  - `routes.go` — rutas chi: home `/` y `/{lang}`, `robots.txt`, `sitemap.xml`, fotolog público, `/login`, `/admin/*` (con `Auth`) y estáticos `/public/*`
  - `middleware.go` — `LoadSession` (siempre activo) y `Auth` (JWT; redirige a `/login?next=…` si no hay token)
- `internal/handlers/` — controladores; `App` es un singleton `*AppRepository` inicializado con `NewRepo()`. `handlers.go` (Home, Robots, Sitemap), `fotolog_public.go`, `fotolog_admin.go` (login, subida, edición, publicar, transcribir), `prompts.go`. Todo se pinta con `app.render()` (renderiza a buffer; error → 500, nunca tumba el proceso). `Contacto`/`Cookies`/`Legal` apuntan a plantillas que no existen y no tienen ruta (código muerto)
- `internal/authsvc/` — cliente HTTP del servicio auth (login → JWT)
- `internal/photos/` — EXIF, derivados con vips, subida a S3
- `internal/openai/` — traducción y transcripción
- `internal/blog/` — carga y renderiza los artículos de `content/blog/`
- `internal/i18n/fotolog.go` — textos del fotolog
- `internal/i18n/i18n.go` — struct `Content` y traducciones es/en/ru de toda la home (meta/SEO incluidos); `i18n.Get(lang)` y `i18n.Default = "es"`
- `internal/models/` — `models.go` define `Page`; `mysqlite/` el acceso a datos
  - `mysqlite/photos.go`, `mysqlite/prompts.go` — modelos del fotolog (crean su esquema con `EnsureSchema`)
  - `mysqlite/pages.go` — `PageModel`: `Insert`, `Update`, `GetBySlug`, `All` sobre tabla `pages` (id, title, slug, published, content, createdAt, updatedAt) → CRUD de páginas
  - `mysqlite/db.go` — conexión
- `internal/tmpl/tmpl.go` — `tmpl.Get("nombre.go.tmpl")` devuelve `(*Template, error)` con partials + layout. Con `tmpl.Cache = true` (se activa si `ENV=prod`) cachea; en dev reparsea en cada petición
- `templates/` — `home.go.tmpl`, `fotolog/`, `admin/`, `layouts/app.go.tmpl`, `pages/` (sin uso), `partials/` (header, sidebar, menu, footer, `brand` = barra lateral compartida portada/fotolog)
- `public/css/main.css` — estáticos
- `database/app.db` — SQLite **no versionado** (en `.gitignore`). Guarda fotos, prompts y sesiones; las tablas se crean al arrancar.

## Desarrollo local

```bash
cp .env.sample .env          # editar: descomentar bloque "Local" si procede
CGO_ENABLED=1 go run ./cmd/web
```

Requiere un compilador C por CGO/SQLite (en macOS: Xcode Command Line Tools). Abrir `http://localhost:$PORT`.

### Variables de entorno (`.env`)

- `PORT` — puerto HTTP (sample: `3335`; **producción real: 3331**)
- `DB_FILE` — ruta SQLite (`./database/app.db`)
- `AUTH_SERVICE` — microservicio de auth
- `JWT_SECRET` — debe coincidir con el del servicio de auth (obligatorio; sin él no arranca)
- `ENV` — `prod` activa cookie `Secure` y caché de plantillas
- `AWS_REGION`, `S3_BUCKET` + credenciales — almacenamiento del fotolog (**local → `azr-fotolog-dev`**, prod → `azr-fotolog`)
- `OPENAI_API_KEY` (opcional), `OPENAI_MODEL`, `OPENAI_TRANSCRIBE_MODEL`

## Build y despliegue

```bash
# Binario Linux (servidor)
env GOOS=linux CGO_ENABLED=1 go build -o app ./cmd/web

# Binario macOS
env GOOS=darwin CGO_ENABLED=1 go build -o bin/app ./cmd/web
```

En el servidor (Ubuntu, host `gosha.vservers.es` / 188.164.198.7 — VPS nuevo tras la migración,
ver `../VPS-MIGRACION-PLAN.md`; el viejo 31.24.41.84 ya no responde):
- Acceso: `ssh anton@188.164.198.7` (de momento por contraseña, sin clave del Mac instalada).
- Código en `/home/anton/www/azr` (clon de `antonrodin/azr_go`), `.env` propio, BD en `database/app.db`.
- Binario `app` **compilado en el propio servidor**; corre en **localhost:3331**.
- Go está en `/usr/local/go/bin` (no en el PATH de sesiones no interactivas).
- **supervisor**: programa `azr` (`/etc/supervisor/conf.d/azr.conf`), logs en `~/www/azr/logs/supervisor.{log,err}`. `supervisorctl` necesita sudo.
- **nginx** hace proxy con TLS (Let's Encrypt) → `https://azr.es` y `www.azr.es`, con `client_max_body_size 100M`.

Despliegue (tras hacer push a `main`):

```bash
ssh anton@188.164.198.7
cd ~/www/azr && git pull
/usr/local/go/bin/go build -o app ./cmd/web
sudo supervisorctl restart azr
```

## Notas y avisos

- **`database/app.db` ya NO está versionado** (se sacó con `git rm --cached` y se añadió a `.gitignore`). El `.db` que había en el repo era **sobras de otro proyecto** ("Circuitos/Nürburgring"): tablas `cars`, `records`, `car_descriptions` (216+ filas de récords de circuito) más datos de demo en `pages`. **No había datos reales de azr.es.**
- `.env.sample` está incompleto (faltan `ENV`, S3 y OpenAI) y aún trae `JWT_SECRET=holacaracola` de ejemplo.
- Los formularios de `/admin` no llevan token CSRF; los protege la cookie `SameSite=Lax`.
- No hay tests.
