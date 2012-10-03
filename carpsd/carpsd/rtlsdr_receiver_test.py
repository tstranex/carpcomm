#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

import rtlsdr_receiver
import config
import tempfile

import unittest
import os


class RTLSDRReceiverTest(unittest.TestCase):

    def Create(self):
        conf = config.GetDefaultConfig()
        section = rtlsdr_receiver.RTLSDRReceiver.__name__
        conf.add_section(section)
        conf.set(section, 'recording_dir', tempfile.mkdtemp())
        conf.set(section, 'device_index', '0')
        conf.set(section, 'tuner_gain_db', '0.1')
        conf.set(section, 'sample_rate_hz', '96000')
        return rtlsdr_receiver.RTLSDRReceiver(conf)

    def testSetHardwareTunerHz(self):
        r = self.Create()
        self.assertTrue(r.SetHardwareTunerHz(1022434))
        self.assertEquals(1022434, r.GetHardwareTunerHz())


if __name__ == '__main__':
    unittest.main()
