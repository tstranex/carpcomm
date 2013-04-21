#!/usr/bin/python

# Author: Timothy Stranex <tstranex@carpcomm.com>
# Copyright 2013 Timothy Stranex

#!/usr/bin/python

"""Amazon SimpleDB utils."""

import base64

def _DecodeItem(item, column):
    column_count = item.get(column + ".v2", None)
    if column_count is None:
        # v1 format
        encoded_data = item[column]
    else:
        # v2 format
        count = int(column_count)
        encoded_data = ''
        for i in range(count):
            encoded_data += item[column + "." + str(i)]
    return base64.b64decode(encoded_data)
