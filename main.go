package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/symmetric-project/symphony/db"
	"github.com/symmetric-project/symphony/model"
	"github.com/symmetric-project/symphony/utils"
)

type SubmissionRow struct {
	Author      string  `json:"author"`
	CreatedUTC  int     `json:"created_utc"`
	Score       int     `json:"score"`
	Subreddit   string  `json:"subreddit"`
	Selftext    string  `json:"selftext"`
	Title       string  `json:"title"`
	Thumbnail   string  `json:"thumbnail"`
	UpvoteRatio float32 `json:"upvote_ratio"`
	URL         string  `json:"url"`
}

func ConvertMarkdownToDraftState(mardown string) (string, error) {
	parts := strings.Fields(`node convert-markdown-to-draft-state.js ` + mardown)
	cmd := exec.Command(parts[0], parts[1:]...)
	output, err := cmd.CombinedOutput()
	log.Println(string(output))
	if err != nil {
		return "", err
	}
	return string(output), nil
}

/* func ConvertDraftStateToRawState(draftState string) (string, error) {
	parts := strings.Fields(`node convert-draft-state-to-raw-state.js ` + draftState)
	cmd := exec.Command(parts[0], parts[1:]...)
	output, err := cmd.CombinedOutput()
	log.Println(string(output))
	if err != nil {
		return "", err
	}
	return string(output), err
} */

/* func ConvertMarkdownToRawState(markdown string) (*string, error) {
	draftState, err := ConvertMarkdownToDraftState(markdown)
	if err != nil {
		utils.StacktraceErrorAndExit(err)
	}
	rawState, err := ConvertDraftStateToRawState(draftState)
	return &rawState, err
} */

func NewPost(subm SubmissionRow, user model.User) model.Post {
	log.Println(user)
	rawState, err := ConvertMarkdownToDraftState(subm.Selftext)
	if err != nil {
		utils.StacktraceErrorAndExit(err)
	}
	id := utils.NewOctid()
	creationTimestamp := utils.CurrentTimestamp()
	post := model.Post{
		ID:                id,
		Title:             subm.Title,
		RawState:          nil,
		NodeName:          subm.Subreddit,
		Slug:              utils.Slugify(subm.Title),
		CreationTimestamp: creationTimestamp,
		AuthorID:          user.ID,
		ThumbnaillURL:     nil,
		ImageURL:          nil,
		Bases:             subm.Score,
	}
	if len(rawState) > 0 {
		post.RawState = &rawState
	}
	if strings.Contains(subm.Thumbnail, ".jpg") {
		post.ThumbnaillURL = &subm.Thumbnail
	}
	if strings.Contains(subm.URL, ".jpg") {
		post.ImageURL = &subm.URL
	}
	return post
}

func NewNode(subm SubmissionRow) model.Node {
	return model.Node{
		Name:      subm.Subreddit,
		Access:    model.NodeAccessPublic,
		CreatorID: "6da8de06-a065-4831-aa82-8015eece9573",
	}
}

func ParseSubreddit(subredditName string) {
	file, err := os.Open("./dumps/RS_2021-06")
	if err != nil {
		utils.StacktraceErrorAndExit(err)
	}
	defer file.Close()

	sc := bufio.NewScanner(file)

	const bufCapacity = 64 * 4096
	buf := make([]byte, bufCapacity)
	sc.Buffer(buf, bufCapacity)

	lines := 0
	for sc.Scan() {
		var subm SubmissionRow
		err := json.Unmarshal(sc.Bytes(), &subm)
		if err != nil {
			utils.StacktraceError(err)
			continue
		}
		if subm.Subreddit != subredditName {
			continue
		}

		var user model.User
		userExists, err := db.GetUserExists(subm.Author)
		if err != nil {
			utils.StacktraceErrorAndExit(err)
		}
		if !userExists {
			user, err = db.AddUser(model.User{Name: subm.Author})
			if err != nil {
				utils.StacktraceErrorAndExit(err)
			} else {
				utils.LogSuccess("New user created " + fmt.Sprint(user))
			}
		} else {
			user, err = db.GetUserByName(subm.Author)
			if err != nil {
				utils.StacktraceErrorAndExit(err)
			}
		}

		/* var node model.Node */
		nodeExists, err := db.GetNodeExists(subm.Subreddit)
		if err != nil {
			utils.StacktraceErrorAndExit(err)
		}
		if !nodeExists {
			creationTimestamp := utils.CurrentTimestamp()
			node := model.Node{
				Name:              subm.Subreddit,
				Access:            model.NodeAccessPublic,
				Nsfw:              false,
				CreationTimestamp: creationTimestamp,
				CreatorID:         "6da8de06-a065-4831-aa82-8015eece9573",
			}
			_, err = db.AddNode(node)
			if err != nil {
				utils.StacktraceErrorAndExit(err)
			} else {
				utils.LogSuccess("New node created " + fmt.Sprint(node))
			}
		}

		/* utils.LogWarning(user) */

		post := NewPost(subm, user)
		err = db.AddPost(post)
		if err != nil {
			utils.StacktraceError(err)
		}
		lines++
		log.Println(lines)
		if lines > 100 {
			break
		}
	}
}

func main() {
	ParseSubreddit("NoNewNormal")
}
