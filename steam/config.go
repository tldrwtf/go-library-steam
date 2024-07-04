package steam

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// envVars holds the environment variables
var envVars map[string]string

// LoadEnv loads environment variables from a .env file
func LoadEnv(filePath string) error {
	envVars = make(map[string]string)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid line in .env file: %s", line)
		}
		envVars[parts[0]] = parts[1]
	}
	return scanner.Err()
}

// GetEnv retrieves the value of the environment variable named by the key
func GetEnv(key string) string {
	return envVars[key]
}

// Config holds the configuration values
var Config struct {
	RateLimit struct {
		RequestsPerSecond int
		Burst             int
	}
}

// LoadConfig loads the configuration from a config.yaml file
func LoadConfig(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid line in config file: %s", line)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "rate_limit.requests_per_second":
			_, err := fmt.Sscanf(value, "%d", &Config.RateLimit.RequestsPerSecond)
			if err != nil {
				return err
			}
		case "rate_limit.burst":
			_, err := fmt.Sscanf(value, "%d", &Config.RateLimit.Burst)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown configuration key: %s", key)
		}
	}
	return scanner.Err()
}

// RateLimiter controls the rate of API requests
var rateLimiter *RateLimiter

// RateLimiter is a simple rate limiter
type RateLimiter struct {
	tokens            chan struct{}
	ticker            *time.Ticker
	requestsPerSecond int
	burst             int
}

// NewRateLimiter creates a new RateLimiter
func NewRateLimiter(requestsPerSecond, burst int) *RateLimiter {
	rl := &RateLimiter{
		tokens:            make(chan struct{}, burst),
		ticker:            time.NewTicker(time.Second / time.Duration(requestsPerSecond)),
		requestsPerSecond: requestsPerSecond,
		burst:             burst,
	}

	go func() {
		for range rl.ticker.C {
			select {
			case rl.tokens <- struct{}{}:
			default:
			}
		}
	}()

	return rl
}

// Wait blocks until a token is available
func (rl *RateLimiter) Wait() {
	<-rl.tokens
}

func init() {
	rateLimiter = NewRateLimiter(1, 5)
}
