package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/google/go-github/github"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
)

var command string
var user string
var path string
var destinationDirectory string
var useCache bool
var debug bool

const tempDirectoryName = "tmp"

func init() {
	flag.StringVar(&command, "command", "", "command to run")
	flag.StringVar(&user, "user", "", "github username")
	flag.StringVar(&path, "output", "", "file output path")
	flag.StringVar(&destinationDirectory, "destination-directory", "", "directory to save geneated markdown post file(s) to")
	flag.BoolVar(&useCache, "cache", true, "cache requests to repo")
	flag.BoolVar(&debug, "debug", false, "print debug information")
}

// RepoPost contents of a post created from a repo
type RepoPost struct {
	Repo             *github.Repository
	Title            string
	Summary          string
	Slug             string
	Tags             []string
	MarkdownBody     string
	PostFileName     string
	PostFileContents string
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func copyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

func copyDirectoryRecursively(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := filepath.Join(src, fd.Name())
		dstfp := filepath.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = copyDirectoryRecursively(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = copyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
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

func getReposForUser(user string, cache bool) ([]*github.Repository, error) {

	cachedReposPath := getCachedReposPathForUser(user)
	if cache {
		if fileExists(cachedReposPath) {
			blob, _ := ioutil.ReadFile(cachedReposPath)
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

	var userRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.List(ctx, user, opt)
		if err != nil {
			fmt.Printf("failed to list repositories for user %s", user)
			return nil, err
		}
		userRepos = append(userRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	if cache {
		os.MkdirAll(tempDirectoryName, os.ModePerm)
		bytes, err := json.Marshal(userRepos)
		if err != nil {
			fmt.Printf("json.Marshal failed")
			return nil, err
		}
		err = ioutil.WriteFile(cachedReposPath, bytes, 0644)
		if err != nil {
			fmt.Printf("ioutil.WriteFile(%s) failed", cachedReposPath)
			return nil, err
		}

	}

	return userRepos, nil
}

func getFilteredRepos(repos []*github.Repository) ([]*github.Repository, error) {
	var filteredRepos []*github.Repository
	for _, repo := range repos {
		var match bool = false
		includeFiltersString := os.Getenv("REPO_NAME_INCLUDE_FILTERS")
		includeFilters := strings.Split(includeFiltersString, ",")

		for _, includeFilter := range includeFilters {
			re := regexp.MustCompile(includeFilter)
			if re.Match([]byte(*repo.Name)) == true {
				match = true
			}
		}

		if match == true {
			filteredRepos = append(filteredRepos, repo)
		}
	}
	return filteredRepos, nil
}

func getPostTitle(repoName string) string {

	indexOf := func(s []string, e string) int {
		for i, a := range s {
			if a == e {
				return i
			}
		}
		return -1
	}

	wordsToCorrectCasing := getEnvAsArray("WORDS_TO_CORRECT_CASING_LIST")
	wordsToCorrectCasingLowerCase := make([]string, 0)

	for _, word := range wordsToCorrectCasing {
		wordsToCorrectCasingLowerCase = append(wordsToCorrectCasingLowerCase, strings.ToLower(word))
	}

	title := strings.Replace(repoName, "playground", "", -1)
	title = strings.Replace(title, "-", " ", -1)
	title = strings.TrimSpace(title)
	title = strings.ToLower(title)
	title = strings.Title(title)

	lowerCaseWords := strings.Split(strings.ToLower(title), " ")

	titleWords := make([]string, 0)

	for _, lowerCaseWord := range lowerCaseWords {
		index := indexOf(wordsToCorrectCasingLowerCase, lowerCaseWord)
		if index != -1 {
			titleWords = append(titleWords, wordsToCorrectCasing[index])
		} else {
			titleWords = append(titleWords, strings.Title(lowerCaseWord))
		}
	}

	return strings.Join(titleWords, " ")
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func getURLResponseCacheDirectory() string {
	return filepath.Join(tempDirectoryName, "url-response-cache")
}

func getURLResponseCacheFilePath(url string) string {
	return filepath.Join(getURLResponseCacheDirectory(), getMD5Hash(url))
}

func getURLResponseBody(url string, cache bool) (string, error) {

	urlResponseCacheFilePath := getURLResponseCacheFilePath(url)
	if cache {
		if fileExists(urlResponseCacheFilePath) {
			b, err := ioutil.ReadFile(urlResponseCacheFilePath)
			if err != nil {
				fmt.Printf("ioutil.ReadFile(%s) failed\n", urlResponseCacheFilePath)
				return "", err
			}
			return string(b), nil
		}
	}

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

	if cache {
		os.MkdirAll(getURLResponseCacheDirectory(), os.ModePerm)
		err := ioutil.WriteFile(urlResponseCacheFilePath, data, os.ModePerm)

		if err != nil {
			fmt.Printf("ioutil.WriteFile(%s) failed\n", urlResponseCacheFilePath)
			return "", err
		}
	}

	return string(data), nil
}

func getPostBodyForRepo(repo *github.Repository) (string, error) {
	url := "https://raw.githubusercontent.com/" + *repo.FullName + "/master/README.md"
	contents, err := getURLResponseBody(url, useCache)
	if err != nil {
		fmt.Printf("getURLContents(%s) failed", url)
		//return "", err
	}

	if err != nil {
		fmt.Printf("failed to getPostBodyForRepo(%s)\n", *repo.Name)
		fmt.Printf("no README.md at \"%s\" setting markdownBody to link to repo(%s)\n", url, *repo.Name)
		contents = "See github repo at [" + *repo.FullName + "](" + *repo.HTMLURL + ")"
	}

	return contents, nil
}

func getPostFileNameForRepo(repo *github.Repository) string {
	return "generated-" + *repo.Name + ".md"
}

func randomSummaryPrefix() string {
	randomSummaryPrefixListString := os.Getenv("RANDOM_SUMMARY_PREFIX_LIST")
	rand.Seed(time.Now().Unix())
	randomSummaryPrefixList := strings.Split(randomSummaryPrefixListString, ",")
	return randomSummaryPrefixList[rand.Intn(len(randomSummaryPrefixList))]
}

func getEnvAsArray(key string) []string {
	result := []string{}
	value := os.Getenv(key)
	if value != "" {
		result = strings.Split(value, ",")
	}
	return result
}

func arrayIntersection(a, b []string) (c []string) {
	m := make(map[string]bool)

	for _, item := range a {
		m[item] = true
	}

	for _, item := range b {
		if _, ok := m[item]; ok {
			c = append(c, item)
		}
	}
	return
}

func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func getPostTags(repo *github.Repository) []string {
	autoTagsIfInRepoName := getEnvAsArray("AUTO_TAGS_IF_IN_REPO_NAME")
	staticTags := getEnvAsArray("STATIC_TAGS")
	repoNameTagMappingsString := os.Getenv("REPO_NAME_TAG_MAPPINGS")
	tagMapJSON := os.Getenv("TAG_MAP_JSON")

	tagMap := make(map[string]string)
	err := json.Unmarshal([]byte(tagMapJSON), &tagMap)
	if err != nil {
		panic(err)
	}

	words := strings.Split(*repo.Name, "-")
	autoTags := arrayIntersection(autoTagsIfInRepoName, words)

	repoNameTagMappings := strings.Split(repoNameTagMappingsString, "|")

	repoMappingTags := []string{}
	for _, repoNameTagMapping := range repoNameTagMappings {
		repoNameToTagsList := strings.Split(repoNameTagMapping, "=")
		if len(repoNameToTagsList) < 2 {
			continue
		}
		repoName := repoNameToTagsList[0]
		tagsListString := repoNameToTagsList[1]
		tags := strings.Split(tagsListString, ",")

		for _, tag := range tags {
			if repoName == *repo.Name {
				repoMappingTags = append(repoMappingTags, tag)
			}
		}

	}

	allPostTags := append(autoTags, append(repoMappingTags, staticTags...)...)

	resultPostTags := make([]string, 0)
	for _, postTag := range allPostTags {

		if tagName, ok := tagMap[postTag]; ok {
			resultPostTags = append(resultPostTags, tagName)
		} else {
			resultPostTags = append(resultPostTags, postTag)
		}
	}

	return unique(resultPostTags)
}

func getPostSlug(repo *github.Repository) string {
	postTitle := getPostTitle(*repo.Name)
	return strings.ToLower(strings.Replace(postTitle, " ", "-", -1))
}
func newRepoPost(repo *github.Repository) (*RepoPost, error) {
	markdownBody, err := getPostBodyForRepo(repo)

	if err != nil {
		fmt.Printf("failed to getPostBodyForRepo(%s)\n", *repo.Name)
		return nil, err
	}

	lines := strings.Split(markdownBody, "\n")

	markdownBody = strings.Join(lines[1:], "\n")

	title := getPostTitle(*repo.Name)
	repoPost := &RepoPost{
		Repo:         repo,
		Title:        title,
		Summary:      randomSummaryPrefix() + " " + title,
		Slug:         getPostSlug(repo),
		Tags:         getPostTags(repo),
		MarkdownBody: markdownBody,
		PostFileName: getPostFileNameForRepo(repo),
	}
	postFileContents, err := getPostFileContents(repoPost)
	if err != nil {
		fmt.Printf("getPostFileContents(%s) failed\n", *repo.Name)
		return nil, err
	}

	repoPost.PostFileContents = postFileContents
	return repoPost, nil
}

func getPostFileContents(repoPost *RepoPost) (string, error) {
	postTemplatePath := filepath.Join("templates", "post.md")

	b, err := ioutil.ReadFile(postTemplatePath)
	contentsAsString := string(b)

	t := template.Must(template.New("hugo-markdown-post-tmpl").Parse(contentsAsString))

	var buf bytes.Buffer

	err = t.Execute(&buf, repoPost)
	if err != nil {
		panic(err)
	}

	result := buf.String()
	return result, nil
}

func getCachedReposPathForUser(user string) string {
	return filepath.Join(tempDirectoryName, "repo-list-"+user+".json")
}

func getAndSaveReposForUser(user string, path string) error {
	result, err := getReposForUser(user, false)
	if err != nil {
		fmt.Printf("getReposForUser(%s) failed\n", user)
		return err
	}
	bytes, err := json.Marshal(result)
	if err != nil {
		fmt.Printf("getAndSaveReposForUser | json.Marshal() failed\n")
		return err
	}
	err = ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		fmt.Printf("getAndSaveReposForUser | ioutil.WriteFile(%s) failed\n", path)
		return err
	}

	return nil
}

func createMarkdownPostFile(repoPost RepoPost, destinationDirectory string) error {

	if err := os.MkdirAll(destinationDirectory, os.ModePerm); err != nil {
		fmt.Printf("os.MkdirAll(%s) failed\n", destinationDirectory)
		return err
	}

	path := filepath.Join(destinationDirectory, repoPost.PostFileName)
	if err := ioutil.WriteFile(path, []byte(repoPost.PostFileContents), 0644); err != nil {
		fmt.Printf("ioutil.WriteFile(%s) failed\n", path)
		return err
	}

	return nil
}

func getFilteredReposForUser(user string) ([]*github.Repository, error) {
	repos, err := getReposForUser(user, true)
	if err != nil {
		fmt.Printf("getReposForUser(%s) failed\n", user)
		return nil, err
	}

	filteredRepos, err := getFilteredRepos(repos)
	if err != nil {
		fmt.Printf("applyIncludeFilterForRepos(repos) failed\n")
		return nil, err
	}

	return filteredRepos, nil
}

func getRepoPosts(username string) ([]RepoPost, error) {
	repoPosts := make([]RepoPost, 0)

	filteredRepos, err := getFilteredReposForUser(user)
	if err != nil {
		fmt.Printf("getFilteredReposForUser(%s) failed\n", user)
		return nil, err
	}

	for _, repo := range filteredRepos {
		repoPost, err := newRepoPost(repo)
		if err != nil {
			fmt.Printf("newRepoPost(%s) failed\n", *repo.Name)
			return nil, err
		}
		repoPosts = append(repoPosts, *repoPost)
	}

	return repoPosts, nil
}

type RepoPostPredicate func(repoPost RepoPost) bool

func getFilteredRepoPosts(username string, fn RepoPostPredicate) ([]RepoPost, error) {
	filteredRepoPosts := make([]RepoPost, 0)

	repoPosts, err := getRepoPosts(username)
	if err != nil {
		fmt.Printf("getRepoPosts(%s) failed\n", username)
		return nil, err
	}

	for _, repoPost := range repoPosts {
		if fn(repoPost) {
			filteredRepoPosts = append(filteredRepoPosts, repoPost)
		}
	}

	return filteredRepoPosts, nil
}

func getRepoPostsWithNoTags(username string) ([]RepoPost, error) {
	filteredRepoPosts, err := getFilteredRepoPosts(username, func(repoPost RepoPost) bool {
		return len(repoPost.Tags) == 0
	})
	if err != nil {
		fmt.Printf("getFilteredRepoPosts(%s) failed\n", username)
		return nil, err
	}
	return filteredRepoPosts, nil
}

func createMarkdownPostFiles(username string, destinationDirectory string) error {
	if debug {
		fmt.Printf("getRepoPosts(%s)\n", username)
	}
	repoPosts, err := getRepoPosts(username)
	if err != nil {
		fmt.Printf("getRepoPosts(%s) failed\n", username)
		return nil
	}

	for _, repoPost := range repoPosts {
		if debug {
			fmt.Printf("createMarkdownPostFile(%s, \"%s\")\n", *repoPost.Repo.Name, destinationDirectory)
		}
		err := createMarkdownPostFile(repoPost, destinationDirectory)
		if err != nil {
			fmt.Printf("generateMarkdownPostFile(%s) failed\n", *repoPost.Repo.Name)
			//return err
		}
	}

	return nil
}

func main() {
	flag.Parse()

	if command == "fetch-and-save-repos-for-user" {
		fmt.Printf("command: %s, user: %s, path: %s\n", command, user, path)
		if err := getAndSaveReposForUser(user, path); err != nil {
			log.Fatal(err)
		}
	}

	if command == "generate-markdown-post-files" {
		fmt.Printf("command: %s, user: %s, destinationDirectory: %s\n", command, user, destinationDirectory)
		if err := createMarkdownPostFiles(user, destinationDirectory); err != nil {
			log.Fatal(err)
		}
	}

}
