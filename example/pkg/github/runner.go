package github

import (
	"context"
	"fmt"

	"github.com/dytlzl/tervi/pkg/tui"
)

func Run(ctx context.Context) error {
	repoView := InitRepoSearchView()
	codeView := InitCodeSearchView()

	// Search Loop
	go func() {
		for {
			select {
			case query := <-repoView.SearchInputCh:
				go SendToChan(SearchRepositories(ctx, query))
			case query := <-repoView.ReadMeInputCh:
				go SendToChan(FetchReadMe(ctx, query))
			case query := <-codeView.SearchInputCh:
				go SendToChan(SearchCode(ctx, query))
			case item := <-codeView.ContentInputCh:
				go SendToChan(FetchContent(ctx, item))
			}
		}
	}()

	err := tui.Run(map[string]tui.View{
		"repo": repoView,
		"code": codeView,
	}, tui.Options{
		DefaultViewName: "repo",
		Style:           tui.CellStyle{F256: 255, B256: 0},
		Footer:          &Footer{},
	}, Channel)
	if err != nil {
		return fmt.Errorf("an error has occured while running tui: %w", err)
	}
	close(Finalizers)
	for finalize := range Finalizers {
		finalize()
	}
	return nil
}
