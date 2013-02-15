LIBS=`pkg-config libusb-1.0 --libs`
CFLAGS=`pkg-config libusb-1.0 --cflags` -I. -Wall

OBJECTS=hid-libusb.o fcd.o

libfcdcarpsd.a: $(OBJECTS)
	ar rcs libfcdcarpsd.a $(OBJECTS)

hid-libusb.o: hid-libusb.c
	gcc $(CFLAGS) -c -o hid-libusb.o hid-libusb.c

fcd.o: fcdpro/fcd.c
	gcc $(CFLAGS) -c -o fcd.o fcdpro/fcd.c

