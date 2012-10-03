#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

import config
import dummy_motor
import unittest

class DummyMotorTest(unittest.TestCase):
    def testNotReadyUntilReset(self):
        m = dummy_motor.DummyMotor(config.GetDefaultConfig())

        self.assertFalse(m.IsReady())

        self.assertFalse(m.IsOn())
        self.assertFalse(
            m.Reset(), 'The power needs to be on before it can be reset.')

        self.assertTrue(m.PowerOn())
        self.assertTrue(m.IsOn())

        self.assertTrue(m.Reset())
        m._FastForwardForTesting()
        self.assertTrue(m.IsReady())

    def testSetAzimuth(self):
        m = dummy_motor.DummyMotor(config.GetDefaultConfig())

        self.assertFalse(m.IsReady())
        self.assertEquals(None, m.GetAzimuthDegrees())
        self.assertFalse(m.SetAzimuthDegrees(12.0),
                         'The azimuth cannot be set until the motor is ready.')

        m.PowerOn()
        m.Reset()
        m._FastForwardForTesting()
        self.assertTrue(m.IsReady())

        self.assertTrue(m.SetAzimuthDegrees(12.0))
        self.assertTrue(m.IsMoving())
        m._FastForwardForTesting()
        self.assertEquals(12.0, m.GetAzimuthDegrees())
        self.assertFalse(m.IsMoving())

        m.PowerOff()
        self.assertFalse(m.SetAzimuthDegrees(13.0))

    def testHalt(self):
        m = dummy_motor.DummyMotor(config.GetDefaultConfig())
        m.PowerOn()
        m.Reset()
        m._FastForwardForTesting()

        self.assertEquals(0.0, m.GetAzimuthDegrees())
        self.assertTrue(m.SetAzimuthDegrees(12.0))
        self.assertTrue(m.Halt())
        self.assertFalse(m.IsMoving())
        self.assertEquals(0.0, m.GetAzimuthDegrees())


if __name__ == '__main__':
    unittest.main()
