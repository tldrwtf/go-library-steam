package steam

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

// TradeOffer represents a trade offer to be sent
type TradeOffer struct {
	PartnerSteamID string `json:"partner_steamid"`
	ItemsToSend    string `json:"items_to_send"`
	ItemsToReceive string `json:"items_to_receive"`
	Message        string `json:"message"`
}

// SendTradeOffer sends a trade offer
// offer: TradeOffer struct containing trade offer details
func (b *Bot) SendTradeOffer(offer TradeOffer) error {
	data := url.Values{}
	data.Set("sessionid", "your_session_id")
	data.Set("partner", offer.PartnerSteamID)
	data.Set("tradeoffermessage", offer.Message)
	data.Set("json_tradeoffer", fmt.Sprintf(`{"newversion":true,"version":2,"me":{"assets":%s,"currency":[],"ready":false},"them":{"assets":%s,"currency":[],"ready":false}}`, offer.ItemsToSend, offer.ItemsToReceive))

	tradeOfferURL := fmt.Sprintf("https://steamcommunity.com/tradeoffer/new/send")

	req, err := http.NewRequest("POST", tradeOfferURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create trade offer request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := b.Session.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send trade offer request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("steam API returned non-OK status: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode trade offer response: %w", err)
	}

	if result["tradeofferid"] == nil {
		return fmt.Errorf("trade offer failed: %v", result)
	}

	log.Printf("Trade offer sent to SteamID: %s\n", offer.PartnerSteamID)
	return nil
}
