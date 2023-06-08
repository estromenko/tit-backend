#!/bin/bash

export PASSWORD="${PASSWORD}"

Xvfb :0 -screen 0 1288x724x24 &

fluxbox -display :0 &
fluxbox_pid=$!

x11vnc -display :0 -forever -passwd "${PASSWORD}" &

websockify 0.0.0.0:8888 localhost:5900 &

wait "${fluxbox_pid}"
