#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

import unittest
import time

import diseqc_motor
import config


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


class DummyDiSEqCController(diseqc_motor.DiSEqCController):
    def __init__(self):
        self.commands = []

    def hard_power_off(self):
        self.commands.append('hard_power_off')

    def hard_power_on(self):
        self.commands.append('hard_power_on')

    def halt(self):
        self.commands.append('halt')

    def drive_east(self):
        self.commands.append('drive_east')

    def drive_west(self):
        self.commands.append('drive_west')

    def goto_stored_position(self, n):
        self.commands.append('goto_stored_position %d' % n)


class DiSEqCMotorTest(unittest.TestCase):

    def create(self):
        conf = config.GetDefaultConfig()
        section = diseqc_motor.DiSEqCMotor.__name__
        conf.add_section(section)
        conf.set(section, 'serial_device', '/dev/null')
        conf.set(section, 'calibrated_rate', '10.0')
        conf.set(section, 'reference_zero', '0.0')
        conf.set(section, 'min_limit', '-45.0')
        conf.set(section, 'max_limit', '45.0')

        c = DummyDiSEqCController()
        m = diseqc_motor.DiSEqCMotor(conf, controller=c)
        m._internal_motor.reset_time = 0.1
        return m, c

    def testAlreadyStopped(self):
        m, c = self.create()
        self.assertEquals(True, m.Stop())

    def testUnknownInfo(self):
        m, c = self.create()
        self.assertEquals({'is_moving': False}, m.GetStateDict())

    def testStart(self):
        m, c = self.create()
        program = [(0.0, 0.0, 30.0),
                   (0.2, 20.0, 30.0),
                   (0.4, 40.0, 30.0),
                   (0.6, 60.0, 30.0),
                   (0.8, 80.0, 30.0)]
        self.assertEquals(True, m.Start(program))

        # This is pretty fragile. :(
        time.sleep(1.0)

        self.assertEquals(
            c.commands,
            ['hard_power_on',
             'goto_stored_position 0',
             'halt',
             'halt',
             'drive_east',
             'halt',
             'drive_east',
             'halt',
             'hard_power_off'])

    def testEmptyProgram(self):
        m, c = self.create()
        self.assertEquals(False, m.Start([]))


if __name__ == '__main__':
    unittest.main()
