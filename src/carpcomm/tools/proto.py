#!/usr/bin/python

# Author: Timothy Stranex <tstranex@carpcomm.com>
# Copyright 2013 Timothy Stranex

from google.protobuf import text_format

import protodb

import optparse
import sys

def GuessProto(args):
    for path in args:
        if path.endswith('.txt'):
            path = path[:-4]
        ext = path.split('.')[-1]
        p = protodb.GetProtoByName(ext)
        if p is not None:
            return p
    return None

p = GuessProto(sys.argv[1:])()

in_file = file(sys.argv[1])
if sys.argv[1].endswith('.txt'):
    text_format.Merge(in_file.read(), p)
else:
    p.ParseFromString(in_file.read())

out_file = sys.stdout
binary = False
if len(sys.argv) > 2:
    path = sys.argv[2]
    out_file = file(path, 'w')
    binary = not path.endswith('.txt')

if binary:
    out_file.write(p.SerializeToString())
else:
    out_file.write(str(p))
