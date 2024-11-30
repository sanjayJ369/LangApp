tests:
	go test ./... -race -count=1 -timeout 30s 

install_pre_commit_hooks:
	pre-commit install
