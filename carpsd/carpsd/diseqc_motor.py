#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

"""DiSEqC motor controller."""

import serial
import time
import threading
import logging

import motor


class Error(Exception):
    pass


class DiSEqCController(object):
    """DiSEqC motor controller over a serial bridge.

    The DisEqC specification can be found at:
    http://www.eutelsat.com/satellites/4_5_5.html
    """

    _DISEQC_FRAMING_MASTER_NOREPLY_FIRST = 0xe0
    
    _DISEQC_ADDRESS_POLAR_POSITIONER = 0x31
    
    _DISEQC_COMMAND_HALT = 0x60
    _DISEQC_COMMAND_DRIVE_EAST = 0x68
    _DISEQC_COMMAND_DRIVE_WEST = 0x69
    _DISEQC_COMMAND_STORE_NN = 0x6a
    _DISEQC_COMMAND_GOTO_NN = 0x6b
    _DISEQC_COMMAND_GOTO_X = 0x6e
    
    def __init__(self, serial_device):
        self.serial = serial.Serial(serial_device, 9600)
        self.hello = self.serial.readline().strip()
        if not self.hello.lower().startswith('hello'):
            raise Error('Device failed to initialize: %s', self.hello)
        logging.info('DiSEqCController initialized. Hello string: %s',
                     self.hello)

    def _send_frame(self, frame):
        self.serial.write(frame)
        result = self.serial.readline().strip()
        if result != 'ok':
            raise Error(result)

    def _send(self, message):
        frame = chr(len(message)) + ''.join(map(chr, message))
        return self._send_frame(frame)

    def hard_power_off(self):
        """Cut off power to the DiSEqC bus completely.

        This is not a DiSEqC command. It's a command to the controller itself.
        """
        return self._send_frame(chr(255))

    def hard_power_on(self):
        """Power on the DiSEqC bus.

        This is not a DiSEqC command. It's a command to the controller itself.
        """
        return self._send_frame(chr(254))

    def halt(self):
        """Make the motor stop.

        This cancels any move operations that are currently underway.
        """
        return self._send([self._DISEQC_FRAMING_MASTER_NOREPLY_FIRST,
                           self._DISEQC_ADDRESS_POLAR_POSITIONER,
                           self._DISEQC_COMMAND_HALT])

    def drive_east(self):
        """Drive the motor east continuously."""
        return self._send([self._DISEQC_FRAMING_MASTER_NOREPLY_FIRST,
                           self._DISEQC_ADDRESS_POLAR_POSITIONER,
                           self._DISEQC_COMMAND_DRIVE_EAST,
                           0x00])

    def drive_west(self):
        """Drive the motor west continuously."""
        return self._send([self._DISEQC_FRAMING_MASTER_NOREPLY_FIRST,
                           self._DISEQC_ADDRESS_POLAR_POSITIONER,
                           self._DISEQC_COMMAND_DRIVE_WEST,
                           0x00])

    def goto_stored_position(self, n):
        """Go to a previously stored position.

        Position 0 corresponds to the reference position.

        Note that some motors lose their memory if they are reset so don't
        rely on this command except for the reference position.
        """
        return self._send([self._DISEQC_FRAMING_MASTER_NOREPLY_FIRST,
                           self._DISEQC_ADDRESS_POLAR_POSITIONER,
                           self._DISEQC_COMMAND_GOTO_NN,
                           n])

    def store_current_position(self, n):
        """Store the current motor position in slot n.

        Note that some motors lose their memory if they are reset so this
        command is not much use in general.
        """
        return self._send([self._DISEQC_FRAMING_MASTER_NOREPLY_FIRST,
                           self._DISEQC_ADDRESS_POLAR_POSITIONER,
                           self._DISEQC_COMMAND_STORE_NN,
                           n])

    @classmethod
    def _encode_azimuth_degrees(cls, azimuth_degrees):
        if azimuth_degrees >= 256:
            reference = 0x1  # +256 degrees
            azimuth_degrees -= 256
        elif azimuth_degrees < 0:
            reference = 0xf  # -256 degrees
            azimuth_degrees += 256
        else:
            reference = 0x0  # 0 degrees

        assert azimuth_degrees >= 0 and azimuth_degrees < 256

        degrees_div_16 = int(azimuth_degrees) / 16
        first_byte = (reference << 4) + degrees_div_16

        degrees = int(azimuth_degrees) % 16

        fraction = azimuth_degrees - int(azimuth_degrees)
        fraction_nibble = 0
        for i in range(4):
            fraction = fraction * 2
            fraction_nibble = fraction_nibble << 1
            if fraction >= 0.5:
                fraction_nibble += 1
                fraction -= 0.5

        second_byte = (degrees << 4) + fraction_nibble

        return first_byte, second_byte

    def goto_x(self, azimuth_degrees):
        """Move the motor to the given angle with respect to the reference pos.

        This command is not supported by all motors.

        azimuth_degrees: degrees clockwise of 'north'.
                         'north' is whatever the motor thinks is north.
        """
        first_byte, second_byte = self._encode_azimuth_degrees(azimuth_degrees)
        return self._send([self._DISEQC_FRAMING_MASTER_NOREPLY_FIRST,
                           self._DISEQC_ADDRESS_POLAR_POSITIONER,
                           self._DISEQC_COMMAND_GOTO_X,
                           first_byte,
                           second_byte])


