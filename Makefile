tests:
	go test ./... -race -count=1 -timeout 30s 

pre_commit_install:
	pip install pre-commit