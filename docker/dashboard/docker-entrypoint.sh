#!/bin/bash

Xvfb :0 &

fluxbox -display :0 &
fluxbox_pid=$!

x11vnc -display :0 -forever &

websockify 0.0.0.0:8888 localhost:5900 &

wait "${fluxbox_pid}"
