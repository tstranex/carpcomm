#!/usr/bin/python

# python -m carpcomm.demod.packet.hrbe_doppler_correction ~/tmp/3351160565435329822 rtl16 out.txt

import numpy
import sys

from carpcomm.demod import format

def FilterIQ(flatdata, fft_size):
    logpower = format.LogPowerSpectrum(flatdata, fft_size)
    logpower = format.RemoveBackgroundStreaks(logpower)

    maxima_pos = logpower.argmax(axis=1)
    maxima_val = logpower.max(axis=1)

    return maxima_pos, maxima_val


def FindBlocks(result_pos, result_val):
    avg = result_val.mean()
    std = result_val.std()

    N = len(result_pos)

    v = result_val
    for i in range(N):
        if v[i] > avg + 2*std:
            v[i] = 1
        else:
            v[i] = 0

    # Find boundaries
    bounds = []
    threshold = 250 #500 #1000
    zero_count = 0
    begin_i = None
    for i in range(N):
        if v[i] == 0:
            if begin_i is None:
                begin_i = i
                zero_count = 0
            zero_count += 1
        else:
            if zero_count > threshold:
                # found end of boundary
                bounds.append((begin_i, i))
            begin_i = None
            zero_count = 0
    if begin_i is not None:
        # The final bound doesn't need to exceed threshold since we assume
        # everything is zero after the pass.
        bounds.append((begin_i, N))

    #print 'bounds:', bounds

    # Invert bounds.
    inverted = []
    for i in range(len(bounds)-1):
        a = bounds[i][1]
        b = bounds[i+1][0]
        inverted.append((a, b))

    #print 'inverted:', inverted
    return inverted


def BlockLinearFit(result_pos, result_val, inverted):

    v = result_val
    N = len(result_pos)

    # do a linear fit to each bound
    fits = []
    for begin, end in inverted:
        X = []
        Y = []

        # for hrbe, use only the first half of the points since it's a
        # single frequency. second half is shifting between mark and space.
        lin_end = max((begin+end)/2, begin+1)
        for i in range(begin, lin_end):
            if v[i] == 1:
                X.append(i)
                Y.append(result_pos[i])

        A = numpy.vstack([X, numpy.ones(len(X))]).T
        m, c = numpy.linalg.lstsq(A, Y)[0]
        fits.append((m, c))

    lin_fit = numpy.zeros(N)
    doppler_points = []
    last_mid = 0
    for i in range(len(inverted)):
        if i+1 >= len(inverted):
            next_begin = N
        else:
            next_begin = inverted[i+1][0]
        mid = (inverted[i][1] + next_begin)/2
        m, c = fits[i]
        for j in range(last_mid, mid):
            lin_fit[j] = m*j + c

        x = 0.5*(inverted[i][0] + inverted[i][1])
        y = m*x + c
        #doppler_points.append((last_mid, m*last_mid + c))
        #doppler_points.append((mid-1, m*(mid-1) + c))
        doppler_points.append((last_mid, y))
        doppler_points.append((mid-1, y))

        last_mid = mid

    return lin_fit, doppler_points


def BlockConstantFit(result_pos, result_val, inverted):

    v = result_val
    N = len(result_pos)

    lin_fit = numpy.zeros(N)
    doppler_points = []
    last_mid = 0
    for i in range(len(inverted)):
        if i+1 >= len(inverted):
            next_begin = N
        else:
            next_begin = inverted[i+1][0]
        mid = (inverted[i][1] + next_begin)/2

        begin, end = inverted[i]

        #c = result_pos[begin:end].mean()
        c = (result_pos[begin:end] * v[begin:end]).sum() / v[begin:end].sum()

        X = []
        for i in range(begin, end):
            if v[i] > 0:
                X.append(result_pos[i])
        s = numpy.array(X).std()

        # remove outliers
        Y = []
        for i in range(begin, end):
            if v[i] > 0 and abs(result_pos[i] - c) < s:
                Y.append(result_pos[i])

        if len(Y) > 0:
            c = sum(Y) / len(Y)

        #c = (result_pos[begin:end] * v[begin:end]).sum() / v[begin:end].sum()

        lin_fit[last_mid:mid] = c

        doppler_points.append((last_mid, c))
        doppler_points.append((mid-1, c))

        last_mid = mid

    return lin_fit, doppler_points


def ConstantDoppler(FFT_SIZE, delta_frac):
    i = FFT_SIZE/2 + delta_frac*FFT_SIZE
    return [(0, i), (10, i)]


def IdentityDoppler(FFT_SIZE):
    return ConstantDoppler(FFT_SIZE, 0)


def LoadAndFindBlocks(input_path, input_format, FFT_SIZE, BLOCK_SIZE):
    f = file(input_path)
    result_pos, result_val = numpy.array([]), numpy.array([])
    while True:
        flatdata = format.LoadIQ(f, input_format, FFT_SIZE*BLOCK_SIZE)
        if len(flatdata) == 0:
            break
        pos, val = FilterIQ(flatdata, FFT_SIZE)
        result_pos = numpy.append(result_pos, pos)
        result_val = numpy.append(result_val, val)

    blocks = FindBlocks(result_pos, result_val)
    return blocks, result_pos, result_val


def OutputDopplerPoints(output_path, doppler_points, FFT_SIZE):
    out = file(output_path, 'w')
    for i, f in doppler_points:
        sample_num = i*FFT_SIZE
        frac = float(f - FFT_SIZE/2) / FFT_SIZE
        print >>out, sample_num, frac


if __name__ == '__main__':
    input_path = sys.argv[1]
    input_format = sys.argv[2]
    output_path = sys.argv[3]
    strategy = sys.argv[4]

    FFT_SIZE = 8192
    BLOCK_SIZE = 512

    blocks, result_pos, result_val = LoadAndFindBlocks(
        input_path, input_format, FFT_SIZE, BLOCK_SIZE)

    if strategy == 'HRBE_LINEAR':
        lin_fit, doppler_points = BlockLinearFit(result_pos, result_val, blocks)
    elif strategy == 'CONSTANT_BURST':
        lin_fit, doppler_points = BlockConstantFit(
            result_pos, result_val, blocks)
    elif strategy == 'SNAPS_FCD_OFFSET':
        doppler_points = ConstantDoppler(FFT_SIZE, 25000.0 / 192000.0)
    elif strategy == 'DISABLED':
        doppler_points = IdentityDoppler(FFT_SIZE)
    else:
        assert False, 'unknown strategy %s' % strategy

    OutputDopplerPoints(output_path, doppler_points, FFT_SIZE)

    #for i in range(len(result_pos)):
    #    print i, result_pos[i], result_val[i], lin_fit[i]


