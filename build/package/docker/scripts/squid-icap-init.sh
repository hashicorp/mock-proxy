#! /usr/bin/env bash
set -euo pipefail

# Background the VCS proxy ICAP protocol.
/vcs-mock-proxy &

# TODO: Is this necessary?
sleep 1

# Start squid in non-daemon mode (i.e., foregrounded)
squid -f /etc/squid/squid.conf -N
