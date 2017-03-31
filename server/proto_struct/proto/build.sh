#!/bin/bash -e
protoc --go_out="/Users/sqwang/Documents/mine/src/server/proto_struct/" --proto_path="./" "./"*.proto
