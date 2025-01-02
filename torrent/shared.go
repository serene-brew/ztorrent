// package torrent

// import (
//     "context"
//     "fmt"
//     "time"
// 	"os"
// 	"net/url"

//     "github.com/anacrolix/torrent"
// 	"github.com/anacrolix/torrent/metainfo"
// )

// func GetPeersFromFile(torrentPath string) error {
//     cfg := torrent.NewDefaultClientConfig()
//     cfg.Seed = false
//     cfg.Debug = false
//     cfg.NoDHT = false
//     cfg.DisablePEX = false // Enable Peer Exchange

//     client, err := torrent.NewClient(cfg)
//     if err != nil {
//         return fmt.Errorf("failed to create client: %v", err)
//     }
//     defer client.Close()

//     // Read the torrent file into a MetaInfo object
//     file, err := os.Open(torrentPath)
//     if err != nil {
//         return fmt.Errorf("failed to open torrent file: %v", err)
//     }
//     defer file.Close()

//     metaInfo, err := metainfo.Load(file)
//     if err != nil {
//         return fmt.Errorf("failed to load torrent file: %v", err)
//     }

//     // Extract the info hash, name, and tracker URLs from the MetaInfo object
//     info, err := metaInfo.UnmarshalInfo()
//     if err != nil {
//         return fmt.Errorf("failed to unmarshal torrent info: %v", err)
//     }
//     infoHash := metaInfo.HashInfoBytes().HexString()
//     name := url.QueryEscape(info.Name)
//     trackers := metaInfo.AnnounceList
//     trackerParams := ""
//     for _, tracker := range trackers {
//         for _, trackerURL := range tracker {
//             trackerParams += "&tr=" + url.QueryEscape(trackerURL)
//         }
//     }

//     // Generate the magnet link using the extracted information
//     magnetLink := fmt.Sprintf("magnet:?xt=urn:btih:%s&dn=%s%s", infoHash, name, trackerParams)

//     // Use the AddMagnet function to add the torrent
//     tor, err := client.AddMagnet(magnetLink)
//     if err != nil {
//         return fmt.Errorf("failed to add magnet: %v", err)
//     }

//     ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute) // Increase timeout
//     defer cancel()

//     select {
//     case <-tor.GotInfo():
//         fmt.Printf("\nTorrent Info:\n")
//         fmt.Printf("Name: %s\n", tor.Name())
//         fmt.Printf("Info Hash: %x\n", tor.InfoHash())
//         fmt.Printf("Total Length: %d bytes\n", tor.Length())
        
//         fmt.Printf("\nPeers:\n")
//         stats := tor.Stats()
//         fmt.Printf("Total Peers: %d\n", stats.TotalPeers)
//         fmt.Printf("Active Peers: %d\n", stats.ActivePeers)
//         fmt.Printf("Pending Peers: %d\n", stats.PendingPeers)

//         activePeers := tor.PeerConns()
//         activeMap := make(map[string]bool)
//         for _, ap := range activePeers {
//             addr := ap.RemoteAddr
//             if addr != nil {
//                 activeMap[addr.String()] = true
//             }
//         }

//         fmt.Printf("\nPeer List:\n")
//         for _, peer := range tor.KnownSwarm() {
//             if peer.Addr != nil {
//                 addr := peer.Addr.String()
//                 status := "[x]" // Inactive
//                 if activeMap[addr] {
//                     status = "[+]" // Active
//                 }
//                 fmt.Printf("- %s %s\n", status, addr)
//             }
//         }

// 		fmt.Println(magnetLink)

//     	case <-ctx.Done():
//         	return fmt.Errorf("timeout waiting for torrent info")
//     }

//     return nil
// }

// func GetPeers(magnetURI string) error {
//     cfg := torrent.NewDefaultClientConfig()
//     cfg.Seed = false
//     cfg.Debug = false
//     cfg.NoDHT = false
// 	cfg.DisablePEX = false

//     client, err := torrent.NewClient(cfg)
//     if err != nil {
//         return fmt.Errorf("failed to create client: %v", err)
//     }
//     defer client.Close()

