default:
	go build -ldflags="-s -w" -o ./bin/golden ./cmd/golden/
	strip ./bin/golden