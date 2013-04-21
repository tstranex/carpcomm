#!/usr/bin/python

# Author: Timothy Stranex <tstranex@carpcomm.com>
# Copyright 2013 Timothy Stranex

import boto
import base64

conn = boto.connect_sdb()
d = conn.get_domain('test-stations')
item = d.get_item('1560574402094221678')
d.delete_item(item)
