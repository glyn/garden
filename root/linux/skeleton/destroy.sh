#!/bin/bash

[ -n "$DEBUG" ] && set -o xtrace
set -o nounset
set -o errexit
shopt -s nullglob

cd $(dirname $0)

source ./etc/config

./net.sh teardown

if [ -f ./run/wshd.pid ]
then
  pid=$(cat ./run/wshd.pid)
  cgroup_path=/tmp/warden/cgroup/instance-$id
  tasks=$cgroup_path/tasks

  if [ -d $cgroup_path ]
  then
    while true
    do
      kill -9 $pid 2> /dev/null || true

      # Wait while there are tasks in one of the instance's cgroups
      if [ -f $tasks ] && [ -n "$(cat $tasks)" ]
      then
        sleep 0.1
      else
        break
      fi
    done
  fi

  # Done, remove pid
  rm -f ./run/wshd.pid

  # Remove cgroups
  if [ -d $cgroup_path ]
  then
    # Remove nested cgroups for nested-warden
    rmdir $cgroup_path/instance* 2> /dev/null || true
    rmdir $cgroup_path
  fi

  exit 0
fi
