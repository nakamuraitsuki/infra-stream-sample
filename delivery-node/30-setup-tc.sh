#!/bin/sh
# Nginxが起動する前に実行される
tc qdisc add dev eth0 root tbf rate 100mbit burst 32kbit latency 50ms || true