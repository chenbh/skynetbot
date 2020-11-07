FILES=$(shell find -type f -name '*.go')
.PHONY: all run deps

all: $(FILES)
	go build -o bot ./cmd

run: all
	./bot

deps:
	sudo apt-get install ffmpeg
	go get -u github.com/bwmarrin/dca/cmd/dca
	# ffmpeg -i test.mp3 -f s16le -ar 48000 -ac 2 pipe:1 | dca > test.dca
