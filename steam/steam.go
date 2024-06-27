package steam

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type Steam struct {
	Key string
}

type VanityResponse struct {
	Vanity Vanity `json:"response"`
}

type Vanity struct {
	SteamID string `json:"steamid"`
}

// ResolveID takes a steam profile URL and returns the steam ID
func (s *Steam) ResolveID(url string) string {
	vanityRegex := regexp.MustCompile(`https:\/\/steamcommunity\.com\/id\/([^\/]+)`)
	vanityMatch := vanityRegex.FindStringSubmatch(url)
	IDRegex := regexp.MustCompile(`https:\/\/steamcommunity\.com\/profiles\/(\d+)`)
	IDMatch := IDRegex.FindStringSubmatch(url)
	steamID := ""

	if len(vanityMatch) > 1 {
		vanity, err := s.resolveVanityURL(vanityMatch[1])
		if err != nil {
			return ""
		}
		steamID = vanity.SteamID
	}
	if len(IDMatch) > 1 {
		steamID = IDMatch[1]
	}

	return steamID
}

func (s *Steam) resolveVanityURL(vanityURL string) (Vanity, error) {
	url := fmt.Sprintf("https://api.steampowered.com/ISteamUser/ResolveVanityURL/v1/?key=%s&vanityurl=%s", s.Key, vanityURL)
	resp, err := http.Get(url)
	if err != nil {
		return Vanity{}, err
	}

	if resp.StatusCode != 200 {
		return Vanity{}, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Vanity{}, err
	}

	VanityResponse := VanityResponse{}
	json.Unmarshal(b, &VanityResponse)
	return VanityResponse.Vanity, nil
}
