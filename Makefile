install:
	go install -v

fmt:
	go fmt
	cd ./http && go fmt
	cd ./driver && go fmt

.PHONY: fmt install
