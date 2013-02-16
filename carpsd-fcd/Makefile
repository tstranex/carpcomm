VERSION=0.01
LIBS=`pkg-config libusb-1.0 --libs`
CFLAGS=`pkg-config libusb-1.0 --cflags` -Ilibfcdcarpsd -Wall

carpsd-fcd: main.cpp libfcdcarpsd/libfcdcarpsd.a
	g++ $(CFLAGS) $(LIBS) -o carpsd-fcd main.cpp libfcdcarpsd/libfcdcarpsd.a

.PHONY : install
install: carpsd-fcd
	mkdir -p $(DESTDIR)/usr/bin/
	cp carpsd-fcd $(DESTDIR)/usr/bin/
	mkdir -p $(DESTDIR)/etc/udev/rules.d/
	cp etc/funcubedongle.rules $(DESTDIR)/etc/udev/rules.d/

PACKAGE=carpsd-fcd_$(VERSION)
.PHONY : deb
deb:
	rm -rf /tmp/$(PACKAGE)
	mkdir /tmp/$(PACKAGE)
	cp -R * /tmp/$(PACKAGE)/
	cd /tmp && tar -czf $(PACKAGE).orig.tar.gz $(PACKAGE)
	cd /tmp/$(PACKAGE) && debuild -us -uc