#!/usr/bin/env python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

"""Capture raw IQ data from a USRP."""

from gnuradio import gr
from gnuradio import uhd
import sys

class top_block(gr.top_block):
    def __init__(self, device_address, center_freq, sample_rate, output_path):
        gr.top_block.__init__(self)

        self.usrp_source = uhd.usrp_source(
            device_addr=device_address,
            stream_args=uhd.stream_args(
                cpu_format="sc16",
                channels=range(1)))
        self.usrp_source.set_samp_rate(sample_rate)
        self.usrp_source.set_center_freq(center_freq, 0)
        self.usrp_source.set_gain(0, 0)
        self.file_sink = gr.file_sink(gr.sizeof_int, output_path)

        self.connect((self.usrp_source, 0), (self.file_sink, 0))

if __name__ == '__main__':
    tb = top_block(
        sys.argv[1], float(sys.argv[2]), float(sys.argv[3]), sys.argv[4])
    tb.run()
