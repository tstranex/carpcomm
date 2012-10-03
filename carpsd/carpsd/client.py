#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

import BaseHTTPServer
import ssl
import urlparse
import socket
import json
import logging
import platform
import random
import time
import StringIO
import signalling

VERSION = '0.17'
CLIENT_NAME = 'carpsd'

PING_TIMEOUT = 3600  # [seconds]
RECONNECT_DELAY = 600  # [seconds]


class _HTTPHandler(BaseHTTPServer.BaseHTTPRequestHandler):
    def do_GET(self):
        url = urlparse.urlparse(self.path)
        params = urlparse.parse_qs(url.query, keep_blank_values=True)
        r = self.server._Dispatch(url.path, params)

        if r == True:
            code, content_type, data = 200, 'text/plain', ''
        elif r == False:
            code, content_type, data = 500, 'text/plain', ''
        else:
            code, content_type, data = r

        self.send_response(code)
        self.send_header('Content-Type', content_type)
        self.send_header('Content-Length', len(data))
        self.end_headers()
        self.wfile.write(data)

    def address_string(self):
        """Overridden to avoid doing a DNS lookup for every request."""
        host, port = self.client_address[:2]
        return host
    

class Client(object):

    def __init__(self, config):
        self._handlers = {}
        self.RegisterHandler('/Identify', self._IdentifyHandler)
        self.RegisterHandler('/Disconnect', self._DisconnectHandler)
        self.RegisterHandler('/Ping', self._PingHandler)

        self._station_id = config.get(Client.__name__, 'id')
        self._secret = config.get(Client.__name__, 'secret')
        host, port = config.get(Client.__name__, 'server').split(':')
        self._server = (host, int(port))
        self._ca_certs = config.get(Client.__name__, 'ca_certificate')

        f = StringIO.StringIO()
        config.write(f)
        self._config = f.getvalue()

    def Connect(self):
        self._disconnected = False
        s = socket.create_connection(self._server, PING_TIMEOUT)
        self._socket = ssl.wrap_socket(
            s,
            ssl_version=ssl.PROTOCOL_TLSv1,
            cert_reqs=ssl.CERT_REQUIRED,
            ca_certs=self._ca_certs)

    def ServeUntilDisconnected(self):
        while not self._disconnected:
            h = _HTTPHandler(self._socket, self._server, self)
            if not h.raw_requestline:
                # FIXME: seems to prematurely timeout on win32.
                # should check the timer to be sure.
                logging.info('Connection closed or timed out.')
                break
        self._socket.close()

    def ConnectAndServeForever(self):
        while True:
            try:
                self.Connect()
                self.ServeUntilDisconnected()
            except socket.error:
                logging.exception('Connection error:')

            # TODO: exponential backoff
            delay = RECONNECT_DELAY + int(RECONNECT_DELAY * random.random())
            logging.info('Reconnecting in %d seconds', delay)
            time.sleep(delay)

    def RegisterHandler(self, path, handler):
        self._handlers[path] = handler

    def _Dispatch(self, path, params):
        if path in self._handlers:
            return self._handlers[path](params)
        else:
            return 404, 'text/plain', ''

    def _IdentifyHandler(self, params):
        # TODO(tstranex): Move this to handlers.py.

        signalling.Get().SignalIdentified()

        # Platform information for stats.
        arch = {
            'machine': platform.machine(),
            'processor': platform.processor(),
            'python_version': platform.python_version(),
            'python_implementation': platform.python_implementation(),
            'system': platform.system(),
            'release': platform.release(),
            'version': platform.version(),
            'win32_ver': platform.win32_ver(),
            'mac_ver': platform.mac_ver(),
            'linux_distribution': platform.linux_distribution(),
            }

        data = {
            'version': VERSION,
            'client': CLIENT_NAME,
            'station_id': self._station_id,
            'secret': self._secret,
            'platform': arch,
            'config': self._config,
            }
        return 200, 'text/plain', json.dumps(data)

    def _DisconnectHandler(self, params):
        logging.error('Disconnected by server: %s', params['reason'][0])
        self._disconnected = True
        return True

    def _PingHandler(self, params):
        signalling.Get().SignalPing()
        return 200, 'text/plain', 'pong'
