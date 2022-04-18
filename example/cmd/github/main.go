package main

import (
	"context"

	"github.com/dytlzl/tervi/example/pkg/github"
)

func main() {
	github.Run(context.Background())
}
