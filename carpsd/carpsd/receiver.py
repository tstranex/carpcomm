#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

class Receiver(object):
    """Interface to a software-defined radio."""

    def SetHardwareTunerHz(self, freq_hz):
        """Set the center frequency of the hardware tuner.

        Returns True on success.
        """
        raise NotImplementedError()

    def GetHardwareTunerHz(self):
        """Returns the center frequency of the hardware tuner."""
        raise NotImplementedError()

    def WaterfallImage(self):
        """Returns a spectrogram Image of the latest data.

        Returns None if there is no new data."""
        raise NotImplementedError()

    def Start(self, stream_url):
        """Start receiving data from the radio.

        Returns True on success."""
        raise NotImplementedError()

    def Stop(self):
        """Stop receiving data."""
        raise NotImplementedError()

    def IsStarted(self):
        """Returns true if we currently receiving data from the radio."""
        raise NotImplementedError()

    def GetStateDict(self):
        """Return a dictionary with info about the current state of the device.

        It contains at least the following keys:
        - hardware_tuner_hz
        - started
        """
        return {
            'hardware_tuner_hz': self.GetHardwareTunerHz(),
            'started': self.IsStarted(),
            }

    def GetInfoDict(self):
        """Return a dictionary with info about the device and its capabilities.

        It contains at least the following keys:
        - driver
        """
        return {
            'driver': self.__class__.__name__,
            }

