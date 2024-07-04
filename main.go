package main

import (
	"go-library-steam/steam"
	"log"
)

func main() {
	// Load environment variables
	err := steam.LoadEnv(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Load configuration
	err = steam.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	apiKey := steam.GetEnv("STEAM_API_KEY")
	if apiKey == "" {
		log.Fatalf("STEAM_API_KEY not set in .env file")
	}

	steamGuard := &steam.SteamGuard{
		SharedSecret:   steam.GetEnv("STEAM_SHARED_SECRET"),
		IdentitySecret: steam.GetEnv("STEAM_IDENTITY_SECRET"),
	}

	bot := steam.NewBot(apiKey, steamGuard)

	username := steam.GetEnv("STEAM_USERNAME")
	password := steam.GetEnv("STEAM_PASSWORD")

	err = bot.Login(username, password)
	if err != nil {
		log.Fatalf("Error logging in: %v", err)
	}

	playerSummaries, err := steam.GetPlayerSummaries(apiKey, "76561197960435530")
	if err != nil {
		log.Fatalf("Error getting player summaries: %v", err)
	}

	log.Printf("Player Summaries: %+v\n", playerSummaries)

	err = bot.AddFriend("76561197960435530")
	if err != nil {
		log.Fatalf("Error adding friend: %v", err)
	}

	playerInventories, err := steam.GetPlayerInventories(apiKey, "76561197960435530", 730, 2)
	if err != nil {
		log.Fatalf("Error getting player inventories: %v", err)
	}

	log.Printf("Player Inventories: %+v\n", playerInventories)

	userStats, err := steam.GetUserStatsForGame(apiKey, "76561197960435530", 730)
	if err != nil {
		log.Fatalf("Error getting user stats for game: %v", err)
	}

	log.Printf("User Stats for Game: %+v\n", userStats)

	ownedGames, err := steam.GetOwnedGames(apiKey, "76561197960435530")
	if err != nil {
		log.Fatalf("Error getting owned games: %v", err)
	}

	log.Printf("Owned Games: %+v\n", ownedGames)

	recentlyPlayedGames, err := steam.GetRecentlyPlayedGames(apiKey, "76561197960435530")
	if err != nil {
		log.Fatalf("Error getting recently played games: %v", err)
	}

	log.Printf("Recently Played Games: %+v\n", recentlyPlayedGames)
}
