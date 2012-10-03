#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

"""Dummy motor controller."""

import logging
import time

import motor

class DummyMotor(motor.Motor):
    """A dummy motor that isn't connected to actual hardware. It's useful for
    testing."""

    def __init__(self, config):
        logging.info('Creating DummyMotor')

        def get(name, cast, default):
            if not config.has_section(DummyMotor.__name__):
                return default
            if not config.has_option(DummyMotor.__name__, name):
                return default
            return cast(config.get(DummyMotor.__name__, name))

        self._limits = (
            get('min_azimuth_limit', float, -50.0),
            get('max_azimuth_limit', float, 45.0))

        self._azimuth = None
        self.PowerOff()

    def PowerOn(self):
        self._Update()
        self._power = True
        return True

    def PowerOff(self):
        self._power = False
        self._end_time = None
        self._command = None
        return True

    def IsOn(self):
        self._Update()
        return self._power

    def _Command(self, delay, command):
        if not self.IsOn():
            return False
        self._end_time = time.time() + delay
        self._command = command
        return True

    def _Update(self):
        t = time.time()
        if t > self._end_time:
            self._ForceUpdate()

    def _ForceUpdate(self):
        if self._end_time is None:
            return
        self._command()
        self._end_time = None
        self._command = None

    def _FastForwardForTesting(self):
        self._ForceUpdate()

    def Reset(self):
        self._Update()
        def NewState():
            self._azimuth = 0.0
        return self._Command(10.0,  NewState)

    def IsReady(self):
        self._Update()
        return self._azimuth is not None

    def IsMoving(self):
        self._Update()
        return self._end_time is not None

    def Halt(self):
        self._Update()
        self._end_time = None
        self._command = None
        return True

    def GetAzimuthDegrees(self):
        self._Update()
        return self._azimuth

    def GetAzimuthLimitsDegrees(self):
        self._Update()
        return self._limits

    def SetAzimuthDegrees(self, azimuth_degrees):
        self._Update()
        if not self.IsReady():
            return False

        def SetAzimuth():
            self._azimuth = azimuth_degrees
        return self._Command(10.0, SetAzimuth)


def Configure(config):
    if config.has_section(DummyMotor.__name__):
        return DummyMotor(config)
