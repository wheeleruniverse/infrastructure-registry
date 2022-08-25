
# https://github.com/aws/aws-lambda-go#for-developers-on-windows

$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"

go build -o main main.go

~\go\bin\build-lambda-zip.exe -o main.zip main

Remove-Item main
