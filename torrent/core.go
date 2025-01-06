package torrent

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func GetDefaultDownloadPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join("/home", os.Getenv("USER"), "Downloads")
	}
	return filepath.Join(homeDir, "Downloads")
}

func DownloadFromMagnet(magnetURI string, downloadPath string) (<-chan ProgressInfo, error) {
	if downloadPath == "" {
		downloadPath = GetDefaultDownloadPath()
	}

	if err := os.MkdirAll(downloadPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create downloads directory: %v", err)
	}

	client, err := createTorrentClient(downloadPath)
	if err != nil {
		return nil, fmt.Errorf("client creation failed: %v", err)
	}

	progress := make(chan ProgressInfo)

	go func() {
		defer client.Close()
		defer close(progress)

		tor, err := client.AddMagnet(magnetURI)
		if err != nil {
			return
		}

		<-tor.GotInfo()
		tor.DownloadAll()

		startTime := time.Now()

		for {
			completed := tor.BytesCompleted()
			total := tor.Length()

			progress <- GetProgressInfo(tor, startTime)

			if completed == total {
				return
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()

	return progress, nil
}

func GetPeers(magnetURI string) (*TorrentInfo, []PeerInfo, error) {
	client, err := createTorrentClient("")
	if err != nil {
		return nil, nil, fmt.Errorf("client creation failed: %v", err)
	}
	defer client.Close()

	tor, err := client.AddMagnet(magnetURI)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to add magnet: %v", err)
	}

	<-tor.GotInfo()

	torrentInfo, err := GetTorrentInfo(tor)
	if err != nil {
		return nil, nil, err
	}

	peerInfo := GetPeerInfo(tor)
	return torrentInfo, peerInfo, nil
}

func GetPeersFromFile(torrentPath string) (*TorrentInfo, []PeerInfo, error) {
	client, err := createTorrentClient("")
	if err != nil {
		return nil, nil, fmt.Errorf("client creation failed: %v", err)
	}
	defer client.Close()

	magnetLink, err := generateMagnetFromFile(torrentPath)
	if err != nil {
		return nil, nil, fmt.Errorf("magnet generation failed: %v", err)
	}

	return GetPeers(magnetLink)
}
