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
		for item := range requestChannel {
			item := item
			go func() {
				SendToChan(func() (any, error) {
					switch typed := item.(type) {
					case RepositorySearchInput:
						return SearchRepositories(ctx, typed)
					case RepositoryWithOrigin:
						return FetchReadMe(ctx, typed)
					case CodeSearchInput:
						return SearchCode(ctx, typed)
					case SearchResultItem:
						return FetchContent(ctx, typed)
					}
					return nil, nil
				}())
			}()
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
			title := ""
			body := func(tui.Size) []tui.Text { return nil }
			subViewTitle := ""
			subViewBody := func(tui.Size) []tui.Text { return nil }
			if mode == "repo" {
				title = "Repository Search"
				if repoView.Result.Query != "" {
					title = fmt.Sprintf("Result of '%s' - %s", repoView.Result.Query, title)
				}
				body = repoView.Body
				subView := repoView.SubView()
				if subView != nil {
					subViewTitle = subView.repo.FullName
					subViewBody = subView.Body
				}
			} else {
				body = codeView.Body
				subView := codeView.SubView()
				title = "Code Search"
				if codeView.Result.Query != "" {
					title = fmt.Sprintf("Result of '%s' - %s", codeView.Result.Query, title)
				}
				if subView != nil {
					subViewTitle = subView.item.Path
					subViewBody = subView.Body
				}
			}
			return tui.VStack(
				tui.ZStack(
					tui.ViewWithRenderer(body).Title(title).Border(tui.Style{F256: 255, B256: 0}),
					tui.HStack(
						tui.Spacer(),
						tui.ViewWithRenderer(subViewBody).Title(subViewTitle).Border(tui.Style{F256: 255, B256: 0}).Hidden(subViewTitle == ""),
					).Padding(2, 2, 2, 2),
				),
				tui.TextView(footerMessage).AbsoluteSize(0, 1).Style(tui.Style{F256: 15, B256: 135}).Padding(0, 1, 0, 1),
			)
		},
		tui.OptionEventHandler(handleEvent),
		tui.OptionStyle(tui.Style{F256: 255, B256: 0}),
		tui.OptionChannel(channel),
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

var channel = make(chan any, 100)

var requestChannel = make(chan any, 100)
