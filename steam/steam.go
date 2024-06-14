package steam

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func (s *Steam) ResolveVanityURL(vanityURL string) (Vanity, error) {
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
	fmt.Println(string(b))

	VanityResponse := VanityResponse{}
	json.Unmarshal(b, &VanityResponse)
	return VanityResponse.Vanity, nil
}
