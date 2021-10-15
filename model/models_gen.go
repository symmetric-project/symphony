// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type Comment struct {
	ID                string `json:"id"`
	PostID            string `json:"postId"`
	PostSlug          string `json:"postSlug"`
	CreationTimestamp int    `json:"creationTimestamp"`
	RawState          string `json:"rawState"`
	AuthorID          string `json:"authorId"`
	Author            *User  `json:"author"`
}

type NewComment struct {
	PostID   string `json:"postId"`
	PostSlug string `json:"postSlug"`
	RawState string `json:"rawState"`
	AuthorID string `json:"authorId"`
}

type NewNode struct {
	Name        string     `json:"name"`
	Tags        []*string  `json:"tags"`
	Access      NodeAccess `json:"access"`
	Nsfw        bool       `json:"nsfw"`
	Description *string    `json:"description"`
}

type NewPost struct {
	Title    string  `json:"title"`
	Link     *string `json:"link"`
	RawState *string `json:"rawState"`
	NodeName string  `json:"nodeName"`
}

type NewUser struct {
	Name *string `json:"name"`
}

type Node struct {
	Name              string     `json:"name"`
	Tags              []*string  `json:"tags"`
	Access            NodeAccess `json:"access"`
	Nsfw              bool       `json:"nsfw"`
	Description       *string    `json:"description"`
	CreationTimestamp int        `json:"creationTimestamp"`
	CreatorID         string     `json:"creatorId"`
}

type Post struct {
	ID                string  `json:"id"`
	Title             string  `json:"title"`
	Link              *string `json:"link"`
	RawState          *string `json:"rawState"`
	NodeName          string  `json:"nodeName"`
	Slug              string  `json:"slug"`
	CreationTimestamp int     `json:"creationTimestamp"`
	AuthorID          string  `json:"authorId"`
	Author            *User   `json:"author"`
	Bases             int     `json:"bases"`
	ThumbnaillURL     *string `json:"thumbnaillUrl"`
	ImageURL          *string `json:"imageUrl"`
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Bases int    `json:"bases"`
}

type NodeAccess string

const (
	NodeAccessPublic     NodeAccess = "PUBLIC"
	NodeAccessRestricted NodeAccess = "RESTRICTED"
	NodeAccessPrivate    NodeAccess = "PRIVATE"
)

var AllNodeAccess = []NodeAccess{
	NodeAccessPublic,
	NodeAccessRestricted,
	NodeAccessPrivate,
}

func (e NodeAccess) IsValid() bool {
	switch e {
	case NodeAccessPublic, NodeAccessRestricted, NodeAccessPrivate:
		return true
	}
	return false
}

func (e NodeAccess) String() string {
	return string(e)
}

func (e *NodeAccess) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = NodeAccess(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid NodeAccess", str)
	}
	return nil
}

func (e NodeAccess) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
