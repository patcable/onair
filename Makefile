TARGETS = onair
VERSION = $(shell bash ./genVersion.sh)

onair: hue.go onair.go watcher.go go.mod go.sum
	GOOS=darwin GOARCH=amd64 go build -ldflags='-extldflags "-sectcreate __TEXT __info_plist $(shell pwd)/Info.plist" -X "main.buildVersion=$(VERSION)"' -o $@

clean: 
	rm $(TARGETS)

.PHONY: sign
sign:
	gon ./gon.hcl
