#!/bin/sh
set -e

/usr/bin/pluto-migrate --config.file=$ConfigFile

/usr/bin/pluto-server --config.file=$ConfigFile