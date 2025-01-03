package torrent

import (
	"fmt"
	"net/url"

	"github.com/anacrolix/torrent"
	bencode "github.com/serene-brew/ztorrent/bencode"
)

// generateMagnetFromFile generates a magnet link from a torrent file
func generateMagnetFromFile(torrentPath string) (string, error){
	torrent, err := bencode.ParseTorrentFile(torrentPath)
	if err != nil{
		fmt.Println("error reading torrent file")
		return "", err
	}
	name := url.QueryEscape(torrent.Info.Name)
	infoHash := torrent.InfoHash
	trackerStub := ""
	for _, trackerArr := range torrent.AnnounceList{
		for _, tracker := range trackerArr{
			trackerStub += "&tr="+url.QueryEscape(tracker)				
		}
	}
	
	magnetLink := fmt.Sprintf("magnet:?xt=urn:btih:%s&dn=%s%s", infoHash, name, trackerStub)
	return magnetLink, nil
	
}

// printPeerInfo prints peer statistics and list
func printPeerInfo(tor *torrent.Torrent) {
    stats := tor.Stats()
    fmt.Printf("\nPeers:\n")
    fmt.Printf("Total Peers: %d\n", stats.TotalPeers)
    fmt.Printf("Active Peers: %d\n", stats.ActivePeers)
    fmt.Printf("Pending Peers: %d\n", stats.PendingPeers)

    activePeers := tor.PeerConns()
    activeMap := make(map[string]bool)
    for _, ap := range activePeers {
        if ap.RemoteAddr != nil {
            activeMap[ap.RemoteAddr.String()] = true
        }
    }

    fmt.Printf("\nPeer List:\n")
    for _, peer := range tor.KnownSwarm() {
        if peer.Addr != nil {
            addr := peer.Addr.String()
            status := "[x]" // Inactive
            if activeMap[addr] {
                status = "[+]" // Active
            }
            fmt.Printf("- %s %s\n", status, addr)
        }
    }
}
