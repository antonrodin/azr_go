# Plan del Fotolog de azr.es

> Documento de traspaso entre fases/sesiones. Cada fase la puede ejecutar un modelo
> distinto: **lee este documento entero y el CLAUDE.md antes de empezar tu fase**, y al
> terminarla marca su checklist y anota decisiones nuevas aquí.

## Qué es

Fotolog personal en azr.es con dos zonas:
- **Pública**: galería y páginas de foto en 3 idiomas (es `/fotolog`, en `/en/fotolog`,
  ru `/ru/fotolog`), SEO-first, licencia **CC BY 4.0** (uso comercial permitido con
  atribución), descargas en varios tamaños incluido el original.
- **Privada**: tras login (servicio auth de microservicios), Anton sube fotos desde el
  móvil, edita títulos/descripciones y publica/despublica.

**Estilo**: seguir el diseño actual de azr.es — negro sobre blanco, CSS propio y mínimo
(NO hay Bootstrap ni frameworks), HTML server-side, tono sarcástico en los textos.
**Keep it simple**: sin ORM (`database/sql` a pelo), sin frameworks JS, URLs amigables,
página ligera.

## Decisiones cerradas (por Anton)

- **BD: SQLite** (la ya montada en el proyecto), sin ORM.
- **GPS**: toggle por foto; en fotos públicas **oculto por defecto**, visible en privadas.
- **RAW**: no por ahora (fotos de móvil JPEG/HEIC). Dejar el diseño abierto para adjuntar
  ficheros extra por foto en el futuro.
- **Zona privada**: auth del microservicio (`POST {AUTH_SERVICE}/api/v1/login` devuelve
  JWT; el middleware `Auth` de `cmd/web/middleware.go` ya lo valida). Solo la usa Anton;
  "logueado = ve todo" es suficiente. Prioridad: funcionalidad de subir/ver antes que
  pulir auth.
- Títulos y descripciones por foto en 3 idiomas como **columnas** (`title_es/en/ru`,
  `desc_es/en/ru`), no tabla de traducciones.
- Traducción automática (fase 3) y dictado por voz (fase 4) con la **API de OpenAI**.

## Infraestructura (Fase 0 — HECHA, 2026-08-04)

- **Bucket S3 `azr-fotolog`** en `eu-west-1`. Prefijos:
  - `public/` — lectura pública por bucket policy (URLs estables para Google Images):
    `https://azr-fotolog.s3.eu-west-1.amazonaws.com/public/...`
  - `private/` — sin acceso público (verificado 403); servir con URLs prefirmadas (~15 min).
  - Publicar/despublicar una foto = copiar sus objetos entre prefijos + borrar origen.
