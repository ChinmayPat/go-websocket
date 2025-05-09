dev-run:
	go run cmd/main.go

templ:
	templ generate

air:
	air --build.cmd "go build -o cmd/main.go"