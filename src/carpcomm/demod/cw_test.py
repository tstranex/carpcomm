#!/usr/bin/python

# Author: Timothy Stranex <tstranex@carpcomm.com>
# Copyright 2013 Timothy Stranex

#!/usr/bin/python

import unittest
import subprocess


class CWTest(unittest.TestCase):

    def decodeCW(self, input_path, satellite_id):
        p = subprocess.Popen(
            ['src/carpcomm/demod/cw_decode.sh',
             input_path,
             satellite_id],
            stdout=subprocess.PIPE)
        return p.communicate()[0].strip()

    def testSwisscubePart0(self):
        self.assertEquals(
            'HB9EG/1',
            self.decodeCW(
                'testdata/carpcomm/demod/swisscube_part0.complex64',
                'swisscube'))

    def testSwisscubePart2(self):
        self.assertEquals(
            'UVAEVU6',
            self.decodeCW(
                'testdata/carpcomm/demod/swisscube_part2.complex64',
                'swisscube'))

    def testSwisscubePart3(self):
        self.assertEquals(
            'VTTTTVTVB',
            self.decodeCW(
                'testdata/carpcomm/demod/swisscube_part3.complex64',
                'swisscube'))

    def testSwisscubeAll(self):
        self.assertEquals(
            'UVAEVU6\nVTTTTVTVB\nHB9EG/1',
            self.decodeCW(
                'testdata/carpcomm/demod/swisscube_all.complex64',
                'swisscube'))

    def testMasat1(self):
        self.assertEquals(
            'HA5MASAT\n408\n05',
            self.decodeCW(
                'testdata/carpcomm/demod/masat1_cut.complex64',
                'masat1'))

    def testHoryu2(self):
        self.assertEquals(
            'JG6YBWHORYU\nD1BABCB589E0FE2',
            self.decodeCW(
                'testdata/carpcomm/demod/horyu2.complex64',
                'horyu2'))


if __name__ == '__main__':
    unittest.main()
