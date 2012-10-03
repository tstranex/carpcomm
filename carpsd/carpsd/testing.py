#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

import config
import client


def GetConfigForTesting():
    conf = config.GetDefaultConfig()

    section = client.Client.__name__
    conf.add_section(section)
    conf.set(section, 'id', 'test_station_id')
    conf.set(section, 'secret', 'test_station_secret')
    conf.set(section, 'ca_certificate', '../etc/ca_cert.pem')

    return conf
