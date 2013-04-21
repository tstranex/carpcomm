#!/usr/bin/python

# Author: Timothy Stranex <tstranex@carpcomm.com>
# Copyright 2013 Timothy Stranex

#!/usr/bin/python


from carpcomm.pb import user_pb2
import dump

dump.Dump('r1-users', 'pb.User', user_pb2.User)
