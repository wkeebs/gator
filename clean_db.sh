#!/bin/zsh
cd ./sql/schema
goose postgres "postgres://will:@localhost:5432/gator" down
goose postgres "postgres://will:@localhost:5432/gator" up