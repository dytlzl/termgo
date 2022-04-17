package github

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/dytlzl/tervi/example/pkg/github"
	"github.com/dytlzl/tervi/pkg/key"
	"github.com/dytlzl/tervi/pkg/tui"
)

const channelSize = 64

type RepoSearchView struct {
	Result             RepositorySearchResult
	SearchInputCh      chan SearchInput
	ReadMeInputCh      chan RepositoryWithOrigin
	IsSearching        bool
	selectedRepository int
	input              string
	position           int
	lastStrokeTime     time.Time
	lastInput          string
	ReadMeMap          map[string]string
	ReadMeRequestMap   map[string]bool
}

func InitRepoSearchView() *RepoSearchView {
	return &RepoSearchView{
		Result:           RepositorySearchResult{},
		SearchInputCh:    make(chan SearchInput, channelSize),
		ReadMeInputCh:    make(chan RepositoryWithOrigin, channelSize),
		ReadMeMap:        map[string]string{},
		ReadMeRequestMap: map[string]bool{},
	}
}

func (m *RepoSearchView) Title() string {
	if m.Result.Query != "" {
		return fmt.Sprintf("Search Result of '%s' - Repository Search", m.Result.Query)
	}
	return "Repository Search"
}

func (m *RepoSearchView) Body(bool, tui.Size) []tui.Text {
	style := tui.CellStyle{F256: 255, B256: 0}
	cursorStyle := tui.CellStyle{F256: 93, B256: style.F256, HasCursor: true}

	slice := make([]tui.Text, 0, len(m.Result.Repositories)+10)
	slice = append(slice, tui.Text{Str: "Query > ", Style: tui.CellStyle{F256: 135, B256: style.B256}})
	if m.position == len(m.input) {
		slice = append(slice, []tui.Text{
			{Str: m.input[:m.position], Style: style},
			{Str: " ", Style: cursorStyle},
		}...)
	} else {
		_, size := utf8.DecodeRuneInString(m.input[m.position:])
		slice = append(slice, []tui.Text{
			{Str: m.input[:m.position], Style: style},
			{Str: m.input[m.position : m.position+size], Style: cursorStyle},
			{Str: m.input[m.position+size:], Style: style},
		}...)
	}
	if m.Result.Query != "" {
		if m.selectedRepository >= len(m.Result.Repositories) {
			m.selectedRepository = len(m.Result.Repositories) - 1
		}
		if m.selectedRepository < 0 {
			m.selectedRepository = 0
		}
		slice = append(slice, tui.Text{Str: "\n\n", Style: style})
		lastOrigin := ""
		for i, repo := range m.Result.Repositories {
			if repo.Origin != lastOrigin {
				slice = append(slice, tui.Text{Str: " " + repo.Origin + ":\n", Style: tui.CellStyle{F256: 8, B256: style.B256}})
				lastOrigin = repo.Origin
			}
			if i == m.selectedRepository {
				slice = append(slice, tui.Text{Str: "> ", Style: tui.CellStyle{F256: 8, B256: style.B256}})
				slice = append(slice, tui.Text{Str: fmt.Sprintf(" #%d", i), Style: tui.CellStyle{F256: 8, B256: 163}})
				slice = append(slice, tui.Text{Str: fmt.Sprintf(" %s \n", repo.FullName), Style: tui.CellStyle{F256: 255, B256: 163}})
			} else {
				slice = append(slice, tui.Text{Str: fmt.Sprintf("   #%d", i), Style: tui.CellStyle{F256: 8, B256: style.B256}})
				slice = append(slice, tui.Text{Str: fmt.Sprintf(" %s \n", repo.FullName), Style: style})
			}
		}
	}
	return slice
}

func (m *RepoSearchView) HandleEvent(event interface{}) interface{} {
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
		m.SearchInputCh <- SearchInput{Query: m.input, CreatedAt: time.Now()}
		m.lastInput = m.input
		Channel <- FooterMessage{Payload: "Searching..."}
		m.IsSearching = true
	}
	if m.IsSearching && m.lastInput == m.Result.Query {
		Channel <- FooterMessage{Payload: ""}
		m.IsSearching = false
	}
	return nil
}

func (m *RepoSearchView) SubViews() []tui.View {
	if m.selectedRepository < len(m.Result.Repositories) {
		repo := m.Result.Repositories[m.selectedRepository]
		if !m.ReadMeRequestMap[repo.HtmlUrl] {
			m.ReadMeInputCh <- repo
			m.ReadMeRequestMap[repo.HtmlUrl] = true
		}
		if repo.Description == "" && m.ReadMeMap[repo.HtmlUrl] == "" {
			return nil
		}
		return []tui.View{
			&RepoSubView{repo: repo, readMe: m.ReadMeMap[repo.HtmlUrl]},
		}
	}
	return nil
}

