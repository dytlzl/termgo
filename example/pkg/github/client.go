package github

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
)

type GithubClient struct {
	Origin     string
	apiAddress string
	token      string
	cli        http.Client
}

func NewClient(origin, apiAddress string) (GithubClient, error) {
	token, err := exec.Command("security", "find-internet-password", "-w", "-s", origin).Output()
	if err != nil {
		return GithubClient{}, fmt.Errorf("could not obtain password of %s from key chain", origin)
	}
	tokenString := strings.TrimRight(string(token), "\n")
	return GithubClient{
		Origin:     origin,
		apiAddress: apiAddress,
		token:      tokenString,
		cli:        http.Client{},
	}, nil
}

func (g GithubClient) Request(ctx context.Context, method string, url string, paramMap map[string]any) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}
	params := req.URL.Query()
	for k, v := range paramMap {
		switch typed := v.(type) {
		case string:
			params.Add(k, typed)
		case int:
			params.Add(k, strconv.Itoa(typed))
		}
	}
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	u, err := user.Current()
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(u.Username, g.token)
	resp, err := g.cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("status code is 403 forbidden, response body: %s", body)
	}
	return body, nil
}

func (g GithubClient) RequestWithEndpoint(ctx context.Context, method string, endpoint string, paramMap map[string]any) ([]byte, error) {
	return g.Request(ctx, method, g.apiAddress+endpoint, paramMap)
}

type CodeSearchResult struct {
	TotalCount int                    `json:"total_count"`
	Items      []CodeSearchResultItem `json:"items"`
}

type RepositoryItemsSearchResult struct {
	TotalCount int          `json:"total_count"`
	Items      []Repository `json:"items"`
}

type IssueSearchResult struct {
	TotalCount int     `json:"total_count"`
	Items      []Issue `json:"items"`
}

type User struct {
	Login string `json:"login"`
}

type Issue struct {
	Title     string `json:"title"`
	HtmlUrl   string `json:"html_url"`
	EventsUrl string `json:"events_url"`
	User      User   `json:"user"`
}

type Event struct {
	Event             string `json:"event"`
	RequestedReviewer User   `json:"requested_reviewer"`
}

type CodeSearchResultItem struct {
	Url        string     `json:"url"`
	Path       string     `json:"path"`
	HtmlUrl    string     `json:"html_url"`
	Repository Repository `json:"repository"`
}

func (i CodeSearchResultItem) Origin() string {
	return strings.Split(i.Url, "/")[2]
}

func (g *GithubClient) FetchSearchResultContents(ctx context.Context, item CodeSearchResultItem) (string, error) {
	res, err := g.Request(ctx, "GET", item.Url, nil)
	if err != nil {
		return "", err
	}
	var jsonMap struct {
		Content string `json:"content"`
	}
	if err := json.Unmarshal(res, &jsonMap); err != nil {
		return "", err
	}
	bytes, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(jsonMap.Content, "\n", ""))
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (g *GithubClient) FetchReadMe(ctx context.Context, fullName string) (string, error) {
	res, err := g.RequestWithEndpoint(ctx, "GET", "/repos/"+fullName+"/readme", nil)
	if err != nil {
		return "", err
	}
	var jsonMap struct {
		Content string `json:"content"`
	}
	if err := json.Unmarshal(res, &jsonMap); err != nil {
		return "", err
	}
	bytes, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(jsonMap.Content, "\n", ""))
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (g *GithubClient) Search(ctx context.Context, query string, page int, per_page int) (CodeSearchResult, error) {
	res, err := g.RequestWithEndpoint(ctx, "GET", "/search/code", map[string]any{
		"q":        query,
		"per_page": per_page,
		"page":     page,
	})
	if err != nil {
		return CodeSearchResult{}, err
	}
	var jsonMap CodeSearchResult
	if err := json.Unmarshal(res, &jsonMap); err != nil {
		return CodeSearchResult{}, fmt.Errorf("failed to unmarshal json: %w: %s", err, string(res))
	}
	return jsonMap, nil
}

func (g *GithubClient) SearchRepositories(ctx context.Context, query string, page int, per_page int) (RepositoryItemsSearchResult, error) {
	res, err := g.RequestWithEndpoint(ctx, "GET", "/search/repositories", map[string]any{
		"q":        query,
		"per_page": per_page,
		"page":     page,
	})
	if err != nil {
		return RepositoryItemsSearchResult{}, err
	}
	var jsonMap RepositoryItemsSearchResult
	if err := json.Unmarshal(res, &jsonMap); err != nil {
		return RepositoryItemsSearchResult{}, fmt.Errorf("failed to unmarshal json: %w: %s", err, string(res))
	}
	return jsonMap, nil
}

func (g *GithubClient) SearchIssues(ctx context.Context, query string, page int, per_page int) (IssueSearchResult, error) {
	res, err := g.RequestWithEndpoint(ctx, "GET", "/search/issues", map[string]any{
		"q":        query,
		"per_page": per_page,
		"page":     page,
	})
	if err != nil {
		return IssueSearchResult{}, err
	}
	var jsonMap IssueSearchResult
	if err := json.Unmarshal(res, &jsonMap); err != nil {
		return IssueSearchResult{}, fmt.Errorf("failed to unmarshal json: %w: %s", err, string(res))
	}
	return jsonMap, nil
}

func (g *GithubClient) FetchEventsFromIssue(ctx context.Context, issue Issue) ([]Event, error) {
	res, err := g.Request(ctx, "GET", issue.EventsUrl, nil)
	if err != nil {
		return nil, err
	}
	var jsonMap []Event
	if err := json.Unmarshal(res, &jsonMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w: %s", err, string(res))
	}
	return jsonMap, nil
}

type Repository struct {
	FullName    string `json:"full_name"`
	HtmlUrl     string `json:"html_url"`
	Description string `json:"description"`
}

const REPOSITORY_NUMBER_PER_PAGE = 100

func (g *GithubClient) fetchRepostories(ctx context.Context, org string, page int) ([]Repository, error) {
	res, err := g.RequestWithEndpoint(ctx, "GET", "/orgs/"+org+"/repos", map[string]any{
		"per_page": REPOSITORY_NUMBER_PER_PAGE,
		"page":     page,
	})
	if err != nil {
		return nil, err
	}
	var jsonMap []Repository
	if err := json.Unmarshal(res, &jsonMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w: %s", err, string(res))
	}
	return jsonMap, nil
}

func (g *GithubClient) FetchAllRepostories(ctx context.Context, org string) ([]Repository, error) {
	repositories := make([]Repository, 0, 512)
	for i := 1; ; i++ {
		pageRepositories, err := g.fetchRepostories(ctx, org, i)
		if err != nil {
			return nil, err
		}
		repositories = append(repositories, pageRepositories...)
		if len(pageRepositories) != REPOSITORY_NUMBER_PER_PAGE {
			return repositories, err
		}
	}
}
