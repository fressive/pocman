MODULE_DIRS := common cli agent server
CMD_AGENT := ./agent/cmd/agent
CMD_SERVER := ./server/cmd/server

.PHONY: all buf-generate tidy build-agent build-server build clean proto

all: build-server build-agent

buf-generate:
	cd common && buf generate

tidy:
	@for mod in $(MODULE_DIRS); do \
		echo "tidying $$mod"; \
		( cd $$mod && go mod tidy ); \
	done

build-agent:
	go build -v -o pocman-agent $(CMD_AGENT)

build-server:
	go build -v -o pocman-server $(CMD_SERVER)

build: all

proto: buf-generate

clean:
	rm -f pocman-agent pocman-server