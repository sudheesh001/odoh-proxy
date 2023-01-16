test:
	go test ./...

all:
	go build -o odoh-proxy

deploy-proxy:
	gcloud app deploy --stop-previous-version proxy.yaml