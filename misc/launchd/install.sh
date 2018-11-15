#!/bin/bash

src_dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
plist=io.udfs.udfs-daemon.plist
dest_dir="$HOME/Library/LaunchAgents"
UDFS_PATH="${UDFS_PATH:-$HOME/.udfs}"
escaped_udfs_path=$(echo $UDFS_PATH|sed 's/\//\\\//g')

UDFS_BIN=$(which udfs || echo udfs)
escaped_udfs_bin=$(echo $UDFS_BIN|sed 's/\//\\\//g')

mkdir -p "$dest_dir"

sed -e 's/{{UDFS_PATH}}/'"$escaped_udfs_path"'/g' \
  -e 's/{{UDFS_BIN}}/'"$escaped_udfs_bin"'/g' \
  "$src_dir/$plist" \
  > "$dest_dir/$plist"

launchctl list | grep udfs-daemon >/dev/null
if [ $? ]; then
  echo Unloading existing udfs-daemon
  launchctl unload "$dest_dir/$plist"
fi

echo Loading udfs-daemon
if (( `sw_vers -productVersion | cut -d'.' -f2` > 9 )); then
  sudo chown root "$dest_dir/$plist"
  sudo launchctl bootstrap system "$dest_dir/$plist"
else
  launchctl load "$dest_dir/$plist"
fi
launchctl list | grep udfs-daemon
