#!/bin/sh

SESSION="slowlorisdb"

tmux kill-session -t $SESSION || true

tmux new-session -d -s $SESSION
tmux split-window -d -t 0 -v

tmux send-keys -t 0 "go run main.go --config config0.yaml start" enter
tmux send-keys -t 1 "go run main.go --config config1.yaml start" enter

tmux attach -t $SESSION

