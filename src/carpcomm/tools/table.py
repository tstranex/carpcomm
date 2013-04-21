#!/usr/bin/python

# Author: Timothy Stranex <tstranex@carpcomm.com>
# Copyright 2013 Timothy Stranex

#!/usr/bin/python

import boto
import struct
import sys

import protodb
from carpcomm.db import recordfile
from carpcomm.db import simpledb


def Dump(domain, column, writer):
    conn = boto.connect_sdb()
    d = conn.get_domain(domain)
    rows = d.select('select * from `%s`' % domain)
    n = 0
    max_size = 0
    for item in rows:
        n += 1
        data = simpledb._DecodeItem(item, column)
        writer.Write(data)
        max_size = max(max_size, len(data))
    print 'Dumped %d records.' % n
    print 'Max record size: %d bytes' % max_size


def Backup(domain, column, out_path):
    w = recordfile.RecordWriter(file(out_path, 'w'))
    Dump(domain, column, w)
    w.Close()


def Lookup(domain, column, name, proto_name):
    conn = boto.connect_sdb()
    d = conn.get_domain(domain)
    item = d.get_item(name)
    data = simpledb._DecodeItem(item, column)
    p = protodb.GetProtoByName(proto_name)()
    p.ParseFromString(data)
    print p


def Verify(in_path, proto_name):
    p = protodb.GetProtoByName(proto_name)()
    r = recordfile.RecordReader(file(in_path))
    n = 0
    while True:
        data = r.Read()
        if data is None:
            break
        n += 1
        try:
            p.ParseFromString(data)
        except Exception, e:
            print n, 'error:', e
    print 'Read %d records.' % n


def main(argv):
    command = argv[1]
    if command == 'backup':
        Backup(*argv[2:])
    elif command == 'verify':
        Verify(*argv[2:])
    elif command == 'lookup':
        Lookup(*argv[2:])
    else:
        print 'Unknown command:', command


if __name__ == '__main__':
    main(sys.argv)
