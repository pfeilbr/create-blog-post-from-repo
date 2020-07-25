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

func getTestdataDirectoryName() string {
	return os.Getenv("TEST_DATA_DIRECTORY_NAME")
}

func getTestRepoNames(username string) []string {
	names := make([]string, 0)
	files, err := ioutil.ReadDir(filepath.Join(getTestdataDirectoryName(), "repos", username))
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
	path := filepath.Join(getTestdataDirectoryName(), "repos", username, name, "expect.json")
	file, _ := ioutil.ReadFile(path)
	data := Expect{}
	_ = json.Unmarshal([]byte(file), &data)
	return data
}

func getReadmeForRepo(username string, name string) string {
	path := filepath.Join(getTestdataDirectoryName(), "repos", username, name, "repo", name, "README.md")

	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}

func mockGetReposForUser(username string) ([]*github.Repository, error) {
	path := filepath.Join(getTestdataDirectoryName(), "repos", username, "repo-list.json")

	blob, _ := ioutil.ReadFile(path)
	var respositoryList []*github.Repository
	if err := json.Unmarshal([]byte(blob), &respositoryList); err != nil {
		fmt.Printf("failed to unmarshall respository list")
		return nil, err
	}
	return respositoryList, nil
}

// func TestAdhoc(t *testing.T) {
// 	username := "pfeilbr"
// 	result, _ := getReposForUser(username)
// 	path := filepath.Join(getTestdataDirectoryName(), "repos", username, "repo-list.json")
// 	bytes, _ := json.Marshal(result)
// 	ioutil.WriteFile(path, bytes, 0644)
// }

func TestListReposForUser(t *testing.T) {
	username := "pfeilbr"
	result, _ := mockGetReposForUser(username)
	if result == nil {
		t.Errorf("no repos. got: %v", result)
	}

	if len(result) == 0 {
		t.Errorf("no repos. got: %v", result)
	}

	t.Logf("repo count: %d", len(result))
}

func TestAllRepos(t *testing.T) {
	username := "pfeilbr"
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
