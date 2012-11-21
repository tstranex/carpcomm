#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

"""Controller for KISS TNCs connected via a serial device."""

import api

import serial
import threading
import logging
import time


SERIAL_READ_TIMEOUT = 0.5


class KISSDecoder(object):
    """State machine that extracts KISS frames from a stream of bytes."""

    FEND = '\xc0'
    FESC = '\xdb'
    TFEND = '\xdc'
    TFESC = '\xdd'

    SIZE_LIMIT = 8192

    def __init__(self):
        self.current_frame = ''
        self.state = 0
        self.frames = []

    def _WriteChar(self, c):
        # We extract KISS frames using a state machine.
        # KISS protocol documentation:
        # http://www.ka9q.net/papers/kiss.html
        # Also relevant is the SLIP protocol:
        # http://tools.ietf.org/html/rfc1055

        if c == self.FEND:
            if len(self.current_frame) > 1:  # Ignore empty frames.
                data_type = self.current_frame[0]
                if data_type != '\x00':
                    logging.info(
                        'Received unknown KISS frame data type from TNC: ' +
                        '%x', ord(data_type))
                self.frames.append(self.current_frame[1:])
            self.state = 0
            self.current_frame = ''
            return

        if self.state == 0:
            if c == self.FESC:
                self.state = 1
            else:
                if len(self.current_frame) < self.SIZE_LIMIT:
                    self.current_frame += c
                self.state = 0
        elif self.state == 1:
            # We're in escape mode.
            if c == self.TFEND:
                if len(self.current_frame) < self.SIZE_LIMIT:
                    self.current_frame += self.FEND
                self.state = 0
            elif c == self.TFESC:
                if len(self.current_frame) < self.SIZE_LIMIT:
                    self.current_frame += self.FESC
                self.state = 0
            else:
                # Error, ignore it.
                self.state = 0
                logging.info(
                    'Received invalid KISS escape char from TNC: %x', ord(c))
        else:
            # Unknown state. This shouldn't happen.
            logging.error('Invalid KISS state %d. Resetting.', self.state)
            self.state = 0

    def Write(self, data):
        for c in data:
            self._WriteChar(c)

    def ReadFrames(self):
        f = self.frames
        self.frames = []
        return f


class _SerialReadThread(threading.Thread):
    """Thread that reads KISS frames from a serial device and uploads them."""

    def __init__(self, serial, api_client, satellite_id):
        threading.Thread.__init__(self)

        self.serial = serial
        self.api_client = api_client
        self.satellite_id = satellite_id

        self.decoder = KISSDecoder()
        self.should_stop = False

    def run(self):
        while not self.should_stop:
            data = self.serial.read()  # Block until there is some data.
            data += self.serial.read(self.serial.inWaiting())

            if data:
                logging.info(
                    '[debug] Recevied serial data: %s', data.encode('hex'))

            self.decoder.Write(data)

            if not self.api_client:
                continue

            timestamp = int(time.time())
            frames = self.decoder.ReadFrames()
            for f in frames:
                ok, status = self.api_client.PostPacket(
                    self.satellite_id, timestamp, f)
                logging.info(
                    '[debug] Posted TNC frame: %s', f.encode('hex'))
                if not ok:
                    host, port = self.api_client.GetServer()
                    logging.info('Error uploading packet to %s:%d: %d, %s',
                                 host, port, status[0], status[1])

        self.serial.close()

    def Stop(self):
        self.should_stop = True

    def GetLatestFrames(self):
        return self.decoder.ReadFrames()


class SerialTNC(object):
    """Controller for KISS TNCs connected via a serial device."""

    def __init__(self, config):
        self._device = config.get(SerialTNC.__name__, 'device')
        self._baud = int(config.get(SerialTNC.__name__, 'baud'))

        self._rtscts = False
        if config.has_option(SerialTNC.__name__, 'rtscts'):
            self._rtscts = config.get(SerialTNC.__name__, 'rtscts') == 'true'

        self._api_client = api.APIClient(config)
        self._thread = None

    def _OpenSerial(self):
        # We need the timeout otherwise the read thread cannot be stopped.
        return serial.Serial(
            self._device,
            self._baud,
            rtscts=self._rtscts,
            timeout=SERIAL_READ_TIMEOUT)

    def Verify(self):
        """Do some quick checks to make sure the configuration works."""
        try:
            s = self._OpenSerial()
            s.close()
        except serial.SerialException, e:
            logging.error('Error opening serial port: %s', e)
            return False
        return True

    def Start(self, api_host, api_port, satellite_id):
        """Start a thread to read packets from the device and upload them.

        If api_host is set to an empty string, then packets are read but
        not uploaded."""

        if self._thread is not None:
            logging.info('Error starting TNC thread: already started')
            return False

        ac = None
        if api_host:
            self._api_client.SetServer(api_host, api_port)
            ac = self._api_client
        else:
            logging.info('TNC frame uploading disabled')

        self._thread = _SerialReadThread(self._OpenSerial(), ac, satellite_id)
        self._thread.start()

        logging.info('Started serial TNC thread for %s', satellite_id)

        return True

    def Stop(self):
        t = self._thread
        if t is None:
            # It's already stopped.
            return True

        t.Stop()
        t.join()

        self._thread = None
        
        if t.isAlive():
            logging.info('Error stopping serial TNC thread. It is still alive.')
            return False

        logging.info('Stopped serial TNC thread.')
        return True

    def GetStateDict(self):
        """Return a dictionary with info about the current state of the device.
        """
        started = self._thread is not None
        return {
            'started': started,
            }

    def GetLatestFrames(self):
        if self._thread is not None:
            return True, self._thread.GetLatestFrames()
        else:
            return False, []


def Configure(config):
    if config.has_section(SerialTNC.__name__):
        return SerialTNC(config)
