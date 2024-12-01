#!/bin/sh
set -e

echo "Installing gqlgen..."
go get github.com/99designs/gqlgen

echo "Generating GraphQL code..."
go run github.com/99designs/gqlgen generate

echo "Done! 🎉" 