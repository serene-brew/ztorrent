package crawler

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

type Torrent struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	InfoHash string `json:"info_hash"`
	Leechers string `json:"leechers"`
	Seeders  string `json:"seeders"`
	NumFiles string `json:"num_files"`
	Size     string `json:"size"`
	Username string `json:"username"`
	Added    string `json:"added"`
	Status   string `json:"status"`
	Category string `json:"category"`
	Imdb     string `json:"imdb"`
}

func GetInfoMediaQuery(query string) ([][]interface{}, error) {
	agent := "Mozilla/5.0 (X11; Linux x86_64; rv:99.0) Gecko/20100101 Firefox/99.0"
	encodedQuery := url.QueryEscape(query)
	apiURL := fmt.Sprintf("https://apibay.org/q.php?q=%s", encodedQuery)

	client := &http.Client{}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", agent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var torrents []Torrent
	if err := json.Unmarshal(body, &torrents); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}
	var result [][]interface{}
	for _, t := range torrents {
		result = append(result, []interface{}{
			t.ID, t.Name, t.InfoHash, t.Seeders, t.Leechers,
			t.NumFiles, t.Size, t.Category,
		})
	}

	return result, nil
}

func ClassifyCategory(categoryID string) string {
	category_ID_string := string(categoryID[0])
	category_ID, _ := strconv.Atoi(category_ID_string)
	if category_ID < len(category)-1 {
		category_ID = category_ID
	} else {
		category_ID = 6
	}

	return category[category_ID]
}

func ConvertSize(sizeBytes int64) string {
	if sizeBytes == 0 {
		return "0B"
	}
	sizeNames := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	i := int(math.Floor(math.Log(float64(sizeBytes)) / math.Log(1024)))
	p := math.Pow(1024, float64(i))
	s := math.Round(float64(sizeBytes)/p*100) / 100

	return fmt.Sprintf("%.2f %s", s, sizeNames[i])
}

func GetMagnet(info_hash string, name string) string {
	return "magnet:?xt=urn:btih:" + info_hash + "&dn=" + name + "&tr=" + GenTrackerStub()
}
