package main

import (
	// "fmt"

	"fmt"

	"github.com/wkeebs/rss-blog-aggregator/internal/config"
)

func main() {
	configStruct, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(configStruct)

	configStruct.SetUser("Will")

	fmt.Println(configStruct)

}
