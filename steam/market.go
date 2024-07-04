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

// MarketItem represents an item to be listed on the Steam market
type MarketItem struct {
	AppID      int    `json:"appid"`
	ContextID  int    `json:"contextid"`
	AssetID    int    `json:"assetid"`
	Price      int    `json:"price"`
	Currency   string `json:"currency"`
	Qty        int    `json:"qty"`
	MarketName string `json:"market_name"`
}

// ListMarketItem lists an item on the Steam market
// item: MarketItem struct containing item details
func (b *Bot) ListMarketItem(item MarketItem) error {
	data := url.Values{}
	data.Set("sessionid", "your_session_id")
	data.Set("appid", fmt.Sprintf("%d", item.AppID))
	data.Set("contextid", fmt.Sprintf("%d", item.ContextID))
	data.Set("assetid", fmt.Sprintf("%d", item.AssetID))
	data.Set("price", fmt.Sprintf("%d", item.Price))
	data.Set("currency", item.Currency)
	data.Set("quantity", fmt.Sprintf("%d", item.Qty))
	data.Set("market_name", item.MarketName)

	listMarketItemURL := fmt.Sprintf("https://steamcommunity.com/market/sellitem/")

	req, err := http.NewRequest("POST", listMarketItemURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create list market item request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := b.Session.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send list market item request: %w", err)
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
		return fmt.Errorf("failed to decode list market item response: %w", err)
	}

	if !result["success"].(bool) {
		return fmt.Errorf("listing market item failed: %v", result)
	}

	log.Printf("Market item listed: %+v\n", item)
	return nil
}
