LIBS=`pkg-config libusb-1.0 --libs`
CFLAGS=`pkg-config libusb-1.0 --cflags` -I. -Wall

OBJECTS=hid-libusb.o fcdpro.o fcdproplus.o

libfcdcarpsd.a: $(OBJECTS)
	ar rcs libfcdcarpsd.a $(OBJECTS)

hid-libusb.o: hid-libusb.c
	gcc $(CFLAGS) -c -o hid-libusb.o hid-libusb.c

fcdpro.o: fcdpro/fcdpro.c
	gcc $(CFLAGS) -c -o fcdpro.o fcdpro/fcdpro.c

fcdproplus.o: fcdproplus/fcdproplus.c
	gcc $(CFLAGS) -c -o fcdproplus.o fcdproplus/fcdproplus.c
