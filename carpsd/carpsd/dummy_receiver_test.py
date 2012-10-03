#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

import config
import dummy_receiver
import unittest
import os.path
#import Image


class DummyReceiverTest(unittest.TestCase):
    def Create(self):
        conf = config.GetDefaultConfig()
        return dummy_receiver.DummyReceiver(conf)

    def testBasic(self):
        r = self.Create()

        self.assertTrue(r.SetHardwareTunerHz(123.0))
        self.assertEqual(123.0, r.GetHardwareTunerHz())

        self.assertFalse(r.IsStarted())

        self.assertTrue(r.Start(None))
        self.assertTrue(r.IsStarted())

        r.Stop()
        self.assertFalse(r.IsStarted())

    def testWaterfallImage(self):
        r = self.Create()

        self.assertFalse(r.IsStarted())
        self.assertEqual(None, r.WaterfallImage())

        self.assertTrue(r.Start(None))

        img = r.WaterfallImage()
        #self.assertTrue(isinstance(img, Image.Image))
        self.assertTrue(img is not None)


if __name__ == '__main__':
    unittest.main()
