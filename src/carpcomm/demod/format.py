#!/usr/bin/python

# Author: Timothy Stranex <tstranex@carpcomm.com>
# Copyright 2013 Timothy Stranex

#!/usr/bin/python

import numpy
import Image

def LoadIQ(f, format, num_samples):
    if format == 'FLOAT32':
        try:
            return numpy.fromfile(f, numpy.complex64, count=num_samples)
        except MemoryError:
            return numpy.fromfile(f, numpy.complex64)
    elif format == 'UINT8':
        try:
            data = numpy.fromfile(f, numpy.byte, count=2*num_samples)
        except MemoryError:
            data = numpy.fromfile(f, numpy.byte)
        data = (data - 127.0) / 128.0
        return data[0::2] + 1j * data[1::2]
    elif format == 'SINT16':
        try:
            data = numpy.fromfile(f, numpy.int16, count=2*num_samples)
        except MemoryError:
            data = numpy.fromfile(f, numpy.int16)
        data = data / 32768.0
        return data[0::2] + 1j * data[1::2]

def LogPowerSpectrum(flatdata, fft_size):
    n = int(len(flatdata) / fft_size)

    frames = flatdata[:n*fft_size].reshape((n, fft_size))
    frames *= numpy.blackman(fft_size)
    ft = numpy.fft.fft(frames)
    ft = numpy.fft.fftshift(ft, axes=(-1,))

    power = numpy.real(numpy.conj(ft) * ft)
    logpower = numpy.log(power + 1e-30)
    return logpower

def RemoveBackgroundStreaks(logpower):
    """Remove continuous background streaks."""
    mean = logpower.mean(axis=0)
    return logpower - mean

def ScaleToImage(arr, min_value, max_value):
    h, w = arr.shape
    flat = arr.reshape(w*h)
    flat = 255.0 * (flat - min_value) / (max_value - min_value)
    data = numpy.array(flat, numpy.int)
    img = Image.new('L', (w, h))
    img.putdata(data)
    return img

def ToImage(arr):
    return ScaleToImage(arr, numpy.min(arr), numpy.max(arr))
