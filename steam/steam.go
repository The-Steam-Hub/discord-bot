package steam

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type Vanity struct {
	SteamID string `json:"steamid"`
}

var (
	VanityURLRegex = `https:\/\/steamcommunity\.com\/id\/([^\/]+)`
	IDURLREgex     = `https:\/\/steamcommunity\.com\/profiles\/(\d+)`
)

func (s Steam) ResolveID(input string) (string, error) {
	if _, err := strconv.ParseUint(input, 10, 64); err == nil {
		return input, nil
	}

	if strings.HasPrefix(input, "[U:1:") {
		return SteamID3ToSteamID64(input)
	}

	if strings.HasPrefix(input, "STEAM_") {
		return SteamIDToSteamID64(input)
	}

	if strings.HasPrefix(input, SteamCommunity) {
		return s.resolveIDFromURL(input), nil
	}

	vanityURL, err := s.resolveVanityURL(input)
	return vanityURL.SteamID, err
}

func (s Steam) resolveIDFromURL(url string) string {
	vanityRegex := regexp.MustCompile(VanityURLRegex)
	vanityMatch := vanityRegex.FindStringSubmatch(url)

	IDRegex := regexp.MustCompile(IDURLREgex)
	IDMatch := IDRegex.FindStringSubmatch(url)

	var steamID string
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

func (s Steam) resolveVanityURL(vanityURL string) (Vanity, error) {
	baseURL, _ := url.Parse(SteamAPIISteamUser)
	baseURL.Path += "ResolveVanityURL/v1"

	params := url.Values{}
	params.Add("key", s.Key)
	params.Add("vanityurl", vanityURL)
	params.Add("format", "json")
	baseURL.RawQuery = params.Encode()

	resp, err := http.Get(baseURL.String())
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

	var response struct {
		Vanity struct {
			SteamID string `json:"steamid"`
		} `json:"response"`
	}

	json.Unmarshal(b, &response)
	return response.Vanity, nil
}

func SteamID64ToSteamID(steamID64 uint64) string {
	universe := (steamID64 >> 56) & 0xFF
	accountID := steamID64 & 0xFFFFFFFF
	authServer := accountID % 2
	accountNumber := (accountID - authServer) / 2
	return fmt.Sprintf("STEAM_%d:%d:%d", universe, authServer, accountNumber)
}

func SteamID64ToSteamID3(steamID64 uint64) string {
	accountID := steamID64 & 0xFFFFFFFF
	return fmt.Sprintf("[U:1:%d]", accountID)
}

func SteamID3ToSteamID64(steamID3 string) (string, error) {
	var accountID uint32
	_, err := fmt.Sscanf(steamID3, "[U:1:%d]", &accountID)
	if err != nil {
		return "", err
	}
	steamID64 := uint64(0x110000100000000) | uint64(accountID)
	return strconv.FormatUint(steamID64, 10), nil
}

func SteamIDToSteamID64(steamID string) (string, error) {
	var universe, authServer, accountNumber uint32
	_, err := fmt.Sscanf(steamID, "STEAM_%d:%d:%d", &universe, &authServer, &accountNumber)
	if err != nil {
		return "", err
	}
	accountID := accountNumber*2 + authServer
	steamID64 := (uint64(universe) << 56) | (1 << 52) | (1 << 32) | uint64(accountID)
	return strconv.FormatUint(steamID64, 10), nil
}
