#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

import api
import config
import testing

import datetime
import unittest


class APIClientTest(unittest.TestCase):

    def create(self):
        conf = testing.GetConfigForTesting()
        c = api.APIClient(conf)
        c.SetServer('host.example.com', 123)
        return c

    def testGetServer(self):
        c = self.create()
        self.assertEquals(('host.example.com', 123), c.GetServer())

    def testPostPacket(self):
        c = self.create()

        def SendRequest(method, path, body):
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

            return True, (200, ''), ''

        c._SendRequest = SendRequest

        ok, status = c.PostPacket('test_satellite', 569, '\xab\x00\x5a')
        self.assertEquals(True, ok)

    def testGetLatestPackets(self):
        c = self.create()

        def SendRequest(method, path, body):
            self.assertEquals(method, 'GET')
            self.assertEquals(
                path,
                '/GetLatestPackets?'
                'station_secret=test_station_secret&'
                'limit=100&'
                'station_id=test_station_id&'
                'satellite_id=test')
            self.assertEquals(body, '')
            return True, (200, ''), (
                '[{"timestamp":1351737353,'
                '"frame_base64":"6C8bpEZ9X9zVTmbtY7fJ4kC+cwBobPwd"},'
                '{"timestamp":1351737340,'
                '"frame_base64":"d3XsBnDyfmep7xC3B1QV3o6X"}]')

        c._SendRequest = SendRequest

        ok, status, packets = c.GetLatestPackets('test', 100)
        self.assertEquals(True, ok)
        self.assertEquals(
            packets,
            [(datetime.datetime(2012, 11, 1, 3, 35, 53),
              '\xe8/\x1b\xa4F}_\xdc\xd5Nf\xedc\xb7\xc9\xe2@'
              '\xbes\x00hl\xfc\x1d'),
             (datetime.datetime(2012, 11, 1, 3, 35, 40),
              'wu\xec\x06p\xf2~g\xa9\xef\x10\xb7\x07T\x15\xde\x8e\x97')])


def liveTest():
    conf = testing.GetConfigForTesting()
    c = api.APIClient(conf)
    c.SetServer('api.carpcomm.com', 5051)

    print c.PostPacket(
        'test_satellite', 569, '\xab\x00\x5a')


if __name__ == '__main__':
    #liveTest()
    unittest.main()
