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
