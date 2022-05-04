
BUILD_FOLDER = build
LINUX_OUT = $(BUILD_FOLDER)/apkg-build-lin64
WINDOWS_OUT = $(BUILD_FOLDER)/apkg-build-win64

all: linux windows

linux:
	GOOS=linux GOARCH=amd64 go build -o $(LINUX_OUT)

windows:
	GOOS=windows GOARCH=amd64 go build -o $(WINDOWS_OUT)

clean:
	rm -rf $(BUILD_FOLDER)