- **Credenciales**: en `.env` (local y servidor `~/www/azr/.env`, chmod 600):
  `AWS_REGION`, `S3_BUCKET`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`.
- **Servidor** (gosha.vservers.es): instalados `vips` 8.12.1 (con HEIF → lee HEIC de
  iPhone) y `exiftool` 12.40. En el Mac: instalados vía brew.
- **Deuda pendiente**:
  - Las credenciales AWS son las del usuario `MacBookAir` (S3-only pero con acceso a
    todos los buckets). El usuario CLI no tiene permisos IAM para crear uno dedicado →
    Anton debe crear en la consola web un usuario `azr-fotolog` (sin consola, solo
    programmatic) con política limitada a `arn:aws:s3:::azr-fotolog[/*]` y rotar el `.env`.
  - `JWT_SECRET=holacaracola` compartido con el servicio auth (hardcodeado allí en
    `cmd/api/helpers.go`). Cambiarlo exige tocar auth y sus consumidores a la vez.
    El código de azr debe leerlo de `.env` (no hardcodear como hace hoy
    `middleware.go:33`).
  - `Cookie.Secure = false` en `main.go` → poner `true` en fase 1 (todo va por HTTPS).

## Arquitectura acordada

- **Pipeline de subida** (todo en el servidor, sin presigned uploads):
  móvil → POST multipart a la app → `exiftool -json` (guardar JSON completo + extraer
  columnas clave) → `vips` genera derivados → subir a S3 → fila en SQLite → borrar temp.
- **Derivados por foto** (claves S3 deterministas, p.ej. `{prefijo}/{id}/{tamaño}.{ext}`):
  WebP 400/1200/2048 para la web (`<picture>` + srcset + lazy) y JPEG 1200/2048 +
  original para descargas.
- **Subida móvil**: página tras login con `<input type="file" multiple accept="image/*">`
  y JS vanilla mínimo (subida secuencial con estado por fichero). Sesión scs: subir
  Lifetime (p.ej. 30 días) para comodidad en el móvil.
- **Esquema orientativo** tabla `photos`: id, slug (del título ES, `gosimple/slug`),
  status (private/public), title_es/en/ru, desc_es/en/ru, taken_at, uploaded_at,
  camera_make, camera_model, lens, focal, aperture, shutter, iso, gps_lat, gps_lng,
  show_gps (bool, default 0), width, height, orig_ext, orig_bytes, exif_json.
- **SEO** (fase 2): hreflang triple + canonical (patrón de la home), OpenGraph,
  schema.org `ImageObject` con `license` + `acquireLicensePage` (badge "Licenciable" en
  Google Images), sitemap de imágenes, paginación `/fotolog/pagina/2`.

## Fases

| # | Fase | Dificultad | Modelo | Estado |
|---|------|-----------|--------|--------|
| 0 | Infraestructura (S3, vips/exiftool, .env) | Baja | — | ✅ Hecha |
| 1 | Núcleo privado: esquema, login, subida+pipeline, admin, galería privada | **Alta** | Fable | ✅ Hecha (2026-08-04) |
| 2 | Zona pública + SEO completo + descargas + licencia | Media | Opus 4.8 | ✅ Hecha (2026-08-04) |
| 3 | Traducción automática ES→EN/RU (OpenAI) al guardar | Baja | Opus 4.8 | ✅ Hecha (2026-08-04) |
| 4 | Dictado por voz (MediaRecorder → transcripción → redacción SEO → traducción) | Media | Opus 4.8 | ✅ Hecha (2026-08-04) |

### Fase 1 — checklist (hecha; testeada e2e en local con foto JPEG+EXIF+GPS y HEIC)
- [x] Tabla `photos` creada al arrancar (`PhotoModel.EnsureSchema`, CREATE IF NOT EXISTS)
- [x] `JWT_SECRET` desde `.env`; `Cookie.Secure` según `ENV=prod`; sesión 30 días;
      `RenewToken` al hacer login (session fixation)
- [x] `GET/POST /login` contra `AUTH_SERVICE` (`internal/authsvc`; servicio "azr" —
      **creado en la BD de auth el 2026-08-04**, usuarios 7 y 10 asignados) + `/admin/logout`
- [x] Rutas `/admin/*` bajo `Auth`: listar, subir, editar, publicar/despublicar, borrar
- [x] Pipeline en `internal/photos/`: `exif.go` (2 llamadas exiftool: volcado legible +
      campos numéricos con -n), `derivatives.go` (vips thumbnail, autorrotate, no
      amplía, strip), `s3.go` (Storage: upload/presign/move/delete)
- [x] Derivados: 400/1200/2048.webp + 1200/2048.jpg + orig.{ext};
      Cache-Control immutable; subida máx. 100 MB
- [x] Slug del título ES al guardar (único; congelado si la foto es pública);
      publicar exige título ES
- [x] Subida móvil: `admin/subir.go.tmpl` con JS vanilla secuencial (Accept: application/json)
      + fallback form clásico sin JS
- [ ] Pendiente de Anton: probar subida desde el móvil real en producción

Notas de implementación (para fases 2-4):
- Vistas admin en `templates/admin/*.go.tmpl`, estilos al final de `public/css/main.css`.
- `photoView` (handlers) = `*models.Photo` + URLs resueltas (presigned si privada,
  `Storage.PublicURL` si pública). Presign TTL 15 min.
- `exif_json` guarda el volcado completo de `exiftool -json -g0` (agrupado, legible) —
  la fase 2 puede pintar la tabla de metadatos desde las columnas y ofrecer el JSON en
  un `<details>`.
- GPS: columnas `gps_lat/gps_lng` NULL si no hay; `show_gps` (0 por defecto) manda en
  lo público.

### Fase 2 — checklist (hecha; testeada e2e en local, incl. privacidad)
- [x] Galería `/fotolog` + `/{lang}/fotolog` + paginación `/fotolog/pagina/{n}` (24/pág)
      y detalle `/fotolog/{slug}` + `/{lang}/fotolog/{slug}`, en es/en/ru
- [x] Textos de interfaz del fotolog en `internal/i18n/fotolog.go` (aparte del home,
      tono guasón); títulos/desc por foto salen de la BD con respaldo a ES
- [x] `<picture>` WebP+JPEG con srcset+sizes, `loading`/`fetchpriority`, lazy en grid
- [x] Descargas mediana/grande/original (con tamaño del original); licencia CC BY 4.0
      con texto humano + atribución copiable + enlace a la licencia
- [x] hreflang triple + x-default + canonical, OpenGraph + Twitter card,
      JSON-LD `ImageObject` con `license` + `acquireLicensePage` (badge Licenciable),
      sitemap de imágenes en `/sitemap.xml` + `/robots.txt`
- [x] GPS solo si `HasGPS && show_gps`; verificado que con `show_gps=0` no se filtra;
      foto privada da 404 en público y no sale en galería
- [x] Enlace al fotolog desde la home (nav, en los 3 idiomas) — enlace interno SEO

Notas de implementación (para fases 3-4):
- **Bootstrap eliminado** del layout (`templates/layouts/app.go.tmpl`): era un CDN de
  ~230 KB que no se usaba y contradecía el objetivo de página rápida. La home y el
  fotolog usan solo `public/css/main.css`. Verificado que la home sigue bien.
- JSON-LD se inyecta con la func de plantilla **`safeJS`** (`internal/tmpl/tmpl.go`),
  que devuelve `template.JS`; con `unescaped`/`template.HTML` html/template lo
  reescapaba dentro de `<script>`. `json.Marshal` ya escapa `< > &`, así que es seguro.
- Rutas: las estáticas (`/fotolog`, `/sitemap.xml`…) van **antes** del comodín
  `/{lang}` en `routes.go`; chi prioriza estáticas, no chocan.
- Consultas públicas en `photos.go`: `CountPublic`, `PublicPage`, `AllPublic`,
  `GetPublicBySlug` (esta filtra por `status='public'`, base de la privacidad).

### ⚠️ Avisos importantes (leer antes de tocar infra)
- **Bucket compartido → separado por entorno.** Producción usa `azr-fotolog`; el `.env`
  **local** apunta a **`azr-fotolog-dev`** (creado 2026-08-04, misma policy public/).
  Motivo: local y prod tienen BDs SQLite independientes, así que los IDs pueden
  colisionar y un test local sobrescribiría ficheros de producción. NO apuntes el .env
  local a `azr-fotolog`.
- **Sesiones en memoria**: `scs` sin store persistente → cada `supervisorctl restart`
  cierra la sesión y hay que volver a hacer login (aunque el Lifetime sea 30 días).
  Mejora opcional futura: `scs/sqlite3store` sobre la misma BD.
- **Anton ya subió una foto real en prod** (id 1, privada, sin título en `~/www/azr`).
  Para publicarla necesita título en español desde `/admin/fotos/1`.

### Fase 3 — checklist (hecha; testeada e2e con llamada real a OpenAI)
- [x] `OPENAI_API_KEY` en `.env` (local y servidor); `OPENAI_MODEL` opcional
      (por defecto `gpt-4o-mini`, cambiable sin tocar código)
- [x] Cliente en `internal/openai/translate.go`: una llamada `chat.completions`
      (JSON mode) devuelve title/desc en EN y RU a partir del ES
- [x] Enganchado en `PhotoEditPost` → `autoTranslate`: al guardar rellena solo los
      campos EN/RU vacíos; no gasta llamada si no hace falta; si la API falla, guarda
      igual y avisa por flash (nunca se pierde lo escrito)
- [x] Casilla "retraducir" en `admin/foto.go.tmpl` para sobrescribir a propósito
      (cumple "nunca pisar texto manual sin pedirlo"); verificado ambos caminos
- [x] Traductor deshabilitado si no hay API key (el resto del admin funciona igual)

Notas para la Fase 4:
- El traductor vive en `app.Translator` (`*openai.Translator`), con `Enabled()`.
- La Fase 4 (voz) puede reutilizar el mismo patrón: transcribir audio → texto ES →
  este mismo `autoTranslate` rellena EN/RU. Falta solo el paso audio→texto
  (`whisper-1` o `gpt-4o-mini-transcribe`, ambos ya disponibles en la cuenta).

### Fase 4 — checklist (hecha; testeada e2e con audio real generado por `say`)
- [x] Botón "🎤 Dictar" junto a la descripción ES en `admin/foto.go.tmpl`, con JS
      vanilla (MediaRecorder; elige webm/mp4/ogg según soporte; graba→para)
- [x] `POST /admin/transcribir`: recibe el audio, `Transcribe` (gpt-4o-mini-transcribe,
      pista `language=es`) → texto, `Polish` (gpt-4o-mini) → descripción SEO en español
- [x] El texto rellena el textarea como **borrador** (no guarda); al pulsar Guardar
      entra el flujo de la Fase 3 y se traduce a EN/RU
- [x] Si el pulido falla, devuelve la transcripción cruda; si no hay OpenAI o el
      navegador no soporta MediaRecorder, el botón se deshabilita y el resto va igual
- [x] Modelo de transcripción configurable con `OPENAI_TRANSCRIBE_MODEL`
      (por defecto `gpt-4o-mini-transcribe`)

Cliente OpenAI unificado en `internal/openai` (renombrado de `Translator` a `Client`):
`Translate` (fase 3), `Transcribe` + `Polish` (fase 4). Campo en `AppRepository.OpenAI`.

## Estado final: fotolog completo (fases 0-4 hechas, 2026-08-04)
Falta solo que Anton use la zona privada con fotos reales. Requisitos de ejecución vivos:
`vips`/`exiftool` en el servidor, `.env` con S3 + `OPENAI_API_KEY` (+ modelos opcionales),
servicio `AZR` en la BD de auth. Ver avisos de infra más arriba (bucket dev vs prod,
sesiones en memoria).
