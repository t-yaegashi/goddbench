#!/bin/bash
PROJECT=$1
if test -z "$PROJECT"; then
        echo "Specify GAE project name"
        exit 1
fi
for i in app/*; do
        curl https://goddbench-${i##*/}-dot-$PROJECT.appspot.com/
done
