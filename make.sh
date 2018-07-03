#!/bin/bash

# make.sh: a script for building the chatter web application

echo ">>> Building the chatter app..."
go build -o chatter
echo ">>> Build done."
