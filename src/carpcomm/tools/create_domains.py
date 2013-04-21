#!/usr/bin/python

# Author: Timothy Stranex <tstranex@carpcomm.com>
# Copyright 2013 Timothy Stranex

#!/usr/bin/python

import boto

conn = boto.connect_sdb()

#conn.delete_domain('test-comments')
print conn.get_all_domains()

#conn.create_domain('test-users')
#conn.create_domain('test-stations')
#conn.create_domain('test-contacts')
conn.create_domain('test-comments')

print conn.get_all_domains()
