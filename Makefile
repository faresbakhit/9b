.POSIX:

GO=go
HURL=hurl
OUT=./9b

default: build

.PHONY: build
build:
	$(GO) build -tags "sqlite_foreign_keys" -o $(OUT)

.PHONY: test
.ONESHELL:
test: clean build
	$(OUT) &
	NINEB_PID=$$!
	sleep 5
	$(MAKE) hurl || true
	kill -2 "$${NINEB_PID}"

.PHONY: hurl
hurl: build
	$(HURL) --very-verbose --variables-file hurl/hurl.vars --test hurl/

.PHONY: clean
clean:
	rm -f $(OUT) ./9b.db
