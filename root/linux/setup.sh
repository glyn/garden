#!/bin/bash

[ -n "$DEBUG" ] && set -o xtrace
set -o nounset
set -o errexit
shopt -s nullglob

cd $(dirname "${0}")

# Check if the old mount point exists, and if so clean it up
if [ -d /dev/cgroup ]
then
  if grep -q /dev/cgroup /proc/mounts
  then
    umount /dev/cgroup
  fi

  rmdir /dev/cgroup
fi

if [ ! -d /tmp/warden/cgroup ]
then
  mkdir -p /tmp/warden

  # Mount tmpfs
  if ! grep "/tmp/warden " /proc/mounts | cut -d' ' -f3 | grep -q tmpfs
  then
    mount -t tmpfs none /tmp/warden
  fi

  mkdir -p /tmp/warden/cgroup

  # Mount cgroup subsystems individually
  mount -t cgroup none /tmp/warden/cgroup
fi

./net.sh setup

# Disable AppArmor if possible
if [ -x /etc/init.d/apparmor ]; then
  /etc/init.d/apparmor teardown
fi

# quotaon(8) exits with non-zero status when quotas are ENABLED
if [ "$DISK_QUOTA_ENABLED" = "true" ] && quotaon -p $CONTAINER_DEPOT_MOUNT_POINT_PATH > /dev/null 2>&1
then
  mount -o remount,usrjquota=aquota.user,grpjquota=aquota.group,jqfmt=vfsv0 $CONTAINER_DEPOT_MOUNT_POINT_PATH
  quotacheck -ugmb -F vfsv0 $CONTAINER_DEPOT_MOUNT_POINT_PATH
  quotaon $CONTAINER_DEPOT_MOUNT_POINT_PATH
elif [ "$DISK_QUOTA_ENABLED" = "false" ] && ! quotaon -p $CONTAINER_DEPOT_MOUNT_POINT_PATH > /dev/null 2>&1
then
  quotaoff $CONTAINER_DEPOT_MOUNT_POINT_PATH
fi
