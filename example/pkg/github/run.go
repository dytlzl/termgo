package github

import (
	"context"
	"fmt"

	"github.com/dytlzl/tervi/pkg/tui"
)

func Run(ctx context.Context) error {
	repoView := NewRepoSearchView()
	codeView := NewCodeSearchView()

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

	mode := "repo"

	footerMessage := ""
	handleEvent := func(event any) any {
		switch typed := event.(type) {
		case FooterMessage:
			footerMessage = typed.Payload
			return nil
		}

		value := func() any {
			if mode == "repo" {
				return repoView.HandleEvent(event)
			} else {
				return codeView.HandleEvent(event)
			}
		}()
		switch typed := value.(type) {
		case string:
			mode = typed
			return nil
		}
		return value
	}

	err := tui.Run(
		func() *tui.View {
			if mode == "repo" {
				return tui.ZStack(
					tui.ZStack(
						tui.ViewWithRenderer(repoView.Body).Title(
							func() string {
								if repoView.Result.Query != "" {
									return fmt.Sprintf("Result of '%s' - Repository Search", repoView.Result.Query)
								} else {
									return "Repository Search"
								}
							}(),
						).Border(tui.Style{F256: 255, B256: 0}),
						func() *tui.View {
							subView := repoView.SubView()
							if subView != nil {
								return tui.ZStack(
									tui.ViewWithRenderer(subView.Body).Title(subView.repo.FullName).RelativeSize(6, 12).Position(6, 0).Border(tui.Style{F256: 255, B256: 0}),
								).Padding(2, 2, 2, 2)
							}
							return nil
						}(),
					).Padding(0, 0, 1, 0),
					tui.ReversedVStack(
						tui.TextView(footerMessage).AbsoluteSize(0, 1).Style(tui.Style{F256: 255, B256: 135}).Padding(0, 1, 0, 1),
					),
				)
			} else {
				return tui.ZStack(
					tui.ZStack(
						tui.ViewWithRenderer(codeView.Body).Title(
							func() string {
								if repoView.Result.Query != "" {
									return fmt.Sprintf("Result of '%s' - Code Search", codeView.Result.Query)
								} else {
									return "Code Search"
								}
							}(),
						).Border(tui.Style{F256: 255, B256: 0}),
						func() *tui.View {
							subView := codeView.SubView()
							if subView != nil {
								return tui.ZStack(
									tui.ViewWithRenderer(subView.Body).Title(subView.item.Path).RelativeSize(6, 12).Position(6, 0).Border(tui.Style{F256: 255, B256: 0}),
								).Padding(2, 2, 2, 2)
							}
							return nil
						}(),
					).Padding(0, 0, 1, 0),
					tui.ReversedVStack(
						tui.TextView(footerMessage).AbsoluteSize(0, 1).Style(tui.Style{F256: 255, B256: 135}).Padding(0, 1, 0, 1),
					),
				)

			}
		},
		tui.OptionEventHandler(handleEvent),
		tui.OptionStyle(tui.Style{F256: 255, B256: 0}),
		tui.OptionChannel(Channel),
	)
	if err != nil {
		return fmt.Errorf("an error has occured while running tui: %w", err)
	}
	close(Finalizers)
	for finalize := range Finalizers {
		finalize()
	}
	return nil
}
