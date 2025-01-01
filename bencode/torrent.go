package bencode

import (
	"crypto/sha1"
	"fmt"
	"os"
	"sort"
)

// parseInfoDictionary parses the `info` dictionary of a torrent file
func parseInfoDictionary(infoDict map[string]interface{}) InfoDictionary {
	info := InfoDictionary{}

	if name, ok := infoDict["name"].(string); ok {
		info.Name = name
	}
	if pieceLength, ok := infoDict["piece length"].(int64); ok {
		info.PieceLength = pieceLength
	}
	/*if pieces, ok := infoDict["pieces"].(string); ok {
		info.Pieces = []byte(pieces)
	}*/
	if length, ok := infoDict["length"].(int64); ok {
		info.Length = length
	}
	if files, ok := infoDict["files"].([]interface{}); ok {
		for _, file := range files {
			fileMap := file.(map[string]interface{})
			fileInfo := FileInfo{}
			if length, ok := fileMap["length"].(int64); ok {
				fileInfo.Length = length
			}
			if md5sum, ok := fileMap["md5sum"].(string); ok {
				fileInfo.MD5Sum = md5sum
			}
			if path, ok := fileMap["path"].([]interface{}); ok {
				for _, p := range path {
					fileInfo.Path = append(fileInfo.Path, p.(string))
				}
			}
			info.Files = append(info.Files, fileInfo)
		}
	}

	return info
}

// ParseTorrentFile parses a .torrent file and calculates the info hash
func ParseTorrentFile(filename string) (Torrent, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return Torrent{}, err
	}

	decoder := NewBencodeDecoder(data)
	torrentDict, err := decoder.Decode()
	if err != nil {
		return Torrent{}, err
	}

	torrent := Torrent{}
	root := torrentDict.(map[string]interface{})

	if announce, ok := root["announce"].(string); ok {
		torrent.Announce = announce
	}
	if announceList, ok := root["announce-list"].([]interface{}); ok {
		for _, tier := range announceList {
			var trackerTier []string
			for _, tracker := range tier.([]interface{}) {
				trackerTier = append(trackerTier, tracker.(string))
			}
			torrent.AnnounceList = append(torrent.AnnounceList, trackerTier)
		}
	}
	if createdBy, ok := root["created by"].(string); ok {
		torrent.CreatedBy = createdBy
	}
	if creationDate, ok := root["creation date"].(int64); ok {
		torrent.CreationDate = creationDate
	}
	if comment, ok := root["comment"].(string); ok {
		torrent.Comment = comment
	}

	if infoDict, ok := root["info"].(map[string]interface{}); ok {
		torrent.Info = parseInfoDictionary(infoDict)

		// Calculate the info hash
		bencodedInfo := Bencode(infoDict)
		hash := sha1.Sum([]byte(bencodedInfo))
		torrent.InfoHash = fmt.Sprintf("%x", hash[:])
	}

	// Calculate total size for multi-file torrents
	if len(torrent.Info.Files) > 0 {
		for _, file := range torrent.Info.Files {
			torrent.TotalSize += file.Length
		}
	} else {
		torrent.TotalSize = torrent.Info.Length
	}

	return torrent, nil
}

// bencode encodes a dictionary into a Bencode string
func Bencode(data map[string]interface{}) string {
	var encode func(interface{}) string

	encode = func(v interface{}) string {
		switch v := v.(type) {
		case string:
			return fmt.Sprintf("%d:%s", len(v), v)
		case int64:
			return fmt.Sprintf("i%de", v)
		case map[string]interface{}:
			keys := make([]string, 0, len(v))
			for k := range v {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			result := "d"
			for _, k := range keys {
				result += encode(k) + encode(v[k])
			}
			result += "e"
			return result
		case []interface{}:
			result := "l"
			for _, item := range v {
				result += encode(item)
			}
			result += "e"
			return result
		default:
			panic("unsupported type")
		}
	}

	return encode(data)
}
