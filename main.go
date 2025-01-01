package main

import (
	"fmt"
	"os"

	bencode "github.com/serene-brew/ztorrent/bencode"
)

func main() {
	torrent, err := bencode.ParseTorrentFile("example.torrent")
	if err != nil {
		fmt.Println("error reading the torrent file")
		os.Exit(1)
	}
	fmt.Println("torrent info: ", torrent)
}
