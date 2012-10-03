#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

import receiver


class TestImage:
    def save(self, f, format):
        # PNG image that says "DummyReceiver".
        f.write('\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x01\x00\x00\x00\x00 \x01\x03\x00\x00\x00\xc7g4\xa7\x00\x00\x00\x01sRGB\x00\xae\xce\x1c\xe9\x00\x00\x00\x06PLTE\xff\xff\xff\x00\x00\x00U\xc2\xd3~\x00\x00\x00\xdeIDAT8\xcbc`\x18\x05\x14\x80\x04\x0c\x11\xfe\x07P\x06\xfb\x07\xec\nx*\xa0\x0c6\x0b\xecf\xf2\xf0\xc0\x14\xf0\xe0P\xc0bp\xb8\xf2\xc1\xbd\x8a\x03m<r\xcc\x0f>\xa4\xf0\xc8\xf1\x1e@U\xc0fp\xf8v\xc1\xb7;\x8em<\xc5|\xc93\xd2x\x8a\xf9\x190\x14\xf0\x19\x1c\xeeIl\x93p\xec`\xe2I\x93p\xec\xc1T\xc0#p\x98%\x91\xfd\xc3\xc1\x19l<\xc9\x12\x07q(0f\xb3\xf8?\x83\xfdO\xb2\xc0\x7f\x14\x05\xcc\x0f\x10\n$\x0eH\xb01$\x1a\x1c\xe0a\xb0C\xf7\x05H\x81\x18\x1b\x0fX\xc1\x86\x03<\xe8\xe1\x00Q \xc7\xc6\xe3\xc0\xc1\xc4\x90\xb8\xc0\x01]\x81\x85$X\x81\x0c\x1b[q_\xf2\x9c\x84\x07\xc5=\xe8q!\x07V\xc0\xc3\xc6&\xcf\xfc\xe0G\xc2\x03y\\!J(\x9e9\x08)\x90 \xa4\xa0`h&p\x00\xf9MA\xa8\xe4\xcf*\xdc\x00\x00\x00\x00IEND\xaeB`\x82')


class DummyReceiver(receiver.Receiver):
    def __init__(self, config):
        self._stream_url = None
        self._freq_hz = None
        self._started = False

    def SetHardwareTunerHz(self, freq_hz):
        self._freq_hz = freq_hz
        return True

    def GetHardwareTunerHz(self):
        return self._freq_hz

    def WaterfallImage(self):
        if not self.IsStarted():
            return None
        return TestImage()

    def Start(self, stream_url):
        self._started = True
        self._stream_url = stream_url
        return True

    def Stop(self):
        if self.IsStarted() and self._stream_url:
            # This is where we'd upload data if we had any real data.
            pass
        self._started = False

    def IsStarted(self):
        return self._started


def Configure(config):
    if config.has_section(DummyReceiver.__name__):
        return DummyReceiver(config)
