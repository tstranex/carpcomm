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


FCD_CONTROL_BINARY = 'carpsd-fcd'
FCD_CONTROL_TIMEOUT = 3
FCD_CONTROL_RETRIES = 5
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
        self._spectrum = None
        self._dir = config.get(FCDReceiver.__name__, 'recording_dir')
        self._alsa_device = config.get(FCDReceiver.__name__, 'alsa_device')

        params = self._GetStatus()
        if not params:
            raise FCDReceiverError('Unable to get FCD device status.')
        
        if 'Pro+' in params['Type']:
            logging.info('Detected FCD Pro+')
            self._sample_rate = 192000
        else:
            logging.info('Detected FCD Pro')
            self._sample_rate = 96000

    def SetHardwareTunerHz(self, freq_hz):
        ok, stdout = self._CommandWithTimeout([
                '--device', '0', '--set_freq_hz', '%d' % freq_hz])
        return ok

    def GetHardwareTunerHz(self):
        params = self._GetStatus()
        if not params:
            return 0
        return int(params['Frequency [Hz]'])

    def _CommandWithTimeout(self, arg):
        for i in range(FCD_CONTROL_RETRIES):
            args = ['timeout', str(FCD_CONTROL_TIMEOUT),
                    FCD_CONTROL_BINARY] + arg
            logging.info('Running command: %s', ' '.join(args))
            p = subprocess.Popen(args, stdout=subprocess.PIPE)
            stdout, stderr = p.communicate()
            if p.returncode == 0:
                return True, stdout
            else:
                logging.error('Error, return code: %d', p.returncode)

        logging.error('Giving up.')
        return False, ''

    def _GetStatus(self):
        ok, stdout = self._CommandWithTimeout(['--device', '0'])
        if not ok:
            return False

        params = {}
        for line in stdout.split('\n'):
            line = line.strip()
            if not line:
                continue
            parts = line.split(': ')
            if len(parts) != 2:
                continue
            params[parts[0]] = parts[1]
        return params

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
