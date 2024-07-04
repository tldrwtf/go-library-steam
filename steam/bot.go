package steam

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

// Bot represents a Steam bot
type Bot struct {
	APIKey     string
	Session    *http.Client
	SteamGuard *Guard
}

// NewBot creates a new Bot instance
// apiKey: Steam Web API key
// steamGuard: Steam Guard instance with shared secret for 2FA
func NewBot(apiKey string, steamGuard *Guard) *Bot {
	return &Bot{
		APIKey:     apiKey,
		Session:    &http.Client{},
		SteamGuard: steamGuard,
	}
}

// AddFriend sends a friend request to a specified SteamID
// steamID: SteamID64 of the user to add as a friend
func (b *Bot) AddFriend(steamID string) error {
	data := url.Values{}
	data.Set("sessionid", "your_session_id")
	data.Set("steamid", steamID)
	addFriendURL := fmt.Sprintf("https://steamcommunity.com/actions/AddFriendAjax")

	req, err := http.NewRequest("POST", addFriendURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create add friend request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := b.Session.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send add friend request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("steam API returned non-OK status: %s", resp.Status)
	}

	log.Printf("Friend request sent to SteamID: %s\n", steamID)
	return nil
}

// RemoveFriend removes a friend from the bot's friend list
// steamID: SteamID64 of the user to remove as a friend
func (b *Bot) RemoveFriend(steamID string) error {
	data := url.Values{}
	data.Set("sessionid", "your_session_id")
	data.Set("steamid", steamID)
	removeFriendURL := fmt.Sprintf("https://steamcommunity.com/actions/RemoveFriendAjax")

	req, err := http.NewRequest("POST", removeFriendURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create remove friend request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := b.Session.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send remove friend request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("steam API returned non-OK status: %s", resp.Status)
	}

	log.Printf("Friend removed with SteamID: %s\n", steamID)
	return nil
}

// AcceptFriendRequest accepts a friend request from a specified SteamID
// steamID: SteamID64 of the user whose friend request to accept
func (b *Bot) AcceptFriendRequest(steamID string) error {
	data := url.Values{}
	data.Set("sessionid", "your_session_id")
	data.Set("steamid", steamID)
	acceptFriendRequestURL := fmt.Sprintf("https://steamcommunity.com/actions/AcceptFriendRequest")

	req, err := http.NewRequest("POST", acceptFriendRequestURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create accept friend request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := b.Session.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send accept friend request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("steam API returned non-OK status: %s", resp.Status)
	}

	log.Printf("Accepted friend request from SteamID: %s\n", steamID)
	return nil
}

// Login handles the login process for the bot
// username: Steam account username
// password: Steam account password
func (b *Bot) Login(username, password string) error {
	return PerformLogin(username, password, b.SteamGuard, b.APIKey)
}
