#!/bin/bash

TAG_ID=$2
GITHUB_API_TOKEN=$1

OWNER="snowdrop"
REPO="generator"
AUTH="Authorization: token $GITHUB_API_TOKEN"
GH_API="https://api.github.com"
GH_REPO="$GH_API/repos/$OWNER/$REPO"


require_clean_go_assets () {
    # Update the index
    git update-index -q --ignore-submodules --refresh
    err=0

    if git diff pkg/template/assets_vfsdata.go | grep -q 'pkg/template/assets_vfsdata.go'; then
      echo >&2 "cannot tag: you have assets unstaged changes."
        git diff-files --name-status -r --ignore-submodules -- pkg/template/assets_vfsdata.go >&2
        echo >&2 -e "\n"
        echo >&2 "Please commit or stash them."
        echo >&2 -e "\n"
        exit 1
    fi

    if git diff --stat --cached origin/master -- pkg/template/assets_vfsdata.go | grep -q 'pkg/template/assets_vfsdata.go'; then
      echo >&2 "cannot tag: your index contains assets uncommitted changes."
        git diff --stat --cached origin/master -- pkg/template/assets_vfsdata.go >&2
        echo >&2 -e "\n"
        echo >&2 "Please push or reset them."
        echo >&2 -e "\n"
        exit 1
    fi
}


require_clean_go_assets

echo "Tagging ..."
git tag -a $TAG_ID -m "$TAG_ID"

echo "Releasing $TAG_ID ..."
JSON='{"tag_name": "'"$TAG_ID"'","target_commitish": "master","name": "'"$TAG_ID"'","body": "'"$TAG_ID"'","draft": false,"prerelease": false}'
curl -H "$AUTH" \
    -H "Content-Type: application/json" \
    -d "$JSON" \
    $GH_REPO/releases
