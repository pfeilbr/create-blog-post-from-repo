package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/github"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
)

func Hello() string {
	return "Hello, world."
}

func getGithubClient() *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return client
}

func getReposForUser(user string) ([]*github.Repository, error) {
	ctx := context.Background()
	client := getGithubClient()

	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 30},
	}
	// get all pages of results
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.List(ctx, user, opt)
		if err != nil {
			fmt.Printf("failed to list repositories for user %s", user)
			return nil, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allRepos, nil
}

func titleForRepo(repoName string) string {
	title := strings.ReplaceAll(repoName, "playground", "")
	title = strings.ReplaceAll(title, "-", " ")
	title = strings.TrimSpace(title)
	title = strings.ToLower(title)
	title = strings.Title(title)
	return title
}

func main() {
	fmt.Printf("%s", Hello())
}
