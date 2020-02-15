.PHONY: build
build: ## Build lancelot binary
	go build -o bin/lancelot main/main.go 
