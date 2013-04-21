#!/usr/bin/python

# Author: Timothy Stranex <tstranex@carpcomm.com>
# Copyright 2013 Timothy Stranex

#!/usr/bin/python

import datetime
import sys

from carpcomm.db import recordfile
from carpcomm.pb import stream_pb2


input_contacts = sys.argv[1]
satellite_id = sys.argv[2]

r = recordfile.RecordReader(file(input_contacts))
p = stream_pb2.Contact()
while True:
    data = r.Read()
    if data is None:
        break
    p.ParseFromString(data)

    if p.satellite_id != satellite_id:
        continue

    frames = []
    for b in p.blob:
        if b.format == stream_pb2.Contact.Blob.FRAME:
            frames.append(b.inline_data)

    if not frames:
        continue

    t = datetime.datetime.utcfromtimestamp(p.start_timestamp).isoformat()
    for f in frames:
        print t, f.encode('hex')
