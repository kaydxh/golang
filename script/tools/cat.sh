#!/bin/sh

if [ "$#" -lt 1 ]; then
    echo "Usage: $0 filename [keyword]"
    exit 1
fi

FILENAME=$1
KEYWORD=$2
LINES=$3

if [ ! -f "$FILENAME" ]; then
    echo "File not found: $FILENAME"
    exit 1
fi

if [ -z "$KEYWORD" ] && [ -z "$LINES" ]; then
    cat "$FILENAME"
elif [ -z "$KEYWORD" ]; then
    head -n "$LINES" "$FILENAME"
elif [ -z "$LINES" ]; then
    grep --color=auto "$KEYWORD" "$FILENAME"
else
    grep --color=auto "$KEYWORD" "$FILENAME" | head -n "$LINES"
fi
