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

func (s *Steam) Friends(ID string) (FriendsResponse, error) {
	url := fmt.Sprintf("http://api.steampowered.com/ISteamUser/GetFriendList/v0001/?key=%s&steamid=%s&relationship=friend", s.Key, ID)
	resp, err := http.Get(url)
	if err != nil {
		return FriendsResponse{}, err
	}

	if resp.StatusCode != 200 {
		return FriendsResponse{}, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return FriendsResponse{}, err
	}

	list := FriendsResponse{}
	json.Unmarshal(b, &list)
	return list, nil
}

func (f FriendsResponse) Oldest() Friend {
	oldest := f.FriendsList.Friends[0]
	for _, current := range f.FriendsList.Friends {
		if current.Relationship == "friend" {
			if current.FriendsSince < oldest.FriendsSince {
				oldest = current
			}
		}
	}

	return oldest
}

func (f FriendsResponse) Newest() Friend {
	newest := f.FriendsList.Friends[0]
	for _, current := range f.FriendsList.Friends {
		if current.Relationship == "friend" {
			if current.FriendsSince > newest.FriendsSince {
				newest = current
			}
		}
	}

	return newest
}

func (f FriendsResponse) Count() int {
	var count int
	for _, friend := range f.FriendsList.Friends {
		if friend.Relationship == "friend" {
			count++
		}
	}
	return count
}
