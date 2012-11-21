#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

import serial_tnc
import testing

import unittest
import threading
import time

BASIC_KISS_FRAME = '\xc0\x00testframe123\xc0'


class KISSDecoderTest(unittest.TestCase):

    def testBasic(self):
        d = serial_tnc.KISSDecoder()
        d.Write(BASIC_KISS_FRAME)
        self.assertEquals(['testframe123'], d.ReadFrames())
        self.assertEquals([], d.ReadFrames())

    def testBasicWithoutPadding(self):
        d = serial_tnc.KISSDecoder()
        d.Write('\x00test1\xc0\x00test2\xc0')
        self.assertEquals(['test1', 'test2'], d.ReadFrames())
        self.assertEquals([], d.ReadFrames())

    def testStartAndEndFraming(self):
        d = serial_tnc.KISSDecoder()
        d.Write('\xc0\x00test1\xc0\xc0\x00test2\xc0')
        self.assertEquals(['test1', 'test2'], d.ReadFrames())
        self.assertEquals([], d.ReadFrames())

    def testIgnoreEmptyFrames(self):
        d = serial_tnc.KISSDecoder()
        d.Write('\xc0\x00\xc0\xc0\x00test\xc0\x00\xc0')
        self.assertEquals(['test'], d.ReadFrames())
        self.assertEquals([], d.ReadFrames())

    def testGarbage(self):
        d = serial_tnc.KISSDecoder()
        d.Write('abjkdfdf\xc0\x00test1\xc0dfdfdfddfdf\xc0\x00test2\xc0sss')
        self.assertEquals(['bjkdfdf', 'test1', 'fdfdfddfdf', 'test2'],
                          d.ReadFrames())
        self.assertEquals([], d.ReadFrames())

    def testEscapeChars(self):
        d = serial_tnc.KISSDecoder()
        d.Write('\xc0\x00test\xdb\xdcabc\xc0')
        self.assertEquals(['test\xc0abc'], d.ReadFrames())
        d.Write('\xc0\x00test\xdb\xddabc\xc0')
        self.assertEquals(['test\xdbabc'], d.ReadFrames())
        self.assertEquals([], d.ReadFrames())

    def testTransposeCharsInNormalMode(self):
        d = serial_tnc.KISSDecoder()
        d.Write('\xc0\x00test\xdc\xddabc\xc0')
        self.assertEquals(['test\xdc\xddabc'], d.ReadFrames())
        self.assertEquals([], d.ReadFrames())

    def testInvalidEscapeIsIgnored(self):
        d = serial_tnc.KISSDecoder()
        d.Write('\xc0\x00test\xdb?abc\xc0')
        self.assertEquals(['testabc'], d.ReadFrames())

    def testUnknownDataFlagIsIgnored(self):
        d = serial_tnc.KISSDecoder()
        d.Write('\xc0\x01test1\xc0 \xc0\xfftest2\xc0')
        self.assertEquals(['test1', 'test2'], d.ReadFrames())

    def testDoubleStart(self):
        d = serial_tnc.KISSDecoder()
        d.Write('\xc0\x00test1\xc0\x00\xc0\xc0\x00test2\xc0\xc0')
        self.assertEquals(['test1', 'test2'], d.ReadFrames())
        self.assertEquals([], d.ReadFrames())

    def testEscapeWithFEND(self):
        d = serial_tnc.KISSDecoder()
        d.Write('\xc0\x00test1\xdb\xc0')
        self.assertEquals(['test1'], d.ReadFrames())
        self.assertEquals([], d.ReadFrames())

    def testTruncateOversize(self):
        large_frame = '1234' + 't' * (2*serial_tnc.KISSDecoder.SIZE_LIMIT)
        data = '\x00test1\xc0\x00' + large_frame + '\xc0\x00test2\xc0'
        truncated_frame = large_frame[:serial_tnc.KISSDecoder.SIZE_LIMIT-1]

        d = serial_tnc.KISSDecoder()
        d.Write(data)
        self.assertEquals(['test1', truncated_frame, 'test2'], d.ReadFrames())
        self.assertEquals([], d.ReadFrames())


class _MockSerial:
    def __init__(self, buf):
        self.close_called = False
        self.buf = buf

    def read(self, n=1):
        if self.buf:
            n = min(n, len(self.buf))
            b = self.buf[:n]
            self.buf = self.buf[n:]
            return b
        else:
            time.sleep(serial_tnc.SERIAL_READ_TIMEOUT)
            return ''

    def inWaiting(self):
        return len(self.buf)

    def close(self):
        self.close_called = True
        

class SerialTNCTest(unittest.TestCase):

    def config(self):
        section = serial_tnc.SerialTNC.__name__
        conf = testing.GetConfigForTesting()
        conf.add_section(section)
        conf.set(section, 'device', '/dev/null')
        conf.set(section, 'baud', '9600')
        return conf

    def create(self):
        return serial_tnc.SerialTNC(self.config())

    def testStartStop(self):
        s = self.create()

        # Mock out the serial device.
        ms = _MockSerial(BASIC_KISS_FRAME)
        s._OpenSerial = lambda: ms

        self.assertTrue(s.Verify())
        self.assertEquals(s.GetStateDict(), {'started': False})
        self.assertTrue(s.Start('', 0, ''))
        self.assertEquals(s.GetStateDict(), {'started': True})
        self.assertTrue(s.Stop())
        self.assertEquals(s.GetStateDict(), {'started': False})
        self.assertTrue(ms.close_called)

    def testGetLatestFrames(self):
        s = self.create()

        # Mock out the serial device.
        ms = _MockSerial(BASIC_KISS_FRAME)
        s._OpenSerial = lambda: ms

        self.assertTrue(s.Verify())
        self.assertTrue(s.Start('', 0, ''))

        # Wait until the data has been read.
        # Surely there is a better way to do this...
        while ms.inWaiting() > 0:
            time.sleep(0.1)
        time.sleep(0.1)

        self.assertEquals((True, ['testframe123']), s.GetLatestFrames())

        self.assertTrue(s.Stop())
        self.assertEquals((False, []), s.GetLatestFrames())

    def testOutOfOrderStop(self):
        s = self.create()
        self.assertTrue(s.Stop())

    def testRTSCTS(self):
        conf = self.config()

        conf.set(serial_tnc.SerialTNC.__name__, 'rtscts', 'true')
        t = serial_tnc.SerialTNC(conf)
        self.assertTrue(t._rtscts)

        conf.set(serial_tnc.SerialTNC.__name__, 'rtscts', 'false')
        t = serial_tnc.SerialTNC(conf)
        self.assertFalse(t._rtscts)


if __name__ == '__main__':
    unittest.main()
