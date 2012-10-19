#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

import fcd_receiver
import config

import unittest
import tempfile


class FCDReceiverTest(unittest.TestCase):

    def testCreate(self):
        conf = config.GetDefaultConfig()
        section = fcd_receiver.FCDReceiver.__name__
        conf.add_section(section)
        conf.set(section, 'recording_dir', tempfile.mkdtemp())
        conf.set(section, 'alsa_device', 'hw:1')
        conf.set(section, 'frequency_correction', '-12')
        return fcd_receiver.FCDReceiver(conf)


if __name__ == '__main__':
    unittest.main()
