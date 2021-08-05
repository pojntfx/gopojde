all: build

backend:
	go build -o out/gopojde-backend/gopojde-backend cmd/gopojde-backend/main.go

manager-browser:
	rm -rf pkg/web/manager out/gopojde-manager-browser
	mkdir -p pkg/web/manager/web out/gopojde-manager-browser
	GOOS=js GOARCH=wasm go build -o pkg/web/manager/web/app.wasm cmd/gopojde-manager/main.go
	BUILDER=true go run cmd/gopojde-manager/main.go -build -out pkg/web/manager
	cp -r web/manager/* pkg/web/manager/web
	cp -r pkg/web/manager/* out/gopojde-manager-browser

manager-wrapper: manager-browser
	go build -o out/gopojde-manager-wrapper/gopojde-manager cmd/gopojde-manager/main.go

build: backend manager-browser manager-wrapper

release-backend:
	go build -a -ldflags '-extldflags "-static"' -o "$(shell [ "$(DST)" = "" ] && echo out/release/gopojde-backend/gopojde-backend.linux-$$(uname -m) || echo $(DST) )" cmd/gopojde-backend/main.go

release-manager-browser: manager-browser
	rm -rf out/release/gopojde-manager
	mkdir -p out/release/gopojde-manager
	cd out/gopojde-manager-browser && tar -czvf ../release/gopojde-manager/gopojde-manager.tar.gz .

release-manager-browser-github-pages: manager-browser
	rm -rf out/release/gopojde-manager-github-pages
	mkdir -p out/release/gopojde-manager-github-pages
	cp -r out/gopojde-manager-browser/* out/release/gopojde-manager-github-pages
	BUILDER=true go run cmd/gopojde-manager/main.go -build -path gopojde -out pkg/web/manager

release-manager-wrapper: manager-browser
	go build -a -ldflags '-extldflags "-static"' -o "$(shell [ "$(DST)" = "" ] && echo out/release/gopojde-manager-wrapper/gopojde-wrapper.linux-$$(uname -m) || echo $(DST) )" cmd/gopojde-manager/main.go

release: release-backend release-manager release-manager-github-pages release-manager-wrapper

install: release-backend release-manager-wrapper
	sudo install out/release/gopojde-backend/gopojde-backend.linux-$$(uname -m) /usr/local/bin/gopojde-backend
	sudo setcap cap_net_bind_service+ep /usr/local/bin/gopojde-backend
	sudo install out/release/gopojde-manager-wrapper/gopojde-manager.linux-$$(uname -m) /usr/local/bin/gopojde-manager
	
dev:
	while [ -z "$$BACKEND_PID" ] || [ -n "$$(inotifywait -q -r -e modify pkg cmd pkg/web/manager/*.css)" ]; do\
		$(MAKE);\
		kill -9 $$BACKEND_PID 2>/dev/null 1>&2;\
		kill -9 $$MANAGER_PID 2>/dev/null 1>&2;\
		wait $$BACKEND_PID $$MANAGER_PID;\
		sudo setcap cap_net_bind_service+ep out/gopojde-backend/gopojde-backend;\
		out/gopojde-backend/gopojde-backend & export BACKEND_PID="$$!";\
		/tmp/gopojde-manager-build -serve & export MANAGER_PID="$$!";\
	done

clean:
	rm -rf out
	rm -rf pkg/api/proto/v1
	rm -rf pkg/web/manager

depend:
	# Setup CLIs
	GO111MODULE=on go get github.com/golang/protobuf/protoc-gen-go@latest
	GO111MODULE=on go get github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
	GO111MODULE=on go get github.com/golang/protobuf/protoc-gen-go@latest
	# Generate bindings
	go generate ./...