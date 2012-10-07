#!/usr/bin/env python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

from distutils.core import setup
from distutils.command.build import build

dependencies=['numpy']

# Custom build class that checks for dependencies before building
class Build(build):
    def run(self):
        if not CheckDependencies():
            print('Unresolved dependencies. Exiting.')
            quit()

        build.run(self)

def CheckDependencies():
    noErrors = True

    for d in dependencies:
        try:
            __import__(d)
        except ImportError as ie:
            print('The module \''+ d + '\' is not installed.')
            noErrors = False

    return noErrors

setup(name='carpsd',
      version='0.17',
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
