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

// FileInfo holds information about a single file
type FileInfo struct {
    Name     string
    Size     int64
    Type     string
    Path     string
    Complete bool
}

// TorrentInfo holds comprehensive torrent information
type TorrentInfo struct {
    Name       string
    InfoHash   string
    TotalSize  int64
    Files      []FileInfo
    FilesByExt map[string][]FileInfo
}

// PeerInfo holds peer connection information
type PeerInfo struct {
    Address  string
    Active   bool
    Stats    PeerStats
}

// PeerStats holds peer statistics
type PeerStats struct {
    TotalPeers   int
    ActivePeers  int
    PendingPeers int
}

// ProgressInfo holds download progress information
type ProgressInfo struct {
    Completed   int64
    Total       int64
    Percentage  float64
    Speed       float64
    TimeElapsed float64
    ETA         float64
}