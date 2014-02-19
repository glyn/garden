#!/bin/bash

[ -n "$DEBUG" ] && set -o xtrace
set -o nounset
set -o errexit
shopt -s nullglob

cd $(dirname $0)/../

source ./lib/common.sh

setup_fs

cp bin/wshd tmp/monkey/sbin/wshd
chmod 700 tmp/monkey/sbin/wshd
