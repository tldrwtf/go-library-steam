# Steam API Library for Go

A Go library for creating Steam bots with features like Steam Guard handling, captcha solving, friend management, trading, market integration, fetching player inventories, and more.

## Features

- **Steam Guard Handling**: Supports generating Steam Guard codes for 2FA.
- **Captcha Handling**: Provides placeholders for handling captcha challenges during login.
- **Friend Management**: Functions for adding, removing, and accepting friend requests.
- **Trading and Market Integration**: Functions for sending trade offers and listing items on the Steam market.
- **Fetching Player Inventories**: Function for fetching player inventories.
- **User Data Retrieval**: Fetches user stats for games, owned games, recently played games, and player summaries.
- **Rate Limiting**: Rate limiting for API requests to prevent IP bans.
- **Logging**: Simple logging using Go's `log` package.
- **Configuration Management**: Config management using a simple file reader.
- **Error Handling and Retries**: Basic error handling for common operations.

## Installation

1. Clone the repository.
2. Ensure you have Go installed.
3. Initialize the Go module:

    ```sh
    go mod init go-library-steam
    ```

## Configuration

Create a `.env` file with the following content:

```env
STEAM_API_KEY=your_steam_api_key
STEAM_SHARED_SECRET=your_shared_secret
STEAM_IDENTITY_SECRET=your_identity_secret
STEAM_USERNAME=your_username
STEAM_PASSWORD=your_password
```

Create a config.yaml file with the following content:

```
rate_limit:
  requests_per_second: 1
  burst: 5
```

## Usage

Here's an example of how to use the library:
```
package main

import (
    "log"

    "go-library-steam/steam"
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
        SharedSecret:  steam.GetEnv("STEAM_SHARED_SECRET"),
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
```

## Example Outputs
The output for the different functionalities might look like this:

Player Summaries:
```
{
    "steamid": "76561197960435530",
    "personaname": "JohnDoe",
    "profileurl": "https://steamcommunity.com/id/johndoe/",
    "avatar": "https://avatars.steamstatic.com/avatar.jpg",
    "avatarmedium": "https://avatars.steamstatic.com/avatarmedium.jpg",
    "avatarfull": "https://avatars.steamstatic.com/avatarfull.jpg"
}
```
Player Inventories:
```
{
    "assets": [
        {
            "appid": 730,
            "contextid": "2",
            "assetid": "1234567890",
            "classid": "111111",
            "instanceid": "0",
            "amount": "1"
        }
    ],
    "descriptions": [
        {
            "appid": 730,
            "classid": "111111",
            "instanceid": "0",
            "name": "AK-47 | Redline",
            "market_hash_name": "AK-47 | Redline",
            "icon_url": "https://steamcommunity-a.akamaihd.net/economy/image/icon.jpg"
        }
    ]
}
```
