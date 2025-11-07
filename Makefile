TARGETS = onair
VERSION = $(shell bash ./genVersion.sh)
RAWVERSION = $(shell bash ./)

onair: ifttt.go hue.go onair.go watcher.go go.mod go.sum
	go build -ldflags='-extldflags "-sectcreate __TEXT __info_plist $(shell pwd)/Info.plist" -X "main.buildVersion=$(VERSION)"' -o $@

clean: 
	rm -rf $(TARGETS) _CodeSignature *.pkg package

.PHONY: sign
sign:
	mkdir -p package
	cp onair package
	codesign -s 'Developer ID Application: Big Technology LLC (FMGF9BLA5F)' -f -v --timestamp --options runtime ./package/onair
	pkgbuild --root package --identifier net.pcable.onair --version $(VERSION:v%=%) --install-location /usr/local/bin onair-$(VERSION)-raw.pkg
	productsign -s 'Developer ID Installer: Big Technology LLC (FMGF9BLA5F)' --timestamp onair-$(VERSION)-raw.pkg onair-$(VERSION).pkg
	xcrun notarytool submit --keychain-profile "BigTech ASC" onair-$(VERSION).pkg

.PHONY: staple
staple:
	xcrun stapler staple onair-$(VERSION).pkg
