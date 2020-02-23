#!/bin/sh
set -e
if [ $AutoMigration == "true" ]; then
  /usr/bin/pluto-migrate --config.file=$ConfigFile
fi

/usr/bin/pluto-server --config.file=$ConfigFile