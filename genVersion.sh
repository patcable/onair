#!/bin/bash

CURTAG=$(git describe --tags --abbrev=0)

if [ -z "$(git status --porcelain)" ]; then 
  echo "${CURTAG}"
else 
  echo "${CURTAG}-dev"
fi
