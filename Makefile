BINARY  := urlshine
VERSION := 2.0.0
LDFLAGS := -ldflags "-X main.version=$(VERSION) -s -w"

.PHONY: all build tidy clean install run-demo windows linux

all: tidy build

build:
	@echo "[BUILD] Compiling $(BINARY) ..."
	go build $(LDFLAGS) -o $(BINARY) .
	@echo "[✔] $(BINARY) binary ready — run 'bash install.sh' to install"

windows:
	@echo "[BUILD] Windows binary ..."
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY).exe .
	@echo "[✔] ./$(BINARY).exe"

linux:
	@echo "[BUILD] Linux binary ..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY)_linux .
	@echo "[✔] ./$(BINARY)_linux"

tidy:
	go mod tidy

install: build
	cp $(BINARY) /usr/local/bin/$(BINARY)
	@echo "[✔] Installed to /usr/local/bin/$(BINARY)"

clean:
	rm -f $(BINARY) $(BINARY).exe $(BINARY)_linux
	rm -rf urlshine_*
