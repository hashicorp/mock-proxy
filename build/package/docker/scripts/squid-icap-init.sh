#! /usr/bin/env bash
set -euo pipefail

# Create a trap that will kill all processes if either exits.
PIDS=()
got_sig_chld=false
trap '
  if ! "$got_sig_chld"; then
    got_sig_chld=true
    ((${#PIDS[@]})) && kill "${PIDS[@]}" 2> /dev/null
  fi
' CHLD

# Background the VCS proxy ICAP protocol.
/vcs-mock-proxy & PIDS+=("$!")

# TODO: Is this necessary?
sleep 1

# Start squid in non-daemon mode, but bash backgrounded. Only use SSL bump if
# an SSL cert has been mounted at /etc/squid/ssl_cert.
if [ -f /etc/squid/ssl_cert/ca.pem ]; then
  squid -f /etc/squid/squid-ssl.conf -N & PIDS+=($!)
else
  squid -f /etc/squid/squid.conf -N & PIDS+=($!)
fi

# Enable "Job Control" mode, then wait:
# https://www.gnu.org/software/bash/manual/html_node/Job-Control.html#Job-Control
set -m
wait
set +m
