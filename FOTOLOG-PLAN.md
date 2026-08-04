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
| 1 | Núcleo privado: esquema, login, subida+pipeline, admin, galería privada | **Alta** | Fable | ⬜ |
| 2 | Zona pública + SEO completo + descargas + licencia | Media | Opus 4.8 | ⬜ |
| 3 | Traducción automática ES→EN/RU (OpenAI) al guardar | Baja | Opus 4.8 | ⬜ |
| 4 | Dictado por voz (MediaRecorder → transcripción → redacción SEO → traducción) | Media | Opus 4.8 | ⬜ |

### Fase 1 — checklist
- [ ] Migración/creación de tabla `photos` (y arranque que la aplique si no existe)
- [ ] `JWT_SECRET` desde `.env`; `Cookie.Secure = true`; sesión 30 días
- [ ] `GET/POST /login` contra `AUTH_SERVICE` (copiar patrón de
      `certifications/services/web/internal/clients/auth.go`) + `/logout`
- [ ] Rutas privadas bajo middleware `Auth`: subir, listar (públicas+privadas), editar
      (títulos/desc/GPS toggle), publicar/despublicar (mueve objetos S3), borrar
- [ ] Pipeline exiftool→vips→S3 (paquete `internal/photos` o similar; shell-out con
      `os/exec`, límite de tamaño de subida, limpieza de temporales, HEIC→web ok)
- [ ] SDK AWS v2 solo S3 (`aws-sdk-go-v2/service/s3` + manager) — presigned GET para privadas
- [ ] Vista privada de galería/foto para testear desde el móvil
- [ ] Probar: subir desde iPhone/Android real → ver foto, EXIF y derivados correctos

### Fase 2 — checklist
- [ ] `/fotolog` + `/en/fotolog` + `/ru/fotolog` + detalle `/fotolog/{slug}` (3 idiomas)
- [ ] Textos de interfaz en `internal/i18n/i18n.go` (mantener tono guasón)
- [ ] `<picture>` WebP+JPEG, srcset, lazy loading, paginación
- [ ] Descargas por tamaño + original; bloque de licencia CC BY 4.0 con atribución
- [ ] hreflang/canonical/OG, schema.org ImageObject con license, sitemap de imágenes
- [ ] GPS visible solo si `show_gps=1`

### Fase 3 — checklist
- [ ] `OPENAI_API_KEY` en `.env`; al guardar desc ES con EN/RU vacíos → traducir
- [ ] Traducciones siempre editables; nunca sobrescribir texto manual sin pedirlo

### Fase 4 — checklist
- [ ] Botón micro (MediaRecorder) en el form de edición móvil
- [ ] Audio → transcripción OpenAI → redacción SEO-friendly ES → fase 3 traduce
- [ ] Todo queda como borrador editable antes de guardar
