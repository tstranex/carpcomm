#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

import client
import config
import handlers
import logging
import optparse
import sys
import os
import ConfigParser
import time

def Daemonize(pidfile):
    pid = os.fork()
    if pid:
        if pidfile:
            print >>file(pidfile, 'w'), pid
        sys.exit(0)

def ParseOptions():
    p = optparse.OptionParser()
    p.add_option('--daemon', dest='daemon', default=False,
                 help='Start in daemon mode')
    p.add_option('--pidfile', dest='pidfile',
                 help='PID file for daemon mode')
    p.add_option('--config', dest='config', default=None,
                 help='Path to a config file')
    p.add_option('--logfile', dest='logfile', default=None,
                 help='Path to write the log info to')
    options, args = p.parse_args()
    return options

def main():
    logging.root.setLevel(logging.INFO)
    options = ParseOptions()
    if options.logfile is not None:
        logging.root.addHandler(logging.FileHandler(options.logfile))
        logging.root.addHandler(logging.StreamHandler(sys.stderr))
    logging.info('CarpSD started: %s', time.ctime())

    conf = config.LoadConfig(options.config)

    if bool(options.daemon):
        Daemonize(options.pidfile)

    try:
        c = client.Client(conf)

        receiver, motor, tnc = config.Configure(conf)
        if receiver:
            handlers.ReceiverHandlers(receiver).RegisterHandlers(c)
        if motor:
            handlers.MotorHandlers(motor).RegisterHandlers(c)
        if tnc:
            handlers.TNCHandlers(tnc).RegisterHandlers(c)
    except ConfigParser.NoOptionError, e:
        logging.error('  Error while reading the configuration file.')
        logging.error('  The "%s" option is missing from the "%s" section.',
                      e.option, e.section)
        logging.error('  - Try editing %s to fill in the missing value.',
                      options.config)
        logging.error('  - Consult http://carpcomm.com/carpsd/ for more info.')
        return

    if tnc is not None:
        if not tnc.Verify():
            logging.error('Error with TNC configuration.')
            return

    c.ConnectAndServeForever()


if __name__ == '__main__':
    main()
