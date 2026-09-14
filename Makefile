GRU_RUNTIME_CURRENT_VERSION := 0.1
GRU_BUILD_DATE := $(shell date -u +%Y-%m-%d)

build:
	go build \
	-ldflags="-X 'github.com/augustofrade/gru-lua/gru.CurrentBuildVersion=$(GRU_RUNTIME_CURRENT_VERSION)' -X 'github.com/augustofrade/gru-lua/gru.CurrentBuildDate=$(GRU_BUILD_DATE)'" \
	-o ./gru.out
	echo "Build complete. Version: $(GRU_RUNTIME_CURRENT_VERSION)"

install-dev: build
	mv ./gru.out ~/.local/bin/gru