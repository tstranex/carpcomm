#!/usr/bin/python

# Author: Timothy Stranex <tstranex@carpcomm.com>
# Copyright 2013 Timothy Stranex

from carpcomm.pb import ranking_pb2
from carpcomm.pb import satellite_pb2
from carpcomm.pb import station_pb2
from carpcomm.pb import user_pb2
from carpcomm.pb import comments_pb2
from carpcomm.pb import stream_pb2

protos = {
    'RankingList': ranking_pb2.RankingList,
    'SatelliteList': satellite_pb2.SatelliteList,
    'Station': station_pb2.Station,
    'User': user_pb2.User,
    'Comment': comments_pb2.Comment,
    'Contact': stream_pb2.Contact,
    }

def GetProtoByName(name):
    return protos.get(name, None)
