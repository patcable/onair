TARGETS = onair
VERSION = $(shell bash ./genVersion.sh)
RAWVERSION = $(shell bash ./)

onair: hue.go onair.go watcher.go go.mod go.sum
	GOOS=darwin GOARCH=amd64 go build -ldflags='-extldflags "-sectcreate __TEXT __info_plist $(shell pwd)/Info.plist" -X "main.buildVersion=$(VERSION)"' -o $@

clean: 
	rm -rf $(TARGETS) _CodeSignature *.pkg package

.PHONY: sign
sign:
	mkdir package
	cp onair package
	codesign -s B6AE7396AD644B78D62F4B970E92E661A7D97B44 -f -v --timestamp --options runtime ./package/onair
	pkgbuild --root package --identifier net.pcable.onair --version $(VERSION:v%=%) --install-location /usr/local/bin onair-$(VERSION)-raw.pkg
	productsign -s 49AE1F5CCA9DAA99F529F88AB06244821D3017AE --timestamp onair-$(VERSION)-raw.pkg onair-$(VERSION).pkg
	xcrun altool --notarize-app --primary-bundle-id net.pcable.onair --username 'pc@pcable.net' --password "@keychain:appleid" --file onair-$(VERSION).pkg

.PHONY: staple
staple:
	xcrun stapler staple onair-$(VERSION).pkg
