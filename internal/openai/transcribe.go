package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

const transcribeEndpoint = "https://api.openai.com/v1/audio/transcriptions"

// Transcribe convierte un audio (voz en español) a texto. filename debe
// llevar una extensión que OpenAI reconozca (webm, mp4, m4a, wav, mp3, ogg…);
// se usa para detectar el formato.
func (t *Client) Transcribe(ctx context.Context, audio io.Reader, filename string) (string, error) {
	if !t.Enabled() {
		return "", fmt.Errorf("OpenAI deshabilitado (sin OPENAI_API_KEY)")
	}

	var body bytes.Buffer
	mw := multipart.NewWriter(&body)

	fw, err := mw.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(fw, audio); err != nil {
		return "", err
	}
	mw.WriteField("model", t.transcribeModel)
	mw.WriteField("language", "es") // pista de idioma: mejora precisión
	mw.WriteField("response_format", "json")
	if err := mw.Close(); err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", transcribeEndpoint, &body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+t.apiKey)

	resp, err := t.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("transcripción OpenAI: %w", err)
	}
	defer resp.Body.Close()

	var out struct {
		Text  string `json:"text"`
		Error *struct {
			Message string `json:"message"`
			Code    string `json:"code"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("respuesta de transcripción ilegible: %w", err)
	}
	if out.Error != nil {
		return "", fmt.Errorf("OpenAI (%s): %s", out.Error.Code, out.Error.Message)
	}

	return strings.TrimSpace(out.Text), nil
}

const polishPrompt = `Eres el editor de un fotolog personal. Te paso una descripción de una foto dictada por voz: puede traer muletillas, repeticiones, errores de transcripción o frases sin rematar.

Reescríbela como una descripción de foto lista para publicar:
- En español, clara, natural y bien redactada.
- SEO-friendly: lenguaje natural, sin repetir palabras clave de forma artificial.
- Conserva el tono y el contenido; corrige lo que haga falta.
- NO inventes datos que no estén en el dictado (lugares, fechas, nombres).
- Devuelve SOLO el texto final, sin comillas ni comentarios.`

// Polish limpia y reescribe en plan SEO una descripción en español (p. ej.
// la salida cruda de una transcripción de voz).
func (t *Client) Polish(ctx context.Context, raw string) (string, error) {
	if !t.Enabled() {
		return "", fmt.Errorf("OpenAI deshabilitado (sin OPENAI_API_KEY)")
	}
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", nil
	}

	reqBody, err := json.Marshal(chatRequest{
		Model: t.model,
		Messages: []chatMessage{
			{Role: "system", Content: polishPrompt},
			{Role: "user", Content: raw},
		},
		Temperature: 0.4,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", chatEndpoint, bytes.NewReader(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+t.apiKey)

	resp, err := t.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("pulido OpenAI: %w", err)
	}
	defer resp.Body.Close()

	var cr chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		return "", fmt.Errorf("respuesta de pulido ilegible: %w", err)
	}
	if cr.Error != nil {
		return "", fmt.Errorf("OpenAI (%s): %s", cr.Error.Code, cr.Error.Message)
	}
	if len(cr.Choices) == 0 {
		return "", fmt.Errorf("OpenAI no devolvió resultado")
	}

	return strings.TrimSpace(cr.Choices[0].Message.Content), nil
}
