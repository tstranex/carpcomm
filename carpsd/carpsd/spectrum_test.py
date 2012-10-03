#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

import unittest
import spectrum
import tempfile
import os


class SpectrumByteTest(unittest.TestCase):
    FRAME_SIZE = spectrum.SpectrumByte.FRAME_SIZE

    # TODO: Remove this once we move to Python > 2.7
    def assertListEqual(self, a, b):
        self.assertEqual(len(a), len(b))
        for i in range(len(a)):
            self.assertEqual(a[i], b[i], 'differ in pos %d' % i)

    def _CreateSpectrum(self, num_frames, skip):
        """Create a Spectrum using a temporary test file."""
        _, path = tempfile.mkstemp()
        f = open(path, 'wb')
        for i in range(num_frames):
            for j in range(self.FRAME_SIZE):
                f.write('%c%c' % (128 + i, 128 - i))
        f.close()

        s = spectrum.SpectrumByte(path, skip)
        s._Open()
        return s, self.FRAME_SIZE

    def testBasic(self):
        s, frame_size = self._CreateSpectrum(11, 2)

        self.assertEquals(5, s._NumFrames())

        frames = s._ReadFrames(1)
        self.assertEqual(1, len(frames))
        self.assertEqual(self.FRAME_SIZE, len(frames[0]))

        self.assertEquals(4, s._NumFrames())

        frames = s._LatestFrames(2)
        self.assertEqual(2, len(frames))
        self.assertEqual(self.FRAME_SIZE, len(frames[0]))
        self.assertEqual(self.FRAME_SIZE, len(frames[1]))

        self.assertEquals(0, s._NumFrames())

    def testLimit(self):
        s, frame_size = self._CreateSpectrum(16, 3)
        s._AdvanceFrames(4)
        frames = s._LatestFrames(2)
        self.assertEqual(1, len(frames))
        self.assertEqual(self.FRAME_SIZE, len(frames[0]))
        self.assertEquals(0, s._NumFrames())

    def testEmptySpectrumFile(self):
        s, frame_size = self._CreateSpectrum(0, 3)
        frames = s._LatestFrames(2)
        self.assertEqual(0, len(frames))

    def testNoSpectrumAvailable(self):
        s, frame_size = self._CreateSpectrum(4, 2)
        s._AdvanceFrames(2)
        frames = s._LatestFrames(2)
        self.assertEqual(0, len(frames))

    def testLazyLoad(self):
        path = 'filedoesntexist'  # File doesn't exist.
        self.assertFalse(os.path.exists(path))
        s = spectrum.SpectrumByte(path, 1)
        self.assertEqual(None, s.LatestImage(5))

    def testImage(self):
        s, frame_size = self._CreateSpectrum(5, 2)
        img = s.LatestImage(2)
        self.assertTrue(img is not None)


class SpectrumInt16Test(SpectrumByteTest):
    FRAME_SIZE = spectrum.SpectrumInt16.FRAME_SIZE

    def _CreateSpectrum(self, num_frames, skip):
        _, path = tempfile.mkstemp()
        f = open(path, 'wb')
        for i in range(num_frames):
            for j in range(self.FRAME_SIZE):
                f.write('%c%c%c%c' % (i, i+1, i+2, i+3))
        f.close()

        s = spectrum.SpectrumInt16(path, skip)
        s._Open()
        return s, self.FRAME_SIZE


if __name__ == '__main__':
    unittest.main()
