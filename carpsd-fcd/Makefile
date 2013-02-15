LIBS=`pkg-config libusb-1.0 --libs`
CFLAGS=`pkg-config libusb-1.0 --cflags` -Ilibfcdcarpsd -Wall

carpsd-fcd: main.cpp libfcdcarpsd/libfcdcarpsd.a
	g++ $(CFLAGS) $(LIBS) -o carpsd-fcd main.cpp libfcdcarpsd/libfcdcarpsd.a

