package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/symmetric-project/symphony/db"
	"github.com/symmetric-project/symphony/model"
)

type SubmissionRow struct {
	Author     string `json:"author"`
	CreatedUTC int    `json:"created_utc"`
	Title      string `json:"title"`
	Subreddit  string `json:"subreddit"`
	Selftext   string `json:"selftext"`
}

func NewPost(row SubmissionRow) model.Post {
	return model.Post{
		Title: row.Title,
	}
}

func main() {
	file, err := os.Open("./dumps/RS_2021-06")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	sc := bufio.NewScanner(file)

	const bufCapacity = 64 * 4096
	buf := make([]byte, bufCapacity)
	sc.Buffer(buf, bufCapacity)

	lines := 0
	for sc.Scan() {
		var row SubmissionRow
		err := json.Unmarshal(sc.Bytes(), &row)
		if err != nil {
			log.Println(err)
		} else {
			post := NewPost(row)
			err = db.AddPost(post)
			if err != nil {
				log.Println(err)
			}
		}
		lines++
		if lines > 100 {
			break
		}
	}
}
