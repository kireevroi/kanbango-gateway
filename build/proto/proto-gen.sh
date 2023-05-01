#!/bin/bash

protoc --go_out=./internal/userproto --go-grpc_out=./internal/userproto internal/protobuf/user.proto