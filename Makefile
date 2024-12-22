install_pre_commit_hooks:
	pre-commit install

build:
	go build ./... > /dev/null

tests:
	go test ./... -race -count=1 -timeout 30s 
