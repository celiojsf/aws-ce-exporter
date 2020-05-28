build:
	go build

install:
	go install

clean:
	go clean
	rm -rf builds

all-platforms:
	env GOOS=linux GOARCH=amd64 go build -o builds/aws-ce-exporter-linux-amd64
	env GOOS=windows GOARCH=amd64 go build -o builds/aws-ce-exporter-win-amd64.exe
	env GOOS=darwin GOARCH=amd64 go build -o builds/aws-ce-exporter-darwin-amd64