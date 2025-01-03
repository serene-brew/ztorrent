package main

import (
//	"fmt"
//	"os"
	entrypoint "github.com/serene-brew/ztorrent/interfaces"
//	bencode "github.com/serene-brew/ztorrent/bencode"
//	crawler "github.com/serene-brew/ztorrent/crawler"
  //  mag "github.com/serene-brew/ztorrent/torrent"
)

func main() {
	//-----------------------------------------------------------------------------
	//torrent file parsing and peers extraction
/*	torrent, err := bencode.ParseTorrentFile("example.torrent")
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
/*	query := "Assassins Creed unity"
	data, err := crawler.GetInfoMediaQuery(query)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, record := range data {
		for _, terms := range record{
			fmt.Println(terms)
		}
	}*/
	//magnet := crawler.GetMagnet(data[0][2].(string), data[0][1].(string))

	//-----------------------------------------------------------------------------
    // crawler magnet URI parsing and peers extraction 
/*	magnetLink := magnet
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

    fmt.Println(magnet)*/

	//-----------------------------------------------------------------------------
	entrypoint.Entrypoint()	
    
}



