package main

import (
	"github.com/wkeebs/gator/internal/config"
	"github.com/wkeebs/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}
