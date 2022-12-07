# Build Go program for multiple OSs and architectures

# Targets
all: darwin_amd64 linux_amd64 linux_arm64

# Architecture flags
ARCH_AMD64 = amd64
ARCH_ARM64 = arm64

# Locations
SRC = aigitmsg/main.go
OUT = bin/

# Go build command
GOBUILD = go build

# Build targets
darwin_amd64:
	GOOS=darwin GOARCH=$(ARCH_AMD64) $(GOBUILD) -o $(OUT)/aigitmsg_darwin_amd64 $(SRC)

linux_amd64:
	GOOS=linux GOARCH=$(ARCH_AMD64) $(GOBUILD) -o $(OUT)aigitmsg_linux_amd64 $(SRC)

linux_arm64:
	GOOS=linux GOARCH=$(ARCH_ARM64) $(GOBUILD) -o $(OUT)aigitmsg_linux_arm64 $(SRC)

windows_amd64:
	GOOS=windows GOARCH=$(ARCH_AMD64) $(GOBUILD) -o $(OUT)aigitmsg_windows_amd64 $(SRC)
