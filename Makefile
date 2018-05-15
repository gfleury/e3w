all: push

TAG = latest
PREFIX = gfleury/e3w

container:
	docker build --pull -t $(PREFIX):$(TAG) .

push: container
	docker push $(PREFIX):$(TAG)

build:
	CGO_ENABLED=0 go build -ldflags='-w' 