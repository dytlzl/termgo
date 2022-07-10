package github

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/dytlzl/tervi/pkg/color"
	"github.com/dytlzl/tervi/pkg/component"
	"github.com/dytlzl/tervi/pkg/key"
	"github.com/dytlzl/tervi/pkg/tui"
)

const channelSize = 64

type SearchView struct {
	Type              string
	Result            searchResult
	IsSearching       bool
	selectedItem      int
	input             string
	position          int
	lastStrokeTime    time.Time
	lastInput         string
	ContentMap        map[string]string
	ContentRequestMap map[string]bool
}

func NewRepoSearchView() *SearchView {
	return &SearchView{
		Type:              "repo",
		ContentMap:        map[string]string{},
		ContentRequestMap: map[string]bool{},
	}
}

func NewCodeSearchView() *SearchView {
	return &SearchView{
		Type:              "code",
		ContentMap:        map[string]string{},
		ContentRequestMap: map[string]bool{},
	}
}

func (m *SearchView) Body() *tui.View {
	slice := make([]*tui.View, 0, len(m.Result.Items)+10)
	if m.Result.Query != "" {
		if m.selectedItem >= len(m.Result.Items) {
			m.selectedItem = len(m.Result.Items) - 1
		}
		if m.selectedItem < 0 {
			m.selectedItem = 0
		}
		slice = append(slice, tui.Break())
		lastOrigin := ""
		width, _, _ := tui.TermSize()
		for i, item := range m.Result.Items {
			if item.Origin != lastOrigin {
				slice = append(slice, tui.String(" "+item.Origin+":\n").FGColor(8))
				lastOrigin = item.Origin
			}
			slice = append(slice, tui.Fmt("%s  #%d ", tui.If(i == m.selectedItem, ">", " "), i).FGColor(8))
			if m.Type == "repo" {
				repo := item.ResultItem.(Repository)
				slice = append(slice, tui.Fmt("%s", repo.FullName).If(i == m.selectedItem, (*tui.View).Underline))
			} else {
				item := item.ResultItem.(CodeSearchResultItem)
				path := item.Path
				if width/2-len(item.Repository.FullName)-15 < 0 {
					path = ""
				} else if len(path) > width/2-len(item.Repository.FullName)-15 {
					for len(path) > width/2-len(item.Repository.FullName)-15 {
						_, size := utf8.DecodeLastRuneInString(path)
						path = path[:len(path)-size]
					}
					path += "..."
				}
				slice = append(slice,
					tui.String(item.Repository.FullName).FGColor(225).If(i == m.selectedItem, (*tui.View).Underline),
					tui.Fmt(" %s", path).If(i == m.selectedItem, (*tui.View).Underline),
				)
			}
			slice = append(slice, tui.Break())
		}
	}
	var title string
	if m.Type == "repo" {
		title = "Repository Search"
	} else {
		title = "Code Search"
	}
	if m.Result.Query != "" {
		title = fmt.Sprintf("Result of '%s' - %s", m.Result.Query, title)
	}
	return tui.VStack(
		tui.HStack(
			tui.String("Query > ").FGColor(color.RGB(200, 0, 200)).AbsoluteSize(8, 1),
			component.TextInput(&m.input, &m.position, func() {
				m.lastStrokeTime = time.Now()
			}),
		).AbsoluteSize(0, 1),
		tui.InlineStack(slice...),
	).Title(title)
}

