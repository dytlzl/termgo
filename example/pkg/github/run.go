package github

import (
	"context"
	"fmt"

	"github.com/dytlzl/tervi/pkg/color"
	"github.com/dytlzl/tervi/pkg/component"
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
					case SearchInput:
						switch typed.Type {
						case "repo":
							return SearchRepositories(ctx, typed)
						case "code":
							return SearchCode(ctx, typed)
						}
					case ResultItemWithOrigin:
						switch typedResultItem := typed.ResultItem.(type) {
						case Repository:
							return FetchReadMe(ctx, typed.Origin, typedResultItem)
						case CodeSearchResultItem:
							return FetchContent(ctx, typed.Origin, typedResultItem)
						}
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

	isConfirmedTermination := false

	err := tui.Run(
		func() *tui.View {
			title := ""
			view := repoView
			if mode == "repo" {
				title = "Repository Search"
				if repoView.Result.Query != "" {
					title = fmt.Sprintf("Result of '%s' - %s", repoView.Result.Query, title)
				}
				view = repoView
			} else {
				title = "Code Search"
				if codeView.Result.Query != "" {
					title = fmt.Sprintf("Result of '%s' - %s", codeView.Result.Query, title)
				}
				view = codeView
			}
			return tui.VStack(
				tui.ZStack(
					view.Body().Title(title).Border(tui.BorderOptionFGColor(color.RGB(100, 100, 100))),
					tui.HStack(
						tui.Spacer(),
						view.SubView().Border(tui.BorderOptionFGColor(color.RGB(100, 100, 100))),
					).Padding(2),
					tui.If(view.IsSearching,
						tui.ZStack(
							tui.String("Searching...").AbsoluteSize(12, 1),
						).Border(tui.BorderOptionFGColor(color.RGB(100, 100, 100))).AbsoluteSize(20, 5),
						nil,
					),
					component.QuitView(&isQuitMenuOpen, &isConfirmedTermination).FGColor(color.RGB(200, 200, 200)).BGColor(color.RGB(100, 0, 100)),
				),
				tui.TextView(footerMessage).AbsoluteSize(0, 1).BGColor(color.RGB(100, 0, 100)).Padding(0, 1),
			)
		},
		tui.OptionEventHandler(handleEvent),
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

var isQuitMenuOpen = false

var channel = make(chan any, 100)

var requestChannel = make(chan any, 100)
