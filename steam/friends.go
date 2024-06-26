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

func (s *Steam) Friends(ID string) (FriendsList, error) {
	url := fmt.Sprintf("http://api.steampowered.com/ISteamUser/GetFriendList/v0001/?key=%s&steamid=%s&relationship=friend", s.Key, ID)
	resp, err := http.Get(url)
	if err != nil {
		return FriendsList{}, err
	}

	if resp.StatusCode != 200 {
		return FriendsList{}, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return FriendsList{}, err
	}

	freindsList := FriendsResponse{}
	json.Unmarshal(b, &freindsList)
	return freindsList.FriendsList, nil
}

func (f FriendsList) Oldest() Friend {
	if len(f.Friends) > 0 {
		oldest := f.Friends[0]
		for _, current := range f.Friends {
			if current.Relationship == "friend" {
				if current.FriendsSince < oldest.FriendsSince {
					oldest = current
				}
			}
		}
		return oldest
	}
	return Friend{}
}

func (f FriendsList) Newest() Friend {
	if len(f.Friends) > 0 {
		newest := f.Friends[0]
		for _, current := range f.Friends {
			if current.Relationship == "friend" {
				if current.FriendsSince > newest.FriendsSince {
					newest = current
				}
			}
		}
		return newest
	}
	return Friend{}
}

func (f FriendsList) Count() int {
	var count int
	for _, friend := range f.Friends {
		if friend.Relationship == "friend" {
			count++
		}
	}
	return count
}
