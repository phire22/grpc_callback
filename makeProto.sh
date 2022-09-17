#!/bin/bash
outdir="./chat"
protoc --go_out=$outdir --go_opt=paths=source_relative --go-grpc_out=$outdir --go-grpc_opt=paths=source_relative chat.proto

outdir="./hook"
protoc --go_out=$outdir --go_opt=paths=source_relative --go-grpc_out=$outdir --go-grpc_opt=paths=source_relative hook.proto

outdir="./register"
protoc --go_out=$outdir --go_opt=paths=source_relative --go-grpc_out=$outdir --go-grpc_opt=paths=source_relative register.proto
