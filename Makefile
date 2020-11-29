FILES=$(shell find -type f -name '*.go')
EXE=bot

.PHONY: run docker

${EXE}: $(FILES)
	go build -o ${EXE} ./cmd

run: ${EXE}
	./${EXE}

docker: Dockerfile $(FILES)
	docker build -t trasacom/skynetbot .
