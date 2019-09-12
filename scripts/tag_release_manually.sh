#!/bin/bash

TAG_ID=$2
GITHUB_API_TOKEN=$1

OWNER="snowdrop"
REPO="generator"
AUTH="Authorization: token $GITHUB_API_TOKEN"
GH_API="https://api.github.com"
GH_REPO="$GH_API/repos/$OWNER/$REPO"

GREEN='\033[0;32m'
NC='\033[0m' # No Color

echo -e "${GREEN}Compiling go templates${NC}"
make assets

echo -e "${GREEN}Commit assets >> git commit -m 'fix: re-generate assets'${NC}"
git add pkg/template/assets_vfsdata.go
git commit -m "fix: re-generate assets" pkg/template/assets_vfsdata.go
echo -e "${GREEN}Pushing assets to master${NC}"
git push origin master

echo -e "${GREEN}Tagging ...${NC}"
git tag -a $TAG_ID -m "$TAG_ID"

echo -e "${GREEN}Releasing $TAG_ID ...${NC}"
JSON='{"tag_name": "'"$TAG_ID"'","target_commitish": "master","name": "'"$TAG_ID"'","body": "'"$TAG_ID"'","draft": false,"prerelease": false}'
curl -H "$AUTH" \
    -H "Content-Type: application/json" \
    -d "$JSON" \
    $GH_REPO/releases
