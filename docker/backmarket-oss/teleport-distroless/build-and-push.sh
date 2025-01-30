#!/usr/bin/env bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
TMP_DIR="$SCRIPT_DIR/.tmp"
CURRENT_BRANCH=$( git rev-parse --abbrev-ref HEAD )

rm -rf "$TMP_DIR"
mkdir -p "$TMP_DIR"
cd "$TMP_DIR" || exit 1

tar -czf teleport-${CURRENT_BRANCH}-linux-amd64-bin.tar.gz -C ../../../../build teleport tctl tbot tsh fdpass-teleport
gh release create ${CURRENT_BRANCH} --generate-notes --repo backmarket-oss/teleport ./teleport-${CURRENT_BRANCH}-linux-amd64-bin.tar.gz

tar xzvf teleport-${CURRENT_BRANCH}-linux-amd64-bin.tar.gz

# To authenticate
# export CR_PAT=xxx
# echo $CR_PAT | docker login ghcr.io -u USERNAME --password-stdin
cd $SCRIPT_DIR
docker buildx build "$SCRIPT_DIR" -f Dockerfile --platform linux/amd64 -t "ghcr.io/backmarket-oss/teleport-distroless:${CURRENT_BRANCH}" --push
docker buildx build "$SCRIPT_DIR" -f Dockerfile --platform linux/amd64 -t "europe-docker.pkg.dev/teleport-prod-sijf/docker/teleport-distroless:${CURRENT_BRANCH}" --push
docker buildx build "$SCRIPT_DIR" -f Dockerfile-debug --platform linux/amd64 -t "ghcr.io/backmarket-oss/teleport-distroless-debug:${CURRENT_BRANCH}" --push
docker buildx build "$SCRIPT_DIR" -f Dockerfile-debug --platform linux/amd64 -t "europe-docker.pkg.dev/teleport-prod-sijf/docker/teleport-distroless-debug:${CURRENT_BRANCH}" --push

rm -rf "$TMP_DIR"
