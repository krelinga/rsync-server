#! /usr/bin/bash
set -e

export DEBIAN_FRONTEND=noninteractive

apt-get update
apt-get install -y --no-upgrade rsync
apt-get autoremove -y
apt-get clean
