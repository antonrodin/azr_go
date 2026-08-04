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

// DefaultComposeInstructions son las instrucciones editables por defecto del
// dictado: convierten la voz cruda en un título y una descripción. NO incluyen
// el contrato de salida (JSON), que se añade aparte.
const DefaultComposeInstructions = `Eres el editor de un fotolog personal. Te paso una descripción de una foto dictada por voz: puede traer muletillas, repeticiones, errores de transcripción o frases sin rematar.

A partir de ella genera dos cosas:
- Un título corto (entre 3 y 6 palabras), claro y con gancho.
- Una descripción lista para publicar: en español, natural, bien redactada y SEO-friendly (sin repetir palabras clave de forma artificial).

Reglas:
- Conserva el tono y el contenido; corrige lo que haga falta.
- NO inventes datos que no estén en el dictado (lugares, fechas, nombres).`

// composeFormat es el contrato de salida (no editable).
const composeFormat = "\n\nDevuelve SOLO un JSON con estas claves exactas: title, description."

// Composition es el resultado del dictado: título y descripción en español.
type Composition struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Compose convierte una transcripción cruda (voz en español) en un título y una
// descripción listos para publicar. instructions son las reglas editables.
func (t *Client) Compose(ctx context.Context, instructions, raw string) (*Composition, error) {
	if !t.Enabled() {
		return nil, fmt.Errorf("OpenAI deshabilitado (sin OPENAI_API_KEY)")
	}
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return &Composition{}, nil
	}
	if instructions == "" {
		instructions = DefaultComposeInstructions
	}

	reqBody, err := json.Marshal(chatRequest{
		Model: t.model,
		Messages: []chatMessage{
			{Role: "system", Content: instructions + composeFormat},
			{Role: "user", Content: raw},
		},
		Temperature:    0.4,
		ResponseFormat: map[string]string{"type": "json_object"},
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", chatEndpoint, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+t.apiKey)

	resp, err := t.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("dictado OpenAI: %w", err)
	}
	defer resp.Body.Close()

	var cr chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		return nil, fmt.Errorf("respuesta de dictado ilegible: %w", err)
	}
	if cr.Error != nil {
		return nil, fmt.Errorf("OpenAI (%s): %s", cr.Error.Code, cr.Error.Message)
	}
	if len(cr.Choices) == 0 {
		return nil, fmt.Errorf("OpenAI no devolvió resultado")
	}

	var out Composition
	if err := json.Unmarshal([]byte(cr.Choices[0].Message.Content), &out); err != nil {
		return nil, fmt.Errorf("JSON de dictado inválido: %w", err)
	}
	out.Title = strings.TrimSpace(out.Title)
	out.Description = strings.TrimSpace(out.Description)
	return &out, nil
}
