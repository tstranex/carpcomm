#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

"""USRP radio Receiver implementation."""

import subprocess
import logging
import os.path
import time

import receiver
import spectrum
import upload


USRP_RECORD_BINARY = ['python', '-m', 'carpsd.usrp_record']


class USRPReceiver(receiver.Receiver):
    """Receiver implementation for the USRP."""

    def __init__(self, config):
        self._stream_url = None
        self._output_path = None
        self._spectrum_path = None
        self._pipe = None
        self._upload_pipe = None
        self._freq_hz = 145500000
        self._spectrum = None

        self._device_address = config.get(
            USRPReceiver.__name__, 'device_address')
        self._sample_rate_hz = int(config.get(
                USRPReceiver.__name__, 'sample_rate_hz'))
        self._dir = config.get(USRPReceiver.__name__, 'recording_dir')

    def SetHardwareTunerHz(self, freq_hz):
        if self._pipe:
            return False
        self._freq_hz = freq_hz
        return True

    def GetHardwareTunerHz(self):
        return self._freq_hz

    def WaterfallImage(self):
        if self._spectrum is None:
            return None
        return self._spectrum.LatestImage(64)

    def Start(self, stream_url):
        self._stream_url = stream_url
        self._output_path = os.path.join(self._dir, str(int(time.time())))
        self._spectrum = None

        args = USRP_RECORD_BINARY + [
            self._device_address,
            '%d' % self._freq_hz,
            '%d' % self._sample_rate_hz,
            self._output_path]
        logging.info('Starting USRP capture: %s', ' '.join(args))
        try:
            self._pipe = subprocess.Popen(args)
        except OSError:
            logging.exception('Error starting USRP capture')
            return False

        if self._pipe.poll() is not None:
            self._pipe = None
            return False

        self._spectrum = spectrum.SpectrumInt16(self._output_path, 10)
        return True

    def Stop(self):
        if self._pipe is None:
            return
        if self._pipe.poll() is None:
            self._pipe.terminate()
        self._pipe.wait()
        self._pipe = None
        self._spectrum = None

        if not self._Upload():
            logging.error('Upload failed')

    def _Upload(self):
        # TODO(tstranex): Start uploading immediately in the background when
        # Start is called.
        if not self._stream_url:
            return True
        return upload.UploadAndDeleteFile(
            self._output_path, self._stream_url, self._sample_rate_hz, 'SINT16')

    def IsStarted(self):
        if self._pipe is None:
            return False
        return self._pipe.poll() is None



def Configure(config):
    if config.has_section(USRPReceiver.__name__):
        return USRPReceiver(config)
