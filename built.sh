#!/bin/bash

BUILT=`git log -1 --format=%cd --date=iso8601-strict`

echo $BUILT

