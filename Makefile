all: build

daemon:
	go build -o out/gopojde-daemon/gopojde-daemon cmd/gopojde-daemon/main.go

manager-web:
	mkdir -p pkg/web/manager/assets/web out/gopojde-manager-web
	cp -ru web/manager/* pkg/web/manager/assets/web
	GOOS=js GOARCH=wasm go build -o pkg/web/manager/assets/web/app.wasm cmd/gopojde-manager/main.go
	BUILDER=true GOOS="" GOARCH="" go run cmd/gopojde-manager/main.go --build --out pkg/web/manager/assets
	cp -ru pkg/web/manager/assets/* out/gopojde-manager-web

manager-native: manager-web
	go build -o out/gopojde-manager-native/gopojde-manager cmd/gopojde-manager/main.go

companion-web:
	mkdir -p pkg/web/companion/assets/web out/gopojde-companion-web
	cp -ru web/companion/* pkg/web/companion/assets/web
	GOOS=js GOARCH=wasm go build -o pkg/web/companion/assets/web/app.wasm cmd/gopojde-companion/main.go
	BUILDER=true GOOS="" GOARCH="" go run cmd/gopojde-companion/main.go --build --out pkg/web/companion/assets
	cp -ru pkg/web/companion/assets/* out/gopojde-companion-web

companion-native: companion-web
	go build -o out/gopojde-companion-native/gopojde-companion cmd/gopojde-companion/main.go

build: daemon manager-web manager-native companion-web companion-native

release-daemon:
	GO386=softfloat CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o "$(shell [ "$(DST)" = "" ] && echo out/release/gopojde-daemon/gopojde-daemon.linux-$$(uname -m) || echo $(DST) )" cmd/gopojde-daemon/main.go

release-manager-web: manager-web
	mkdir -p out/release/gopojde-manager-web
	cd out/gopojde-manager-web && tar -czvf ../release/gopojde-manager-web/gopojde-manager.tar.gz .

release-manager-web-github-pages: manager-web
	mkdir -p out/release/gopojde-manager-web-github-pages
	cp -ru out/gopojde-manager-web/* out/release/gopojde-manager-web-github-pages
	BUILDER=true GOOS="" GOARCH="" go run cmd/gopojde-manager/main.go --build --path gopojde --out out/release/gopojde-manager-web-github-pages

release-manager-native: manager-web
	GO386=softfloat CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o "$(shell [ "$(DST)" = "" ] && echo out/release/gopojde-manager-native/gopojde-manager.linux-$$(uname -m) || echo $(DST) )" cmd/gopojde-manager/main.go

release-companion-web: companion-web
	mkdir -p out/release/gopojde-companion-web
	cd out/gopojde-companion-web && tar -czvf ../release/gopojde-companion-web/gopojde-companion.tar.gz .

release-companion-web-github-pages: companion-web
	mkdir -p out/release/gopojde-companion-web-github-pages
	cp -ru out/gopojde-companion-web/* out/release/gopojde-companion-web-github-pages
	BUILDER=true GOOS="" GOARCH="" go run cmd/gopojde-companion/main.go --build --path gopojde --out out/release/gopojde-companion-web-github-pages

release-companion-native: companion-web
	GO386=softfloat CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o "$(shell [ "$(DST)" = "" ] && echo out/release/gopojde-companion-native/gopojde-companion.linux-$$(uname -m) || echo $(DST) )" cmd/gopojde-companion/main.go

release: release-daemon release-manager-web release-manager-web-github-pages release-manager-native release-companion-web release-companion-web-github-pages release-companion-native

install: release-daemon release-manager-native release-companion-native
	sudo install out/release/gopojde-daemon/gopojde-daemon.linux-$$(uname -m) /usr/local/bin/gopojde-daemon
	sudo install out/release/gopojde-manager-native/gopojde-manager.linux-$$(uname -m) /usr/local/bin/gopojde-manager
	sudo install out/release/gopojde-companion-native/gopojde-companion.linux-$$(uname -m) /usr/local/bin/gopojde-companion
	
dev:
	while [ -z "$$DAEMON_PID" ] || [ -n "$$(inotifywait -q -r -e modify pkg cmd web/manager web/companion)" ]; do\
		$(MAKE);\
		kill -9 $$DAEMON_PID 2>/dev/null 1>&2;\
		kill -9 $$MANAGER_PID 2>/dev/null 1>&2;\
		kill -9 $$COMPANION_PID 2>/dev/null 1>&2;\
		wait $$DAEMON_PID $$MANAGER_PID $$COMPANION_PID;\
		out/gopojde-daemon/gopojde-daemon & export DAEMON_PID="$$!";\
		BUILDER=true go run cmd/gopojde-manager/main.go --serve & export MANAGER_PID="$$!";\
		out/gopojde-companion-native/gopojde-companion & export COMPANION_PID="$$!";\
	done

clean:
	rm -rf out
	rm -rf pkg/api/proto/v1
	rm -rf pkg/web/*/assets

depend:
	# Setup CLIs
	GO111MODULE=on go get github.com/golang/protobuf/protoc-gen-go@latest
	GO111MODULE=on go get github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
	GO111MODULE=on go get github.com/golang/protobuf/protoc-gen-go@latest
	# Generate bindings
	go generate ./...