#!/usr/bin/python

# Author: Timothy Stranex <tstranex@carpcomm.com>
# Copyright 2013 Timothy Stranex

#!/usr/bin/python

"""Simple binary record file format."""

import struct

RECORD_WRITER_V0_HEADER = 'RecordWriter000'


class RecordWriter(object):
    def __init__(self, f):
        self.f = f
        self.f.write(RECORD_WRITER_V0_HEADER)

    def Write(self, record):
        self.f.write(struct.pack('!L', len(record)))
        self.f.write(record)

    def Close(self):
        self.f.close()


class RecordReader(object):
    def __init__(self, f):
        self.f = f
        h = f.read(len(RECORD_WRITER_V0_HEADER))
        assert h == RECORD_WRITER_V0_HEADER

    def Read(self):
        s = self.f.read(4)
        if not s:
            return None
        size = struct.unpack('!L', s)[0]
        return self.f.read(size)
