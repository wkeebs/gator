#!/bin/zsh
go build -C src -o ../build/gator
./build/gator "$@"