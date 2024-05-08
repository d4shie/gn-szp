.ONESHELL:

build:
	go generate ./...
ifeq ($(OS),Windows_NT)
#TODO: Determine whether user is running from CMD.EXE or PowerShell. Powershell is currently assumed.
	@echo Detected that Make is being run on Windows: Be aware that the compiled binary will only be for linux/amd64.
	@pwsh -Command '$$env:GOOS="linux"; $$env:GOARCH="amd64"; go build -o bin/gn-szp -trimpath -ldflags="-s -w" ./cmd/gn-szp; $$env:GOOS=$$null; $$env:GOARCH=$$null'
else
	GOOS="linux" GOARCH="amd64" go build -o bin/gn-szp -trimpath -ldflags="-s -w" ./cmd/gn-szp
endif