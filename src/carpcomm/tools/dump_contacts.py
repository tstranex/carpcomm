#!/usr/bin/python

# Author: Timothy Stranex <tstranex@carpcomm.com>
# Copyright 2013 Timothy Stranex

#!/usr/bin/python

from carpcomm.pb import stream_pb2
import dump

dump.Dump('test-contacts', 'pb.Contact', stream_pb2.Contact)
