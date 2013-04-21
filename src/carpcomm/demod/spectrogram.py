#!/usr/bin/python

# Author: Timothy Stranex <tstranex@carpcomm.com>
# Copyright 2013 Timothy Stranex

#!/usr/bin/python

import Image
import ImageDraw
import sys

import format


FFT_SIZE = 2048 #1024
BLOCK_SIZE = 1024 #2048


def _GuessRange(sample_rate, sample_type):
    if sample_type == 'UINT8':
        # Probably from RTL-SDR
        return -11.0, 14.0
    elif sample_type == 'FLOAT32':
        return -14.0, 12.0
    else:
        # Probably from FUNcube Dongle
        return -25.0, 5.0


def Spectrogram(path, sample_rate, sample_type, title):
    min_value, max_value = _GuessRange(sample_rate, sample_type)

    blocks = []
    height = 0

    f = file(path)
    i = 0
    while True:
        flatdata = format.LoadIQ(f, sample_type, FFT_SIZE*BLOCK_SIZE)
        if len(flatdata) == 0:
            break
        logpower = format.LogPowerSpectrum(flatdata, FFT_SIZE)
        print logpower.min(), logpower.max()
        img = format.ScaleToImage(logpower, min_value, max_value)
        height += img.size[1]

        p = path + '_spectrogram_%05d.png' % i
        i += 1
        img.save(p)

        blocks.append(p)

    xpad = 50
    ypad = 40

    result = Image.new('L', (FFT_SIZE + 2*xpad, height + 2*ypad))
    draw = ImageDraw.Draw(result)
    draw.rectangle([(0, 0), (FFT_SIZE + 2*xpad, height + 2*ypad)], fill='#fff')

    y = ypad
    for p in blocks:
        img = Image.open(p)
        result.paste(img, (xpad, y))
        y += img.size[1]

    # Frequency marks
    delta_f = 10000.0
    n = int(sample_rate / delta_f)
    for i in range(-n/2,n/2+1):
        x = xpad + FFT_SIZE/2 + i * (FFT_SIZE * delta_f / sample_rate)
        f = i*delta_f/1000
        #draw.line([(x, ypad), (x, ypad-5)])
        draw.line([(x, ypad-5), (x, height+ypad+5)])
        draw.text((x-25, ypad-20), '%.0f kHz' % f)

    # Time marks
    n = int(height * FFT_SIZE / float(sample_rate))
    for i in range(0, n, 2):
        y = ypad + i * sample_rate / FFT_SIZE
        #draw.line([(xpad+FFT_SIZE, y), (xpad+FFT_SIZE+5, y)])
        draw.line([(xpad-5, y), (xpad+FFT_SIZE+5, y)])
        draw.text((xpad+FFT_SIZE + 15, y-5), '%d s' % i)

    # Title
    draw.text((5, 0), title)

    return result


#import ephem
import math
import datetime
def OverlayDoppler(initial_img, fft_size, sample_rate):
    img = Image.new('RGB', initial_img.size)
    img.paste(initial_img, (0, 0))

    obs = ephem.Observer()
    obs.lat = math.radians(47.355063)
    obs.lon = math.radians(8.559047)
    obs.elevation = 418.0

    body = ephem.readtle(
        'FITSAT-1 (NIWAKA)',
        '1 38853U 98067CP  12295.88729863  .00054031  00000-0  85976-3 0   180',
        '2 38853  51.6495 211.8234 0014897 192.8078 167.2604 15.52629142  2661')

    begin_time = ephem.Date(datetime.datetime.utcfromtimestamp(1350666289))

    xpad = 50
    ypad = 40
    for i in xrange(img.size[1] - 2*ypad):
        t_offset = i * fft_size / float(sample_rate)
        obs.date = ephem.Date(begin_time + ephem.second * t_offset)
        body.compute(obs)

        deltaf = -body.range_velocity * 437.25e6 / 2.99e8
        dx = deltaf * fft_size / sample_rate  - 4*121 - 10
        x = xpad + fft_size/2 + int(dx)
        y = ypad + i
        img.putpixel((x+10, y), (0, 255, 0))
        img.putpixel((x+11, y), (0, 255, 0))
        #img.putpixel((x-10, y), (0, 255, 0))
        #img.putpixel((x-11, y), (0, 255, 0))

    return img


if __name__ == '__main__':
    sample_rate = int(sys.argv[2])
    title = sys.argv[5]
    img = Spectrogram(sys.argv[1], sample_rate, sys.argv[3], title)
    #img = OverlayDoppler(img, FFT_SIZE, sample_rate)
    img.save(sys.argv[4])
