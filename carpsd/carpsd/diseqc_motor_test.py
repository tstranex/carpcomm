#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

import diseqc_motor
import unittest


class DiSEqCControllerTest(unittest.TestCase):
    def testAzimuthEncoding(self):
        # Based on Table 3 in the DiSEqC Positioner Application Note, page 12.
        encode = diseqc_motor.DiSEqCController._encode_azimuth_degrees

        self.assertEquals(encode(0.0), (0x0, 0x0))
        self.assertEquals(encode(90.0), (0x05, 0xa0))
        self.assertEquals(encode(180.0), (0x0b, 0x40))

        self.assertEquals(encode(-180.0), (0xf4, 0xc0))
        self.assertEquals(encode(-90.0), (0xfa, 0x60))

        self.assertEquals(encode(270.0), (0x10, 0xe0))
        self.assertEquals(encode(360.0), (0x16, 0x80))
        self.assertEquals(encode(450.0), (0x1c, 0x20))


if __name__ == '__main__':
    unittest.main()
