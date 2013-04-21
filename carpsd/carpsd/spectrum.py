#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

"""Spectrum
"""

import Image
import numpy
import os


class SpectrumByte:
    ITEM_TYPE = numpy.byte
    FRAME_SIZE = 1024

    def __init__(self, path, skip):
        self.path = path
        self.skip = skip
        self.f = None
        self.frame_bytes = 2 * self.ITEM_TYPE().nbytes * self.FRAME_SIZE

    def _Open(self):
        if self.f is not None:
            return self.f
        if os.path.exists(self.path):
            self.f = open(self.path, 'rb')
        return self.f

    def _NumFrames(self):
        try:
            available_bytes = os.path.getsize(self.path) - self.f.tell()
            return available_bytes / self.frame_bytes / self.skip
        except OSError:
            return 0

    def _AdvanceFrames(self, num_frames):
        self.f.seek(num_frames * self.frame_bytes * self.skip, 1)

    def _ReadIQSamples(self, f, num_samples):
        data = numpy.fromfile(f, dtype=self.ITEM_TYPE, count=2*num_samples)
        return (data - 127.0) / 128.0

    def _ReadFrames(self, num_frames):
        frames = []
        for i in range(num_frames):
            data = self._ReadIQSamples(self.f, self.FRAME_SIZE)
            iq = data[0::2] + 1j * data[1::2]
            frames.append(iq)
            self.f.seek(self.frame_bytes * (self.skip - 1), 1)
        return numpy.array(frames)

    def _LatestFrames(self, max_frames):
        available = self._NumFrames()
        num_to_read = min(available, max_frames)
        self._AdvanceFrames(available - num_to_read)
        return self._ReadFrames(num_to_read)

    def LatestImage(self, max_frames):
        """Return the latest spectrum Image with at max_frames rows.

        Returns None if no data is available.
        """

        if self._Open() is None:
            return None

        frames = self._LatestFrames(max_frames)
        w = self.FRAME_SIZE
        h = len(frames)
        if h == 0:
            return None

        # Flip the image so that newer rows are on top.
        frames = numpy.flipud(frames)

        # Compute the power spectrum
        power = []
        for i in range(h):
            ft = numpy.fft.fftshift(numpy.fft.fft(frames[i]))
            power.append(numpy.real(numpy.conj(ft) * ft))
        power = numpy.array(power)

        n = w * h
        power += 1e-30  # To avoid -inf after taking the log.
        logpower = numpy.log(power).reshape(n)
        max_value = numpy.max(logpower)
        min_value = numpy.min(logpower)
        data = 255.0 * (logpower - min_value) / (max_value - min_value)
        data = numpy.array(data, dtype=numpy.uint8)

        img = Image.new('L', (w, h))
        img.putdata(data)

        return img


class SpectrumInt16(SpectrumByte):
    ITEM_TYPE = numpy.int16

    def _ReadIQSamples(self, f, num_samples):
        data = numpy.fromfile(f, dtype=self.ITEM_TYPE, count=2*num_samples)
        return data / 32768.0
