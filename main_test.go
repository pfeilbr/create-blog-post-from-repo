package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-github/github"
	_ "github.com/joho/godotenv/autoload"
)

type Expect struct {
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

func getTestRepoNames(username string) []string {
	names := make([]string, 0)
	files, err := ioutil.ReadDir(filepath.Join(testdataDirectoryName, "repos", username))
	if err != nil {
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

func getExpectedResultsForRepo(username string, name string) Expect {
	path := filepath.Join(testdataDirectoryName, "repos", username, name, "expect.json")
	file, _ := ioutil.ReadFile(path)
	data := Expect{}
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

func mockGetReposForUser(username string) ([]*github.Repository, error) {
	path := filepath.Join(testdataDirectoryName, "repos", username, "repo-list.json")

	blob, _ := ioutil.ReadFile(path)
	var respositoryList []*github.Repository
	if err := json.Unmarshal([]byte(blob), &respositoryList); err != nil {
		fmt.Printf("failed to unmarshall respository list")
		return nil, err
	}
	return respositoryList, nil
}

func TestListReposForUser(t *testing.T) {
	username := githubUsername
	result, _ := mockGetReposForUser(username)
	if result == nil {
		t.Errorf("no repos. got: %v", result)
	}

	if len(result) == 0 {
		t.Errorf("no repos. got: %v", result)
	}
}

func TestApplyIncludeFilterForRepos(t *testing.T) {
	username := githubUsername
	repos, _ := mockGetReposForUser(username)
	filteredRepos, err := applyIncludeFilterForRepos(repos)
	if err != nil {
		t.Error(err)
	}

	filteredReposCount := len(filteredRepos)
	if filteredReposCount == 0 {
		t.Errorf("expected >0 filtered repos.  got %d", filteredReposCount)
	}

	t.Logf("filteredReposCount: %d", filteredReposCount)
	// t.Logf("filteredRepos:\n%v", Map(filteredRepos, func(repo *github.Repository) string {
	// 	return *repo.Name
	// }))
}

func TestAllRepos(t *testing.T) {
	username := githubUsername
	names := getTestRepoNames(username)

	for _, name := range names {
		expect := getExpectedResultsForRepo(username, name)

		t.Run("titleForRepo/"+name, func(t *testing.T) {
			want := expect.Title
			result := titleForRepo(name)

			if result != want {
				t.Errorf("got %s, want %s", result, want)
			}
		})

		t.Run("getReadmeForRepo/"+name, func(t *testing.T) {
			result := getReadmeForRepo(username, name)
			if len(result) == 0 {
				t.Errorf("empty README contents. got %s", result)
			}
		})

	}

}

func TestCreatePostString(t *testing.T) {
	username := githubUsername
	repos, _ := mockGetReposForUser(username)
	filteredRepos, err := applyIncludeFilterForRepos(repos)
	if err != nil {
		t.Error(err)
	}

	repo := filteredRepos[0]
	postString, err := createPostString(repo)
	if err != nil {
		t.Error(err)
	}
	t.Log(postString)
}
