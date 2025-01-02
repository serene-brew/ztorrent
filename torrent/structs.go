package torrent

import "time"

// Config holds the configuration for torrent operations
type Config struct {
    Timeout      time.Duration
    Debug        bool
    ShowProgress bool
    Seed         bool
}

// DefaultConfig returns default configuration values
func DefaultConfig() Config {
    return Config{
        Timeout:      2 * time.Minute,
        Debug:        false,
        ShowProgress: true,
        Seed:        false,
    }
}