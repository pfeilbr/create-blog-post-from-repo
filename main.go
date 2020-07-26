package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/google/go-github/github"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
)

var command string
var user string
var path string
var destinationDirectory string

const tempDirectoryName = "tmp"

func init() {
	flag.StringVar(&command, "command", "fetch-and-save-repo-metadata-list-for-user", "command to run")
	flag.StringVar(&user, "user", "pfeilbr", "github username")
	flag.StringVar(&path, "output", "repo-list.json", "file output path")
	flag.StringVar(&destinationDirectory, "destination-directory", "", "directory to save geneated markdown post file(s) to")
}

type RepoMetadata struct {
	Repo             *github.Repository
	Title            string
	MarkdownBody     string
	PostFileName     string
	PostFileContents string
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

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func getReposForUser(user string, useCache bool) ([]*github.Repository, error) {

	cacheRepoMetadataListPath := cacheRepoMetadataListPathForUser(user)
	if useCache {
		if fileExists(cacheRepoMetadataListPath) {
			blob, _ := ioutil.ReadFile(cacheRepoMetadataListPath)
			var respositoryList []*github.Repository
			if err := json.Unmarshal([]byte(blob), &respositoryList); err != nil {
				fmt.Printf("failed to unmarshall respository list")
				return nil, err
			}
			return respositoryList, nil
		}
	}

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

	if useCache {
		os.MkdirAll(tempDirectoryName, os.ModePerm)
		bytes, _ := json.Marshal(allRepos)
		ioutil.WriteFile(cacheRepoMetadataListPath, bytes, 0644)
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

func getURLContents(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Read body: %v", err)
	}

	return string(data), nil
}

func getMarkdownBodyForRepo(repo *github.Repository) (string, error) {
	url := "https://raw.githubusercontent.com/" + *repo.FullName + "/master/README.md"
	contents, err := getURLContents(url)
	if err != nil {
		fmt.Printf("getURLContents(%s) failed", url)
		return "", err
	}
	return contents, nil
}

func getPostFileNameForRepo(repo *github.Repository) string {
	return "generated-" + *repo.Name + ".md"
}

func newRepoMetadata(repo *github.Repository) (*RepoMetadata, error) {
	markdownBody, err := getMarkdownBodyForRepo(repo)

	if err != nil {
		fmt.Printf("failed to getMarkdownBodyForRepo(%s)", *repo.Name)
		return nil, err
	}
	repoMetadata := &RepoMetadata{
		Repo:         repo,
		Title:        titleForRepo(*repo.Name),
		MarkdownBody: markdownBody,
		PostFileName: getPostFileNameForRepo(repo),
	}
	postFileContents, err := getPostFileContents(repoMetadata)
	if err != nil {
		fmt.Printf("getPostFileContents(%s) failed\n", *repo.Name)
		return nil, err
	}

	repoMetadata.PostFileContents = postFileContents
	return repoMetadata, nil
}

func getPostFileContents(repoMetadata *RepoMetadata) (string, error) {
	postTemplatePath := filepath.Join("templates", "post.md")

	b, err := ioutil.ReadFile(postTemplatePath)
	contentsAsString := string(b)

	t := template.Must(template.New("hugo-markdown-post-tmpl").Parse(contentsAsString))

	var buf bytes.Buffer

	err = t.Execute(&buf, repoMetadata)
	if err != nil {
		panic(err)
	}

	result := buf.String()
	return result, nil
}

func cacheRepoMetadataListPathForUser(username string) string {
	return filepath.Join(tempDirectoryName, "repo-metadata-list-"+username+".json")
}

func fetchAndSaveRepoMetadataListForUser(username string, path string) error {
	result, err := getReposForUser(username, false)
	if err != nil {
		fmt.Printf("getReposForUser(%s) failed\n", username)
		return err
	}
	bytes, _ := json.Marshal(result)
	ioutil.WriteFile(path, bytes, 0644)
	return nil
}

func generateMarkdownPostFile(repo *github.Repository, destinationDirectory string) error {
	repoMetadata, err := newRepoMetadata(repo)

	if err != nil {
		fmt.Printf("newRepoMetadata(%s) failed\n", *repo.Name)
		return err
	}

	os.MkdirAll(destinationDirectory, os.ModePerm)
	path := filepath.Join(destinationDirectory, repoMetadata.PostFileName)
	if err := ioutil.WriteFile(path, []byte(repoMetadata.PostFileContents), 0644); err != nil {
		fmt.Printf("ioutil.WriteFile(%s) failed\n", path)
		return err
	}

	return nil
}

func generateMarkdownPostFiles(user string, destinationDirectory string) error {
	repos, err := getReposForUser(user, true)
	if err != nil {
		fmt.Printf("getReposForUser(%s) failed\n", user)
		return err
	}

	filteredRepos, err := applyIncludeFilterForRepos(repos)
	if err != nil {
		fmt.Printf("applyIncludeFilterForRepos(repos) failed\n")
		return err
	}

	for _, repo := range filteredRepos {
		err := generateMarkdownPostFile(repo, destinationDirectory)
		if err != nil {
			fmt.Printf("generateMarkdownPostFile(%s) failed\n", *repo.Name)
			return err
		}

	}

	return nil
}

func main() {
	flag.Parse()

	if command == "fetch-and-save-repo-metadata-list-for-user" {
		fmt.Printf("command: %s, user: %s, path: %s\n", command, user, path)
		if err := fetchAndSaveRepoMetadataListForUser(user, path); err != nil {
			log.Fatal(err)
		}
	}

	if command == "generate-markdown-post-files" {
		fmt.Printf("command: %s, user: %s, destinationDirectory: %s\n", command, user, destinationDirectory)
		if err := generateMarkdownPostFiles(user, destinationDirectory); err != nil {
			log.Fatal(err)
		}
	}

}
