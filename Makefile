run:
	go run cmd/main.go

test:
	go test ./... -v

build:
	go build -o tui_proxy_client cmd/main.go

build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o tui_proxy_client cmd/main.go

clean:
	rm -f tui_proxy_client
	rm -rf AppDir
	rm -f *.AppImage

appimage: build-linux
	mkdir -p AppDir/usr/bin
	mkdir -p AppDir/usr/share/applications
	cp tui_proxy_client AppDir/usr/bin/
	cp appimage/tui_proxy_client.desktop AppDir/usr/share/applications/
	cp appimage/tui_proxy_client.desktop AppDir/
	cp appimage/tui_proxy_client.png AppDir/ 2>/dev/null || true
	cp appimage/AppRun AppDir/
	chmod +x AppDir/AppRun
	chmod +x AppDir/usr/bin/tui_proxy_client
	./appimagetool AppDir tui_proxy_client-x86_64.AppImage