#!/usr/bin/env python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

from distutils.core import setup
from distutils.command.build import build

dependencies = ['numpy', 'Image']

class Build(build):
    """"Custom build class that checks for dependencies before building."""
    def run(self):
        if not check_dependencies():
            print('Unresolved dependencies. Exiting.')
            quit()

        build.run(self)

def check_dependencies():
    no_errors = True

    for d in dependencies:
        try:
            __import__(d)
        except ImportError:
            print("The module '%s' is not installed." % d)
            no_errors = False

    return no_errors

setup(name='carpsd',
      version='0.19',
      description='Carpcomm Station Daemon',
      author='Carpcomm GmbH',
      author_email='info@carpcomm.com',
      url='http://carpcomm.com/carpsd/',
      license='Apache-2.0',
      packages=['carpsd'],
      scripts=['scripts/carpsd'],
      data_files=[('/etc/init.d', ['etc/init.d/carpsd']),
                  ('/etc/carpsd', ['etc/carpsd.conf']),
                  ('/etc/carpsd', ['etc/ca_cert.pem'])],
      cmdclass={'build': Build}
      )
