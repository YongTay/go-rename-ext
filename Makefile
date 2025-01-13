# Makefile

# Go source file
SRC = main.go

# Output directory
OUT_DIR = bin

# Executable names
WIN_EXEC = $(OUT_DIR)/rename-ext.exe
LINUX_EXEC = $(OUT_DIR)/rename-ext-linux
MAC_EXEC = $(OUT_DIR)/rename-ext-mac
LINUX_ARM_EXEC = $(OUT_DIR)/rename-ext-linux-arm
LINUX_ARM64_EXEC = $(OUT_DIR)/rename-ext-linux-arm64
MAC_ARM64_EXEC = $(OUT_DIR)/rename-ext-mac-arm64

# Default target
all: windows linux mac linux-arm linux-arm64 mac-arm64

# Compile for Windows
windows:
	GOOS=windows GOARCH=amd64 go build -o $(WIN_EXEC) $(SRC)

# Compile for Linux
linux:
	GOOS=linux GOARCH=amd64 go build -o $(LINUX_EXEC) $(SRC)

# Compile for macOS
mac:
	GOOS=darwin GOARCH=amd64 go build -o $(MAC_EXEC) $(SRC)

# Compile for Linux ARM
linux-arm:
	GOOS=linux GOARCH=arm go build -o $(LINUX_ARM_EXEC) $(SRC)

# Compile for Linux ARM64
linux-arm64:
	GOOS=linux GOARCH=arm64 go build -o $(LINUX_ARM64_EXEC) $(SRC)

# Compile for macOS ARM64
mac-arm64:
	GOOS=darwin GOARCH=arm64 go build -o $(MAC_ARM64_EXEC) $(SRC)

# Clean up
clean:
	rm -rf $(OUT_DIR)/*
