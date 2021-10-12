package db

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/symmetric-project/symphony/env"
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

func AddSubscription(subscription Subscription) error {
	builder := SQ.Insert(`"subscription"`).Columns(`"endpoint"`, `"auth"`, `"p256dh"`).Values(subscription.Endpoint, subscription.Auth, subscription.P256dh).Suffix(`RETURNING *`)
	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}
	_, err = DB.Exec(context.Background(), query, args...)
	return err
}
