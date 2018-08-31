#!/bin/bash

TAG_ID="release-"$1
GIT_TOKEN=$2

echo "Tagging ..."
git tag -a $TAG_ID -m "$TAG_ID"

echo "Releasing $TAG_ID ..."
JSON='{"tag_name": "'"$TAG_ID"'","target_commitish": "master","name": "'"$TAG_ID"'","body": "'"$TAG_ID"'-release","draft": false,"prerelease": false}'
curl -H "$AUTH" \
    -H "Content-Type: application/json" \
    -d "$JSON" \
    $GH_REPO/releases