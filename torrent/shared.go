// package torrent

// import (
//     "context"
//     "fmt"
//     "time"
//     "github.com/anacrolix/torrent"
//     "strings"
//     "os"
// )

// // createDefaultClient creates a new torrent client with default configuration
// func createDefaultClient() (*torrent.Client, error) {
//     cfg := torrent.NewDefaultClientConfig()
//     cfg.Seed = false
//     cfg.Debug = false
//     cfg.NoDHT = false
//     cfg.DisablePEX = false
//     cfg.ListenPort = 0 // System assigned port as per available range

//     // Add retry logic
//     var client *torrent.Client
//     var err error
//     for retries := 0; retries < 3; retries++ {
//         client, err = torrent.NewClient(cfg)
//         if err == nil {
//             return client, nil
//         }
//         time.Sleep(time.Second)
//     }
//     return nil, fmt.Errorf("failed to create client after retries: %v", err)
// }

// // GetPeersFromFile retrieves peer information from a torrent file
// func GetPeersFromFile(torrentPath string) error {
//     client, err := createDefaultClient()
//     if err != nil {
//         return fmt.Errorf("failed to create client: %v", err)
//     }
//     defer client.Close()

//     magnetLink, err := generateMagnetFromFile(torrentPath)
//     if err != nil {
//         return err
//     }

//     tor, err := client.AddMagnet(magnetLink)
//     if err != nil {
//         return fmt.Errorf("failed to add magnet: %v", err)
//     }

//     ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
//     defer cancel()

//     // Start progress monitoring
//     go monitorProgress(ctx, tor)

//     select {
//     case <-tor.GotInfo():
//         fmt.Printf("\nTorrent Info:\n")
//         fmt.Printf("Name: %s\n", tor.Name())
//         fmt.Printf("Info Hash: %x\n", tor.InfoHash())
//         fmt.Printf("Total Length: %d bytes\n", tor.Length())

//         printFileInfo(tor)  // Add this line
//         printPeerInfo(tor)
//         fmt.Printf("\nMagnet Link: %s\n", magnetLink)
//     case <-ctx.Done():
//         return fmt.Errorf("timeout waiting for torrent info")
//     }

//     return nil
// }

// // GetPeers retrieves peer information from a magnet link
// func GetPeers(magnetURI string) error {
//     client, err := createDefaultClient()
//     if err != nil {
//         return fmt.Errorf("failed to create client: %v", err)
//     }
//     defer client.Close()

//     tor, err := client.AddMagnet(magnetURI)
//     if err != nil {
//         return fmt.Errorf("failed to add magnet: %v", err)
//     }

//     ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
//     defer cancel()

//     // Start progress monitoring
//     go monitorProgress(ctx, tor)

//     select {
//     case <-tor.GotInfo():
//         fmt.Printf("\nTorrent Info:\n")
//         fmt.Printf("Name: %s\n", tor.Name())
//         fmt.Printf("Info Hash: %x\n", tor.InfoHash())
//         fmt.Printf("Total Length: %d bytes\n", tor.Length())

//         printDetailedFileInfo(tor)  // Add this line
//         printPeerInfo(tor)
//     case <-ctx.Done():
//         return fmt.Errorf("timeout waiting for torrent info")
//     }

//     return nil
// }

// // monitorProgress monitors and displays torrent progress
// func monitorProgress(ctx context.Context, tor *torrent.Torrent) {
//     ticker := time.NewTicker(1 * time.Second)
//     defer ticker.Stop()

//     for {
//         select {
//         case <-ticker.C:
//             stats := tor.Stats()
//             if stats.TotalPeers > 0 {
//                 fmt.Printf("\rPeers: %d Active, %d Total",
//                     stats.ActivePeers, stats.TotalPeers)
//             }
//         case <-ctx.Done():
//             return
//         }
//     }
// }

// func printFileInfo(tor *torrent.Torrent) {
//     info := tor.Info()
//     if info == nil {
//         fmt.Println("No file information available")
//         return
//     }

//     fmt.Printf("\nFiles Available:\n")
//     if len(info.Files) == 0 {
//         // Single file torrent
//         fmt.Printf("- %s (%d bytes)\n", info.Name, info.Length)
//         return
//     }

//     // Multiple files torrent
//     for _, file := range info.Files {
//         path := append([]string{info.Name}, file.Path...)
//         fmt.Printf("- %s (%d bytes)\n",
//             joinPath(path),
//             file.Length)
//     }
// }

// func joinPath(parts []string) string {
//     return strings.Join(parts, "/")
// }

// func humanReadableSize(bytes int64) string {
//     const unit = 1024
//     if bytes < unit {
//         return fmt.Sprintf("%d B", bytes)
//     }
//     div, exp := int64(unit), 0
//     for n := bytes / unit; n >= unit; n /= unit {
//         div *= unit
//         exp++
//     }
//     return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
// }

// func getFileExtension(filename string) string {
//     parts := strings.Split(filename, ".")
//     if len(parts) > 1 {
//         return parts[len(parts)-1]
//     }
//     return "unknown"
// }

