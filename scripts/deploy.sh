#!/bin/bash

# Check if version is provided
if [ -z "$1" ]
  then
    echo "No version number provided"
    exit 1
fi

# Check if docker is installed
if ! [ -x "$(command -v docker)" ]; then
  echo 'Error: docker is not installed.' >&2
  exit 1
fi

echo "Building container for Rincon v$1"

# Build the docker container
docker build -t bk1031/rincon:"$1" -t bk1031/rincon:latest --platform linux/amd64,linux/arm64,linux/arm/v7,linux/arm/v6 --push --progress=plain .

echo "Container deployed successfully"

# Check if GitHub CLI is installed
if ! command -v gh &> /dev/null
then
    echo "GitHub CLI (gh) is not installed. Please install it to proceed."
    exit 1
fi

# Create a release tag
git tag v$version -m "Release version $version"
git push origin v$version

# Create a release
gh release create v$version --generate-notes

echo "Package released successfully for version $version"