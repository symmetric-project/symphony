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
	Author     string `json:"author"`
	CreatedUTC int    `json:"created_utc"`
	Title      string `json:"title"`
	Subreddit  string `json:"subreddit"`
	Selftext   string `json:"selftext"`
}

func NewPost(submissionRow SubmissionRow) model.Post {
	return model.Post{
		Title: submissionRow.Title,
	}
}

func ConvertMarkdownToDraftJSRawState() {
	parts := strings.Fields("node draft.js")
	cmd := exec.Command(parts[0], parts[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
	}
	log.Println(string(output))
}

func main() {
	ConvertMarkdownToDraftJSRawState()
	return
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
		userExists, err := db.GetUserExists(subm.Author)
		if err != nil {
			utils.StacktraceError(err)
		}

		var user model.User
		if !userExists {
			user, err = db.AddUser(model.User{Name: subm.Author})
			if err != nil {
				utils.StacktraceErrorAndExit(err)
			} else {
				utils.LogSuccess("user created " + fmt.Sprint(user))
			}
		}
		user, err = db.GetUserByName(subm.Author)
		if err != nil {
			utils.StacktraceErrorAndExit(err)
		}

		/* utils.LogWarning(user) */

		post := NewPost(subm)
		post.AuthorID = user.ID
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
