VERSION=1.0.0
BUILD_TIME=`date +%FT%T%z`

LDFLAGS=-ldflags "-X github.com/alexgear/suplat/config.Version=${VERSION} -X github.com/alexgear/suplat/config.BuildTime=${BUILD_TIME}"

mac:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build ${LDFLAGS}
linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ${LDFLAGS}
deploy:
	ssh -t admin@battleship 'sudo systemctl stop suplat'
	scp suplat admin@battleship:/opt/suplat/
	scp config.toml admin@battleship:/opt/suplat/
	ssh -t admin@battleship 'sudo systemctl start suplat'
