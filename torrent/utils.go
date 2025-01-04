package torrent

import (
    "fmt"
    "net/url"
    "strings"
    "time"
    // "context"

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

func HumanReadableSize(bytes int64) string {
    const unit = 1024
    if bytes < unit {
        return fmt.Sprintf("%d B", bytes)
    }
    div, exp := int64(unit), 0
    for n := bytes / unit; n >= unit; n /= unit {
        div *= unit
        exp++
    }
    return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func getFileExtension(filename string) string {
    parts := strings.Split(filename, ".")
    if len(parts) > 1 {
        return parts[len(parts)-1]
    }
    return "unknown"
}

func joinPath(parts []string) string {
    return strings.Join(parts, "/")
}

// DEEZ FUCKERIES CAN BE USED LATER IF NECESSARY OR WILL BE ERADICATED FROM EXISTENCE
// // Progress tracking and monitoring
// func monitorDownloadProgress(tor *torrent.Torrent) error {
//     ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
//     defer cancel()

//     for {
//         select {
//         case <-ctx.Done():
//             return fmt.Errorf("download timeout")
//         default:
//             completed := tor.BytesCompleted()
//             total := tor.Length()
            
//             if completed == total {
//                 fmt.Printf("\nDownload completed successfully\n")
//                 return nil
//             }
//             time.Sleep(500 * time.Millisecond)
//         }
//     }
// }

// // downloadTorrent handles the torrent download process
// func downloadTorrent(client *torrent.Client, magnetURI string) error {
//     tor, err := client.AddMagnet(magnetURI)
//     if err != nil {
//         return fmt.Errorf("failed to add magnet: %v", err)
//     }

//     <-tor.GotInfo()
//     tor.DownloadAll()
    
//     return monitorDownloadProgress(tor)
// }

// GetTorrentInfo extracts comprehensive torrent information
func GetTorrentInfo(tor *torrent.Torrent) (*TorrentInfo, error) {
    info := tor.Info()
    if info == nil {
        return nil, fmt.Errorf("no torrent info available")
    }

    torInfo := &TorrentInfo{
        Name:       info.Name,
        InfoHash:   tor.InfoHash().String(),
        TotalSize:  info.TotalLength(),
        FilesByExt: make(map[string][]FileInfo),
    }

    if len(info.Files) == 0 {
        // Single file torrent
        torInfo.Files = append(torInfo.Files, FileInfo{
            Name: info.Name,
            Size: info.Length,
            Type: getFileExtension(info.Name),
            Path: info.Name,
        })
    } else {
        // Multiple files torrent
        for _, file := range info.Files {
            path := append([]string{info.Name}, file.Path...)
            fullPath := joinPath(path)
            fileInfo := FileInfo{
                Name: file.Path[len(file.Path)-1],
                Size: file.Length,
                Type: getFileExtension(file.Path[len(file.Path)-1]),
                Path: fullPath,
            }
            torInfo.Files = append(torInfo.Files, fileInfo)
            
            // Group by extension
            ext := fileInfo.Type
            torInfo.FilesByExt[ext] = append(torInfo.FilesByExt[ext], fileInfo)
        }
    }

    return torInfo, nil
}

// GetPeerInfo extracts peer connection information
func GetPeerInfo(tor *torrent.Torrent) []PeerInfo {
    var peers []PeerInfo
    stats := tor.Stats()
    
    activePeers := tor.PeerConns()
    activeMap := make(map[string]bool)
    for _, ap := range activePeers {
        if ap.RemoteAddr != nil {
            activeMap[ap.RemoteAddr.String()] = true
        }
    }

    // THIS SHIT IS SO ASS I DONT UNDERSTAND SHIT, ACCORDING TO DOCUMENTATION KnownSwarm returns "KNOWN subset of peers (active, inactive, half-open, full-open) THOUGH IT DOESNT MAKE SENSE HOW IT WORKS"
    for _, peer := range tor.KnownSwarm() {
        if peer.Addr != nil {
            addr := peer.Addr.String()
            peers = append(peers, PeerInfo{
                Address: addr,
                Active:  activeMap[addr],
                Stats: PeerStats{
                    TotalPeers:   stats.TotalPeers,
                    ActivePeers:  stats.ActivePeers,
                    PendingPeers: stats.PendingPeers,
                },
            })
        }
    }

    return peers
}

// GetProgressInfo calculates current download progress
func GetProgressInfo(tor *torrent.Torrent, startTime time.Time) ProgressInfo {
    completed := tor.BytesCompleted()
    total := tor.Length()
    elapsed := time.Since(startTime).Seconds()
    speed := float64(completed) / elapsed
    
    var eta float64
    if speed > 0 {
        eta = float64(total-completed) / speed
    }

    return ProgressInfo{
        Completed:   completed,
        Total:       total,
        Percentage:  float64(completed) * 100 / float64(total),
        Speed:       speed,
        TimeElapsed: elapsed,
        ETA:         eta,
    }
}