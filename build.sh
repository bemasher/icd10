#!/bin/bash

COMMIT=`git rev-parse --short HEAD`
GTag=""

go install -v -ldflags "-X main.commit=${COMMIT} -X main.gtag=${GTag}"