#!/usr/bin/python

import numpy
import Image
import sys
from carpcomm.demod import format


BLOCK_SIZE = 2048
BANDWIDTH = 2000.0


def FilterSignal(logpower, filter_frac, fft_size):
    """Filter to the signal only."""

    delta = int(0.5 * fft_size * filter_frac)

    # assume the signal is not near the edge
    # FIXME: we should use a computed frequency hint instead
    boundary = 0.2*fft_size
    center = logpower.std(axis=0)[boundary:-boundary].argmax() + boundary
    max_variation = max(50, delta)

    block = logpower[:,center-max_variation:center+max_variation]
    s = block.std(axis=0)
    signal_pos = min(max(s.argmax(), delta), len(block[0])-delta)
    band = block[:,signal_pos-delta:signal_pos+delta]
    return band


def Output1D(logpower):
    for i in range(len(logpower)):
        print i, logpower[i]


def FilterIQ(flatdata, filter_frac, fft_size):
    logpower = format.LogPowerSpectrum(flatdata, fft_size)
    #format.ToImage(logpower).save('spectrum.png')
    logpower = format.RemoveBackgroundStreaks(logpower)
    #format.ToImage(logpower).save('no_streaks.png')
    logpower = FilterSignal(logpower, filter_frac, fft_size)
    #format.ToImage(logpower).save('active_band.png')
    logpower = logpower.max(axis=1)

    #raw_input('enter')

    return logpower



def ToImage1D(arr):
    return ToImage(numpy.array([arr]*8).transpose())


if __name__ == '__main__':
    sample_rate = float(sys.argv[3])
    fft_size = int(sys.argv[4])

    filter_frac = BANDWIDTH / sample_rate

    f = file(sys.argv[1])
    result = numpy.array([])
    while True:
        flatdata = format.LoadIQ(f, sys.argv[2], fft_size*BLOCK_SIZE)
        if len(flatdata) == 0:
            break
        filtered = FilterIQ(flatdata, filter_frac, fft_size)
        result = numpy.append(result, filtered)
    Output1D(result)
    result.tofile(sys.argv[5])
