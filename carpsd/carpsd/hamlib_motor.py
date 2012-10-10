#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

"""Motor controller using the hamlib backend."""

import threading
import logging
import time
import subprocess

import motor


ROTCTL_BINARY = 'rotctl'

class _HamlibRotator(object):

    def __init__(self, model, device, params):
        self.model = model
        self.device = device
        self.conf = ','.join(['%s=%s' % p for p in params.items()])

    def _check_output(self, args):
        return subprocess.check_output(args)

    def _RunCommand(self, command):
        args = [ROTCTL_BINARY,
                '--model=%s' % self.model,
                '--rot-file=%s' % self.device]
        if self.conf:
            args.append('--set-conf=%s' % self.conf)
        args += command
        logging.info('Sending hamlib rotator command: %s', ' '.join(args))
        try:
            result = self._check_output(args)
        except subprocess.CalledProcessError, e:
            logging.info('Error calling hamlib rotator controller: %s',
                         e.output)
            return False, e.output
        return True, result

    def SetAzimuthElevation(self, az, el):
        ok, output = self._RunCommand(['P', str(az), str(el)])
        return ok

    def GetAzimuthElevation(self):
        ok, output = self._RunCommand(['p'])
        if not ok:
            return False, (None, None)
        az, el = map(float, output.split())
        return True, (az, el)

    def Stop(self):
        ok, output = self._RunCommand(['S'])
        return ok

    def GetInfo(self):
        return self._RunCommand(['_'])


class _MotorThread(threading.Thread):

    def __init__(self, program, rotator):
        threading.Thread.__init__(self)

        self.program = program
        self.rotator = rotator
        self.start_time = time.time()
        self.should_stop = False

    def run(self):
        for t, az, el in self.program:
            dt = self.start_time + t - time.time()
            if dt < 0.0:
                continue
            time.sleep(dt)
            if self.should_stop:
                return
            self.rotator.SetAzimuthElevation(az, el)

    def Stop(self):
        self.should_stop = True


class HamlibMotor(motor.Motor):
    """Motor controller for hamlib rotors."""

    def __init__(self, config):
        model = config.get(HamlibMotor.__name__, 'model')
        device = config.get(HamlibMotor.__name__, 'device')

        params = {}
        prefix = 'hamlib_param_'
        for option in config.options(HamlibMotor.__name__):
            if not option.startswith(prefix):
                continue
            k = option[len(prefix):]
            params[k] = config.get(HamlibMotor.__name__, option)

        self.rotator = _HamlibRotator(model, device, params)
        self.thread = None

        logging.info('HamlibMotor created.')

    def Start(self, program):
        if not program:
            return False
        if not self.Stop():
            return False

        self.thread = _MotorThread(program, self.rotator)
        self.thread.start()
        logging.info('Started motor control thread.')
        return True

    def Stop(self):
        if self.thread is not None:
            self.thread.Stop()
            self.thread = None
        return self.rotator.Stop()

    def GetStateDict(self):
        result = {}
        ok, (az, el) = self.rotator.GetAzimuthElevation()
        if az is not None:
            result['azimuth_degrees'] = az
        if el is not None:
            result['elevation_degrees'] = el
        return result

    def GetInfoDict(self):
        d = {'driver': HamlibMotor.__name__}
        ok, info = self.rotator.GetInfo()
        if ok:
            d['hamlib_info'] = info
        return d


def Configure(config):
    if config.has_section(HamlibMotor.__name__):
        return HamlibMotor(config)
