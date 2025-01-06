# Load variables from .env file (if exists).
-include .env

install_pre_commit_hooks:
	pre-commit install

build:
	go build ./... > /dev/null

tests: infra-start run-tests

run-tests:
	go test ./... -race -count=1 -timeout 30s 

cleanup:
	go mod tidy

# Ollama language model with default fallback if not set.
OLLAMA_LLM := $(or $(OLLAMA_LLM), "aya-expanse:8b")

infra-start:
	docker volume create ollama
	docker compose up --detach
	docker exec ollama ollama run $(OLLAMA_LLM)

infra-stop:
	docker compose down
