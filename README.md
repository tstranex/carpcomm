carpcomm
========

Carpcomm Open Source projects

http://carpcomm.com/


Directories
-----------

api: Go API client
carpsd-fcd: FCD control utility for carpsd
carpsd: CarpSD ground station control software
telemetry: Telemetry decoder framework
src: Server-side carpcomm system (fe, streamer, muxd, schedd, etc)


Installation
------------

git clone git://github.com/tstranex/carpcomm.git
cd carpcomm
export GOPATH=`pwd`
go get launchpad.net/goamz
go get code.google.com/p/goauth2
go get code.google.com/p/goprotobuf
