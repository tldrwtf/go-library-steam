package steam

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// LoginResponse represents the response from the Steam login API call
type LoginResponse struct {
	Success           bool   `json:"success"`
	RequiresTwoFactor bool   `json:"requires_twofactor"`
	CaptchaNeeded     bool   `json:"captcha_needed"`
	CaptchaGid        string `json:"captcha_gid"`
	Message           string `json:"message"`
}

// PerformLogin performs the login process
// username: Steam account username
// password: Steam account password
// steamGuard: Steam Guard instance with shared secret for 2FA
// apiKey: Steam Web API key
func PerformLogin(username, password string, steamGuard *Guard, apiKey string) error {
	// Prepare the form data for the login request
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)

	// If Steam Guard is enabled, generate the Steam Guard code
	if steamGuard != nil {
		code, err := GenerateSteamGuardCode(steamGuard.SharedSecret)
		if err != nil {
			return fmt.Errorf("failed to generate Steam Guard code: %w", err)
		}
		data.Set("twofactorcode", code)
	}

	// Send the login request
	resp, err := http.PostForm("https://steamcommunity.com/login/dologin", data)
	if err != nil {
		return fmt.Errorf("failed to perform login: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// Check for non-OK HTTP status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("steam login returned non-OK status: %s", resp.Status)
	}

	// Decode the JSON response
	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return fmt.Errorf("failed to decode login response: %w", err)
	}

	// Handle various login scenarios
	if !loginResp.Success {
		if loginResp.RequiresTwoFactor {
			return fmt.Errorf("two-factor authentication required")
		}
		if loginResp.CaptchaNeeded {
			return fmt.Errorf("captcha required: %s", loginResp.CaptchaGid)
		}
		return fmt.Errorf("login failed: %s", loginResp.Message)
	}

	return nil
}
