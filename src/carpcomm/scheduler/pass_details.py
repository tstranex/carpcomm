#!/usr/bin/python

# Author: Timothy Stranex <tstranex@carpcomm.com>
# Copyright 2013 Timothy Stranex

#!/usr/bin/python

"""
Example input:

1349814979.0
120.0
47.4
8.5
400.0
1998-067CQ
1 38854U 98067CP  12283.07336473  .00054398  00000-0  89879-3 0    14
2 38854  51.6473 275.4897 0014651 145.1372 215.0744 15.51582271   671
1.0
"""

import ephem
import math
import sys
import datetime

def PassDetails(begin_timestamp, duration_seconds,
                latitude_degrees, longitude_degrees, elevation_metres,
                tle, resolution_seconds):
    obs = ephem.Observer()
    obs.lat = math.radians(latitude_degrees)
    obs.lon = math.radians(longitude_degrees)
    obs.elevation = elevation_metres

    body = ephem.readtle(*tle)

    begin_time = ephem.Date(datetime.datetime.utcfromtimestamp(begin_timestamp))

    n = int(duration_seconds/resolution_seconds)
    print n
    for i in xrange(n):
        t_offset = i*resolution_seconds
        timestamp = begin_timestamp + t_offset
        obs.date = ephem.Date(begin_time + ephem.second*t_offset)
        body.compute(obs)

        print '%f %f %f %f %f %f %f %f %d' % (
            timestamp,
            math.degrees(body.az),
            math.degrees(body.alt),
            float(body.range), # [m]
            float(body.range_velocity), # [m/s]
            math.degrees(body.sublat),
            math.degrees(body.sublong),
            float(body.elevation), # [m]
            bool(body.eclipsed))


def main():
    PassDetails(float(sys.stdin.readline()),
                float(sys.stdin.readline()),
                float(sys.stdin.readline()),
                float(sys.stdin.readline()),
                float(sys.stdin.readline()),
                [sys.stdin.readline(),
                 sys.stdin.readline(),
                 sys.stdin.readline()],
                float(sys.stdin.readline()))
                

if __name__ == '__main__':
    main()
