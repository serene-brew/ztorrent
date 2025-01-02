package torrent

import (
    "context"
    "fmt"
    "time"

    "github.com/anacrolix/torrent"
)

func GetPeers(magnetURI string) error {
    cfg := torrent.NewDefaultClientConfig()
    cfg.Seed = false
    cfg.Debug = false
    cfg.NoDHT = false

    client, err := torrent.NewClient(cfg)
    if err != nil {
        return fmt.Errorf("failed to create client: %v", err)
    }
    defer client.Close()

    tor, err := client.AddMagnet(magnetURI)
    if err != nil {
        return fmt.Errorf("failed to add magnet: %v", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    select {
    case <-tor.GotInfo():
        fmt.Printf("\nTorrent Info:\n")
        fmt.Printf("Name: %s\n", tor.Name())
        fmt.Printf("Info Hash: %x\n", tor.InfoHash())
        fmt.Printf("Total Length: %d bytes\n", tor.Length())
        
        fmt.Printf("\nPeers:\n")
        stats := tor.Stats()
        fmt.Printf("Total Peers: %d\n", stats.TotalPeers)
        fmt.Printf("Active Peers: %d\n", stats.ActivePeers)
        fmt.Printf("Pending Peers: %d\n", stats.PendingPeers)

        activePeers := tor.PeerConns()
        activeMap := make(map[string]bool)
        for _, ap := range activePeers {
            addr := ap.RemoteAddr
            if addr != nil {
                activeMap[addr.String()] = true
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

    case <-ctx.Done():
        return fmt.Errorf("timeout waiting for torrent info")
    }

    return nil
}