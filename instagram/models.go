package instagram

import (
	"fmt"
)

type Data struct {
	EntryData EntryData `json:"entry_data"`
}

func (d Data) UserProfile() User {
	return d.EntryData.ProfilePage[0].User
}

type EntryData struct {
	ProfilePage []UserContainer `json:"ProfilePage"`
}

type UserContainer struct {
	User User `json:"user"`
}

type User struct {
	Username   string `json:"username"`
	ProfilePic string `json:"profile_pic_url"`

	Media Media `json:"media"`
}

func (u User) MediaNodes(num int) []Node {
	nodes := u.Media.Nodes

	if len(nodes) > num {
		return nodes[:num]
	}

	return nodes
}

func (u User) ProfileUrl() string {
	return fmt.Sprintf("%s/%s", baseUrl, u.Username)
}

type Media struct {
	Nodes []Node `json:"nodes"`
}

type Node struct {
	Id   string `json:"id"`
	Date int64  `json:"date"`
	Code string `json:"code"`
	Src  string `json:"display_src"`
}

func (n Node) DetailUrl() string {
	return fmt.Sprintf("%s/p/%s", baseUrl, n.Code)
}
