#!/bin/bash

GOCMD="go build -o"
cursor=">>>"

message="$cursor Building "

echo "$message domainfinder..."
$GOCMD domainfinder

echo "$message available..."
cd ../available
$GOCMD ../domainfinder/lib/available

echo "$message coolify..."
cd ../coolify
$GOCMD ../domainfinder/lib/coolify

echo "$message domainify..."
cd ../domainify
$GOCMD ../domainfinder/lib/domainify

echo "$message sprinkle..."
cd ../sprinkle
$GOCMD ../domainfinder/lib/sprinkle

echo "$message synonyms..."
cd ../synonyms
$GOCMD ../domainfinder/lib/synonyms

cd ../build
echo "$message Building done."
