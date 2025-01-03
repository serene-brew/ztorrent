package main

import (
	"fmt"
	"os"
	"strings"
	// entrypoint "github.com/serene-brew/ztorrent/interfaces"
	bencode "github.com/serene-brew/ztorrent/bencode"
	crawler "github.com/serene-brew/ztorrent/crawler"
   	mag "github.com/serene-brew/ztorrent/torrent"
)

func main() {
	//-----------------------------------------------------------------------------
	//torrent file parsing and peers extraction
	torrent, err := bencode.ParseTorrentFile("example.torrent")
	if err != nil {
		fmt.Println("error reading the torrent file")
		os.Exit(1)
	}
	fmt.Println("[-] info hash: ", torrent.InfoHash)
	fmt.Println("[-] trackers Array: ", torrent.AnnounceList)
	fmt.Println("[-] torrent file name: ", torrent.Info.Name)
	fmt.Println("[-] torrent files: ", torrent.Info.Files)
	fmt.Println("[-] total download size: ", torrent.TotalSize)
	
	err = mag.GetPeersFromFile("example.torrent")
    if err != nil {
        fmt.Println("Error getting peers:", err)
        os.Exit(1)
    }
	//-----------------------------------------------------------------------------
    // crawler test script*/ 
	query := "terminator 1"
	data, err := crawler.GetInfoMediaQuery(query)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, record := range data {
		for _, terms := range record{
			fmt.Println(terms)
		}
	}
	magnet := crawler.GetMagnet(data[0][2].(string), data[0][1].(string))

	//-----------------------------------------------------------------------------
    // crawler magnet URI parsing and peers extraction 
	magnetLink := magnet
    metadata, err := bencode.ParseMagnetLink(magnetLink)
    if err != nil {
        fmt.Println("Error parsing magnet link:", err)
        os.Exit(1)
    }

    for _, tracker := range metadata.UDPTrackers {
        fmt.Printf("  %s\n", tracker)
    }

    err = mag.GetPeers(magnetLink)
    if err != nil {
        fmt.Println("Error getting peers:", err)
        os.Exit(1)
    }

	// Ask user for file selection
    fmt.Println("\nEnter file extensions to download (comma-separated, e.g., mp4,mkv):")
    var extInput string
    fmt.Scanln(&extInput)
    
    extensions := []string{}
    if extInput != "" {
        extensions = strings.Split(extInput, ",")
    }

    fmt.Println("Enter maximum file size in MB (0 for no limit):")
    var maxSizeMB int64
    fmt.Scanln(&maxSizeMB)

	selection := mag.FileSelection{
        Extensions: extensions,
		MaxSize:    maxSizeMB * 1024 * 1024,
    }

    err = mag.DownloadSelectedFilesFromMagnet(magnet, selection)
    if err != nil {
        fmt.Println("Error downloading files:", err)
        os.Exit(1)
    }

    fmt.Println(magnet)

	//-----------------------------------------------------------------------------
	// entrypoint.Entrypoint()	
    
}



