all: build

daemon:
	go build -o out/gopojde-daemon/gopojde-daemon cmd/gopojde-daemon/main.go

manager-web:
	mkdir -p pkg/web/manager/web out/gopojde-manager-web
	GOOS=js GOARCH=wasm go build -o pkg/web/manager/web/app.wasm cmd/gopojde-manager/main.go
	BUILDER=true GOOS="" GOARCH="" go run cmd/gopojde-manager/main.go --build --out pkg/web/manager
	cp -rn web/manager/* pkg/web/manager/web
	cp -rn pkg/web/manager/* out/gopojde-manager-web

manager-native: manager-web
	go build -o out/gopojde-manager-native/gopojde-manager cmd/gopojde-manager/main.go

build: daemon manager-web manager-native

release-daemon:
	go build -a -ldflags '-extldflags "-static"' -o "$(shell [ "$(DST)" = "" ] && echo out/release/gopojde-daemon/gopojde-daemon.linux-$$(uname -m) || echo $(DST) )" cmd/gopojde-daemon/main.go

release-manager-web: manager-web
	mkdir -p out/release/gopojde-manager
	cd out/gopojde-manager-web && tar -czvf ../release/gopojde-manager/gopojde-manager.tar.gz .

release-manager-web-github-pages: manager-web
	mkdir -p out/release/gopojde-manager-github-pages
	cp -rn out/gopojde-manager-web/* out/release/gopojde-manager-github-pages
	BUILDER=true GOOS="" GOARCH="" go run cmd/gopojde-manager/main.go --build --path gopojde --out out/release/gopojde-manager-github-pages

release-manager-native: manager-web
	go build -a -ldflags '-extldflags "-static"' -o "$(shell [ "$(DST)" = "" ] && echo out/release/gopojde-manager-native/gopojde-native.linux-$$(uname -m) || echo $(DST) )" cmd/gopojde-manager/main.go

release: release-daemon release-manager-web release-manager-web-github-pages release-manager-native

install: release-daemon release-manager-native
	sudo install out/release/gopojde-daemon/gopojde-daemon.linux-$$(uname -m) /usr/local/bin/gopojde-daemon
	sudo setcap cap_net_bind_service+ep /usr/local/bin/gopojde-daemon
	sudo install out/release/gopojde-manager-native/gopojde-manager.linux-$$(uname -m) /usr/local/bin/gopojde-manager
	
dev:
	while [ -z "$$DAEMON_PID" ] || [ -n "$$(inotifywait -q -r -e modify pkg cmd pkg/web/manager/*.css)" ]; do\
		$(MAKE);\
		kill -9 $$DAEMON_PID 2>/dev/null 1>&2;\
		kill -9 $$MANAGER_PID 2>/dev/null 1>&2;\
		wait $$DAEMON_PID $$MANAGER_PID;\
		sudo setcap cap_net_bind_service+ep out/gopojde-daemon/gopojde-daemon;\
		out/gopojde-daemon/gopojde-daemon & export DAEMON_PID="$$!";\
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