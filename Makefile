
build:
	go build

run:
	go run *.go

files:
	@echo "package main" > autoGenStaticFiles.go
	@echo 'import "encoding/base64"' >> autoGenStaticFiles.go
	@echo 'var b64decode = func(encoded string) (raw []byte) {' >> autoGenStaticFiles.go
	@echo '	raw, _ = base64.StdEncoding.DecodeString(encoded)' >> autoGenStaticFiles.go
	@echo '	return' >> autoGenStaticFiles.go
	@echo '}' >> autoGenStaticFiles.go
	@echo 'var staticFiles = map[string][]byte{' >> autoGenStaticFiles.go
	@for F in `cd statics && find . -type f | sed 's/^\.\///' | grep -v "\.git"`; do \
		C=`base64 statics/$$F`; \
		echo "\"$$F\": b64decode(\`$$C\`)," >> autoGenStaticFiles.go ; \
	done
	@echo '}' >> autoGenStaticFiles.go

