APP=seed_db
BIN_DIR=bin

build: linux windows

linux:
	mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/$(APP)-linux .

windows:
	mkdir -p $(BIN_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(BIN_DIR)/$(APP)-windows.exe .

clean:
	rm -rf $(BIN_DIR)