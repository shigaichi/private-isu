#!/bin/bash

find /home/isucon/image/storage/ -type f -regextype egrep -regex ".*/[0-9]+\.(png|jpg|gif)" | while read file; do
  filename=$(basename "$file")
  number="${filename%.*}"
  if [[ "$number" -gt 10000 ]]; then
    echo "Deleting $file"
    rm -v "$file"
  fi
done
