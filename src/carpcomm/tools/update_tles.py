#!/usr/bin/python

# Author: Timothy Stranex <tstranex@carpcomm.com>
# Copyright 2013 Timothy Stranex

#!/usr/bin/python

import sys
import urllib2

from google.protobuf import text_format
from carpcomm.pb import satellite_pb2


urls = [
    'http://mstl.atl.calpoly.edu/~ops/keps/kepler.txt',
    'http://www.celestrak.com/NORAD/elements/cubesat.txt',
    'http://www.celestrak.com/NORAD/elements/amateur.txt',
    'http://www.celestrak.com/NORAD/elements/noaa.txt',
    'http://www.celestrak.com/NORAD/elements/engineering.txt',
    'http://www.celestrak.com/NORAD/elements/tle-new.txt',
    'http://www.celestrak.com/NORAD/elements/stations.txt',
    ]
#    'http://www.amsat.org/amsat/ftp/keps/current/nasa.all'


def LoadTLEs(tle_file, label_to_sat):
    ids_updated = set()
    lines = tle_file.readlines()
    for i in range(0, len(lines), 3):
        label = lines[i].strip()
        if label not in label_to_sat:
            continue
        sat = label_to_sat[label]
        sat.tle = '\n'.join([lines[i].strip(),
                             lines[i+1].strip(),
                             lines[i+2].strip()])
        ids_updated.add(sat.id)
    return ids_updated


# TODO: This is only needed for the Stanford ground station. We can remove it
# once CarpSD supports rotators.
def OutputTLEs(sat_list, path):
    f = file(path, 'w')
    for sat in sat_list.satellite:
        tle = sat.tle
        if not tle:
            continue
        if sat.disable_tracking:
            continue
        print >>f, sat.tle


def main(in_path, out_path, out_tle_path):
    sat_list = satellite_pb2.SatelliteList()
    text_format.Merge(open(in_path).read(), sat_list)

    label_to_sat = {}
    all_ids = set()
    for sat in sat_list.satellite:
        all_ids.add(sat.id)
        if not sat.celestrak_tle_label:
            print >>sys.stderr, 'Warning: missing tle label:', sat.id
        else:
            label_to_sat[sat.celestrak_tle_label] = sat

    ids_updated = set()
    for url in urls:
        print >>sys.stderr, 'Loading', url
        f = urllib2.urlopen(url)
        ids_updated.update(LoadTLEs(f, label_to_sat))

    file(out_path, 'w').write(sat_list.SerializeToString())

    ids_not_updated = all_ids.difference(ids_updated)
    for sat_id in ids_not_updated:
        print >>sys.stderr, 'Warning: %s not updated' % sat_id

    print >>sys.stderr, 'Num updated:', len(ids_updated)
    print >>sys.stderr, 'Num not updated:', len(ids_not_updated)

    OutputTLEs(sat_list, out_tle_path)


if __name__ == '__main__':
    main(sys.argv[1], sys.argv[2], sys.argv[3])
