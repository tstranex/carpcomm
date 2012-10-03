#!/usr/bin/env python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

from distutils.core import setup
import py2exe

setup(name='carpsd',
      version='0.17',
      description='Carpcomm Station Daemon',
      author='Carpcomm GmbH',
      author_email='info@carpcomm.com',
      url='http://www.carpcomm.com/',
      license='Apache-2.0',
      packages=['carpsd'],
      console=['scripts/carpsd'],
      py_modules = ['carpsd'],
      data_files=[#('/etc/init.d', ['etc/init.d/carpsd']),
                  ('', ['etc/carpsd.conf']),
                  ('', ['etc/ca_cert.pem'])],
      )
