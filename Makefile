9b:
	go build -tags "sqlite_foreign_keys" -o ./9b

.PHONY: watch
watch:
	air --build.cmd "go build -tags \"sqlite_foreign_keys\" -o ./9b" \
	    --build.bin "./9b" \
	    --build.exclude_dir "templates,public,art" \
