#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

"""FUNcube Dongle radio Receiver implementation."""

import subprocess
import logging
import os.path
import time

import receiver
import spectrum
import upload


QTHID_BINARY = 'qthid-cli'
ARECORD_BINARY = 'arecord'  # We could also try using 'parec' for pulseaudio.


class FCDReceiverError(Exception):
    pass


class FCDReceiver(receiver.Receiver):
    """Receiver implementation for the FUNcube Dongle."""

    def __init__(self, config):
        self._stream_url = None
        self._output_path = None
        self._spectrum_path = None
        self._fcd_pipe = None
        self._upload_pipe = None
        self._freq_hz = 145500000
        self._spectrum = None
        self._dir = config.get(FCDReceiver.__name__, 'recording_dir')
        self._alsa_device = config.get(FCDReceiver.__name__, 'alsa_device')
        self._frequency_correction = int(config.get(
                FCDReceiver.__name__, 'frequency_correction'))

        self._sample_rate = 96000
        if config.has_option(FCDReceiver.__name__, 'model'):
            model = config.get(FCDReceiver.__name__, 'model')
            if model == 'pro':
                self._sample_rate = 96000
            elif model == 'proplus':
                self._sample_rate = 192000
            else:
                raise FCDReceiverError('Unknown FCDReceiver model: %s' % model)

    def SetHardwareTunerHz(self, freq_hz):
        args = [QTHID_BINARY, '--set_freq_hz', '%d' % freq_hz]
        logging.info('Setting FCD frequency: %s', ' '.join(args))
        try:
            subprocess.check_call(args)
        except subprocess.CalledProcessError:
            logging.exception('Error tuning FCD')
            return False
        self._freq_hz = freq_hz
        return True

    def GetHardwareTunerHz(self):
        return self._freq_hz

    def WaterfallImage(self):
        if self._spectrum is None:
            return None
        return self._spectrum.LatestImage(256)

    def Start(self, stream_url):
        self._stream_url = stream_url
        self._output_path = os.path.join(self._dir, str(int(time.time())))
        self._spectrum = None

        args = [ARECORD_BINARY,
                '-D', self._alsa_device,
                '-f', 'S16_LE',
                '-r', '%d' % self._sample_rate,
                '-c', '2',
                '-t', 'raw',
                self._output_path]
        logging.info('Starting FCD capture: %s', ' '.join(args))
        try:
            self._fcd_pipe = subprocess.Popen(args)
        except OSError:
            logging.exception('Error starting FCD')
            return False

        if self._fcd_pipe.poll() is not None:
            self._fcd_pipe = None
            return False

        self._spectrum = spectrum.SpectrumInt16(self._output_path, 10)

        return True

    def Stop(self):
        if self._fcd_pipe is None:
            return
        if self._fcd_pipe.poll() is None:
            self._fcd_pipe.terminate()
        self._fcd_pipe.wait()
        self._fcd_pipe = None
        self._spectrum = None

        if not self._Upload():
            logging.error('Upload failed')

    def _Upload(self):
        # TODO(tstranex): Start uploading immediately in the background when
        # Start is called.
        if not self._stream_url:
            return True
        return upload.UploadAndDeleteFile(
            self._output_path, self._stream_url, self._sample_rate, 'SINT16')

    def IsStarted(self):
        if self._fcd_pipe is None:
            return False
        return self._fcd_pipe.poll() is None


def Configure(config):
    if config.has_section(FCDReceiver.__name__):
        return FCDReceiver(config)
