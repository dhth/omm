#!/usr/bin/env bash

if [ "$#" -ne 1 ]; then
    echo "Usage: ./checknumtasks.sh <expected_num_tasks>"
    exit 1
fi

num_tasks=$(/var/tmp/omm_head tasks | wc -l | xargs)
if [ "$num_tasks" -ne "$1" ]; then
    echo "Number of tasks: $num_tasks; expected: $1"
    exit 1
fi
