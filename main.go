package main

import (
	"fmt"
	"os"

	bencode "github.com/serene-brew/ztorrent/bencode"
	crawler "github.com/serene-brew/ztorrent/crawler"
    mag "github.com/serene-brew/ztorrent/torrent"
)

func main() {
	torrent, err := bencode.ParseTorrentFile("../example.torrent")
	if err != nil {
		fmt.Println("error reading the torrent file")
		os.Exit(1)
	}
	fmt.Println("torrent info: ", torrent)


    fmt.Println(crawler.GenTrackerStub())

	query := "Assassins Creed unity"
	data, err := crawler.GetInfoMediaQuery(query)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the result
	for _, record := range data {
		fmt.Println(record)
		for i := range record{
			if i == 7 {
				fmt.Println(crawler.ClassifyCategory(record[i].(string)))
			} 
		}
	}
	magnet := crawler.GetMagnet(data[0][2].(string), data[0][1].(string))
	fmt.Println(magnet)


	// magnetLink := "magnet:?xt=urn:btih:4ED8248102AC7DA2578B106C95B708B2648F176F&dn=DeadPool+and+Wolverine+2024+1080p+V1+Clean+HDTS+H264+COLLECTiVE&tr=http%3A%2F%2Fp4p.arenabg.com%3A1337%2Fannounce&tr=udp%3A%2F%2F47.ip-51-68-199.eu%3A6969%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2780%2Fannounce&tr=udp%3A%2F%2F9.rarbg.to%3A2710%2Fannounce&tr=udp%3A%2F%2F9.rarbg.to%3A2730%2Fannounce&tr=udp%3A%2F%2F9.rarbg.to%3A2920%2Fannounce&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Fopentracker.i2p.rocks%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.cyberia.is%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.dler.org%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.openbittorrent.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337&tr=udp%3A%2F%2Ftracker.pirateparty.gr%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.tiny-vps.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451%2Fannounce"
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

    fmt.Println(magnet)

}



