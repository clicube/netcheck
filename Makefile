fmt:
	go fmt ./...

build:
	mkdir -p ./bin
	go build -o ./bin/netcheckd ./
	sudo setcap cap_net_raw=+ep bin/netcheckd

clean:
	go clean ./...
	rm -rf ./bin

test:
	go test ./...

remote-build:
	mkdir -p ./bin
	GOOS=linux GOARCH=arm GOARM=6 go build -o ./bin/netcheckd ./

remote-copy: remote-build
	ssh raspi "mkdir -p services/netcheck/bin"
	scp ./bin/netcheckd raspi:services/netcheck/bin/netcheckd
	ssh raspi "sudo setcap cap_net_raw=+ep services/netcheck/bin/netcheckd"

remote-run: remote-copy
	ssh raspi "cd services/netcheck/bin && ./netcheckd"
