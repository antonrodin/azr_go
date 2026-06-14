# CLAUDE.md — azr_go

Web personal de **azr.es** (sitio de Anton Zekeriev Rodin). Aplicación web en **Go** con renderizado server-side de HTML.

> **Estado actual:** la web es una landing "en construcción" — solo muestra nombre, mensaje, contacto y enlaces (Github/Youtube/LinkedIn/Blog). El handler `Home` renderiza `home.go.tmpl` con `nil`: **no usa la base de datos**. El CRUD de `pages` y los modelos existen en el código pero no están cableados a rutas.

> El `README.md` está copiado de otro proyecto ("Circuitos / Nürburgring") y se usa como plantilla de despliegue; ignora sus textos, no son de este sitio. El stack y los comandos sí aplican.

## Stack

- **Go 1.23.4**
- **Router:** `go-chi/chi v5` sobre `net/http`
- **Base de datos:** **SQLite** vía `mattn/go-sqlite3` (requiere **CGO**), fichero en `DB_FILE` (`./database/app.db`)
- **Sesiones:** `alexedwards/scs/v2` (cookies, vida 24h) → respuestas con `Vary: Cookie`
- **Auth:** `golang-jwt/jwt/v5` validando un token JWT guardado en sesión, contra un microservicio externo `AUTH_SERVICE` (`https://auth.microservices.azr.es`)
- **Utilidades:** `gosimple/slug` (URLs), `joho/godotenv` (`.env`)

## Estructura

- `cmd/web/` — entrypoint
  - `main.go` — `run()`: carga `.env`, abre SQLite, configura sesiones, monta `http.Server` en `:PORT`
  - `routes.go` — rutas chi; sirve estáticos en `/public/*` desde `./public/`
  - `middleware.go` — `LoadSession` (siempre activo) y `Auth` (JWT; redirige a `/login` si no hay token)
- `internal/handlers/` — controladores; `App` es un singleton `*AppRepository` inicializado con `NewRepo()`. Handlers: `Home`, `Contacto`, `Cookies`, `Legal`
- `internal/models/` — `models.go` define `Page`; `mysqlite/` el acceso a datos
  - `mysqlite/pages.go` — `PageModel`: `Insert`, `Update`, `GetBySlug`, `All` sobre tabla `pages` (id, title, slug, published, content, createdAt, updatedAt) → CRUD de páginas
  - `mysqlite/db.go` — conexión
- `internal/tmpl/tmpl.go` — `tmpl.Get("nombre.go.tmpl")` carga y ejecuta plantillas
- `templates/` — `home.go.tmpl`, `layouts/app.go.tmpl`, `pages/` (index/create/edit/show), `partials/` (header, sidebar, menu, footer)
- `public/css/main.css` — estáticos
- `database/app.db` — SQLite **no versionado** (en `.gitignore`). Se crea/usa solo en ejecución; la web actual no lo necesita.

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
- `JWT_SECRET` — debe coincidir con el del servicio de auth

## Build y despliegue

```bash
# Binario Linux (servidor)
env GOOS=linux CGO_ENABLED=1 go build -o app ./cmd/web

# Binario macOS
env GOOS=darwin CGO_ENABLED=1 go build -o bin/app ./cmd/web
```

En el servidor (Ubuntu 22.04, host `gosha` / 31.24.41.84):
- Código en `/home/anton/www/azr`, binario `app` corre en **localhost:3331**
- Gestionado por **supervisor** (patrón igual al de `circuitos`: `/etc/supervisor/conf.d/`, `supervisorctl restart`)
- **nginx** hace proxy con TLS (Let's Encrypt) → `https://azr.es` y `www.azr.es`

## Notas y avisos

- **`database/app.db` ya NO está versionado** (se sacó con `git rm --cached` y se añadió a `.gitignore`). El `.db` que había en el repo era **sobras de otro proyecto** ("Circuitos/Nürburgring"): tablas `cars`, `records`, `car_descriptions` (216+ filas de récords de circuito) más datos de demo en `pages`. **No había datos reales de azr.es** — la web ni lee la BD.
- El `JWT_SECRET` por defecto (`holacaracola`) está **hardcodeado** en `middleware.go` (línea 33) además de en `.env.sample`. Debería leerse de entorno y no estar en el repo.
- `Cookie.Secure = false` en `main.go`: en producción (HTTPS) lo correcto sería `true`.
- El middleware `Auth` y rutas como `/login` están definidos pero el router actual (`routes.go`) solo expone `/` y estáticos; el resto de handlers/rutas pueden estar pendientes o desactivados.
