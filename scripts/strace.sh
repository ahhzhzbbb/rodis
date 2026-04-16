#!/bin/bash

set -e

MAIN_DIR="/home/hoang/projects/rodis/cmd/rodis-server"
PROJECT_DIR="/home/hoang/projects/rodis"
LOG_DIR="$PROJECT_DIR/logs"
PORT="6379"
NUM_CONNS="100"
NUM_REQS="2000000"

mkdir -p "$LOG_DIR"

cd "$MAIN_DIR"
go build -o "$PROJECT_DIR/rodis-server" .

cd "$PROJECT_DIR"
taskset -c 0,1 ./rodis-server >> /dev/null 2> "$LOG_DIR/error.log" &
PID=$!

# đảm bảo server không chết ngay
sleep 1
if ! kill -0 $PID 2>/dev/null; then
    echo "rodis-server died early"
    exit 1
fi

# cleanup nếu script exit
trap "kill $PID $STRACE_PID 2>/dev/null" EXIT

# chờ server warm-up
sleep 1

sudo strace -c -p $PID 2> "$LOG_DIR/strace.log" &
STRACE_PID=$!

echo "RUNNING..."
redis-benchmark -t set,get,ping -p $PORT -c $NUM_CONNS -n $NUM_REQS > "$LOG_DIR/bench.log" 2>&1

# giống Ctrl+C
sudo kill -INT $STRACE_PID
# wait $STRACE_PID

echo "=== FUTEX ==="
awk '/futex/ {printf "time=%s%% calls=%s\n", $1, $4}' "$LOG_DIR/strace.log"

kill $PID

tail -6 "$LOG_DIR/bench.log"