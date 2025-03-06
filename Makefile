.PHONY: build
build:
	go build -o ./build/mreminder ./cmd/medication-reminder/main.go

run_app: build
	./build/mreminder

run_app_cfgpath: build
	./build/mreminder -config-path=./config/config.yml