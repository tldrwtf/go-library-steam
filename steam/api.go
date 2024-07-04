package steam

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// PlayerSummariesResponse represents the response from the GetPlayerSummaries API call
type PlayerSummariesResponse struct {
	Response struct {
		Players []struct {
			SteamID      string `json:"steamid"`
			PersonaName  string `json:"personaname"`
			ProfileURL   string `json:"profileurl"`
			Avatar       string `json:"avatar"`
			AvatarMedium string `json:"avatarmedium"`
			AvatarFull   string `json:"avatarfull"`
		} `json:"players"`
	} `json:"response"`
}

// PlayerInventoryResponse represents the response from the GetPlayerInventories API call
type PlayerInventoryResponse struct {
	Assets []struct {
		AppID      int    `json:"appid"`
		ContextID  string `json:"contextid"`
		AssetID    string `json:"assetid"`
		ClassID    string `json:"classid"`
		InstanceID string `json:"instanceid"`
		Amount     string `json:"amount"`
	} `json:"assets"`
	Descriptions []struct {
		AppID          int    `json:"appid"`
		ClassID        string `json:"classid"`
		InstanceID     string `json:"instanceid"`
		Name           string `json:"name"`
		MarketHashName string `json:"market_hash_name"`
		IconURL        string `json:"icon_url"`
	} `json:"descriptions"`
}

// UserStatsForGameResponse represents the response from the GetUserStatsForGame API call
type UserStatsForGameResponse struct {
	PlayerStats struct {
		SteamID  string `json:"steamID"`
		GameName string `json:"gameName"`
		Stats    []struct {
			Name  string `json:"name"`
			Value int    `json:"value"`
		} `json:"stats"`
		Achievements []struct {
			Name     string `json:"name"`
			Achieved int    `json:"achieved"`
		} `json:"achievements"`
	} `json:"playerstats"`
}

// OwnedGamesResponse represents the response from the GetOwnedGames API call
type OwnedGamesResponse struct {
	Response struct {
		GameCount int `json:"game_count"`
		Games     []struct {
			AppID           int    `json:"appid"`
			Name            string `json:"name"`
			PlaytimeForever int    `json:"playtime_forever"`
			ImgIconURL      string `json:"img_icon_url"`
			ImgLogoURL      string `json:"img_logo_url"`
		} `json:"games"`
	} `json:"response"`
}

// RecentlyPlayedGamesResponse represents the response from the GetRecentlyPlayedGames API call
type RecentlyPlayedGamesResponse struct {
	Response struct {
		TotalCount int `json:"total_count"`
		Games      []struct {
			AppID           int    `json:"appid"`
			Name            string `json:"name"`
			Playtime2Weeks  int    `json:"playtime_2weeks"`
			PlaytimeForever int    `json:"playtime_forever"`
			ImgIconURL      string `json:"img_icon_url"`
			ImgLogoURL      string `json:"img_logo_url"`
		} `json:"games"`
	} `json:"response"`
}

// GetPlayerSummaries fetches player summaries from Steam API
// apiKey: Steam Web API key
// steamID: SteamID64 of the player
func GetPlayerSummaries(apiKey, steamID string) (*PlayerSummariesResponse, error) {
	rateLimiter.Wait()

	// Build the URL for the API request
	url := fmt.Sprintf("https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v2/?key=%s&steamids=%s", apiKey, steamID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// Check for non-OK HTTP status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("steam API returned non-OK status: %s", resp.Status)
	}

	// Decode the JSON response
	var result PlayerSummariesResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return &result, nil
}

// GetPlayerInventories fetches player inventories from Steam API
// apiKey: Steam Web API key
// steamID: SteamID64 of the player
// appID: Application ID (e.g., 730 for CS:GO)
// contextID: Context ID (e.g., 2 for CS:GO)
func GetPlayerInventories(apiKey, steamID string, appID, contextID int) (*PlayerInventoryResponse, error) {
	rateLimiter.Wait()

	// Build the URL for the API request
	url := fmt.Sprintf("https://steamcommunity.com/inventory/%s/%d/%d?l=english&count=5000", steamID, appID, contextID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// Check for non-OK HTTP status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("steam API returned non-OK status: %s", resp.Status)
	}

	// Decode the JSON response
	var result PlayerInventoryResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return &result, nil
}

// GetUserStatsForGame fetches user stats for a specific game from the Steam API
func GetUserStatsForGame(apiKey, steamID string, appID int) (*UserStatsForGameResponse, error) {
	rateLimiter.Wait()

	// Build the URL for the API request
	url := fmt.Sprintf("https://api.steampowered.com/ISteamUserStats/GetUserStatsForGame/v2/?key=%s&steamid=%s&appid=%d", apiKey, steamID, appID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// Check for non-OK HTTP status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("steam API returned non-OK status: %s", resp.Status)
	}

	// Decode the JSON response
	var result UserStatsForGameResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return &result, nil
}

// GetOwnedGames fetches the list of games owned by a user from the Steam API
func GetOwnedGames(apiKey, steamID string) (*OwnedGamesResponse, error) {
	rateLimiter.Wait()

	// Build the URL for the API request
	url := fmt.Sprintf("https://api.steampowered.com/IPlayerService/GetOwnedGames/v1/?key=%s&steamid=%s&include_appinfo=true&include_played_free_games=true", apiKey, steamID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// Check for non-OK HTTP status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("steam API returned non-OK status: %s", resp.Status)
	}

	// Decode the JSON response
	var result OwnedGamesResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return &result, nil
}

// GetRecentlyPlayedGames fetches the list of recently played games from the Steam API
func GetRecentlyPlayedGames(apiKey, steamID string) (*RecentlyPlayedGamesResponse, error) {
	rateLimiter.Wait()

	// Build the URL for the API request
	url := fmt.Sprintf("https://api.steampowered.com/IPlayerService/GetRecentlyPlayedGames/v1/?key=%s&steamid=%s", apiKey, steamID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// Check for non-OK HTTP status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("steam API returned non-OK status: %s", resp.Status)
	}

	// Decode the JSON response
	var result RecentlyPlayedGamesResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return &result, nil
}
