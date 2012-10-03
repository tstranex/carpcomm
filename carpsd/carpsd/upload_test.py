#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

import upload
import config
import testing

import unittest


class APIClientTest(unittest.TestCase):

    def create(self):
        conf = testing.GetConfigForTesting()
        c = upload.APIClient(conf)
        c.SetServer('host.example.com', 123)
        return c

    def testGetServer(self):
        c = self.create()
        self.assertEquals(('host.example.com', 123), c.GetServer())

    def testPostPacket(self):
        c = self.create()

        method, path, body = c._GetPostPacketRequest(
            'test_satellite', 569, '\xab\x00\x5a')
        self.assertEquals(method, 'POST')
        self.assertEquals(path, '/PostPacket')
        self.assertEquals(
            body, 
            '{"format": "FRAME", '
            '"timestamp": 569, '
            '"station_id": "test_station_id", '
            '"station_secret": "test_station_secret", '
            '"satellite_id": "test_satellite", '
            '"frame_base64": "qwBa"}')


def liveTest():
    conf = testing.GetConfigForTesting()
    c = upload.APIClient(conf)
    c.SetServer('localhost', 5051)

    print c.PostPacket(
        'test_satellite', 569, '\xab\x00\x5a')


if __name__ == '__main__':
    #liveTest()
    unittest.main()
