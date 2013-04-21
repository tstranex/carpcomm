#!/usr/bin/python

# Author: Timothy Stranex <tstranex@carpcomm.com>
# Copyright 2013 Timothy Stranex

#!/usr/bin/python

import boto
import base64

def Dump(domain, column, proto):
    conn = boto.connect_sdb()
    d = conn.get_domain(domain)

    rows = d.select('select * from `%s`' % domain)
    n = 0
    for item in rows:
        if column not in item:
            print 'invalid: ', item.name, item
            print
            continue

        p = proto()
        p.ParseFromString(base64.b64decode(item[column]))
        print p
        n += 1

    print 'total items:', n
