package torrent

import (
    "fmt"
    "net/url"

    "github.com/anacrolix/torrent"
    "github.com/anacrolix/torrent/metainfo"
)

// generateMagnetFromFile generates a magnet link from a torrent file
func generateMagnetFromFile(torrentPath string) (string, error) {
    metaInfo, err := metainfo.LoadFromFile(torrentPath)
    if err != nil {
        return "", fmt.Errorf("failed to load torrent file: %v", err)
    }

    info, err := metaInfo.UnmarshalInfo()
    if err != nil {
        return "", fmt.Errorf("failed to unmarshal torrent info: %v", err)
    }

    infoHash := metaInfo.HashInfoBytes().HexString()
    name := url.QueryEscape(info.Name)
    
    // Handle trackers
    trackerParams := ""
    if len(metaInfo.AnnounceList) == 0 && metaInfo.Announce != "" {
        trackerParams = "&tr=" + url.QueryEscape(metaInfo.Announce)
    } else {
        for _, tracker := range metaInfo.AnnounceList {
            for _, trackerURL := range tracker {
                trackerParams += "&tr=" + url.QueryEscape(trackerURL)
            }
        }
    }

    return fmt.Sprintf("magnet:?xt=urn:btih:%s&dn=%s%s", infoHash, name, trackerParams), nil
}

// printTorrentInfo prints torrent metadata and peer information
func printTorrentInfo(tor *torrent.Torrent) {
    fmt.Printf("\nTorrent Info:\n")
    fmt.Printf("Name: %s\n", tor.Name())
    fmt.Printf("Info Hash: %x\n", tor.InfoHash())
    fmt.Printf("Total Length: %d bytes\n", tor.Length())
    
    printPeerInfo(tor)
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