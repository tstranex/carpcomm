#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

"""Dummy motor controller."""

import logging
import time

import motor

class DummyMotor(motor.Motor):
    """A dummy motor that isn't connected to actual hardware.

    It's useful for testing."""

    def __init__(self, config):
        pass

    def Start(self, program):
        if not program:
            return False
        return True

    def Stop(self):
        return True

    def GetStateDict(self):
        return {
            'azimuth_degrees': 0.0,
            'elevation_degrees': 0.0,
            'is_moving': False
            }


def Configure(config):
    if config.has_section(DummyMotor.__name__):
        return DummyMotor(config)
