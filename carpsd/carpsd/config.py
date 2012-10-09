#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

import ConfigParser
import StringIO
import logging

import diseqc_motor
import dummy_motor
import hamlib_motor
import dummy_receiver
import fcd_receiver
import rtlsdr_receiver
import serial_tnc
import signalling


def GetDefaultConfig():
    return ConfigParser.SafeConfigParser()


def LoadConfig(config_path):
    config = GetDefaultConfig()
    if config_path is not None:
        logging.info('Reading config file: %s', config_path)
        config.read(config_path)
    return config


def _ConfigureReceiver(config):
    r = fcd_receiver.Configure(config)
    if r:
        return r
    r = rtlsdr_receiver.Configure(config)
    if r:
        return r
    r = dummy_receiver.Configure(config)
    if r:
        return r
    return None


def _ConfigureMotor(config):
    m = hamlib_motor.Configure(config)
    if m:
        return m
    m = diseqc_motor.Configure(config)
    if m:
        return m
    m = dummy_motor.Configure(config)
    if m:
        return m
    return None


def _ConfigureTNC(config):
    t = serial_tnc.Configure(config)
    if t:
        return t
    return None


def Configure(config):
    f = StringIO.StringIO()
    config.write(f)
    logging.info('Configuring with this config:\n%s', f.getvalue())

    signalling.Configure(config)

    receiver = _ConfigureReceiver(config)
    motor = _ConfigureMotor(config)
    tnc = _ConfigureTNC(config)
    return receiver, motor, tnc
