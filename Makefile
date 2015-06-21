
VERSION=$(shell git describe --tags)
OPTS=-ldflags '-X main.VERSION $(VERSION)'

version:
	go run $(OPTS) idok.go -version

all: _prepare linux32 linux64 darwin freebsd64 freebsd32 windows pack

_prepare: clean
	mkdir dist
	sed 's/@VERSION@/$(VERSION)/' install-idok.sh > dist/install-idok.sh

darwin:
	GOOS=darwin go build $(OPTS) -o idok-darwin idok.go

linux32:
	GOARCH=386 go build $(OPTS) -o idok-i686 idok.go
	strip idok-i686

linux64:
	go build $(OPTS) -o idok-x86_64 idok.go
	strip idok-x86_64

freebsd32:
	GOOS=freebsd GOARCH=386 go build $(OPTS) -o idok-freebsd32 idok.go

freebsd64:
	GOOS=freebsd go build $(OPTS) -o idok-freebsd64 idok.go

windows:
	GOOS=windows go build $(OPTS) idok.go

pack:
	mv idok-* idok.exe dist
	cd dist && \
	gzip idok-i686 && \
	gzip idok-x86_64 && \
	gzip idok-darwin && \
	gzip idok-freebsd32 && \
	gzip idok-freebsd64 && \
	zip idok idok.exe && \
	rm idok.exe

clean:
	rm -rf dist

deploy:
	scp -C idok*.zip idok*.gz root@metal3d.org:/var/www/dists/