// func printDetailedFileInfo(tor *torrent.Torrent) {
//     info := tor.Info()
//     if info == nil {
//         fmt.Println("No file information available")
//         return
//     }

//     fmt.Printf("\n=== Files Information ===\n")

//     if len(info.Files) == 0 {
//         // Single file torrent
//         ext := getFileExtension(info.Name)
//         fmt.Printf("Single File Torrent:\n")
//         fmt.Printf("- Name: %s\n", info.Name)
//         fmt.Printf("- Size: %s\n", humanReadableSize(info.Length))
//         fmt.Printf("- Type: %s\n", ext)
//         return
//     }

//     // Multiple files torrent
//     fmt.Printf("Multiple Files Torrent (%d files):\n", len(info.Files))
//     var totalSize int64

//     // Group files by extension
//     filesByExt := make(map[string][]string)
//     sizeByExt := make(map[string]int64)

//     for _, file := range info.Files {
//         path := append([]string{info.Name}, file.Path...)
//         fullPath := joinPath(path)
//         ext := getFileExtension(file.Path[len(file.Path)-1])

//         filesByExt[ext] = append(filesByExt[ext], fullPath)
//         sizeByExt[ext] += file.Length
//         totalSize += file.Length

//         fmt.Printf("\n- File: %s\n", fullPath)
//         fmt.Printf("  Size: %s\n", humanReadableSize(file.Length))
//         fmt.Printf("  Type: %s\n", ext)
//     }

//     // Print summary
//     fmt.Printf("\n=== Summary ===\n")
//     fmt.Printf("Total Size: %s\n", humanReadableSize(totalSize))
//     fmt.Printf("File Types:\n")
//     for ext, files := range filesByExt {
//         fmt.Printf("- %s: %d files (%s)\n",
//             ext,
//             len(files),
//             humanReadableSize(sizeByExt[ext]))
//     }
// }

// func DownloadFromMagnet(magnetURI string) error {
//     downloadDir := "downloads"
//     if err := os.MkdirAll(downloadDir, 0755); err != nil {
//         return fmt.Errorf("failed to create downloads directory: %v", err)
//     }

//     cfg := torrent.NewDefaultClientConfig()
//     cfg.Seed = false
//     cfg.Debug = false
//     cfg.DataDir = downloadDir

//     client, err := torrent.NewClient(cfg)
//     if err != nil {
//         return fmt.Errorf("failed to create client: %v", err)
//     }
//     defer client.Close()

//     tor, err := client.AddMagnet(magnetURI)
//     if err != nil {
//         return fmt.Errorf("failed to add magnet: %v", err)
//     }

//     <-tor.GotInfo()

//     tor.DownloadAll()

//     startTime := time.Now()
//     ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
//     defer cancel()

//     for {
//         select {
//         case <-ctx.Done():
//             return fmt.Errorf("download timeout")
//         default:
//             completed := tor.BytesCompleted()
//             total := tor.Length()

//             printProgress(completed, total, startTime)

//             if tor.BytesCompleted() == tor.Length() {
//                 fmt.Printf("\nDownload completed successfully\n")
//                 return nil
//             }
//             time.Sleep(500 * time.Millisecond)
//         }
//     }
// }

// func printProgress(completed, total int64, startTime time.Time) {
//     width := 50
//     percentage := float64(completed) * 100 / float64(total)
//     filled := int(float64(width) * float64(completed) / float64(total))

//     // Calculate speed
//     elapsed := time.Since(startTime).Seconds()
//     speed := float64(completed) / elapsed // bytes per second

//     // Calculate ETA
//     var eta float64
//     if speed > 0 {
//         eta = float64(total-completed) / speed
//     }

//     // Create progress bar
//     bar := strings.Repeat("=", filled) + strings.Repeat(" ", width-filled)

//     // Format speed and ETA
//     var speedStr, etaStr string
//     if speed < 1024 {
//         speedStr = fmt.Sprintf("%.0f B/s", speed)
//     } else if speed < 1024*1024 {
//         speedStr = fmt.Sprintf("%.1f KB/s", speed/1024)
//     } else {
//         speedStr = fmt.Sprintf("%.1f MB/s", speed/1024/1024)
//     }

//     if eta > 0 {
//         etaStr = fmt.Sprintf("ETA: %ds", int(eta))
//     } else {
//         etaStr = "ETA: --"
//     }

//     fmt.Printf("\r[%s] %.1f%% %s %s", bar, percentage, speedStr, etaStr)
// }

package torrent

import (
	"fmt"
	"github.com/anacrolix/torrent"
	"time"
)

func createTorrentClient(dataDir string) (*torrent.Client, error) {
	cfg := torrent.NewDefaultClientConfig()
	cfg.Seed = false
	cfg.Debug = false
	cfg.NoDHT = false
	cfg.DisablePEX = false
	cfg.ListenPort = 0
	if dataDir != "" {
		cfg.DataDir = dataDir
	}

	var client *torrent.Client
	var err error
	for retries := 0; retries < 3; retries++ {
		client, err = torrent.NewClient(cfg)
		if err == nil {
			return client, nil
		}
		time.Sleep(time.Second)
	}
	return nil, fmt.Errorf("failed to create client after retries: %v", err)
}
