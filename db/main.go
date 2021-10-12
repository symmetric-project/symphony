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

func GetUserExists(userId string) (bool, error) {
	var exists bool
	var user model.User
	builder := SQ.Select(`"user"`).Where("id = $1", userId).Suffix("RETURNING id")
	query, args, err := builder.ToSql()
	if err != nil {
		return exists, err
	}
	err = pgxscan.Get(context.Background(), DB, &user, query, args...)
	if user.ID == userId {
		exists = true
	}
	return exists, err
}

func AddUser(user model.User) error {
	builder := SQ.Insert(`"user"`).Columns(`"name"`).Values(user.Name)
	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}
	err = pgxscan.Get(context.Background(), DB, &user, query, args...)
	return err
}

func AddPost(post model.Post) error {
	id := utils.NewOctid()
	slug := utils.Slugify(post.Title)
	creationTimestamp := utils.CurrentTimestamp()

	var link *string
	var deltaOps *string

	if post.Link != nil {
		link = post.Link
	} else {
		deltaOps = post.DeltaOps
	}

	builder := SQ.Insert(`post`).Columns(`id`, `title`, `link`, `delta_ops`, `node_name`, `slug`, `creation_timestamp`, `author_id`).Values(id, post.Title, link, deltaOps, post.NodeName, slug, creationTimestamp, "asddas")
	query, args, err := builder.ToSql()

	if err != nil {
		return err
	}
	err = pgxscan.Get(context.Background(), DB, &post, query, args...)
	return err
}
