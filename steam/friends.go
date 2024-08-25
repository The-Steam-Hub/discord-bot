package steam

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Friend struct {
	ID           string `json:"steamid"`
	Relationship string `json:"relationship"`
	FriendsSince int64  `json:"friend_since"`
}

func (s Steam) FriendsList(ID string) (*[]Friend, error) {
	baseURL, _ := url.Parse(SteamWebAPIISteamUser)
	baseURL.Path += "GetFriendList/v0001"

	params := url.Values{}
	params.Add("key", s.Key)
	params.Add("steamid", ID)
	params.Add("format", "json")
	params.Add("relationship", "friend")
	baseURL.RawQuery = params.Encode()

	resp, err := http.Get(baseURL.String())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response struct {
		FriendsList struct {
			Friends []Friend `json:"friends"`
		} `json:"friendslist"`
	}

	json.Unmarshal(b, &response)
	return &response.FriendsList.Friends, nil
}

func FriendsSort(friends []Friend) []Friend {
	length := len(friends)
	for i := 0; i < length-1; i++ {
		for j := i + 1; j < length; j++ {
			left := friends[i].FriendsSince
			right := friends[j].FriendsSince
			if left > right {
				friends[i], friends[j] = friends[j], friends[i]
			}
		}
	}
	return friends
}

func FriendIDs(friends []Friend) []string {
	IDs := make([]string, len(friends))
	for k, v := range friends {
		IDs[k] = v.ID
	}
	return IDs
}
