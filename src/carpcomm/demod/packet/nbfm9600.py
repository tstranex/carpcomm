#!/usr/bin/env python

from gnuradio import blks2
from gnuradio import gr
from gnuradio.gr import firdes
import sys

class top_block(gr.top_block):

	def __init__(self, input_path, sample_rate, output_path):
		gr.top_block.__init__(self)

		# We don't use the existing NBFM demodulator block because it
		# contains a lowpass output filter which is unsuitable for 9600
		# GMSK (it's designed for voice).

		self.source = gr.file_source(
			gr.sizeof_gr_complex*1, input_path, False)
		self.low_pass_filter = gr.fir_filter_ccf(4, firdes.low_pass(
			1, sample_rate, 15000, 100, firdes.WIN_HAMMING, 6.76))

		# High pass filter to remove the DC component. This is important
		# when the signal is near the SDR's local oscillator.
		# NOTE(tstranex): Disabled since we are now shifting the FCD
		# center frequency instead.
		#self.high_pass_filter = gr.fir_filter_ccf(1, firdes.high_pass(
		#	1, sample_rate/4, 100, 100, firdes.WIN_HAMMING, 6.76))

		self.quadrature_demod = gr.quadrature_demod_cf(
			sample_rate/4/(2*3.14*3000))
		self.fm_deemph = blks2.fm_deemph(fs=sample_rate/4, tau=75e-6)
		self.boost_volume = gr.multiply_const_vff((1.52, ))
		self.sink = gr.wavfile_sink(output_path, 1, sample_rate/4, 16)


		self.connect((self.source, 0), (self.low_pass_filter, 0))
		#self.connect((self.low_pass_filter, 0), (self.high_pass_filter, 0))
		#self.connect((self.high_pass_filter, 0), (self.quadrature_demod, 0))
		self.connect((self.low_pass_filter, 0), (self.quadrature_demod, 0))
		
		self.connect((self.quadrature_demod, 0), (self.fm_deemph, 0))
		self.connect((self.fm_deemph, 0), (self.boost_volume, 0))
		self.connect((self.boost_volume, 0), (self.sink, 0))


if __name__ == '__main__':
	sample_rate = int(float(sys.argv[2]))
	tb = top_block(sys.argv[1], sample_rate, sys.argv[3])
	tb.run()
