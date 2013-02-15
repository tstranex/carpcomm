LIBS=`pkg-config libusb-1.0 --libs`
CFLAGS=`pkg-config libusb-1.0 --cflags` -I. -Isrc -Isrc/hidapi -Wall -std=gnu99

OBJECTS=hid-libusb.o fcd_pro.o fcd_pro_plus.o fcdcarpsd.o

libfcdcarpsd.a: $(OBJECTS)
	ar rcs libfcdcarpsd.a $(OBJECTS)

hid-libusb.o: src/hidapi/hid-libusb.c
	gcc $(CFLAGS) -c -o hid-libusb.o src/hidapi/hid-libusb.c

fcd_pro.o: src/fcd_pro.c
	gcc $(CFLAGS) -c -o fcd_pro.o src/fcd_pro.c

fcd_pro_plus.o: src/fcd_pro_plus.c
	gcc $(CFLAGS) -c -o fcd_pro_plus.o src/fcd_pro_plus.c

fcdcarpsd.o: src/fcdcarpsd.c
	gcc $(CFLAGS) -c -o fcdcarpsd.o src/fcdcarpsd.c