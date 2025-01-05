package main

import (
	"fmt"
	bencode "github.com/serene-brew/ztorrent/bencode"
	crawler "github.com/serene-brew/ztorrent/crawler"
	mag "github.com/serene-brew/ztorrent/torrent"
	"os"
	"strings"
)

func main() {
	// Torrent file parsing and peers extraction
	torrent, err := bencode.ParseTorrentFile("example.torrent")
	if err != nil {
		fmt.Println("Error reading the torrent file")
		os.Exit(1)
	}

	// Display torrent metadata
	fmt.Println("[-] info hash: ", torrent.InfoHash)
	fmt.Println("[-] trackers Array: ", torrent.AnnounceList)
	fmt.Println("[-] torrent file name: ", torrent.Info.Name)
	fmt.Println("[-] torrent files: ", torrent.Info.Files)
	fmt.Println("[-] total download size: ", torrent.TotalSize)

	// Get peers from torrent file
	torInfo, peers, err := mag.GetPeersFromFile("example.torrent")
	if err != nil {
		fmt.Println("Error getting peers:", err)
		os.Exit(1)
	}

	// Display torrent information
	fmt.Printf("\n=== Torrent Information ===\n")
	fmt.Printf("Name: %s\n", torInfo.Name)
	fmt.Printf("Info Hash: %s\n", torInfo.InfoHash)
	fmt.Printf("Total Size: %s\n", mag.HumanReadableSize(torInfo.TotalSize))

	// Display files
	fmt.Printf("\nFiles:\n")
	for _, file := range torInfo.Files {
		fmt.Printf("- %s (%s)\n", file.Path, mag.HumanReadableSize(file.Size))
	}

	// Display peer information
	fmt.Printf("\n=== Peer Information ===\n")
	fmt.Printf("Total Peers: %d\n", len(peers))
	for _, peer := range peers {
		status := "inactive"
		if peer.Active {
			status = "active"
		}
		fmt.Printf("- %s (%s)\n", peer.Address, status)
	}

	// Crawler test script
	query := "terminator 1"
	data, err := crawler.GetInfoMediaQuery(query)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Display crawler results
	for _, record := range data {
		for _, terms := range record {
			fmt.Println(terms)
		}
	}

	// Get magnet link from crawler data
	magnet := crawler.GetMagnet(data[0][2].(string), data[0][1].(string))

	// Parse magnet link
	metadata, err := bencode.ParseMagnetLink(magnet)
	if err != nil {
		fmt.Println("Error parsing magnet link:", err)
		os.Exit(1)
	}

	// Display trackers
	for _, tracker := range metadata.UDPTrackers {
		fmt.Printf("  %s\n", tracker)
	}

	// Get peers from magnet
	torInfo, peers, err = mag.GetPeers(magnet)
	if err != nil {
		fmt.Println("Error getting peers:", err)
		os.Exit(1)
	}

	// Display magnet torrent information
	fmt.Printf("\n=== Magnet Torrent Information ===\n")
	fmt.Printf("Name: %s\n", torInfo.Name)
	fmt.Printf("Total Size: %s\n", mag.HumanReadableSize(torInfo.TotalSize))
	fmt.Printf("Files: %d\n", len(torInfo.Files))
	fmt.Printf("Peers: %d\n", len(peers))

	// Start download with monitoring deez nuts
	fmt.Printf("\nStarting download...\n")
	progress, err := mag.DownloadFromMagnet(magnet)
	if err != nil {
		fmt.Println("Error starting download:", err)
		os.Exit(1)
	}

	// Monitor download progress
	for p := range progress {
		fmt.Printf("\r[%s] %.1f%% %.1f MB/s ETA: %.0fs",
			getProgressBar(p.Percentage),
			p.Percentage,
			p.Speed/1024/1024,
			p.ETA)
	}
	fmt.Println("\nDownload completed!")
}

// TEMPORARY DRIVER FUNCTION JUST TO CHECK IF SHIT WORKSSSS
func getProgressBar(percentage float64) string {
	width := 50
	filled := int(float64(width) * percentage / 100)
	return strings.Repeat("=", filled) + strings.Repeat(" ", width-filled)
}
