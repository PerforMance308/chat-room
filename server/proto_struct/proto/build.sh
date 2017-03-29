#!/bin/bash -e
protoc --go_out="/Users/sqwang/Documents/mine/src/proto_struct/" --proto_path="./" "./"*.proto
