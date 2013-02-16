#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

import fcd_receiver
import config

import unittest
import tempfile


class TestFCDReceiver(fcd_receiver.FCDReceiver):
    def _GetStatus(self):
        return {'Type': 'FUNcube Dongle Pro+',
                'Frequency [Hz]': '144000000'}


class FCDReceiverTest(unittest.TestCase):

    def testCreate(self):
        conf = config.GetDefaultConfig()
        section = fcd_receiver.FCDReceiver.__name__
        conf.add_section(section)
        conf.set(section, 'recording_dir', tempfile.mkdtemp())
        conf.set(section, 'alsa_device', 'hw:1')

        r = TestFCDReceiver(conf)
        self.assertEquals(r._sample_rate, 192000)
        self.assertEquals(r.GetHardwareTunerHz(), 144000000)
        

if __name__ == '__main__':
    unittest.main()
