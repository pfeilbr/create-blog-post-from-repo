# create-blog-post-from-repo

create markdown ([hugo](https://gohugo.io/)) blog post from structured repository's `README.md` file

## Session

```sh
# run tests recursively on change
make watch-test

# run main on change
make watch-run
```

## TODO

* unit tests
  * create github api json response for all user repos
  * git clone test repo to local test fixtures
  * mock GetRepoReadme to point to local test repo README.md
  * create post.yaml corresponding to single playground repo for unit test
    * ensure repo README.md is sufficiently complex.  Code, images, links to repo files, etc.
    * remove first line / h1 `#` from `README.md`
* fetch all repos for user
* filter by *-playground and allow manual repo adds that don't have the *-playground naming convension
* GetRepoReadme - use regular raw URL instead of github API to minimize potential rate limiting
  * e.g. <https://raw.githubusercontent.com/pfeilbr/heroku-node-worker-playground/master/README.md>
* CreatePostTitleFromRepoName
* repoName.split('-').join(' ').TitleCase.replace('playground', '')
* <https://github.com/google/go-github>
