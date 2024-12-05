package main

import (
	"github.com/wkeebs/rss-blog-aggregator/internal/config"
	"github.com/wkeebs/rss-blog-aggregator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}
