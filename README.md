# create-blog-post-from-repo

create markdown ([hugo](https://gohugo.io/)) blog post from structured repository's `README.md` file

## Session

```sh
# fetch test data
make fetch-test-data

# run tests recursively on change
make watch-test

# run main on change
make watch-run
```

## TODO

* make relative references in README.md absolute references to the resource in github
* manual tag mappings for a repo.  can it be put in the yaml front matter of `README.md` and not display.  if not put in `.env`
* add "see corresponding github repo for this post @ ..."
* make repo post desciptions fixed for a given post.  don't want them changing between repo -> post is regenerated
* in site search.  type ahead search
  * use hugo to generate site-metadata.json which can be used by react search component
* verify google indexes all pages

## Completed

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

## Scratch

```sh
"name": "serverless-plugin-cloudfront-lambda-edge-playground",
"full_name": "pfeilbr/serverless-plugin-cloudfront-lambda-edge-playground",
"created_at": "2019-09-10T21:55:07Z",
"html_url": "https://github.com/pfeilbr/serverless-plugin-cloudfront-lambda-edge-playground",
"language": "CSS",
"full_name": "pfeilbr/16-games-in-c--sfml",


"https://raw.githubusercontent.com/" + *repo.FullName + "/master/README.md"

```
