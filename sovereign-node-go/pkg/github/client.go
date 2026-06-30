package github

import (
	"context"
	"github.com/google/go-github/v60/github"
)

type Client struct {
	GitHub *github.Client
}

func NewClient(token string) *Client {
	return &Client{
		GitHub: github.NewClient(nil).WithAuthToken(token),
	}
}

func (c *Client) GetRepoInfo(ctx context.Context, owner, repo string) (*github.Repository, error) {
	repository, _, err := c.GitHub.Repositories.Get(ctx, owner, repo)
	return repository, err
}

func (c *Client) ListIssues(ctx context.Context, owner, repo string) ([]*github.Issue, error) {
	opts := &github.IssueListByRepoOptions{
		State: "all",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	issues, _, err := c.GitHub.Issues.ListByRepo(ctx, owner, repo, opts)
	return issues, err
}
