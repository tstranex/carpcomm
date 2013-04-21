#!/usr/bin/env python

from gnuradio import gr
from gnuradio.gr import firdes
import sys

FREQ_SPACE = 2200

class top_block(gr.top_block):

	def __init__(self, input_path, sample_rate, output_path):
		gr.top_block.__init__(self)

		self.source = gr.file_source(
			gr.sizeof_gr_complex, input_path, False)
		self.lowpass_and_decimate = gr.fir_filter_ccf(
			8, firdes.low_pass(1, sample_rate, FREQ_SPACE, 100))
		self.lsb_tune = gr.freq_xlating_fir_filter_ccc(
			1, (1,), -FREQ_SPACE, sample_rate/8)
		self.boost_volume = gr.multiply_const_vcc((10, ))
		self.complex_to_real = gr.complex_to_real(1)
		self.sink = gr.wavfile_sink(output_path, 1, sample_rate/8, 16)

		self.connect((self.source, 0), (self.lowpass_and_decimate, 0))
		self.connect((self.lowpass_and_decimate, 0), (self.lsb_tune, 0))
		self.connect((self.lsb_tune, 0), (self.boost_volume, 0))
		self.connect((self.boost_volume, 0), (self.complex_to_real, 0))
		self.connect((self.complex_to_real, 0), (self.sink, 0))


if __name__ == '__main__':
	sample_rate = int(float(sys.argv[2]))
	tb = top_block(sys.argv[1], sample_rate, sys.argv[3])
	tb.run()
