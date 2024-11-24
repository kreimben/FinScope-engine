# /bin/sh -c

GOOS=linux GOARCH=arm64 go build -o bootstrap main.go

chmod +x bootstrap

zip bootstrap.zip bootstrap

rm bootstrap

open .