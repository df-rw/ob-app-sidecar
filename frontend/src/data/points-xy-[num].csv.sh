#!/bin/bash

MAX_POINTS=100
MAX_RADIUS=10
MAX_DIM=10000
ID=0

for i in "$@"; do
    case $i in
        --num=*)
            MAX_POINTS="${i#*=}"
            shift;
            shift;
            ;;
        *)
            ;;
    esac
done

echo "r,x,y"

while [ $ID -lt $MAX_POINTS ]; do
    R=$(($RANDOM % $MAX_RADIUS))
    X=$(($RANDOM % $MAX_DIM))
    Y=$(($RANDOM % $MAX_DIM))
    echo "${R},${X},${Y}"
    ID=$((ID+1))
done
