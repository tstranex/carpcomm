#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

import unittest
import time

import hamlib_motor
import config


class HamlibMotorTest(unittest.TestCase):

    def create(self):
        conf = config.GetDefaultConfig()
        section = hamlib_motor.HamlibMotor.__name__
        conf.add_section(section)
        conf.set(section, 'model', '1')
        conf.set(section, 'device', '/dev/null')
        conf.set(section, 'hamlib_param_min_az', '10.0')
        conf.set(section, 'hamlib_param_max_az', '80.0')

        commands = []
        def check_output(args):
            commands.append(args)
            return '124.0\n64.0\n'

        m = hamlib_motor.HamlibMotor(conf)
        m.rotator._check_output = check_output
        return m, commands

    def expectCommands(self, actual, expected):
        base = ['rotctl',
                '--model=1',
                '--rot-file=/dev/null',
                '--set-conf=max_az=80.0,min_az=10.0']
        self.assertEquals(len(actual), len(expected))
        for a, e in zip(actual, expected):
            self.assertEquals(a, base + [e])        

    def testAlreadyStopped(self):
        m, commands = self.create()
        self.assertEquals(True, m.Stop())
        self.expectCommands(commands, ['S'])

    def testGetStateDict(self):
        m, commands = self.create()
        self.assertEquals(
            {'azimuth_degrees': 124.0,
             'elevation_degrees': 64.0},
            m.GetStateDict())
        self.expectCommands(commands, ['p'])

    def testStartEmpty(self):
        m, commands = self.create()
        self.assertEquals(False, m.Start([]))

    def testStart(self):
        m, commands = self.create()
        program = [(0.0, 0.0, 30.0),
                   (0.2, 20.0, 30.0),
                   (0.4, 40.0, 20.0),
                   (0.6, 60.0, 15.0),
                   (0.8, 80.0, 10.0)]
        self.assertEquals(True, m.Start(program))

        self.expectCommands(commands, ['S'])

        # This is pretty fragile. :(
        time.sleep(1.0)

        self.expectCommands(
            commands,
            ['S',
             'P 20.000000 30.000000',
             'P 40.000000 20.000000',
             'P 60.000000 15.000000',
             'P 80.000000 10.000000',
             'S'])


if __name__ == '__main__':
    unittest.main()