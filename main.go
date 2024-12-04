package main

import (
	// "fmt"

	"fmt"

	"github.com/wkeebs/rss-blog-aggregator/internal/config"
)

const CONFIG_FILENAME = "gatorconfig.json"

func main() {
	configData, err := config.Read(CONFIG_FILENAME)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(configData)
}
