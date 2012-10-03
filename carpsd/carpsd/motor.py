#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

"""Motor controller interface."""


class Motor(object):
    """Abstract interface for antenna motors."""

    def PowerOn(self):
        """Switch on the motor.

        Returns True if successful.
        """
        raise NotImplementedError()

    def PowerOff(self):
        """Switch off the motor.

        Returns True if successful.
        """
        raise NotImplementedError()

    def IsOn(self):
        """Returns True if the motor is switched on."""
        raise NotImplementedError()

    def Reset(self):
        """Move the motor to a known reference position.

        Returns True on success.
        """
        raise NotImplementedError()

    def IsReady(self):
        """Returns true if the position is known and the motor is ready to
        accept commands.
        """
        raise NotImplementedError()

    def IsMoving(self):
        """Returns true if the motor is currently moving."""
        raise NotImplementedError()

    def Halt(self):
        """Cancel any underway commands and bring the motor to a halt.

        Returns True on success.
        """
        raise NotImplementedError()

    def GetAzimuthDegrees(self):
        """Return the current azimuth or None if it's unknown.

        The angle is measured in degrees clockwize from true north.
        """
        raise NotImplementedError()

    def GetAzimuthLimitsDegrees(self):
        """Return a tuple (min, max) of the possible azimuth rotation.

        The angles are measured in degrees clockwize from true north.
        """
        raise NotImplementedError()

    def IsAllowedAzimuthDegrees(self, azimuth_degrees):
        """Return true if azimuth_degrees is within the limits returned by
        GetAzimuthDegreesLimits.
        """
        min_limit, max_limit = self.GetAzimuthLimitsDegrees()
        return min_limit <= azimuth_degrees and azimuth_degrees <= max_limit

    def SetAzimuthDegrees(self, azimuth_degrees):
        """Command the motor to move to the given azimuth.

        The angle is measured in degrees clockwize from true north.

        Returns true if successful.
        """
        raise NotImplementedError()

    def GetStateDict(self):
        """Return a dictionary with info about the current state of the device.

        It contains at least the following keys:
        - power
        - ready
        - azimuth_degrees
        - moving
        """
        return  {
            'power': self.IsOn(),
            'ready': self.IsReady(),
            'azimuth_degrees': self.GetAzimuthDegrees(),
            'moving': self.IsMoving(),
            }

    def GetInfoDict(self):
        """Return a dictionary with info about the device and its capabilities.

        It contains at least the following keys:
        - driver
        - azimuth_limits
        """
        return {
            'driver': self.__class__.__name__,
            'azimuth_limits': self.GetAzimuthLimitsDegrees(),
            }
