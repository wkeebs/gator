#!/usr/bin/zsh
rm build/*
go build -C src -o ../build/aggregator 
./build/aggregator