tests:
	go test ./... -race -count=1

pre_commit_install:
	pip install pre-commit