default: install

generate:
	go generate ./...

install:
	go install .

test:
	go test -count=1 -parallel=4 ./...

testacc:
	TF_ACC=1 go test -count=1 -parallel=4 -timeout 10m -v ./...

doc:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name omglol

upgrade-client:
	GOPROXY=https://proxy.golang.org GO111MODULE=on go get github.com/ejstreet/omglol-client-go@latest
	go mod tidy