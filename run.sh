#!/bin/bash

SERVICE=$1
PROTO_DIR="./proto"
OUTPUT_DIR="$PROTO_DIR/$SERVICE"

# Create output directory if it doesn't exist
mkdir -p $OUTPUT_DIR

# Generate Go code from proto files
find $PROTO_DIR -name "$SERVICE*.proto" -exec protoc \
  --proto_path=$PROTO_DIR \
  --go_out=$OUTPUT_DIR \
  --go_opt=paths=source_relative \
  --go-grpc_out=$OUTPUT_DIR \
  --go-grpc_opt=paths=source_relative \
  {} \;

echo "✅ Generated code for $SERVICE"
