#!/usr/bin/env bash

set -e
set -x

source /etc/profile
# work within the current docker working dir
if [ ! -f "./hcnet-core.cfg" ]; then
   cp /hcnet-core.cfg ./
fi   

echo "using config:"
cat hcnet-core.cfg

# initialize new db
hcnet-core new-db

if [ "$1" = "standalone" ]; then
  # initialize for new history archive path, remove any pre-existing on same path from base image
  rm -rf ./history
  hcnet-core new-hist vs

  # serve history archives to aurora on port 1570
  pushd ./history/vs/
  python3 -m http.server 1570 &
  popd
fi

exec hcnet-core run
