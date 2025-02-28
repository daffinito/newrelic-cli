#!/bin/bash

COLOR_RED='\033[0;31m'
COLOR_NONE='\033[0m'
CURRENT_GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD)

if [ $CURRENT_GIT_BRANCH != 'main' ]; then
  printf "\n"
  printf "${COLOR_RED} Error: The release.sh script must be run while on the main branch. \n ${COLOR_NONE}"
  printf "\n"

  exit 1
fi

if [ $# -ne 1 ]; then
  printf "\n"
  printf "${COLOR_RED} Error: Release version argument required. \n\n ${COLOR_NONE}"
  printf " Example: \n\n    ./tools/release.sh 0.9.0 \n\n"
  printf "  Example (make): \n\n    make release version=0.9.0 \n"
  printf "\n"

  exit 1
fi

RELEASE_VERSION=$1
GIT_USER=$(git config user.email)

echo "Generating release for v${RELEASE_VERSION} using system user git user ${GIT_USER}"

git checkout -b release/v${RELEASE_VERSION}

# Auto-generate CLI documentation
NATIVE_OS=$(go version | awk -F '[ /]' '{print $4}')
if [ -x "bin/${NATIVE_OS}/newrelic" ]; then
   rm -rf docs/cli/*
   mkdir -p docs/cli
   bin/${NATIVE_OS}/newrelic documentation --outputDir docs/cli/ --format markdown
   git add docs/cli/*
   git commit -m "chore(docs): Regenerate CLI docs for v${RELEASE_VERSION}"
fi

# Auto-generate CHANGELOG updates
git-chglog --next-tag v${RELEASE_VERSION} -o CHANGELOG.md --sort semver
# Fix any spelling issues in the CHANGELOG
misspell -source text -w CHANGELOG.md

# Commit CHANGELOG updates
git add CHANGELOG.md
git commit -m "chore(changelog): Update CHANGELOG for v${RELEASE_VERSION}"
git push origin release/v${RELEASE_VERSION}
