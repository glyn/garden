#!/bin/bash

[ -f etc/config ] && source etc/config


function overlay_directory_in_rootfs() {
  # Skip if exists
  if [ ! -d tmp/rootfs/$1 ]
  then
    if [ -d tmp/monkey/$1 ]
    then
      cp -r tmp/monkey/$1 tmp/rootfs/
    else
      mkdir -p tmp/rootfs/$1
    fi
  fi

  mount -n --bind tmp/rootfs/$1 tmp/monkey/$1
  mount -n --bind -o remount,$2 tmp/rootfs/$1 tmp/monkey/$1
}

function setup_fs_other() {
  mkdir -p $rootfs_path/proc

  mount -n --bind $rootfs_path tmp/monkey
  mount -n --bind -o remount,ro $rootfs_path tmp/monkey

  overlay_directory_in_rootfs /dev rw
  overlay_directory_in_rootfs /etc rw
  overlay_directory_in_rootfs /home rw
  overlay_directory_in_rootfs /sbin rw
  overlay_directory_in_rootfs /var rw

  mkdir -p tmp/rootfs/tmp
  chmod 777 tmp/rootfs/tmp
  overlay_directory_in_rootfs /tmp rw
}

function get_mountpoint() {
  df -P $1 | tail -1 | awk '{print $NF}'
}

function current_fs() {
  mountpoint=$(get_mountpoint $1)

  local mp
  local fs

  while read _ mp fs _; do
    if [ "$fs" = "rootfs" ]; then
      continue
    fi

    if [ "$mp" = "$mountpoint" ]; then
      echo $fs
    fi
  done < /proc/mounts
}

function should_use_overlayfs() {
  # load it so it's in /proc/filesystems
  modprobe -q overlayfs >/dev/null 2>&1 || true

  # cannot mount overlayfs in aufs
  if [ "$(current_fs tmp/rootfs)" == "aufs" ]; then
    return 1
  fi

  # cannot mount overlayfs in overlayfs; whiteout not supported
  if [ "$(current_fs tmp/rootfs)" == "overlayfs" ]; then
    return 1
  fi

  # check if it's a known filesystem
  grep -q overlayfs /proc/filesystems
}

function should_use_aufs() {
  # load it so it's in /proc/filesystems
  modprobe -q aufs >/dev/null 2>&1 || true

  # don't use aufs for nested warden as neither overlayfs nor aufs can mount
  # on it
  if [ "$allow_nested_warden" == "true" ]; then
    return 1
  fi

  # cannot mount aufs in aufs
  if [ "$(current_fs tmp/rootfs)" == "aufs" ]; then
    return 1
  fi

  # cannot mount aufs in overlayfs
  if [ "$(current_fs tmp/rootfs)" == "overlayfs" ]; then
    return 1
  fi

  # check if it's a known filesystem
  grep -q aufs /proc/filesystems
}

function setup_fs() {
  mkdir -p tmp/rootfs tmp/monkey

  if should_use_aufs; then
    mount -n -t aufs -o br:tmp/rootfs=rw:$rootfs_path=ro+wh none tmp/monkey
  elif should_use_overlayfs; then
    mount -n -t overlayfs -o rw,upperdir=tmp/rootfs,lowerdir=$rootfs_path none tmp/monkey
  else
    setup_fs_other
  fi
}

function teardown_fs() {
  umount tmp/monkey
}
