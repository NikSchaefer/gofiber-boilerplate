.PHONY: run test clean
default: run

# Rebuild the Golang service inside the long running container. 
run: run-base
	docker exec -it helloworld-service go run /go/src/app/main.go

# Run all the tests
test: run-base
	docker exec -it helloworld-service go test -cover ./...

# Remove the container.
clean:
	docker stop helloworld-service && docker rm helloworld-service || exit 0

# run-base: Create a long running container in the background.
run-base:
	if [ "$(shell docker ps --filter=name=helloworld-service -q)" = "" ]; then \
		docker build --target builder -t helloworld-service-base .; \
		docker run \
			-p 0.0.0.0:8080:8080/tcp \
			--mount type=bind,source=$(shell pwd),target=/go/src/app \
			--name helloworld-service \
			-td helloworld-service-base; \
	fi;