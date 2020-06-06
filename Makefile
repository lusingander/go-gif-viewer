.PHONY: all run build bundle clean credits

BINARY_NAME=go-gif-viewer

all: build

run:
	go run *.go

build: bundle
	go build -o $(BINARY_NAME)

bundle:
	fyne bundle ./resource/icons/svg/ > ./resource.go

clean:
	rm $(BINARY_NAME)

credits:
	fyne-credits-generator > ./credits.go