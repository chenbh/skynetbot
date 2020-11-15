#!/bin/bash

files=$(find audio/* -not -name "*.opus")

set -eux
set pipefail

for f in $files; do
  ffmpeg -i $f -ar 48000 -ac 2 -b:a 64K "${f%.*}.opus"
  rm $f
done
