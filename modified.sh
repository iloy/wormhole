#!/bin/bash

MODIFIED=`git status --porcelain`
if [ -n "$MODIFIED" ]; then
    MODIFIED="(modified)"
else
    CHERRY=`git cherry`
    if [ -n "$CHERRY" ]; then
        MODIFIED="(modified)"
    else
        MODIFIED="(not-modified)"
    fi
fi

echo $MODIFIED

