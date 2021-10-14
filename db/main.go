package db

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/symmetric-project/symphony/env"
	"github.com/symmetric-project/symphony/model"
	"github.com/symmetric-project/symphony/utils"
)

type ParsedStampPrice struct {
	Raw float32 `json:"raw"`
}

type ParsedStampVolume struct {
	Raw float32 `json:"raw"`
}

type ParsedStamp struct {
	Symbol             string           `json:"symbol"`
	PriceChange        ParsedStampPrice `json:"regularMarketPrice"`
	PriceChangePercent ParsedStampPrice `json:"regularMarketChangePercent"`
	Timestamp          int              `json:"timestamp"`
}

type Stamp struct {
	ID                    int32   `json:"id"`
	Symbol                string  `json:"symbol"`
	PriceChange           float32 `json:"priceChange"`
	PriceChangePercentage float32 `json:"priceChangePercentage"`
	Timestamp             int32   `json:"timestamp"`
}

type Subscription struct {
	Endpoint string `json:"endpoint"`
	Auth     string `json:"auth"`
	P256dh   string `json:"p256dh"`
}

type User struct {
	ID               int32  `json:"id"`
	PushSubscription string `json:"pushSubscription"`
	PushInterval     int32  `json:"pushInterval"`
}

var DB *pgxpool.Pool

var SQ sq.StatementBuilderType

func init() {
	SQ = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	var err error
	DB, err = pgxpool.Connect(context.Background(), env.CONFIG.DATABASE_URL)
	if err != nil {
		utils.StacktraceErrorAndExit(err)
	}
}

func GetUserExists(userName string) (bool, error) {
	var exists bool
	err := pgxscan.Get(context.Background(), DB, &exists, `SELECT EXISTS(SELECT 1 from "user" WHERE "name" = $1)`, userName)
	return exists, err
}

func GetUserByName(userName string) (model.User, error) {
	var user model.User
	err := pgxscan.Get(context.Background(), DB, &user, `SELECT * from "user" WHERE "name" = $1`, userName)
	return user, err
}

func GetNodeExists(nodeName string) (bool, error) {
	var exists bool
	err := pgxscan.Get(context.Background(), DB, &exists, `SELECT EXISTS(SELECT 1 from "node" WHERE "name" = $1)`, nodeName)
	return exists, err
}

func AddNode(newNode model.Node) (model.Node, error) {
	var node model.Node
	creationTimestamp := utils.CurrentTimestamp()
	builder := SQ.Insert(`node`).Columns(`name`, `access`, `nsfw`, `creation_timestamp`, `creator_id`).Values(newNode.Name, newNode.Access, newNode.Nsfw, creationTimestamp, "8d05b23f-daff-4cc2-9bac-743c38e2b88e").Suffix(`RETURNING *`)

	query, args, err := builder.ToSql()
	if err != nil {
		return node, err
	}
	err = pgxscan.Get(context.Background(), DB, &node, query, args...)
	if err != nil {
		return node, err
	}
	return node, err
}

func AddUser(newUser model.User) (model.User, error) {
	var user model.User
	builder := SQ.Insert(`"user"`).Columns(`"name"`).Values(newUser.Name).Suffix("RETURNING *")
	query, args, err := builder.ToSql()
	if err != nil {
		return user, err
	}
	err = pgxscan.Get(context.Background(), DB, &user, query, args...)
	return user, err
}

func AddPost(post model.Post) error {
	id := utils.NewOctid()
	slug := utils.Slugify(post.Title)
	creationTimestamp := utils.CurrentTimestamp()

	builder := SQ.Insert(`post`).Columns(`id`, `title`, `link`, `raw_state`, `node_name`, `slug`, `creation_timestamp`, `author_id`).Values(id, post.Title, post.Link, post.RawState, post.NodeName, slug, creationTimestamp, post.AuthorID)
	query, args, err := builder.ToSql()

	if err != nil {
		return err
	}
	_, err = DB.Exec(context.Background(), query, args...)
	return err
}
