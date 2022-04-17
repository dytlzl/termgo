package github

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/dytlzl/tervi/example/pkg/github"
	"github.com/dytlzl/tervi/pkg/tui"
)

var Finalize = func() {}

var Servers = []github.GithubClient{}

var ServerMap = map[string]github.GithubClient{}

func init() {
	for _, value := range []struct {
		origin     string
		apiAddress string
	}{
		{
			origin:     "github.com",
			apiAddress: "https://api.github.com",
		},
	} {
		server, err := github.NewClient(value.origin, value.apiAddress)
		if err != nil {
			panic(err)
		}
		Servers = append(Servers, server)
		ServerMap[value.origin] = server
	}
}

var Channel = make(chan interface{}, 100)

type SearchInput struct {
	Query     string
	CreatedAt time.Time
}

type RepositorySearchResult struct {
	SearchInput
	Repositories []RepositoryWithOrigin
}

type CodeSearchResult struct {
	SearchInput
	Items []github.SearchResultItem
}

type ContentResult struct {
	Url     string
	Content string
}

type ReadMeResult struct {
	HtmlUrl string
	ReadMe  string
}

type RepositoryWithOrigin struct {
	github.Repository
	Origin string
}

type FooterMessage struct {
	Payload string
}

func RepositoryPath(url string) string {
	segments := strings.Split(url, "/")
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	rootPath := u.HomeDir + "/" + "ghq"
	return rootPath + strings.Join(segments[1:], "/")
}

func CloneRepository(dirPath, url string) error {
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory where repository is cloned: %w", err)
	}
	prevDir, _ := os.Getwd()
	os.Chdir(dirPath)
	defer os.Chdir(prevDir)
	_, err = exec.Command("git", "clone", url).Output()
	if err != nil {
		return fmt.Errorf("failed to execute git clone %s: %w", url, err)
	}
	return nil
}

func OpenRepository(url string) error {
	repoPath := RepositoryPath(url)
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		Channel <- FooterMessage{"Cloning " + url + "..."}
		dirPath := filepath.Dir(repoPath)
		err = CloneRepository(dirPath, url)
		if err != nil {
			return fmt.Errorf("failed to clone %s: %w", url, err)
		}
		Channel <- FooterMessage{"Cloning " + url + "..." + " Done."}
	} else {
		Channel <- FooterMessage{url + " already exists locally."}
	}
	err := exec.Command("open", "-a", "Visual Studio Code", repoPath).Start()
	if err != nil {
		return fmt.Errorf("failed to open Visual Studio Code: %w", err)
	}
	return nil
}

func OpenUrl(url string) error {
	return exec.Command("open", url).Start()
}

func SearchCode(ctx context.Context, input SearchInput, out chan interface{}) {
	items := make([]github.SearchResultItem, 0, 10)
	for _, server := range Servers {
		result, err := server.Search(ctx, input.Query, 1, 10)
		if err != nil {
			terminateWithError(out, err)
			return
		}
		items = append(items, result.Items...)
	}
	out <- CodeSearchResult{
		SearchInput: input,
		Items:       items,
	}
}

func FetchContent(ctx context.Context, item github.SearchResultItem, out chan interface{}) {
	server := ServerMap[item.Origin()]
	result, err := server.FetchSearchResultContents(ctx, item)
	if err != nil {
		terminateWithError(out, err)
		return
	}
	out <- ContentResult{
		Url:     item.Url,
		Content: result,
	}
}

func SearchRepositories(ctx context.Context, input SearchInput, out chan interface{}) {
	repositories := make([]RepositoryWithOrigin, 0, 20)

	for _, server := range Servers {
		res, err := server.SearchRepositories(ctx, input.Query, 1, 10)
		if err != nil {
			terminateWithError(out, err)
			return
		}
		for _, v := range res.Items {
			repositories = append(repositories, RepositoryWithOrigin{v, server.Origin})
		}
	}
	out <- RepositorySearchResult{
		SearchInput:  input,
		Repositories: repositories,
	}
}

func FetchReadMe(ctx context.Context, repo RepositoryWithOrigin, out chan interface{}) {
	server := ServerMap[repo.Origin]
	result, err := server.FetchReadMe(ctx, repo.FullName)
	if err != nil {
		terminateWithError(out, err)
		return
	}
	out <- ReadMeResult{
		HtmlUrl: repo.HtmlUrl,
		ReadMe:  result,
	}
}

func terminateWithError(out chan interface{}, err error) {
	Finalize = func() {
		fmt.Println(err)
	}
	out <- tui.Terminate{}
}
