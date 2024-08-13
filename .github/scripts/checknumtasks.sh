#!/usr/bin/env bash

if [ "$#" -ne 2 ]; then
    echo "Usage: ./checknumtasks.sh <actual_num_tasks> <expected_num_tasks> "
    exit 1
fi

if [ "$1" -ne "$2" ]; then
    echo "Number of tasks: $1; expected: $2"
    exit 1
fi
