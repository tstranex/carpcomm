#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

import upload

import unittest
import subprocess
import tempfile
import os.path


class UploadTest(unittest.TestCase):

    def testPipeWaitThread(self):
        _, path = tempfile.mkstemp()
        pipe = subprocess.Popen(['echo'])

        self.assertTrue(os.path.exists(path))

        t = upload._PipeWaitThread()
        t.SetPipe(pipe)
        t.SetPathToDelete(path)
        t.run()

        self.assertFalse(os.path.exists(path))


if __name__ == '__main__':
    unittest.main()
