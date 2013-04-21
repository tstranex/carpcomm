#!/usr/bin/python

# Author: Timothy Stranex <tstranex@carpcomm.com>
# Copyright 2013 Timothy Stranex

#!/usr/bin/python

from carpcomm.pb import station_pb2
import dump

dump.Dump('r1-stations', 'pb.Station', station_pb2.Station)