func (m *RepoSearchView) Options() tui.ViewOptions {
	return tui.ViewOptions{
		Title:    m.Title(),
		SubViews: m.SubViews(),
	}
}

type RepoSubView struct {
	tui.DefaultView
	repo   RepositoryWithOrigin
	readMe string
}

func (m *RepoSubView) Title() string {
	return m.repo.FullName
}

func (m *RepoSubView) Body(bool, tui.Size) []tui.Text {
	style := tui.CellStyle{F256: 255, B256: 0}
	keyStyle := tui.CellStyle{F256: 135, B256: style.B256}
	slice := make([]tui.Text, 0, 5)

	if m.repo.Description != "" {
		slice = append(slice, tui.Text{Str: "Description: \n ", Style: keyStyle})
		slice = append(slice, tui.Text{Str: m.repo.Description + "\n\n", Style: style})

	}
	if m.readMe != "" {
		slice = append(slice, tui.Text{Str: "README: \n ", Style: keyStyle})
		slice = append(slice, tui.Text{Str: m.readMe + "\n", Style: style})
	}
	return slice
}

func (m *RepoSubView) Options() tui.ViewOptions {
	return tui.ViewOptions{
		Title: m.Title(),
		Width: tui.NewFraction(2, 3),
	}
}

type CodeSearchView struct {
	Result            CodeSearchResult
	SearchInputCh     chan SearchInput
	ContentInputCh    chan github.SearchResultItem
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

func InitCodeSearchView() *CodeSearchView {
	return &CodeSearchView{
		Result:            CodeSearchResult{},
		SearchInputCh:     make(chan SearchInput, channelSize),
		ContentInputCh:    make(chan github.SearchResultItem, channelSize),
		ContentMap:        map[string]string{},
		ContentRequestMap: map[string]bool{},
	}
}

func (m *CodeSearchView) Title() string {
	if m.Result.Query != "" {
		return fmt.Sprintf("Search Result of '%s' - Code Search", m.Result.Query)
	}
	return "Code Search"
}

func (m *CodeSearchView) Body(_ bool, size tui.Size) []tui.Text {
	style := tui.CellStyle{F256: 255, B256: 0}
	cursorStyle := tui.CellStyle{F256: 93, B256: style.F256, HasCursor: true}

	slice := make([]tui.Text, 0, len(m.Result.Items)+10)
	slice = append(slice, tui.Text{Str: "Query > ", Style: tui.CellStyle{F256: 135, B256: style.B256}})
	if m.position == len(m.input) {
		slice = append(slice, []tui.Text{
			{Str: m.input[:m.position], Style: style},
			{Str: " ", Style: cursorStyle},
		}...)
	} else {
		_, size := utf8.DecodeRuneInString(m.input[m.position:])
		slice = append(slice, []tui.Text{
			{Str: m.input[:m.position], Style: style},
			{Str: m.input[m.position : m.position+size], Style: cursorStyle},
			{Str: m.input[m.position+size:], Style: style},
		}...)
	}
	if m.Result.Query != "" {
		if m.selectedItem >= len(m.Result.Items) {
			m.selectedItem = len(m.Result.Items) - 1
		}
		if m.selectedItem < 0 {
			m.selectedItem = 0
		}
		slice = append(slice, tui.Text{Str: "\n\n", Style: style})
		lastOrigin := ""
		for i, item := range m.Result.Items {
			if item.Origin() != lastOrigin {
				slice = append(slice, tui.Text{Str: " " + item.Origin() + ":\n", Style: tui.CellStyle{F256: 8, B256: style.B256}})
				lastOrigin = item.Origin()
			}
			path := item.Path
			if size.Width/3-len(item.Repository.FullName)-15 < 0 {
				path = ""
			} else if len(path) > size.Width/3-len(item.Repository.FullName)-15 {
				for len(path) > size.Width/3-len(item.Repository.FullName)-15 {
					_, size := utf8.DecodeLastRuneInString(path)
					path = path[:len(path)-size]
				}
				path += "..."
			}
			if i == m.selectedItem {
				slice = append(slice, tui.Text{Str: "> ", Style: tui.CellStyle{F256: 8, B256: style.B256}})
				slice = append(slice, tui.Text{Str: fmt.Sprintf(" #%d ", i), Style: tui.CellStyle{F256: 8, B256: 163}})
				slice = append(slice, tui.Text{Str: item.Repository.FullName, Style: tui.CellStyle{F256: 225, B256: 163}})
				slice = append(slice, tui.Text{Str: fmt.Sprintf(" %s \n", path), Style: tui.CellStyle{F256: 255, B256: 163}})
			} else {
				slice = append(slice, tui.Text{Str: fmt.Sprintf("   #%d ", i), Style: tui.CellStyle{F256: 8, B256: style.B256}})
				slice = append(slice, tui.Text{Str: item.Repository.FullName, Style: tui.CellStyle{F256: 225, B256: style.B256}})
				slice = append(slice, tui.Text{Str: fmt.Sprintf(" %s \n", path), Style: style})
			}
		}
	}
	return slice
}

func (m *CodeSearchView) HandleEvent(event interface{}) interface{} {
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
		m.SearchInputCh <- SearchInput{Query: m.input, CreatedAt: time.Now()}
		m.lastInput = m.input
		Channel <- FooterMessage{Payload: "Searching..."}
		m.IsSearching = true
	}
	if m.IsSearching && m.lastInput == m.Result.Query {
		Channel <- FooterMessage{Payload: ""}
		m.IsSearching = false
	}
	return nil
}

