package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/didit-pub/rt-gcli/internal/config"
)

// Client contiene la lógica principal del cliente RT
type Client struct {
	baseURL    string
	httpClient *http.Client
	auth       *AuthConfig
	debug      bool
}

// AuthConfig contiene la configuración de autenticación
type AuthConfig struct {
	Username string
	Password string
	Token    string
}

// ClientOption define un tipo para las opciones del cliente
type ClientOption func(*Client)

// WithDebug habilita el modo debug
func WithDebug(debug bool) ClientOption {
	return func(c *Client) {
		c.debug = debug
	}
}

// NewClient crea una nueva instancia del cliente
func NewClient(cfg *config.Config, opts ...ClientOption) (*Client, error) {
	if cfg.APIURL == "" {
		return nil, fmt.Errorf("RT API URL is required")
	}

	// Crear el cliente HTTP con timeout
	httpClient := &http.Client{
		Timeout: cfg.Timeout,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	// Crear el cliente base
	client := &Client{
		baseURL:    cfg.APIURL,
		httpClient: httpClient,
		auth: &AuthConfig{
			Username: cfg.Username,
			Password: cfg.Password,
			Token:    cfg.Token,
		},
		debug: cfg.Debug,
	}

	// Aplicar opciones adicionales
	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}

// doRequest realiza una petición HTTP y maneja la respuesta
func (c *Client) doRequest(method, endpoint string, body interface{}, params map[string]string) ([]byte, error) {
	var err error
	// Construir URL completa
	baseURL := fmt.Sprintf("%s/%s", c.baseURL, endpoint)

	// Crear URL con parámetros
	requestURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %w", err)
	}

	// Añadir parámetros si existen
	if len(params) > 0 {
		q := requestURL.Query()
		for key, value := range params {
			q.Add(key, value)
		}
		requestURL.RawQuery = q.Encode()
	}

	var req *http.Request
	if body != nil {
		jsonBody, errb := json.Marshal(body)
		if errb != nil {
			return nil, fmt.Errorf("error marshaling request body: %w", errb)
		}

		// Debug logging del body si está habilitado
		if c.debug {
			fmt.Printf("Request body:\n%s\n", string(jsonBody))
		}

		req, err = http.NewRequest(method, requestURL.String(), bytes.NewBuffer(jsonBody))
	} else {
		req, err = http.NewRequest(method, requestURL.String(), nil)
	}

	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	// Añadir headers
	req.Header.Set("Content-Type", "application/json")

	// Configurar autenticación
	if c.auth.Token != "" {
		req.Header.Set("Authorization", "token "+c.auth.Token)
	} else {
		req.SetBasicAuth(c.auth.Username, c.auth.Password)
	}

	// Debug logging
	if c.debug {
		fmt.Printf("Request: %s %s\n", method, requestURL.String())
	}

	// Realizar la petición
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// Leer la respuesta
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Verificar el código de estado
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s",
			resp.StatusCode, string(respBody))
	}
	if c.debug {
		fmt.Printf("Response: %s\n", string(respBody))
	}
	return respBody, nil
}
