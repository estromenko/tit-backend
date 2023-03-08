#!/bin/bash

RANDOM_STRING="$(tr -dc A-Za-z0-9 </dev/urandom | head -c 10 ; echo '')"

export PASSWORD="${RANDOM_STRING}"

echo "${PASSWORD}" > /password.txt

Xvfb :0 &

fluxbox -display :0 &
fluxbox_pid=$!

x11vnc -display :0 -forever -passwd "${PASSWORD}" &

websockify 0.0.0.0:8888 localhost:5900 &

wait "${fluxbox_pid}"