func (m *CodeSearchView) SubViews() []tui.View {
	if m.selectedItem < len(m.Result.Items) {
		item := m.Result.Items[m.selectedItem]
		if !m.ContentRequestMap[item.Url] {
			m.ContentInputCh <- item
			m.ContentRequestMap[item.Url] = true
		}
		return []tui.View{
			&CodeSubView{item: item, content: m.ContentMap[item.Url], query: m.Result.Query, runeMode: m.runeMode},
		}
	}
	return nil
}

func (m *CodeSearchView) Options() tui.ViewOptions {
	return tui.ViewOptions{
		Title:    m.Title(),
		SubViews: m.SubViews(),
	}
}

type CodeSubView struct {
	tui.DefaultView
	content  string
	query    string
	item     github.SearchResultItem
	runeMode bool
}

func (m *CodeSubView) Title() string {
	return m.item.Path
}

func (m *CodeSubView) Body(_ bool, size tui.Size) []tui.Text {
	if m.runeMode {
		return m.Body_RuneMode(false)
	}
	if m.content == "" {
		return []tui.Text{{Str: "Loading...", Style: tui.CellStyle{F256: 135, B256: 0}}}
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
	beginRow := row - size.Height/2
	endRow := row + size.Height/2
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
	slice := make([]tui.Text, 0, 64)
	lineNumberWidth := len(strconv.Itoa(endRow + 1))
	for i := beginRow; i <= endRow; i++ {
		slice = append(slice, tui.Text{Str: fmt.Sprintf(fmt.Sprintf("%%%dd ", lineNumberWidth), i+1), Style: tui.CellStyle{F256: 135, B256: 0}})
		index := strings.Index(lines[i], m.query)
		if index == -1 {
			slice = append(slice, tui.Text{Str: lines[i] + "\n", Style: tui.DefaultStyle})
		} else {
			slice = append(slice, tui.Text{Str: lines[i][:index], Style: tui.DefaultStyle})
			slice = append(slice, tui.Text{Str: lines[i][index : index+len(m.query)], Style: tui.CellStyle{B256: 163}})
			slice = append(slice, tui.Text{Str: lines[i][index+len(m.query):] + "\n", Style: tui.DefaultStyle})
		}
	}
	return slice
}

func (m *CodeSubView) Body_RuneMode(bool) []tui.Text {
	if m.content == "" {
		return []tui.Text{{Str: "Loading...", Style: tui.CellStyle{F256: 135, B256: 0}}}
	}
	slice := make([]tui.Text, 0, 2048)
	lineNumber := 0
	slice = append(slice, tui.Text{Str: fmt.Sprintf("%3d ", lineNumber), Style: tui.CellStyle{F256: 135, B256: 0}})
	for _, r := range m.content {
		slice = append(slice, tui.Text{Str: " " + strconv.Itoa(int(r)), Style: tui.DefaultStyle})
		if r == '\n' {
			lineNumber++
			slice = append(slice, tui.Text{Str: "\n", Style: tui.DefaultStyle})
			slice = append(slice, tui.Text{Str: fmt.Sprintf("%3d ", lineNumber), Style: tui.CellStyle{F256: 135, B256: 0}})
		}
	}
	return slice
}

func (m *CodeSubView) Options() tui.ViewOptions {
	return tui.ViewOptions{
		Title: m.Title(),
		Width: tui.NewFraction(2, 3),
	}
}

type Footer struct {
	message        string
	MessageChannel chan string
}

func (*Footer) Style() tui.CellStyle {
	return tui.CellStyle{B256: 135}
}

func (f *Footer) Text() []tui.Text {
	return []tui.Text{{Str: f.message, Style: f.Style()}}
}

func (f *Footer) HandleEvent(event interface{}) interface{} {
	switch typed := event.(type) {
	case FooterMessage:
		f.message = typed.Payload
	}
	return nil
}
