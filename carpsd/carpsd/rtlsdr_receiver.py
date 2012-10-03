#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

"""RTL SDR radio Receiver implementation."""

import os.path
import time
import subprocess
import logging

import receiver
import spectrum
import upload


RTL_SDR_BINARY = 'rtl_sdr'


class RTLSDRReceiver(receiver.Receiver):
    """Receiver implementation for RTL SDR."""

    def __init__(self, config):
        self._stream_url = None
        self._output_path = None
        self._spectrum_path = None
        self._pipe = None
        self._upload_pipe = None
        self._freq_hz = 145500000
        self._tuner_gain_db = float(config.get(
                RTLSDRReceiver.__name__, 'tuner_gain_db'))
        self._sample_rate_hz = int(config.get(
                RTLSDRReceiver.__name__, 'sample_rate_hz'))
        self._spectrum = None
        self._dir = config.get(RTLSDRReceiver.__name__, 'recording_dir')
        self._device_index = int(config.get(
                RTLSDRReceiver.__name__, 'device_index'))

    def SetHardwareTunerHz(self, freq_hz):
        if self._pipe:
            # TODO(tstranex): Implement this.
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

        args = [RTL_SDR_BINARY,
                self._output_path,
                '-f %d' % self._freq_hz,
                '-s %d' % self._sample_rate_hz,
                '-d %d' % self._device_index,
                '-g %f' % self._tuner_gain_db]
        logging.info('Starting rtl-sdr capture: %s', ' '.join(args))
        try:
            self._pipe = subprocess.Popen(args)
        except OSError:
            logging.exception('Error starting rtl-sdr')
            return False

        if self._pipe.poll() is not None:
            self._pipe = None
            return False

        self._spectrum = spectrum.SpectrumByte(self._output_path, 20)

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
        return upload.UploadAndDeleteFile(self._output_path, self._stream_url)

    def IsStarted(self):
        if self._pipe is None:
            return False
        return self._pipe.poll() is None


def Configure(config):
    if config.has_section(RTLSDRReceiver.__name__):
        return RTLSDRReceiver(config)