func (m *SearchView) HandleEvent(event any) any {
	switch typed := event.(type) {
	case searchResult:
		if typed.CreatedAt.UnixMicro() > m.Result.CreatedAt.UnixMicro() {
			m.Result = typed
		}
	case ReadMeResult:
		m.ContentMap[typed.HtmlUrl] = typed.ReadMe
	case ContentResult:
		m.ContentMap[typed.Url] = typed.Content
	case rune:
		switch typed {
		case key.CtrlS:
			if m.Type == "repo" {
				return "code"
			} else {
				return "repo"
			}
		case key.CtrlL:
			if m.Type == "repo" {
				if m.selectedItem < len(m.Result.Items) {
					go OpenRepository(m.Result.Items[m.selectedItem].ResultItem.(Repository).HtmlUrl)
				}
			}
		case key.Enter:
			if m.selectedItem < len(m.Result.Items) {
				if m.Type == "repo" {
					go OpenUrl(m.Result.Items[m.selectedItem].ResultItem.(Repository).HtmlUrl)
				} else {
					go OpenUrl(m.Result.Items[m.selectedItem].ResultItem.(CodeSearchResultItem).HtmlUrl)

				}
			}
		case key.ArrowUp:
			m.selectedItem--
		case key.ArrowDown:
			m.selectedItem++
		case key.Esc:
			isQuitMenuOpen = true
		}
	}
	if m.input != m.lastInput && m.lastStrokeTime.UnixMilli()+50 < time.Now().UnixMilli() {
		requestChannel <- SearchInput{Type: m.Type, Query: m.input, CreatedAt: time.Now()}
		m.lastInput = m.input
		m.IsSearching = true
	}
	if m.IsSearching && m.lastInput == m.Result.Query {
		m.IsSearching = false
	}
	return nil
}

func (m *SearchView) SubView() *tui.View {
	if m.selectedItem >= len(m.Result.Items) {
		m.selectedItem = len(m.Result.Items) - 1
	}
	if m.selectedItem < 0 {
		m.selectedItem = 0
	}
	if m.selectedItem < len(m.Result.Items) {
		if m.Type == "repo" {
			repo := m.Result.Items[m.selectedItem].ResultItem.(Repository)
			if !m.ContentRequestMap[repo.HtmlUrl] {
				requestChannel <- m.Result.Items[m.selectedItem]
				m.ContentRequestMap[repo.HtmlUrl] = true
			}
			if repo.Description == "" && m.ContentMap[repo.HtmlUrl] == "" {
				return nil
			}
			return tui.InlineStack(
				tui.If(repo.Description != "",
					tui.InlineStack(
						tui.String("Description: \n ").FGColor(8),
						tui.String(repo.Description+"\n\n"),
					),
					nil,
				),
				tui.If(m.ContentMap[repo.HtmlUrl] != "",
					tui.InlineStack(
						tui.String("README: \n ").FGColor(8),
						tui.String(m.ContentMap[repo.HtmlUrl]+"\n"),
					),
					nil,
				),
			).Title(repo.FullName)
		} else {
			item := m.Result.Items[m.selectedItem].ResultItem.(CodeSearchResultItem)
			if !m.ContentRequestMap[item.Url] {
				requestChannel <- m.Result.Items[m.selectedItem]
				m.ContentRequestMap[item.Url] = true
			}
			if m.ContentMap[item.Url] == "" {
				return tui.String("Loading...").FGColor(8)
			}
			content := strings.ReplaceAll(m.ContentMap[item.Url], string(rune(9)), "    ")
			lines := strings.Split(content, "\n")
			col := -1
			row := -1
			for number, line := range lines {
				index := strings.Index(strings.ToUpper(line), strings.ToUpper(m.Result.Query))
				if index != -1 {
					col = index
					row = number
					break
				}
			}
			if col == -1 {
				row = 0
			}
			_, height, _ := tui.TermSize()
			beginRow := row - height/2
			endRow := row + height/2
			if beginRow < 0 {
				endRow -= beginRow
				beginRow = 0
			}
			if endRow > len(lines)-1 {
				beginRow -= endRow - (len(lines) - 1)
				if beginRow < 0 {
					beginRow = 0
				}
				endRow = len(lines) - 1
			}
			lineNumberWidth := len(strconv.Itoa(endRow + 1))
			return tui.InlineMapN(endRow-beginRow+1, func(i int) *tui.View {
				rowNumber := beginRow + i
				return tui.InlineStack(
					tui.Fmt(fmt.Sprintf("%%%dd ", lineNumberWidth), rowNumber+1).FGColor(135),
					codeLineView(lines[rowNumber], m.Result.Query),
					tui.Break(),
				)
			}).Title(item.Path)
		}
	}
	return nil
}

func codeLineView(line, pattern string) *tui.View {
	index := strings.Index(strings.ToUpper(line), strings.ToUpper(pattern))
	if index == -1 {
		return tui.String(line)
	} else {
		return tui.InlineStack(
			tui.String(line[:index]),
			tui.String(line[index:index+len(pattern)]).BGColor(135),
			tui.String(line[index+len(pattern):]),
		)
	}
}
