#!/bin/sh

# https://stackoverflow.com/questions/53031035/generate-build-timestamp-in-go

clear

PKG_PATH='github.com/easymirror/easymirror-backend/internal/build'
BUILD_TIME=$(date +"%Y-%m-%d %H:%M:%S %Z")
CommitHash=N/A
GoVersion=N/A
GitTag=N/A

if [[ $(go version) =~ [0-9]+\.[0-9]+\.[0-9]+ ]];
then
    GoVersion=${BASH_REMATCH[0]}
fi

GV=$(git tag || echo 'N/A')
if [[ $GV =~ [^[:space:]]+ ]];
then
    GitTag=${BASH_REMATCH[0]}
fi

GH=$(git log -1 --pretty=format:%h || echo 'N/A')
if [[ GH =~ 'fatal' ]];
then
    CommitHash=N/A
else
    CommitHash=$GH
fi


# Get version from `VERSION` file
echo -e "Getting build version from file..."
VERSION=`cat version`


# Build docker image
echo -e "Building production Dockerfile..."
docker build \
--build-arg PKG_PATH="$PKG_PATH" \
--build-arg BUILD_TIME="$BUILD_TIME" \
--build-arg CommitHash="$CommitHash" \
--build-arg GoVersion="$GoVersion" \
--build-arg GitTag="$GitTag" \
--build-arg VERSION="$VERSION" \
--file Dockerfile.prod \
-t easymirror-platform-api:"$VERSION" .