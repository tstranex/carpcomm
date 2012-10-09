#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

"""Motor controller interface."""


class Motor(object):
    """Abstract interface for antenna motors."""

    def Start(self, program):
        """Activate the motor and follow the given rotation program.

        program is a list of tuples:
        (time_seconds, azimuth_degrees, elevation_degrees).
        They should be sorted in ascending order by time.
        The times are relative to when the method is called.

        azimuth_degrees >= 0 and is measured from clockwise from North.
        elevation_degrees >= 0 and is measured upward from the horizon.

        The azimuth and elevation must vary continuously over the pass.
        E.g. a north to south pass which passes directly overhead should have
        azimuth_degrees = 0 and elevation_degrees varying from 0 to 180.
        azimuth_degrees should not dicontinuously switch to 180 in this case.

        Returns True if successful.
        """

    def Stop(self):
        """Stop the motor.

        This will also halt any ongoing rotation program.

        Returns True if successful.
        """

    def GetStateDict(self):
        """Return a dictionary with info about the current state of the device.

        It contains the following keys:
        - azimuth_degrees (float)
        - elevation_degrees (float)
        - is_moving (bool)

        If azimuth_degrees or elevation_degrees are unknown, they will not be
        present in the dictionary.
        """

    def GetInfoDict(self):
        """Return a dictionary with info about the device and its capabilities.

        It contains at least the following keys:
        - driver
        """
        return {
            'driver': self.__class__.__name__,
            }
