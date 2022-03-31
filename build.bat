set GOOS=windows
set GOARCH=amd64
go build -o out/awm-amd64.exe .
set GOOS=linux
set GOARCH=amd64
go build -o out/awm-linux-amd64 .
