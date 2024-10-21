#!/bin/bash
echo "ts,x,y"
awk -F, 'NR>1 { print $1","($2/2)","($3/2) }' ./src/data/rand-xy.csv 
