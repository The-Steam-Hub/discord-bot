package steam

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type FriendsResponse struct {
	FriendsList FriendsList `json:"friendslist"`
}

type FriendsList struct {
	Friends []Friend `json:"friends"`
}

type Friend struct {
	ID           string `json:"steamid"`
	Relationship string `json:"relationship"`
	FriendsSince int64  `json:"friend_since"`
}

func (s *Steam) GetFriendsList(ID string) ([]Friend, error) {
	url := fmt.Sprintf("http://api.steampowered.com/ISteamUser/GetFriendList/v0001/?key=%s&steamid=%s&relationship=friend", s.Key, ID)
	resp, err := http.Get(url)
	if err != nil {
		return []Friend{}, err
	}

	if resp.StatusCode != 200 {
		return []Friend{}, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Friend{}, err
	}

	freindsList := FriendsResponse{}
	json.Unmarshal(b, &freindsList)
	return freindsList.FriendsList.Friends, nil
}

func SortFriends(friends []Friend) []Friend {
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

func GetFriendsIDs(friends []Friend) []string {
	IDs := make([]string, len(friends))
	for k, v := range friends {
		IDs[k] = v.ID
	}
	return IDs
}
