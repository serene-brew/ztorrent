package bencode

import (
    "encoding/hex"
    "fmt"
    "net/url"
    "strings"
)

func ExtractUDPTrackers(trackers []string) []string {
    var udpTrackers []string
    for _, tracker := range trackers {
        if strings.HasPrefix(tracker, "udp://") {
            cleaned := strings.TrimPrefix(tracker, "udp://")
            cleaned = strings.TrimSuffix(cleaned, "/announce")
            udpTrackers = append(udpTrackers, cleaned)
        }
    }
    return udpTrackers
}

func ParseMagnetLink(magnetURI string) (*MagnetMetadata, error) {
    uri, err := url.Parse(magnetURI)
    if err != nil {
        return nil, fmt.Errorf("invalid magnet URI: %v", err)
    }

    params := uri.Query()
    xt := params.Get("xt")
    if !strings.HasPrefix(xt, "urn:btih:") {
        return nil, fmt.Errorf("missing or invalid xt parameter")
    }

    infoHashHex := strings.TrimPrefix(xt, "urn:btih:")
    infoHash, err := hex.DecodeString(infoHashHex)
    if err != nil {
        return nil, fmt.Errorf("invalid info hash: %v", err)
    }

    trackers := params["tr"]
    udpTrackers := ExtractUDPTrackers(trackers)

    metadata := &MagnetMetadata{
        InfoHash:    infoHash,
        DisplayName: params.Get("dn"),
        UDPTrackers: udpTrackers,
    }

    return metadata, nil
}