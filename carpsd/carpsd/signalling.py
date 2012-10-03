#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

"""Interface for signalling the user about events."""

import logging


class Signaller(object):

    def SignalIdentified(self):
        logging.info('Station Identified')

    def SignalPing(self):
        logging.info('Ping')

    def SignalReceiverStart(self):
        logging.info('Receiver Start')

    def SignalReceiverStop(self):
        logging.info('Recevier Stop')

    def SignalUploadStart(self):
        logging.info('Upload Start')

    def SignalUploadStop(self):
        logging.info('Upload Stop')


class LEDSignaller(Signaller):

    def _Signal(self, data):
        f = file('/tmp/led_signal', 'w')
        f.write(data)
        f.close()

    def SignalIdentified(self):
        Signaller.SignalIdentified(self)
        self._Signal('-')

    def SignalPing(self):
        Signaller.SignalPing(self)
        self._Signal('-')

    def SignalReceiverStart(self):
        Signaller.SignalReceiverStart(self)
        self._Signal('[')

    def SignalReceiverStop(self):
        Signaller.SignalReceiverStop(self)
        self._Signal(']')

    def SignalUploadStart(self):
        Signaller.SignalUploadStart(self)
        self._Signal('(')

    def SignalUploadStop(self):
        Signaller.SignalUploadStop(self)
        self._Signal(')')



_global_signaller = Signaller()

def Get():
    return _global_signaller


def Configure(config):
    global _global_signaller
    if config.has_section(LEDSignaller.__name__):
        _global_signaller = LEDSignaller()
        logging.info('Installed LEDSignaller')