# TODO(tstranex): Need to implement locking.
class DiSEqCMotor(motor.Motor):

    CALIBRATED_RATE = 1.45  # [degrees / second]
    REFERENCE_ZERO = -30.0  # [degrees] east of true north
    MIN_LIMIT = -80.0  # [degrees] wrt reference
    MAX_LIMIT = +80.0  # [degrees] wrt reference
    SERIAL_DEVICE = '/dev/ttyUSB0'

    def __init__(self, config):

        def get(name, cast, default):
            if not config.has_section(DiSEqCMotor.__name__):
                return default
            if not config.has_option(DiSEqCMotor.__name__, name):
                return default
            return cast(config.get(DiSEqCMotor.__name__, name))

        self.CALIBRATED_RATE = get(
            'calibrated_rate', float, self.CALIBRATED_RATE)
        self.REFERENCE_ZERO = get(
            'reference_zero', float, self.REFERENCE_ZERO)
        self.MIN_LIMIT = get('min_limit', float, self.MIN_LIMIT)
        self.MAX_LIMIT = get('max_limit', float, self.MAX_LIMIT)
        self.SERIAL_DEVICE = get('serial_device', str, self.SERIAL_DEVICE)

        self.controller = DiSEqCController(self.SERIAL_DEVICE)
        self.current = None
        self.target = None
        self.timer = None
        self.last_update = None
        self.power = False

    def PowerOn(self):
        self.controller.hard_power_on()
        logging.info('Powered on motor')
        self.power = True
        return True

    def PowerOff(self):
        self.Halt()
        self.controller.hard_power_off()
        logging.info('Powered off motor')
        self.power = False
        return True

    def IsOn(self):
        return self.power

    def Reset(self):
        return self._Command(self.REFERENCE_ZERO, 60.0,
                             self.controller.goto_stored_position, 0)

    def IsReady(self):
        self._UpdateCurrent()
        return self.current is not None

    def IsMoving(self):
        self._UpdateCurrent()
        return self.target is not None

    def _CommandFinished(self):
        self.current = self.target
        self.target = None
        self.end_time = None
        self.Halt()

    def _Command(self, target, duration, func, *args):
        if not self.IsOn():
            return False
        if self.target is not None:
            return False
        self.timer = threading.Timer(duration, self._CommandFinished)
        now = time.time()
        self.target = target
        self.last_update = now
        self.end_time = now + duration
        func(*args)
        logging.debug('Sent motor command: %s', func.__name__)
        self.timer.start()
        return True

    def Halt(self):
        if self.timer:
            self.timer.cancel()
            self.timer = None

        self.controller.halt()
        logging.debug('Sent motor halt command.')
        self._UpdateCurrent()
        self.target = None
        self.end_time = None
        return True

    def _UpdateCurrent(self):
        if self.target is None:
            return
        if self.current is None:
            return

        now = time.time()
        if now <= self.end_time:
            dt = now - self.last_update
            f = dt / (self.end_time - self.last_update)
            self.current += f * (self.target - self.current)
        self.last_update = now

    def GetAzimuthDegrees(self):
        self._UpdateCurrent()
        return self.current

    def GetAzimuthLimitsDegrees(self):
        return (self.MIN_LIMIT + self.REFERENCE_ZERO,
                self.MAX_LIMIT + self.REFERENCE_ZERO)

    def SetAzimuthDegrees(self, azimuth_degrees):
        if self.current is None:
            return False
        if not self.IsAllowedAzimuthDegrees(azimuth_degrees):
            return False

        # Don't bother if it's less than 1 second resolution.
        def CloseEnough(a):
            if a is None:
                return False
            return abs(a - azimuth_degrees) < self.CALIBRATED_RATE
        if self.target is None:
            if CloseEnough(self.current):
                return True
        else:
            if CloseEnough(self.target):
                return True

        self.Halt()

        delta = azimuth_degrees - self.current
        if delta > 0:
            return self._Command(azimuth_degrees,
                                 delta / self.CALIBRATED_RATE,
                                 self.controller.drive_east)
        else:
            return self._Command(azimuth_degrees,
                                 -delta / self.CALIBRATED_RATE,
                                 self.controller.drive_west)


def Configure(config):
    if config.has_section(DiSEqCMotor.__name__):
        return DiSEqCMotor(config)
