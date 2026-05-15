build:
	go build -gcflags -S -o ./gru.out

install-dev: build
	mv ./gru.out ~/.local/bin/gru