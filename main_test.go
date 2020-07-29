package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-github/github"
	_ "github.com/joho/godotenv/autoload"
)

type ExpectedTestResults struct {
	Title string `json: "title"`
}

var testdataDirectoryName string
var githubUsername string

func init() {
	testdataDirectoryName = os.Getenv("TEST_DATA_DIRECTORY_NAME")
	githubUsername = os.Getenv("GITHUB_USERNAME")
}

func Map(vs []*github.Repository, f func(*github.Repository) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func getTestRepoNames(user string) []string {
	names := make([]string, 0)
	path := filepath.Join(testdataDirectoryName, "repos", user)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Printf("ioutil.ReadDir(\"%s\") failed\n", path)
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		if !file.IsDir() {
			continue
		}
		names = append(names, file.Name())
	}
	return names
}

func getExpectedResultsForRepo(username string, name string) ExpectedTestResults {
	path := filepath.Join(testdataDirectoryName, "repos", username, name, "expect.json")
	file, _ := ioutil.ReadFile(path)
	data := ExpectedTestResults{}
	_ = json.Unmarshal([]byte(file), &data)
	return data
}

func getReadmeForRepo(username string, name string) string {
	path := filepath.Join(testdataDirectoryName, "repos", username, name, "repo", name, "README.md")

	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}

// func mockGetReposForUser(username string) ([]*github.Repository, error) {
// 	path := filepath.Join(testdataDirectoryName, "repos", username, "repo-list.json")

// 	blob, _ := ioutil.ReadFile(path)
// 	var respositoryList []*github.Repository
// 	if err := json.Unmarshal([]byte(blob), &respositoryList); err != nil {
// 		log.Printf("failed to unmarshall respository list")
// 		return nil, err
// 	}
// 	return respositoryList, nil
// }

func TestGetReposForUser(t *testing.T) {
	user := githubUsername
	result, _ := getReposForUser(user, true)
	if result == nil {
		t.Errorf("no repos. got: %v", result)
	}

	if len(result) == 0 {
		t.Errorf("no repos. got: %v", result)
	}
}

func TestGetFilteredRepos(t *testing.T) {
	user := githubUsername
	repos, _ := getReposForUser(user, true)
	filteredRepos, err := getFilteredRepos(repos)
	if err != nil {
		t.Error(err)
	}

	filteredReposCount := len(filteredRepos)
	if filteredReposCount == 0 {
		t.Errorf("expected >0 filtered repos.  got %d", filteredReposCount)
	}

	//t.Logf("filteredReposCount: %d", filteredReposCount)
	// t.Logf("filteredRepos:\n%v", Map(filteredRepos, func(repo *github.Repository) string {
	// 	return *repo.Name
	// }))
}

func TestAllRepos(t *testing.T) {
	user := githubUsername
	names := getTestRepoNames(user)

	for _, name := range names {
		expect := getExpectedResultsForRepo(user, name)

		t.Run("titleForRepo/"+name, func(t *testing.T) {
			want := expect.Title
			result := getPostTitle(name)

			if result != want {
				t.Errorf("got %s, want %s", result, want)
			}
		})

		t.Run("getReadmeForRepo/"+name, func(t *testing.T) {
			result := getReadmeForRepo(user, name)
			if len(result) == 0 {
				t.Errorf("empty README contents. got %s", result)
			}
		})

	}

}

func TestCreateMarkdownPostFiles(t *testing.T) {
	useCache = true
	destinationDirectory := "tmp/posts"
	os.RemoveAll(destinationDirectory)
	os.MkdirAll(destinationDirectory, 0777)
	user := githubUsername
	if err := createMarkdownPostFiles(user, destinationDirectory); err != nil {
		log.Printf("createMarkdownPostFiles(%s, %s). failed\n", user, destinationDirectory)
		t.Error(err)
	}

	destinationCopyToDirectoryPath := "/Users/pfeilbr/Dropbox/mac01/Users/brianpfeil/projects/personal-website/content/post"
	copyDirectoryRecursively(destinationDirectory, destinationCopyToDirectoryPath)
}

func TestGetFilteredReposForUser(t *testing.T) {
	t.SkipNow()
	username := githubUsername
	repoNames := make([]string, 0)
	t.Logf("starting ...")
	filteredRepos, err := getFilteredReposForUser(username)
	if err != nil {
		log.Printf("getFilteredReposForUser(%s) failed\n", username)
		return
	}

	for _, repo := range filteredRepos {
		repoNames = append(repoNames, *repo.Name)
	}
	repoNamesString := strings.Join(repoNames, "\n")
	t.Logf(repoNamesString)
}

func TestGetRepoPostsWithNoTags(t *testing.T) {
	//t.SkipNow()
	username := githubUsername
	repoPosts, err := getRepoPostsWithNoTags(username)
	if err != nil {
		t.Error(err)
	}

	repoNames := make([]string, 0)
	for _, repoPost := range repoPosts {
		repoNames = append(repoNames, *repoPost.Repo.Name)
	}
	repoNamesString := strings.Join(repoNames, "=|")
	t.Log(repoNamesString)
}
