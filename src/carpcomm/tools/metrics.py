#!/usr/bin/python

# Author: Timothy Stranex <tstranex@carpcomm.com>
# Copyright 2013 Timothy Stranex

#!/usr/bin/python

import csv
import datetime
import sys
from carpcomm.db import recordfile
from carpcomm.pb import stream_pb2


def CountRecords(path):
    r = recordfile.RecordReader(file(path))
    n = 0
    while True:
        data = r.Read()
        if data is None:
            break
        n += 1
    return n

def CountContacts(path):
    num_contacts = 0
    num_contacts_with_frame = 0
    num_contacts_with_datum = 0
    num_frame_bytes = 0

    r = recordfile.RecordReader(file(path))
    p = stream_pb2.Contact()
    while True:
        data = r.Read()
        if data is None:
            break
        p.ParseFromString(data)
        num_contacts += 1

        has_frame = 0
        has_datum = 0

        for b in p.blob:
            if b.format == stream_pb2.Contact.Blob.FRAME:
                num_frame_bytes += len(b.inline_data)
                has_frame = 1
            elif b.format == stream_pb2.Contact.Blob.DATUM:
                has_datum = 1

        num_contacts_with_frame += has_frame
        num_contacts_with_datum += has_datum

    row = {}
    row['num_contacts'] = num_contacts
    row['num_contacts_with_frame'] = num_contacts_with_frame
    row['num_contacts_with_datum'] = num_contacts_with_datum
    row['num_frame_bytes'] = num_frame_bytes
    return row



backup_dir = sys.argv[1]
output_path = 'data/metrics.txt'

r = csv.reader(file(output_path))
fields = r.next()

row = {}
row['date'] = backup_dir.split('/')[-1].split('_')[0]
row['num_users'] = CountRecords(backup_dir + '/test-users.User.rec')
row['num_comments'] = CountRecords(backup_dir + '/test-comments.Comment.rec')
row['num_stations'] = CountRecords(backup_dir + '/test-stations.Station.rec')
row.update(CountContacts(backup_dir + '/test-contacts.Contact.rec'))

w = csv.DictWriter(file(output_path, 'a'), fields)
w.writerow(row)

for k in fields:
    print k, row[k]
