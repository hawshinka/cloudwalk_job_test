in := qgames.log
out := qgames.json

unit-test:
	# Running unit tests...
	go test -race -tags=unit ./...

unit-test-coverage:
	# Generating coverage from tests...
	go test -race -tags=unit ./... -v -coverprofile coverage.txt
	go tool cover -html=coverage.txt -o coverage.html
	rm coverage.txt

validate-in-file:
	@ if [ -z "${in}" ]; then echo "Error: 'in' variable is not set. Defaulting to qgames.log."; fi

validate-out-file:
	@ if [ -z "${out}" ]; then echo "Error: 'out' variable is not set. Defaulting to qgames.json."; fi

run: validate-in-file validate-out-file
	# Running application...
	go run main.go -in=${in} -out=${out}
