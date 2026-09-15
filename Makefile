GRU_RUNTIME_CURRENT_VERSION := $(patsubst v%,%,$(shell git describe --tags --abbrev=0)) dev build
GRU_BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ) dev build

build:
	go build \
	-ldflags="-X 'github.com/augustofrade/gru-lua/gru.CurrentBuildVersion=$(GRU_RUNTIME_CURRENT_VERSION)' -X 'github.com/augustofrade/gru-lua/gru.CurrentBuildDate=$(GRU_BUILD_DATE)'" \
	-o ./gru.out
	echo "Build complete. Version: $(GRU_RUNTIME_CURRENT_VERSION)"

install-dev: build
	mv ./gru.out ~/.local/bin/gru