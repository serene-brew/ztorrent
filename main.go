package main

// Standard library imports
import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	// local package imports
	bencode "github.com/serene-brew/ztorrent/bencode"
	crawler "github.com/serene-brew/ztorrent/crawler"
	mag "github.com/serene-brew/ztorrent/torrent"
	// interfaces "github.com/serene-brew/ztorrent/interfaces"
)

// main serves as a test harness for the torrent functionality for now
// It demonstrates the following capabilities:
// 1. Torrent file parsing
// 2. Peer discovery
// 3. Magnet link handling
// 4. Download functionality
// Note: This will be replaced with a TUI interface in the final version
// the TUI entrypoint will be coded later into main.go
func main() {
	// SECTION 1: Torrent File Processing
	// Parse a local torrent file to extract metadata
	torrent, err := bencode.ParseTorrentFile("example.torrent")
	if err != nil {
		fmt.Println("Error reading the torrent file")
		os.Exit(1)
	}

	// Display basic torrent metadata for verification
	fmt.Println("[-] info hash: ", torrent.InfoHash)            // Unique identifier for the torrent
	fmt.Println("[-] trackers Array: ", torrent.AnnounceList)   // List of tracker URLs
	fmt.Println("[-] torrent file name: ", torrent.Info.Name)   // Name of the torrent
	fmt.Println("[-] torrent files: ", torrent.Info.Files)      // List of files in the torrent
	fmt.Println("[-] total download size: ", torrent.TotalSize) // Total size of all files

	// SECTION 2: Peer Discovery
	// Get peer information from the torrent file
	torInfo, peers, err := mag.GetPeersFromFile("example.torrent")
	if err != nil {
		fmt.Println("Error getting peers:", err)
		os.Exit(1)
	}

	// Display detailed torrent information
	fmt.Printf("\n=== Torrent Information ===\n")
	fmt.Printf("Name: %s\n", torInfo.Name)
	fmt.Printf("Info Hash: %s\n", torInfo.InfoHash)
	fmt.Printf("Total Size: %s\n", mag.HumanReadableSize(torInfo.TotalSize))

	// Display individual file information
	fmt.Printf("\nFiles:\n")
	for _, file := range torInfo.Files {
		fmt.Printf("- %s (%s)\n", file.Path, mag.HumanReadableSize(file.Size))
	}

	// Display peer connection information
	fmt.Printf("\n=== Peer Information ===\n")
	fmt.Printf("Total Peers: %d\n", len(peers))
	for _, peer := range peers {
		status := "inactive"
		if peer.Active {
			status = "active"
		}
		fmt.Printf("- %s (%s)\n", peer.Address, status)
	}

	// SECTION 3: Magnet Link Processing
	// Test the crawler functionality with a search query
	query := "terminator 1"
	data, err := crawler.GetInfoMediaQuery(query)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Display raw crawler results
	for _, record := range data {
		for _, terms := range record {
			fmt.Println(terms)
		}
	}

	// Generate magnet link from crawler data
	// data[0][2] contains the info hash
	// data[0][1] contains the name
	magnet := crawler.GetMagnet(data[0][2].(string), data[0][1].(string))

	// Parse and validate the magnet link
	metadata, err := bencode.ParseMagnetLink(magnet)
	if err != nil {
		fmt.Println("Error parsing magnet link:", err)
		os.Exit(1)
	}

	// Display available trackers from magnet link
	for _, tracker := range metadata.UDPTrackers {
		fmt.Printf("  %s\n", tracker)
	}

	// SECTION 4: Download Process
	// Get peer information from magnet link
	torInfo, peers, err = mag.GetPeers(magnet) // core.go line 64 GetPeers function
	if err != nil {
		fmt.Println("Error getting peers:", err)
		os.Exit(1)
	}

	// Display magnet torrent information
	fmt.Printf("\n=== Magnet Torrent Information ===\n")
	fmt.Printf("Name: %s\n", torInfo.Name)
	fmt.Printf("Total Size: %s\n", mag.HumanReadableSize(torInfo.TotalSize)) // utils.go line 35 humanreadablesize function
	fmt.Printf("Files: %d\n", len(torInfo.Files))
	fmt.Printf("Peers: %d\n", len(peers))

	// Start download process
	fmt.Printf("\nStarting download...\n")

	// Set custom download path in user's home directory
	customPath := filepath.Join(mag.GetDefaultDownloadPath(), "torrents") // core.go line 10 GetDefaultDownloadPath function

	// Initialize download and get progress channel
	progress, err := mag.DownloadFromMagnet(magnet, customPath) // for default path, keep second parameter as ""
	if err != nil {
		fmt.Println("Error starting download:", err)
		os.Exit(1)
	}

	// Display download location
	fmt.Println(customPath) // just for showing it

	// Monitor and display download progress
	for p := range progress {
		fmt.Printf("\r[%s] %.1f%% %.1f MB/s ETA: %s",
			getProgressBar(p.Percentage),
			p.Percentage,
			p.Speed/1024/1024,
			formatETA(p.ETA))
	}
	fmt.Println("\nDownload completed!")

	// interfaces.Entrypoint()
}

// getProgressBar generates a visual progress bar string
// percentage: download progress percentage (0-100)
// returns: a string representing the progress bar [=====    ]
func getProgressBar(percentage float64) string {
	width := 50 // Width of the progress bar in characters
	filled := int(float64(width) * percentage / 100)
	return strings.Repeat("=", filled) + strings.Repeat(" ", width-filled)
}

// some abomination which can be used to format ETA (ETA is float64 by default but this function changes that to string)
// formatETA converts seconds to a human readable duration string
func formatETA(seconds float64) string {
	if seconds < 0 {
		return "calculating..."
	}

	secs := int(seconds)

	// Break down into time units
	months := secs / (30 * 24 * 3600) // Approximate months
	secs %= (30 * 24 * 3600)

	weeks := secs / (7 * 24 * 3600)
	secs %= (7 * 24 * 3600)

	days := secs / (24 * 3600)
	secs %= (24 * 3600)

	hours := secs / 3600
	secs %= 3600

	minutes := secs / 60
	secs %= 60

	// Format based on largest unit
	switch {
	case months > 0:
		return fmt.Sprintf("%dmo%dw", months, weeks)
	case weeks > 0:
		return fmt.Sprintf("%dw%dd", weeks, days)
	case days > 0:
		return fmt.Sprintf("%dd%dh", days, hours)
	case hours > 0:
		return fmt.Sprintf("%dh%02dm", hours, minutes)
	case minutes > 0:
		return fmt.Sprintf("%dm%02ds", minutes, secs)
	default:
		return fmt.Sprintf("%ds", secs)
	}
}
