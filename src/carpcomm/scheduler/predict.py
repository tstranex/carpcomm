#!/usr/bin/python

# Author: Timothy Stranex <tstranex@carpcomm.com>
# Copyright 2013 Timothy Stranex

#!/usr/bin/python

"""
Example input:

1338799998.2630811
47.4
8.5
400.0
20.0
270.0
90.0
SWISSCUBE               
1 35932U 09051B   12110.66765508  .00000638  00000-0  15500-3 0  5172
2 35932  98.3348 213.8703 0006768 284.4795  75.6141 14.52927878136365
"""

import ephem
import math
import sys
import time
import datetime

class SatellitePass(object):
    def __init__(self,
                 start_time, end_time,
                 start_azimuth, end_azimuth,
                 max_altitude):
        self.start_time = start_time
        self.end_time = end_time
        self.start_azimuth = start_azimuth
        self.end_azimuth = end_azimuth
        self.max_altitude = max_altitude

    def Serialize(self):
        def Timestamp(t):
            return time.mktime(ephem.localtime(t).timetuple())

        return '%f %f %f %f %f' % (
            Timestamp(self.start_time),
            Timestamp(self.end_time),
            math.degrees(self.start_azimuth),
            math.degrees(self.end_azimuth),
            math.degrees(self.max_altitude))


def IsVisible(body, min_altitude, min_azimuth, max_azimuth):
    if float(body.alt) < min_altitude:
        return False

    assert float(body.az) >= 0.0
    assert float(body.az) < 2*math.pi

    if min_azimuth <= max_azimuth:
        return min_azimuth <= body.az and body.az < max_azimuth
    else:
        return body.az < max_azimuth or body.az >= min_azimuth


def FindVisibleTransition(body, t_lower_bound, t_upper_bound, obs, is_visible):
    """Run a binary search to find the time when the visibility changes."""

    def VisibleAt(t):
        obs.date = t
        body.compute(obs)
        return is_visible(body)

    v1 = VisibleAt(t_lower_bound)
    v2 = VisibleAt(t_upper_bound)
    if v1 == v2:
        obs.date = t_lower_bound
        body.compute(obs)
        return t_lower_bound, body.az

    desired_resolution = 1.0 * ephem.second
    while t_upper_bound - t_lower_bound > desired_resolution:
        t = 0.5 * (t_upper_bound + t_lower_bound)
        v = VisibleAt(t)
        if v == v1:
            t_lower_bound = t
        else:
            t_upper_bound = t

    t = ephem.Date(t_upper_bound)
    obs.date = t
    body.compute(obs)
    return t, body.az


def FindMaximumAltitude(body, t_lower_bound, t_upper_bound, obs):
    def GetAltitude(t):
        obs.date = t
        body.compute(obs)
        return body.alt

    alt_lower = GetAltitude(t_lower_bound)
    alt_upper = GetAltitude(t_upper_bound)

    # Altitude is a unimodal function over this time interval.

    desired_resolution = 1.0 * ephem.second
    while t_upper_bound - t_lower_bound > desired_resolution:
        dt = t_upper_bound - t_lower_bound
        t1 = t_lower_bound + 0.3 * dt
        t2 = t_lower_bound + 0.7 * dt
        a1 = GetAltitude(t1)
        a2 = GetAltitude(t2)

        if a1 > a2:
            t_upper_bound = t2
        else:
            t_lower_bound = t1

    return GetAltitude(t_lower_bound)


def NextPasses(body, base, seconds, obs, is_visible):
    resolution = 60.0  # [s]

    passes = []
    start_time = None
    start_azimuth = None
    i = 0
    while i < seconds/resolution:
        tlast = ephem.Date(base + i*ephem.second*resolution)
        i += 1
        t = ephem.Date(base + i*ephem.second*resolution)
        obs.date = t
        body.compute(obs)

        visible = is_visible(body)

        if start_time is None:
            if not visible:
                continue
            else:
                # There was a visibility transition.
                start_time, start_azimuth = FindVisibleTransition(
                    body, tlast, t, obs, is_visible)
        else:
            if not visible:
                # There was a visibility transition.
                end_time, end_azimuth = FindVisibleTransition(
                    body, tlast, t, obs, is_visible)
                max_altitude = FindMaximumAltitude(
                    body, start_time, end_time, obs)

                passes.append(SatellitePass(
                        start_time, end_time,
                        start_azimuth, end_azimuth,
                        max_altitude))
                start_time = None
                start_azimuth = None
                max_altitude = 0.0
                continue

    # FIXME: extend the loop time to complete the final pass
    if start_time is not None:
        max_altitude = FindMaximumAltitude(
            body, start_time, t, obs)
        passes.append(SatellitePass(
                start_time, t,
                start_azimuth, body.az,
                max_altitude))

    return passes


def main():
    begin_timestamp = float(sys.stdin.readline())
    duration_seconds = float(sys.stdin.readline())
    latitude_degrees = float(sys.stdin.readline())
    longitude_degrees = float(sys.stdin.readline())
    elevation_metres = float(sys.stdin.readline())
    min_altitude_degrees = float(sys.stdin.readline())
    min_azimuth_degrees = float(sys.stdin.readline())
    max_azimuth_degrees = float(sys.stdin.readline())
    tle = [sys.stdin.readline(), sys.stdin.readline(), sys.stdin.readline()]

    obs = ephem.Observer()
    obs.lat = math.radians(latitude_degrees)
    obs.lon = math.radians(longitude_degrees)
    obs.elevation = elevation_metres

    body = ephem.readtle(*tle)

    def is_visible(body):
        return IsVisible(body,
                         math.radians(min_altitude_degrees),
                         math.radians(min_azimuth_degrees),
                         math.radians(max_azimuth_degrees))

    begin_time = ephem.Date(datetime.datetime.utcfromtimestamp(begin_timestamp))
    passes = NextPasses(body, begin_time, duration_seconds, obs, is_visible)
    print len(passes)
    for p in passes:
        print p.Serialize()


if __name__ == '__main__':
    main()
