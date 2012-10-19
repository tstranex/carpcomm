#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

import usrp_receiver
import config
import tempfile

import unittest


class USRPReceiverTest(unittest.TestCase):

    def Create(self):
        conf = config.GetDefaultConfig()
        section = usrp_receiver.USRPReceiver.__name__
        conf.add_section(section)
        conf.set(section, 'recording_dir', tempfile.mkdtemp())
        conf.set(section, 'device_address', 'test')
        conf.set(section, 'sample_rate_hz', '32000')
        return usrp_receiver.USRPReceiver(conf)

    def testSetHardwareTunerHz(self):
        r = self.Create()
        self.assertTrue(r.SetHardwareTunerHz(1022434))
        self.assertEquals(1022434, r.GetHardwareTunerHz())


if __name__ == '__main__':
    unittest.main()
