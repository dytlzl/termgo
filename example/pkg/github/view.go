package github

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/dytlzl/tervi/pkg/key"
	"github.com/dytlzl/tervi/pkg/tui"
)

const channelSize = 64

type RepoSearchView struct {
	Result             RepositorySearchResult
	IsSearching        bool
	selectedRepository int
	input              string
	position           int
	lastStrokeTime     time.Time
	lastInput          string
	ReadMeMap          map[string]string
	ReadMeRequestMap   map[string]bool
}

func NewRepoSearchView() *RepoSearchView {
	return &RepoSearchView{
		Result:           RepositorySearchResult{},
		ReadMeMap:        map[string]string{},
		ReadMeRequestMap: map[string]bool{},
	}
}

func (m *RepoSearchView) Body() *tui.View {
	style := tui.Style{F256: 255, B256: 0}
	cursorStyle := tui.Style{F256: 93, B256: style.F256, HasCursor: true}

	slice := make([]*tui.View, 0, len(m.Result.Repositories)+10)
	slice = append(slice, tui.Span("Query > ").Style(tui.Style{F256: 135, B256: style.B256}))
	if m.position == len(m.input) {
		slice = append(slice,
			tui.Span(m.input[:m.position]),
			tui.Span(" ").Style(cursorStyle),
		)
	} else {
		_, size := utf8.DecodeRuneInString(m.input[m.position:])
		slice = append(slice,
			tui.Span(m.input[:m.position]),
			tui.Span(m.input[m.position:m.position+size]).Style(cursorStyle),
			tui.Span(m.input[m.position+size:]),
		)
	}
	if m.Result.Query != "" {
		if m.selectedRepository >= len(m.Result.Repositories) {
			m.selectedRepository = len(m.Result.Repositories) - 1
		}
		if m.selectedRepository < 0 {
			m.selectedRepository = 0
		}
		slice = append(slice, tui.Break(), tui.Break())
		lastOrigin := ""
		for i, repo := range m.Result.Repositories {
			if repo.Origin != lastOrigin {
				slice = append(slice, tui.Span(" "+repo.Origin+":\n").Style(tui.Style{F256: 8, B256: style.B256}))
				lastOrigin = repo.Origin
			}
			if i == m.selectedRepository {
				slice = append(slice,
					tui.Span("> ").Style(tui.Style{F256: 8, B256: style.B256}),
					tui.Fmt(" #%d", i).Style(tui.Style{F256: 8, B256: 163}),
					tui.Fmt(" %s \n", repo.FullName).Style(tui.Style{F256: 255, B256: 163}),
				)
			} else {
				slice = append(slice,
					tui.Fmt("   #%d", i).Style(tui.Style{F256: 8, B256: style.B256}),
					tui.Fmt(" %s \n", repo.FullName),
				)
			}
		}
	}
	return tui.P(slice...).Style(style)
}

func (m *RepoSearchView) HandleEvent(event any) any {
	switch typed := event.(type) {
	case RepositorySearchResult:
		if typed.CreatedAt.UnixMicro() > m.Result.CreatedAt.UnixMicro() {
			m.Result = typed
		}
	case ReadMeResult:
		m.ReadMeMap[typed.HtmlUrl] = typed.ReadMe
	case rune:
		switch typed {
		case key.CtrlS:
			return "code"
		case key.CtrlL:
			if m.selectedRepository < len(m.Result.Repositories) {
				go OpenRepository(m.Result.Repositories[m.selectedRepository].HtmlUrl)
			}
		case key.CtrlV:
			if m.selectedRepository < len(m.Result.Repositories) {
				go OpenRepository(m.Result.Repositories[m.selectedRepository].HtmlUrl)
			}
		case key.Enter:
			if m.selectedRepository < len(m.Result.Repositories) {
				go OpenUrl(m.Result.Repositories[m.selectedRepository].HtmlUrl)
			}
		case key.ArrowLeft:
			if m.position > 0 {
				_, size := utf8.DecodeLastRuneInString(m.input[:m.position])
				m.position -= size
			}
		case key.ArrowRight:
			if m.position < len(m.input) {
				_, size := utf8.DecodeRuneInString(m.input[m.position:])
				m.position += size
			}
		case key.ArrowUp:
			m.selectedRepository--
		case key.ArrowDown:
			m.selectedRepository++
		case key.Del:
			if m.input != "" {
				_, size := utf8.DecodeLastRuneInString(m.input[:m.position])
				m.input = m.input[:m.position-size] + m.input[m.position:]
				m.position -= size
			}
			m.lastStrokeTime = time.Now()
		case key.Esc:
			return tui.Terminate
		default:
			m.input = m.input[:m.position] + string(typed) + m.input[m.position:]
			m.position += utf8.RuneLen(typed)
			m.lastStrokeTime = time.Now()
		}
	}
	if m.input != m.lastInput && m.lastStrokeTime.UnixMilli()+10 < time.Now().UnixMilli() {
		requestChannel <- RepositorySearchInput{Query: m.input, CreatedAt: time.Now()}
		m.lastInput = m.input
		channel <- FooterMessage{Payload: "Searching..."}
		m.IsSearching = true
	}
	if m.IsSearching && m.lastInput == m.Result.Query {
		channel <- FooterMessage{Payload: ""}
		m.IsSearching = false
	}
	return nil
}

