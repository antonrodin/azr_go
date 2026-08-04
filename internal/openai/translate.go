// Package openai integra los servicios de OpenAI que usa el fotolog:
// traducción de textos (Fase 3) y transcripción de voz + pulido (Fase 4).
package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const chatEndpoint = "https://api.openai.com/v1/chat/completions"

// DefaultModel es el modelo de chat por defecto (OPENAI_MODEL).
const DefaultModel = "gpt-4o-mini"

// DefaultTranscribeModel es el modelo de transcripción por defecto
// (OPENAI_TRANSCRIBE_MODEL).
const DefaultTranscribeModel = "gpt-4o-mini-transcribe"

// Client llama a la API de OpenAI. Si apiKey está vacío, se considera
// deshabilitado (el llamante debe comprobar Enabled()).
type Client struct {
	apiKey          string
	model           string
	transcribeModel string
	http            *http.Client
}

func NewClient(apiKey, model, transcribeModel string) *Client {
	if model == "" {
		model = DefaultModel
	}
	if transcribeModel == "" {
		transcribeModel = DefaultTranscribeModel
	}
	return &Client{
		apiKey:          apiKey,
		model:           model,
		transcribeModel: transcribeModel,
		// Traducir es rápido; transcribir un audio puede tardar más. 60s de
		// margen para no cortar una respuesta legítima.
		http: &http.Client{Timeout: 60 * time.Second},
	}
}

// Enabled indica si hay API key configurada.
func (t *Client) Enabled() bool { return t != nil && t.apiKey != "" }

// Translation son los textos traducidos. Un campo vacío en el original
// produce vacío en sus traducciones.
type Translation struct {
	TitleEN string `json:"title_en"`
	DescEN  string `json:"desc_en"`
	TitleRU string `json:"title_ru"`
	DescRU  string `json:"desc_ru"`
}

const systemPrompt = `Eres el traductor de un fotolog personal. Traduce el título y la descripción de una foto del español a inglés y ruso.

Reglas:
- Traducción natural y fluida, como la escribiría un nativo, no literal.
- Conserva el tono del original (puede ser desenfadado, irónico o con humor).
- Texto claro y SEO-friendly: natural, sin repetir palabras clave de forma artificial.
- No añadas comillas, ni prefijos, ni comentarios; solo el texto traducido.
- Si un campo de origen viene vacío, devuelve cadena vacía en sus traducciones.

Devuelve SOLO un JSON con estas claves exactas: title_en, desc_en, title_ru, desc_ru.`

type chatRequest struct {
	Model          string            `json:"model"`
	Messages       []chatMessage     `json:"messages"`
	Temperature    float64           `json:"temperature"`
	ResponseFormat map[string]string `json:"response_format,omitempty"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"error"`
}

// Translate traduce título y descripción en español. Una sola llamada
// devuelve las cuatro traducciones (EN y RU).
func (t *Client) Translate(ctx context.Context, titleES, descES string) (*Translation, error) {
	if !t.Enabled() {
		return nil, fmt.Errorf("traductor deshabilitado (sin OPENAI_API_KEY)")
	}

	input, err := json.Marshal(map[string]string{"title_es": titleES, "desc_es": descES})
	if err != nil {
		return nil, err
	}

	reqBody, err := json.Marshal(chatRequest{
		Model: t.model,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: string(input)},
		},
		Temperature:    0.3,
		ResponseFormat: map[string]string{"type": "json_object"},
	})
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", chatEndpoint, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+t.apiKey)

	resp, err := t.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("llamada a OpenAI: %w", err)
	}
	defer resp.Body.Close()

	var cr chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		return nil, fmt.Errorf("respuesta de OpenAI ilegible: %w", err)
	}
	if cr.Error != nil {
		return nil, fmt.Errorf("OpenAI (%s): %s", cr.Error.Code, cr.Error.Message)
	}
	if len(cr.Choices) == 0 {
		return nil, fmt.Errorf("OpenAI no devolvió resultado")
	}

	var out Translation
	if err := json.Unmarshal([]byte(cr.Choices[0].Message.Content), &out); err != nil {
		return nil, fmt.Errorf("JSON de traducción inválido: %w", err)
	}

	return &out, nil
}
