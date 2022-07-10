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

var Clients = []GithubClient{}

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
		Clients = append(Clients, client)
		clientMap[value.Origin] = client
	}
}

type SearchInput struct {
	Type      string
	Query     string
	CreatedAt time.Time
}

type searchResult struct {
	SearchInput
	Items []ResultItemWithOrigin
}

type ContentResult struct {
	Url     string
	Content string
}

type ReadMeResult struct {
	HtmlUrl string
	ReadMe  string
}

type ResultItemWithOrigin struct {
	ResultItem any
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

func CloneRepository(repoPath, url string) error {
	dirPath := filepath.Dir(repoPath)
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory where repository is cloned: %w", err)
	}
	_, err = exec.Command("git", "clone", url, repoPath).Output()
	if err != nil {
		return fmt.Errorf("failed to clone %s: %w", url, err)
	}
	return nil
}

func OpenRepository(url string) error {
	repoPath := RepositoryPath(url)
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		channel <- FooterMessage{"Cloning " + url + "..."}
		err = CloneRepository(repoPath, url)
		if err != nil {
			return fmt.Errorf("failed to clone %s: %w", url, err)
		}
		channel <- FooterMessage{"Cloning " + url + "..." + " Done."}
	} else {
		channel <- FooterMessage{url + " already exists locally."}
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

func SendToChan(data any, err error) {
	if err != nil {
		Finalizers <- func() {
			fmt.Fprintln(os.Stderr, err)
		}
		channel <- tui.Terminate
		return
	}
	channel <- data
}

func SearchCode(ctx context.Context, input SearchInput) (searchResult, error) {
	items := make([]ResultItemWithOrigin, 0, 10)
	for _, client := range Clients {
		result, err := client.Search(ctx, input.Query, 1, 10)
		if err != nil {
			return searchResult{}, err
		}
		for _, v := range result.Items {
			items = append(items, ResultItemWithOrigin{v, client.Origin})
		}
	}
	return searchResult{
		SearchInput: input,
		Items:       items,
	}, nil
}

func FetchContent(ctx context.Context, origin string, item CodeSearchResultItem) (ContentResult, error) {
	client := clientMap[origin]
	result, err := client.FetchSearchResultContents(ctx, item)
	if err != nil {
		return ContentResult{}, err
	}
	return ContentResult{
		Url:     item.Url,
		Content: result,
	}, nil
}

func SearchRepositories(ctx context.Context, input SearchInput) (searchResult, error) {
	repositories := make([]ResultItemWithOrigin, 0, 20)
	for _, client := range Clients {
		res, err := client.SearchRepositories(ctx, input.Query, 1, 10)
		if err != nil {
			return searchResult{}, nil
		}
		for _, v := range res.Items {
			repositories = append(repositories, ResultItemWithOrigin{v, client.Origin})
		}
	}
	return searchResult{
		SearchInput: input,
		Items:       repositories,
	}, nil
}

func FetchReadMe(ctx context.Context, origin string, repo Repository) (ReadMeResult, error) {
	client := clientMap[origin]
	result, err := client.FetchReadMe(ctx, repo.FullName)
	if err != nil {
		return ReadMeResult{}, nil
	}
	return ReadMeResult{
		HtmlUrl: repo.HtmlUrl,
		ReadMe:  result,
	}, nil
}
