#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

import motor
import unittest


class TestMotor(motor.Motor):
    def GetAzimuthLimitsDegrees(self):
        return (-45.0, 80.0)


class MotorTest(unittest.TestCase):
    def testAzimuthLimits(self):
        m = TestMotor()

        self.assertTrue(m.IsAllowedAzimuthDegrees(0.0))
        self.assertTrue(m.IsAllowedAzimuthDegrees(-20.0))
        self.assertTrue(m.IsAllowedAzimuthDegrees(60.0))

        self.assertFalse(m.IsAllowedAzimuthDegrees(-60.0))
        self.assertFalse(m.IsAllowedAzimuthDegrees(120.0))


if __name__ == '__main__':
    unittest.main()
