package main

import (
	"context"
	"fmt"

	"github.com/dytlzl/tervi/example/internal/github"
	"github.com/dytlzl/tervi/pkg/tui"
)

func main() {
	Run(context.Background())
}

func Run(ctx context.Context) error {
	repoView := github.InitRepoSearchView()
	codeView := github.InitCodeSearchView()

	// Search Loop
	go func() {
		for {
			select {
			case query := <-repoView.SearchInputCh:
				go github.SearchRepositories(ctx, query, github.Channel)
			case query := <-repoView.ReadMeInputCh:
				go github.FetchReadMe(ctx, query, github.Channel)
			case query := <-codeView.SearchInputCh:
				go github.SearchCode(ctx, query, github.Channel)
			case item := <-codeView.ContentInputCh:
				go github.FetchContent(ctx, item, github.Channel)
			}
		}
	}()

	err := tui.Run(map[string]tui.View{
		"repo": repoView,
		"code": codeView,
	}, tui.Options{
		DefaultViewName: "repo",
		Style:           tui.CellStyle{F256: 255, B256: 0},
		Footer:          &github.Footer{},
	}, github.Channel)
	if err != nil {
		return fmt.Errorf("an error has occured while running tui: %w", err)
	}
	close(github.Finalizers)
	for finalize := range github.Finalizers {
		finalize()
	}
	return nil
}
