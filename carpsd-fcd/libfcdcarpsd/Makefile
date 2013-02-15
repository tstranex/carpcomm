LIBS=`pkg-config libusb-1.0 --libs`
CFLAGS=`pkg-config libusb-1.0 --cflags` -I. -Wall -std=gnu99

OBJECTS=hid-libusb.o fcd_pro.o fcd_pro_plus.o fcdcarpsd.o

libfcdcarpsd.a: $(OBJECTS)
	ar rcs libfcdcarpsd.a $(OBJECTS)

hid-libusb.o: hid-libusb.c
	gcc $(CFLAGS) -c -o hid-libusb.o hid-libusb.c

fcd_pro.o: fcdpro/fcd_pro.c
	gcc $(CFLAGS) -c -o fcd_pro.o fcdpro/fcd_pro.c

fcd_pro_plus.o: fcdproplus/fcd_pro_plus.c
	gcc $(CFLAGS) -c -o fcd_pro_plus.o fcdproplus/fcd_pro_plus.c

fcdcarpsd.o: fcdcarpsd.c
	gcc $(CFLAGS) -c -o fcdcarpsd.o fcdcarpsd.c