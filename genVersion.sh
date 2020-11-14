#!/bin/bash

CURTAG=$(git describe --tags --abbrev=0)
if [ -z "$(git status --porcelain)" ]; then 
  NUMBER_FILES_CHANGED=$(git diff --name-only HEAD $CURTAG | wc -l)
  if [[ $NUMBER_FILES_CHANGED -eq 0 ]]; then
    # Directory is clean - No commits since last tag
    echo "${CURTAG}"
  else
    # Directory is clean - Commits since last tag
    SHORT=$(git rev-parse --short HEAD)
    echo "${CURTAG}-${SHORT}"
  fi
else 
  # Directory is not clean
  echo "${CURTAG}-DEV"
fi