func (m *RepoSearchView) SubView() *RepoSubView {
	if m.selectedRepository >= len(m.Result.Repositories) {
		m.selectedRepository = len(m.Result.Repositories) - 1
	}
	if m.selectedRepository < 0 {
		m.selectedRepository = 0
	}
	if m.selectedRepository < len(m.Result.Repositories) {
		repo := m.Result.Repositories[m.selectedRepository]
		if !m.ReadMeRequestMap[repo.HtmlUrl] {
			requestChannel <- repo
			m.ReadMeRequestMap[repo.HtmlUrl] = true
		}
		if repo.Description == "" && m.ReadMeMap[repo.HtmlUrl] == "" {
			return nil
		}
		return &RepoSubView{repo: repo, readMe: m.ReadMeMap[repo.HtmlUrl]}
	}
	return nil
}

type RepoSubView struct {
	repo   RepositoryWithOrigin
	readMe string
}

func (m *RepoSubView) Body() *tui.View {
	style := tui.Style{F256: 255, B256: 0}
	keyStyle := tui.Style{F256: 135, B256: style.B256}
	return tui.P(
		tui.If(m.repo.Description != "",
			tui.P(
				tui.Span("Description: \n ").Style(keyStyle),
				tui.Span(m.repo.Description+"\n\n").Style(style),
			),
			nil,
		),
		tui.If(m.readMe != "",
			tui.P(
				tui.Span("README: \n ").Style(keyStyle),
				tui.Span(m.readMe+"\n").Style(style),
			),
			nil,
		),
	)
}

type CodeSearchView struct {
	Result            CodeSearchResult
	IsSearching       bool
	ContentMap        map[string]string
	ContentRequestMap map[string]bool
	selectedItem      int
	input             string
	position          int
	lastStrokeTime    time.Time
	lastInput         string
	runeMode          bool
}

func NewCodeSearchView() *CodeSearchView {
	return &CodeSearchView{
		ContentMap:        map[string]string{},
		ContentRequestMap: map[string]bool{},
	}
}

func (m *CodeSearchView) Body() *tui.View {
	style := tui.Style{F256: 255, B256: 0}
	cursorStyle := tui.Style{F256: 93, B256: style.F256, HasCursor: true}

	slice := make([]*tui.View, 0, len(m.Result.Items)+10)
	slice = append(slice, tui.Span("Query > ").Style(tui.Style{F256: 135, B256: style.B256}))
	if m.position == len(m.input) {
		slice = append(slice,
			tui.Span(m.input[:m.position]),
			tui.Span(" ").Style(cursorStyle),
		)
	} else {
		_, size := utf8.DecodeRuneInString(m.input[m.position:])
		slice = append(slice,
			tui.Span(m.input[:m.position]),
			tui.Span(m.input[m.position:m.position+size]).Style(cursorStyle),
			tui.Span(m.input[m.position+size:]),
		)
	}
	if m.Result.Query != "" {
		if m.selectedItem >= len(m.Result.Items) {
			m.selectedItem = len(m.Result.Items) - 1
		}
		if m.selectedItem < 0 {
			m.selectedItem = 0
		}
		slice = append(slice, tui.Break(), tui.Break())
		lastOrigin := ""
		width, _, _ := tui.TermSize()
		for i, item := range m.Result.Items {
			if item.Origin() != lastOrigin {
				slice = append(slice, tui.Span(" "+item.Origin()+":\n").Style(tui.Style{F256: 8, B256: style.B256}))
				lastOrigin = item.Origin()
			}
			path := item.Path
			if width/3-len(item.Repository.FullName)-15 < 0 {
				path = ""
			} else if len(path) > width/3-len(item.Repository.FullName)-15 {
				for len(path) > width/3-len(item.Repository.FullName)-15 {
					_, size := utf8.DecodeLastRuneInString(path)
					path = path[:len(path)-size]
				}
				path += "..."
			}
			if i == m.selectedItem {
				slice = append(slice,
					tui.Span("> ").Style(tui.Style{F256: 8, B256: style.B256}),
					tui.Fmt(" #%d ", i).Style(tui.Style{F256: 8, B256: 163}),
					tui.Span(item.Repository.FullName).Style(tui.Style{F256: 225, B256: 163}),
					tui.Fmt(" %s \n", path).Style(tui.Style{F256: 255, B256: 163}),
				)
			} else {
				slice = append(slice,
					tui.Fmt("   #%d ", i).Style(tui.Style{F256: 8, B256: style.B256}),
					tui.Span(item.Repository.FullName).Style(tui.Style{F256: 225, B256: style.B256}),
					tui.Fmt(" %s \n", path),
				)
			}
		}
	}
	return tui.P(slice...).Style(style)
}

