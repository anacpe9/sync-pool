format:
	gofmt -w .

test:
	GO_ENV="test" go test ./...

test-coverage-template:
	GO_ENV="test" go test ./... -coverprofile coverage/coverage.out

test-coverage: test-coverage-template
	GO_ENV="test" go tool cover -func         coverage/coverage.out

test-coverage-html-ci: test-coverage-template
	GO_ENV="test" go tool cover -html         coverage/coverage.out -o coverage/coverage.html

test-coverage-html: test-coverage-html-ci
	open                                                               coverage/coverage.html

# https://github.com/nikolaydubina/go-cover-treemap
# go install github.com/nikolaydubina/go-cover-treemap@latest
test-coverage-treemap-ci: test-coverage-template
	go-cover-treemap            -coverprofile coverage/coverage.out > coverage/coverage.svg

test-coverage-treemap: test-coverage-treemap-ci
	open                                                              coverage/coverage.svg

# https://github.com/oligot/go-mod-upgrade
list-upgrades-available:
	go list -u -f '{{if (and (not (or .Main .Indirect)) .Update)}}{{.Path}}: {{.Version}} -> {{.Update.Version}}{{end}}' -m all 2> /dev/null

list-directly-modules-upgradable:
	go list -u -m $(go list -m -f '{{.Indirect}} {{.}}' all | grep '^false' | cut -d ' ' -f2) | grep '\['

upgrade-all-packages:
	go get -u ./...
	make requirements

requirements:
	go mod tidy

clean-packages:
	go clean -modcache

# go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
# find . -type f -name "*.model.go" -print0 | xargs -0 realpath | xargs maligned
# find . -type f -name "*.go" -print0 | tr '\n' '\0' | xargs -0 -n1 -I % sh -c 'echo "maligned %";  maligned "%"; echo ""'
# find . -type f -name "*.model.go" -print0 | xargs -0 realpath | xargs fieldalignment -fix
fix-struct:
	find . -type f -name "*.go" -print0 | tr '\n' '\0' | xargs -0 -n1 -I % sh -c 'echo "fieldalignment -fix %";  fieldalignment -fix "%"; echo ""'

fix-struct-v2:
	find . -type d ! -path '**/.*' | tr '\n' '\0' | xargs -0 -n1 -I % sh -c 'echo "fieldalignment -fix %/*.go";  fieldalignment -fix "%/*.go"; echo ""'
# find . -type d ! -path '**/.*' -print0 | tr '\n' '\0' | xargs -0 -n1 -I % sh -c 'echo "fieldalignment -fix %/*.go" && fieldalignment -fix "%/*.go"; echo ""'

