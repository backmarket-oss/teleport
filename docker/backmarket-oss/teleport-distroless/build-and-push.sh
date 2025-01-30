#!/usr/bin/env bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
TMP_DIR="$SCRIPT_DIR/.tmp"

rm -rf "$TMP_DIR"
mkdir -p "$TMP_DIR"
cd "$TMP_DIR" || exit 1
wget https://github.com/backmarket-oss/teleport/releases/download/v17.2.1-bm-alloydb/teleport-v17.2.1-bm-alloydb-linux-amd64-bin.tar.gz -O teleport-v17.2.1-bm-alloydb-linux-amd64-bin.tar.gz

tar xzvf teleport-v17.2.1-bm-alloydb-linux-amd64-bin.tar.gz

# To authenticate
# export CR_PAT=xxx
# echo $CR_PAT | docker login ghcr.io -u USERNAME --password-stdin
docker buildx build "$SCRIPT_DIR" --platform linux/amd64 -t "ghcr.io/backmarket-oss/teleport-distroless:v17.2.1-bm-alloydb" --push
docker buildx build "$SCRIPT_DIR" --platform linux/amd64 -t "europe-docker.pkg.dev/teleport-prod-sijf/docker/teleport-distroless:v17.2.1-bm-alloydb" --push

rm -rf "$TMP_DIR"