func (m *CodeSearchView) HandleEvent(event any) any {
	switch typed := event.(type) {
	case CodeSearchResult:
		if typed.CreatedAt.UnixMicro() > m.Result.CreatedAt.UnixMicro() {
			m.Result = typed
		}
	case ContentResult:
		m.ContentMap[typed.Url] = typed.Content
	case rune:
		switch typed {
		case key.CtrlS:
			return "repo"
		case key.Enter:
			if m.selectedItem < len(m.Result.Items) {
				go OpenUrl(m.Result.Items[m.selectedItem].HtmlUrl)
			}
		case key.CtrlR:
			m.runeMode = !m.runeMode
		case key.ArrowLeft:
			if m.position > 0 {
				_, size := utf8.DecodeLastRuneInString(m.input[:m.position])
				m.position -= size
			}
		case key.ArrowRight:
			if m.position < len(m.input) {
				_, size := utf8.DecodeRuneInString(m.input[m.position:])
				m.position += size
			}
		case key.ArrowUp:
			m.selectedItem--
		case key.ArrowDown:
			m.selectedItem++
		case key.Del:
			if m.input != "" {
				_, size := utf8.DecodeLastRuneInString(m.input[:m.position])
				m.input = m.input[:m.position-size] + m.input[m.position:]
				m.position -= size
			}
			m.lastStrokeTime = time.Now()
		case key.Esc:
			return tui.Terminate
		default:
			m.input = m.input[:m.position] + string(typed) + m.input[m.position:]
			m.position += utf8.RuneLen(typed)
			m.lastStrokeTime = time.Now()
		}
	}
	if m.input != m.lastInput && m.lastStrokeTime.UnixMilli()+50 < time.Now().UnixMilli() {
		requestChannel <- CodeSearchInput{Query: m.input, CreatedAt: time.Now()}
		m.lastInput = m.input
		channel <- FooterMessage{Payload: "Searching..."}
		m.IsSearching = true
	}
	if m.IsSearching && m.lastInput == m.Result.Query {
		channel <- FooterMessage{Payload: ""}
		m.IsSearching = false
	}
	return nil
}

func (m *CodeSearchView) SubView() *CodeSubView {
	if m.selectedItem >= len(m.Result.Items) {
		m.selectedItem = len(m.Result.Items) - 1
	}
	if m.selectedItem < 0 {
		m.selectedItem = 0
	}
	if m.selectedItem < len(m.Result.Items) {
		item := m.Result.Items[m.selectedItem]
		if !m.ContentRequestMap[item.Url] {
			requestChannel <- item
			m.ContentRequestMap[item.Url] = true
		}
		return &CodeSubView{item: item, content: m.ContentMap[item.Url], query: m.Result.Query}
	}
	return nil
}

type CodeSubView struct {
	content string
	query   string
	item    SearchResultItem
}

func (m *CodeSubView) Body() *tui.View {
	if m.content == "" {
		return tui.Span("Loading...").Style(tui.Style{F256: 135, B256: 0})
	}
	content := strings.ReplaceAll(m.content, string(rune(9)), "    ")
	lines := strings.Split(content, "\n")
	col := -1
	row := -1
	for number, line := range lines {
		index := strings.Index(line, m.query)
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
	return tui.PMapN(endRow-beginRow+1, func(i int) *tui.View {
		rowNumber := beginRow + i
		index := strings.Index(lines[rowNumber], m.query)
		if index == -1 {
			return tui.P(
				tui.Fmt(fmt.Sprintf("%%%dd ", lineNumberWidth), rowNumber+1).Style(tui.Style{F256: 135, B256: 0}),
				tui.Span(lines[rowNumber]+"\n"),
			)
		}
		return tui.P(
			tui.Fmt(fmt.Sprintf("%%%dd ", lineNumberWidth), rowNumber+1).Style(tui.Style{F256: 135, B256: 0}),
			tui.P(
				tui.Span(lines[rowNumber][:index]),
				tui.Span(lines[rowNumber][index:index+len(m.query)]).BGColor(163),
				tui.Span(lines[rowNumber][index+len(m.query):]+"\n"),
			),
		)
	})
}
