package main

import (
	"context"

	"github.com/dytlzl/tervi/example/pkg/github"
)

func main() {
	github.SetAPIs([]github.API{
		{
			Origin:  "github.com",
			Address: "https://api.github.com",
		},
	})
	github.Run(context.Background())
}
