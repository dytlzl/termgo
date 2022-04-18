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

	"github.com/dytlzl/tervi/pkg/tui"
)

var Finalizers = make(chan func(), 16)

var clients = []GithubClient{}

var clientMap = map[string]GithubClient{}

type API struct {
	Origin  string
	Address string
}

func SetAPIs(apis []API) {
	for _, value := range apis {
		client, err := NewClient(value.Origin, value.Address)
		if err != nil {
			panic(err)
		}
		clients = append(clients, client)
		clientMap[value.Origin] = client
	}
}

func init() {
	SetAPIs([]API{
		{
			Origin:  "github.com",
			Address: "https://api.github.com",
		},
	})
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
	Items []SearchResultItem
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
	Repository
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
	items := make([]SearchResultItem, 0, 10)
	for _, client := range clients {
		result, err := client.Search(ctx, input.Query, 1, 10)
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

func FetchContent(ctx context.Context, item SearchResultItem, out chan interface{}) {
	client := clientMap[item.Origin()]
	result, err := client.FetchSearchResultContents(ctx, item)
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

	for _, client := range clients {
		res, err := client.SearchRepositories(ctx, input.Query, 1, 10)
		if err != nil {
			terminateWithError(out, err)
			return
		}
		for _, v := range res.Items {
			repositories = append(repositories, RepositoryWithOrigin{v, client.Origin})
		}
	}
	out <- RepositorySearchResult{
		SearchInput:  input,
		Repositories: repositories,
	}
}

func FetchReadMe(ctx context.Context, repo RepositoryWithOrigin, out chan interface{}) {
	client := clientMap[repo.Origin]
	result, err := client.FetchReadMe(ctx, repo.FullName)
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
	Finalizers <- func() {
		fmt.Println(err)
	}
	out <- tui.Terminate
}
