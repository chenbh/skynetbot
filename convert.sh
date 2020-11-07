#!/bin/bash

files=$(find audio/* -not -name "*.dsa")

set -eux
set pipefail

for f in $files; do
  ffmpeg -i $f -f s16le -ar 48000 -ac 2 pipe:1 | dca > "${f%.*}.dsa"
  mv $f archive/
done
