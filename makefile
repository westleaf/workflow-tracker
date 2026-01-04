APP_NAME=workflow-tracker
WIN_BIN=$(APP_NAME).exe
LINUX_BIN=$(APP_NAME)

.PHONY: win run-win linux run clean

# Build Windows binary
win:
	GOOS=windows GOARCH=amd64 go build -o /mnt/c/Users/joaki/Desktop/$(WIN_BIN) .

# Run Windows binary (from Windows, not WSL)
run-win: win
	"/mnt/c/Users/joaki/Desktop/$(WIN_BIN)"

# Build Linux binary (for WSL or real Linux)
linux:
	GOOS=linux GOARCH=amd64 go build -o $(LINUX_BIN) .

# Run Linux binary (WSL)
run: linux
	./$(LINUX_BIN)

clean:
	rm -f $(WIN_BIN) $(LINUX_BIN)
