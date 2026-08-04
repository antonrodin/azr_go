// Package authsvc habla con el servicio de auth de microservicios
// (AUTH_SERVICE) para hacer login. El JWT devuelto se guarda en la sesión
// y lo valida el middleware Auth en cada petición.
package authsvc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var ErrInvalidCredentials = errors.New("credenciales inválidas")

type Client struct {
	baseURL string
	service string
	http    *http.Client
}

// New crea el cliente. service es el nombre del servicio en auth ("azr").
func New(baseURL, service string) *Client {
	return &Client{
		baseURL: baseURL,
		service: service,
		http:    &http.Client{Timeout: 10 * time.Second},
	}
}

// Login devuelve el JWT si las credenciales son válidas.
func (c *Client) Login(ctx context.Context, email, password string) (string, error) {
	body, err := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
		"service":  c.service,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/v1/login", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("servicio auth no disponible: %w", err)
	}
	defer resp.Body.Close()

	var out struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
		Data    struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("respuesta de auth ilegible: %w", err)
	}

	if out.Error || out.Data.Token == "" {
		// Cualquier fallo de login se trata igual de cara al usuario.
		return "", ErrInvalidCredentials
	}

	return out.Data.Token, nil
}
