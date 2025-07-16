package recaptcha

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
)

const verifyURL = "https://www.google.com/recaptcha/api/siteverify"

// Response es la estructura de la respuesta del servidor de verificación de reCAPTCHA de Google.
type Response struct {
	Success     bool     `json:"success"`
	Score       float64  `json:"score"`
	Action      string   `json:"action"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes,omitempty"`
}

var ErrRecaptchaNotConfigured = errors.New("la clave secreta de reCAPTCHA no está configurada")

// Verify envía el token de reCAPTCHA a Google para su verificación.
// Devuelve true si el token es válido y la puntuación está por encima del umbral.
func Verify(token, secretKey string) (bool, error) {
	if secretKey == "" {
		log.Println("ADVERTENCIA: La clave secreta de reCAPTCHA no está configurada en el servidor. La verificación de reCAPTCHA fallará.")
		return false, ErrRecaptchaNotConfigured
	}

	resp, err := http.PostForm(verifyURL, url.Values{"secret": {secretKey}, "response": {token}})
	if err != nil {
		log.Printf("Error al contactar el servidor de reCAPTCHA: %v", err)
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error al leer el cuerpo de la respuesta de reCAPTCHA: %v", err)
		return false, err
	}

	var recaptchaResponse Response
	if err := json.Unmarshal(body, &recaptchaResponse); err != nil {
		log.Printf("Error al decodificar la respuesta de reCAPTCHA: %v", err)
		return false, err
	}

	const scoreThreshold = 0.5 // Puedes ajustar este umbral según tus necesidades.
	return recaptchaResponse.Success && recaptchaResponse.Score >= scoreThreshold, nil
}
