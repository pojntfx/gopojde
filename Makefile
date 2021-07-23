all: build

backend:
	go build -o out/gopojde-backend/gopojde-backend cmd/gopojde-backend/main.go

frontend:
	rm -f web/app.wasm
	GOOS=js GOARCH=wasm go build -o web/app.wasm cmd/gopojde-frontend/main.go
	go build -o /tmp/gopojde-frontend-build cmd/gopojde-frontend/main.go
	rm -rf out/gopojde-frontend
	/tmp/gopojde-frontend-build -build
	cp -r web/* out/gopojde-frontend/web

build: backend frontend

release-backend:
	CGO_ENABLED=0 go build -ldflags="-extldflags=-static" -tags netgo -o out/release/gopojde-backend/gopojde-backend.linux-$$(uname -m) cmd/gopojde-backend/main.go

release-frontend: frontend
	rm -rf out/release/gopojde-frontend
	mkdir -p out/release/gopojde-frontend
	cd out/gopojde-frontend && tar -czvf ../release/gopojde-frontend/gopojde-frontend.tar.gz .

release-frontend-github-pages: frontend
	rm -rf out/release/gopojde-frontend-github-pages
	mkdir -p out/release/gopojde-frontend-github-pages
	/tmp/gopojde-frontend-build -build -path gopojde -out out/release/gopojde-frontend-github-pages
	cp -r web/* out/release/gopojde-frontend-github-pages/web

release: release-backend release-frontend release-frontend-github-pages

install: release-backend
	sudo install out/release/gopojde-backend/gopojde-backend.linux-$$(uname -m) /usr/local/bin/gopojde-backend
	
dev:
	while [ -z "$$BACKEND_PID" ] || [ -n "$$(inotifywait -q -r -e modify pkg cmd web/*.css)" ]; do\
		$(MAKE);\
		kill -9 $$BACKEND_PID 2>/dev/null 1>&2;\
		kill -9 $$FRONTEND_PID 2>/dev/null 1>&2;\
		wait $$BACKEND_PID $$FRONTEND_PID;\
		out/gopojde-backend/gopojde-backend & export BACKEND_PID="$$!";\
		/tmp/gopojde-frontend-build -serve & export FRONTEND_PID="$$!";\
	done

clean:
	rm -rf out
	rm -rf pkg/api/proto/v1
	rm -rf ~/.local/share/gopojde

depend:
	# Setup CLIs
	GO111MODULE=on go get github.com/golang/protobuf/protoc-gen-go@latest
	GO111MODULE=on go get github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
	GO111MODULE=on go get github.com/golang/protobuf/protoc-gen-go@latest
	# Generate bindings
	go generate ./...