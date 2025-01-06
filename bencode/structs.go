package bencode

import (
	"bytes"
	// "time"
)

// Torrent represents the structure of a torrent file
type Torrent struct {
	Announce     string
	AnnounceList [][]string
	CreatedBy    string
	CreationDate int64
	Comment      string
	Info         InfoDictionary
	InfoHash     string
	TotalSize    int64
}

// InfoDictionary represents the `info` section of a torrent file
type InfoDictionary struct {
	Name        string
	Length      int64
	MD5Sum      string
	PieceLength int64
	Pieces      []byte
	Files       []FileInfo
}

// FileInfo represents individual file information in multi-file torrents
type FileInfo struct {
	Length int64
	MD5Sum string
	Path   []string
}

// BencodeDecoder decodes Bencoded data
type BencodeDecoder struct {
	reader *bytes.Reader
}

// NewBencodeDecoder creates a new BencodeDecoder
func NewBencodeDecoder(data []byte) *BencodeDecoder {
	return &BencodeDecoder{reader: bytes.NewReader(data)}
}

func (d *BencodeDecoder) readByte() (byte, error) {
	return d.reader.ReadByte()
}

func (d *BencodeDecoder) unreadByte() error {
	return d.reader.UnreadByte()
}

// TorrentMetadata holds parsed magnet link information
type MagnetMetadata struct {
	InfoHash    []byte
	DisplayName string
	UDPTrackers []string
}

// Protocol constants as defined in BEP 15 (will be used for connecting to peers using the trackers, will be implemented later)
// const (
//     connectionID     = int64(0x41727101980)     // Magic constant for connection request
//     actionConnect    = int32(0)                 // Identifies a connect request/response
//     actionAnnounce   = int32(1)                 // Identifies an announce request/response
//     timeoutDuration  = 15 * time.Second         // Timeout for UDP responses
//     maxRetries       = 1                        // Maximum number of retry attempts
// )
