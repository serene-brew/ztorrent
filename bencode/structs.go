package bencode

import (
	"bytes"
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
