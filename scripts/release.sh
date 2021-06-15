#!/usr/bin/env bash

set -e

docker login -u $DOCKER_USERNAME -p $DOCKER_PASSWORD

tag=$(git describe --abbrev=0 --tags)

function release() {
  name="$1"
  tag="$2"
  image="gotwaygateway/$name"
  platform="linux/amd64,linux/arm64,linux/arm"

  echo "🏗    Building '$image'"
  docker buildx create --name "$name" --use --append
  docker buildx build --platform "$platform" --build-arg SERVICE="$name" -t "$image:$tag" -t "$image:latest" --push .
  docker buildx imagetools inspect "$image:latest"
}

release "gotway" "$tag" .

for ms in $(ls -d cmd/*); do
  name=$(basename "$ms")
  release "$name" "$tag" .
done
