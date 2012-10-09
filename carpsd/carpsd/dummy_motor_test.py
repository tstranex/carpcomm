#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

import config
import dummy_motor
import unittest

class DummyMotorTest(unittest.TestCase):

    def testStart(self):
        m = dummy_motor.DummyMotor(config.GetDefaultConfig())
        program = [[0.0, 10.0, 20.0],
                   [1.0, 20.0, 30.0]]
        self.assertEquals(True, m.Start(program))

    def testStop(self):
        m = dummy_motor.DummyMotor(config.GetDefaultConfig())
        self.assertEquals(True, m.Stop())

    def testGetStateDict(self):
        m = dummy_motor.DummyMotor(config.GetDefaultConfig())
        self.assertTrue(isinstance(m.GetStateDict(), dict))


if __name__ == '__main__':
    unittest.main()
