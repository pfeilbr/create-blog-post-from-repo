package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/google/go-github/github"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
)

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

func applyIncludeFilterForRepos(repos []*github.Repository) ([]*github.Repository, error) {
	var result []*github.Repository
	for _, repo := range repos {
		var match bool = false
		repoNameIncludeFiltersString := os.Getenv("REPO_NAME_INCLUDE_FILTERS")
		repoNameIncludeFilters := strings.Split(repoNameIncludeFiltersString, ",")

		for _, includeFilter := range repoNameIncludeFilters {
			re := regexp.MustCompile(includeFilter)
			if re.Match([]byte(*repo.Name)) == true {
				match = true
			}
		}

		if match == true {
			result = append(result, repo)
		}
	}
	return result, nil
}

func titleForRepo(repoName string) string {
	title := strings.ReplaceAll(repoName, "playground", "")
	title = strings.ReplaceAll(title, "-", " ")
	title = strings.TrimSpace(title)
	title = strings.ToLower(title)
	title = strings.Title(title)
	return title
}

func fetchRepoMetadataListForUser(username string, path string) error {
	result, err := getReposForUser(username)
	if err != nil {
		fmt.Printf("getReposForUser(%s) failed\n", username)
		return err
	}
	bytes, _ := json.Marshal(result)
	ioutil.WriteFile(path, bytes, 0644)
	return nil
}

var command string
var user string
var path string

func init() {
	flag.StringVar(&command, "command", "fetch-repo-metadata-list-for-user", "command to run")
	flag.StringVar(&user, "user", "pfeilbr", "github username")
	flag.StringVar(&path, "output", "repo-list.json", "file output path")
}

func main() {
	flag.Parse()

	if command == "fetch-repo-metadata-list-for-user" {
		fmt.Printf("command: %s, user: %s, path: %s\n", command, user, path)
		if err := fetchRepoMetadataListForUser(user, path); err != nil {
			log.Fatal(err)
		}
	}
}
