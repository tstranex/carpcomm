#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

"""Methods for uploading streams."""

import client

import logging
import threading
import subprocess
import os
import base64
import json
import httplib
import socket
import ssl
import signalling
import urllib
import datetime


class _PipeWaitThread(threading.Thread):
    """Thread that calls wait() on a pipe."""

    def SetPipe(self, pipe):
        self.pipe = pipe

    def SetPathToDelete(self, path):
        self.path = path

    def run(self):
        signalling.Get().SignalUploadStart()
        self.pipe.wait()
        logging.info('Upload complete.')
        signalling.Get().SignalUploadStop()
        os.remove(self.path)


def UploadAndDeleteFile(path, stream_url, rate, dtype):
    """Upload the finalized file in another process."""

    query = '?rate=%d&type=%s' % (rate, dtype)
    url = stream_url + query

    args = ['curl', '--upload-file', path, url]
    logging.info('Starting upload: %s', ' '.join(args))
    try:
        pipe = subprocess.Popen(args)
    except OSError:
        logging.exception('Error starting upload')
        return False

    if pipe.poll() is not None:
        return False

    # Although the upload is running in another process, we need to wait()
    # for it to terminate to avoid leaving defunct processes around.
    t = _PipeWaitThread()
    t.SetPipe(pipe)
    t.SetPathToDelete(path)
    t.start()

    return True


# TODO(tstranex): Consider moving this this to a new file called api.py.
class APIClient:
    def __init__(self, config):
        self._station_id = config.get(client.Client.__name__, 'id')
        self._secret = config.get(client.Client.__name__, 'secret')
        self._ca_certs = config.get(client.Client.__name__, 'ca_certificate')

        # We should consider setting default values here.
        self.SetServer(None, None)

    def SetServer(self, host, port):
        self.host = host
        self.port = port

    def GetServer(self):
        return self.host, self.port

    def _SendRequest(self, method, path, body):
        s = socket.create_connection((self.host, self.port))
        s = ssl.wrap_socket(
            s,
            ssl_version=ssl.PROTOCOL_TLSv1,
            cert_reqs=ssl.CERT_REQUIRED,
            ca_certs=self._ca_certs)

        c = httplib.HTTPSConnection(
            self.host, self.port, None, self._ca_certs)
        # UGLY. We have to open the socket ourselves because HTTPSConnection
        # doesn't allow us to specify all SSL parameters.
        c.sock = s
        c.request(method, path, body)
        r = c.getresponse()
        if r.status != httplib.OK:
            return False, (r.status, r.reason), None

        return True, (r.status, r.reason), r.read()

    def PostPacket(self, satellite_id, timestamp, frame):
        req = {
            'station_id': self._station_id,
            'station_secret': self._secret,
            'timestamp': timestamp,
            'satellite_id': satellite_id,
            'format': 'FRAME',
            'frame_base64': base64.b64encode(frame),
            }
        body = json.dumps(req)
        ok, status, body = self._SendRequest('POST', '/PostPacket', body)
        return ok, status

    def GetLatestPackets(self, satellite_id, limit):
        req = {
            'group_id': self._station_id,
            'group_secret': self._secret,
            'satellite_id': satellite_id,
            'limit': limit,
            }
        params = urllib.urlencode(req)
        ok, status, body = self._SendRequest(
            'GET', '/GetLatestPackets?%s' % params, '')
        packets = []
        if ok:
            for p in json.loads(body):
                timestamp = datetime.datetime.fromtimestamp(int(p['timestamp']))
                frame = base64.b64decode(p['frame_base64'])
                packets.append((timestamp, frame))
        return ok, status, packets
