#!/bin/bash

TAG_ID=$1
GITHUB_API_TOKEN=$2

OWNER="snowdrop"
REPO="generator"
AUTH="Authorization: token $GITHUB_API_TOKEN"
GH_API="https://api.github.com"
GH_REPO="$GH_API/repos/$OWNER/$REPO"
GH_TAGS="$GH_REPO/releases/tags/$TAG_ID"

# Read asset tags.
response=$(curl -sH "$AUTH" $GH_TAGS)

# Get ID of the asset based on given filename.
eval $(echo "$response" | grep -m 1 "id.:" | grep -w id | tr : = | tr -cd '[[:alnum:]]=')
[ "$id" ] || { echo "Error: Failed to get release id for tag: $tag"; echo "$response" | awk 'length($0)<100' >&2; exit 1; }
echo "ID : $id"

# Delete the Github Release
curl -X DELETE \
     -H "$AUTH" \
     $GH_REPO/releases/$id
