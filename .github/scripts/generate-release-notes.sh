#!/bin/sh

TAGS=$(cat CHANGELOG.md | grep -o "##.*\d\.\d\.\d")
CURR_TAG=$(echo "$TAGS" | sed -n '1p')
PREV_TAG=$(echo "$TAGS" | sed -n '2p')

cat CHANGELOG.md | tr '\n' '\a' | grep -o "$CURR_TAG.*$PREV_TAG" | tr '\a' '\n' | tail -n +2 | sed '$d' | awk 'NF' > changelog-temp.md
