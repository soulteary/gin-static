#!/bin/bash

OWNER=soulteary
REPO=gin-static

command -v curl >/dev/null 2>&1 || { echo "需要安装 curl"; exit 1; }
command -v jq >/dev/null 2>&1 || { echo "需要安装 jq"; exit 1; }

API_URL="https://api.github.com/repos/$OWNER/$REPO/releases/latest"
RESPONSE=$(curl -s $API_URL)

if [ $? -ne 0 ]; then
    echo "Get latest release info failed"
    exit 1
fi

VERSION=$(echo $RESPONSE | jq -r .tag_name)
PUBLISHED_AT=$(echo $RESPONSE | jq -r .published_at)
DOWNLOAD_URL=$(echo $RESPONSE | jq -r .tarball_url)

if [ "$VERSION" = "null" ]; then
    echo "Can't get latest release info"
    exit 1
fi

echo "Latest Version: $VERSION"
echo "Publish Date: $PUBLISHED_AT"
echo "Download URL: $DOWNLOAD_URL"

echo "Publishing $VERSION to golang.org package registry"
GOPROXY=proxy.golang.org go list -m github.com/$OWNER/$REPO@$VERSION
echo "Done"