//     tor, err := client.AddMagnet(magnetURI)
//     if err != nil {
//         return fmt.Errorf("failed to add magnet: %v", err)
//     }

//     ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//     defer cancel()

//     select {
//     case <-tor.GotInfo():
//         fmt.Printf("\nTorrent Info:\n")
//         fmt.Printf("Name: %s\n", tor.Name())
//         fmt.Printf("Info Hash: %x\n", tor.InfoHash())
//         fmt.Printf("Total Length: %d bytes\n", tor.Length())
        
//         fmt.Printf("\nPeers:\n")
//         stats := tor.Stats()
//         fmt.Printf("Total Peers: %d\n", stats.TotalPeers)
//         fmt.Printf("Active Peers: %d\n", stats.ActivePeers)
//         fmt.Printf("Pending Peers: %d\n", stats.PendingPeers)

//         activePeers := tor.PeerConns()
//         activeMap := make(map[string]bool)
//         for _, ap := range activePeers {
//             addr := ap.RemoteAddr
//             if addr != nil {
//                 activeMap[addr.String()] = true
//             }
//         }

//         fmt.Printf("\nPeer List:\n")
//         for _, peer := range tor.KnownSwarm() {
//             if peer.Addr != nil {
//                 addr := peer.Addr.String()
//                 status := "[x]" // Inactive
//                 if activeMap[addr] {
//                     status = "[+]" // Active
//                 }
//                 fmt.Printf("- %s %s\n", status, addr)
//             }
//         }

//     case <-ctx.Done():
//         return fmt.Errorf("timeout waiting for torrent info")
//     }

//     return nil
// }


package torrent

import (
    "context"
    "fmt"
    "time"
    "github.com/anacrolix/torrent"
)

// createDefaultClient creates a new torrent client with default configuration
func createDefaultClient() (*torrent.Client, error) {
    cfg := torrent.NewDefaultClientConfig()
    cfg.Seed = false
    cfg.Debug = false
    cfg.NoDHT = false
    cfg.DisablePEX = false
    return torrent.NewClient(cfg)
}

// GetPeersFromFile retrieves peer information from a torrent file
func GetPeersFromFile(torrentPath string) error {
    client, err := createDefaultClient()
    if err != nil {
        return fmt.Errorf("failed to create client: %v", err)
    }
    defer client.Close()

    magnetLink, err := generateMagnetFromFile(torrentPath)
    if err != nil {
        return err
    }

    tor, err := client.AddMagnet(magnetLink)
    if err != nil {
        return fmt.Errorf("failed to add magnet: %v", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
    defer cancel()

    // Start progress monitoring
    go monitorProgress(ctx, tor)

    select {
    case <-tor.GotInfo():
        printTorrentInfo(tor)
        fmt.Printf("\nMagnet Link: %s\n", magnetLink)
    case <-ctx.Done():
        return fmt.Errorf("timeout waiting for torrent info")
    }

    return nil
}

// GetPeers retrieves peer information from a magnet link
func GetPeers(magnetURI string) error {
    client, err := createDefaultClient()
    if err != nil {
        return fmt.Errorf("failed to create client: %v", err)
    }
    defer client.Close()

    tor, err := client.AddMagnet(magnetURI)
    if err != nil {
        return fmt.Errorf("failed to add magnet: %v", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
    defer cancel()

    // Start progress monitoring
    go monitorProgress(ctx, tor)

    select {
    case <-tor.GotInfo():
        printTorrentInfo(tor)
    case <-ctx.Done():
        return fmt.Errorf("timeout waiting for torrent info")
    }

    return nil
}

// monitorProgress monitors and displays torrent progress
func monitorProgress(ctx context.Context, tor *torrent.Torrent) {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            stats := tor.Stats()
            if stats.TotalPeers > 0 {
                fmt.Printf("\rPeers: %d Active, %d Total", 
                    stats.ActivePeers, stats.TotalPeers)
            }
        case <-ctx.Done():
            return
        }
    }
}