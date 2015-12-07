# Build for my raspberry
echo "Compiling ..."
GOARCH=arm GOOS=linux go build
echo "Copying new binary to raspberry ..."
scp ./http_files_client pi@raspberry:/usr/local/bin/httpfilesclient

