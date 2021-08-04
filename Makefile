all: build

backend:
	go build -o out/gopojde-backend/gopojde-backend cmd/gopojde-backend/main.go

companion:
	go build -o out/gopojde-companion/gopojde-companion cmd/gopojde-companion/main.go

frontend:
	rm -f web/app.wasm
	GOOS=js GOARCH=wasm go build -o web/app.wasm cmd/gopojde-frontend/main.go
	go build -o /tmp/gopojde-frontend-build cmd/gopojde-frontend/main.go
	rm -rf out/gopojde-frontend
	/tmp/gopojde-frontend-build -build
	cp -r web/* out/gopojde-frontend/web

build: backend companion frontend

release-backend:
	go build -a -ldflags '-extldflags "-static"' -o "$(shell [ "$(DST)" = "" ] && echo out/release/gopojde-backend/gopojde-backend.linux-$$(uname -m) || echo $(DST) )" cmd/gopojde-backend/main.go

release-companion:
	go build -a -ldflags '-extldflags "-static"' -o "$(shell [ "$(DST)" = "" ] && echo out/release/gopojde-companion/gopojde-companion.linux-$$(uname -m) || echo $(DST) )" cmd/gopojde-companion/main.go

release-frontend: frontend
	rm -rf out/release/gopojde-frontend
	mkdir -p out/release/gopojde-frontend
	cd out/gopojde-frontend && tar -czvf ../release/gopojde-frontend/gopojde-frontend.tar.gz .

release-frontend-github-pages: frontend
	rm -rf out/release/gopojde-frontend-github-pages
	mkdir -p out/release/gopojde-frontend-github-pages
	/tmp/gopojde-frontend-build -build -path gopojde -out out/release/gopojde-frontend-github-pages
	mkdir -p out/release/gopojde-frontend-github-pages

release: release-backend release-companion release-frontend release-frontend-github-pages

install: release-backend release-companion
	sudo install out/release/gopojde-backend/gopojde-backend.linux-$$(uname -m) /usr/local/bin/gopojde-backend
	sudo setcap cap_net_bind_service+ep /usr/local/bin/gopojde-backend
	sudo install out/release/gopojde-companion/gopojde-companion.linux-$$(uname -m) /usr/local/bin/gopojde-companion
	
dev:
	while [ -z "$$BACKEND_PID" ] || [ -n "$$(inotifywait -q -r -e modify pkg cmd web/*.css)" ]; do\
		$(MAKE);\
		kill -9 $$BACKEND_PID 2>/dev/null 1>&2;\
		kill -9 $$COMPANION_PID 2>/dev/null 1>&2;\
		kill -9 $$FRONTEND_PID 2>/dev/null 1>&2;\
		wait $$BACKEND_PID $$COMPANION_PID $$FRONTEND_PID;\
		sudo setcap cap_net_bind_service+ep out/gopojde-backend/gopojde-backend;\
		out/gopojde-backend/gopojde-backend & export BACKEND_PID="$$!";\
		/tmp/gopojde-frontend-build -serve & export FRONTEND_PID="$$!";\
		out/gopojde-companion/gopojde-companion & export COMPANION_PID="$$!";\
	done

clean:
	rm -rf out
	rm -rf pkg/api/proto/v1

depend:
	# Setup CLIs
	GO111MODULE=on go get github.com/golang/protobuf/protoc-gen-go@latest
	GO111MODULE=on go get github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
	GO111MODULE=on go get github.com/golang/protobuf/protoc-gen-go@latest
	# Generate bindings
	go generate ./...