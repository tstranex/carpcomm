#!/usr/bin/python

# Copyright 2012 Carpcomm GmbH
# Author: Timothy Stranex <tstranex@carpcomm.com>

import time
import json
import cStringIO
import base64
import signalling


OK = 200
OK_NO_DATA = 204


class ReceiverHandlers(object):

    def __init__(self, receiver):
        self.r = receiver

    def RegisterHandlers(self, client):
        client.RegisterHandler('/ReceiverGetInfo', self.ReceiverGetInfo)
        client.RegisterHandler('/ReceiverGetState', self.ReceiverGetState)
        client.RegisterHandler('/ReceiverStart', self.ReceiverStart)
        client.RegisterHandler('/ReceiverStop', self.ReceiverStop)
        client.RegisterHandler('/ReceiverSetFrequency',
                               self.ReceiverSetFrequency)
        client.RegisterHandler('/ReceiverWaterfallPNG',
                               self.ReceiverWaterfallPNG)

    def ReceiverGetInfo(self, params):
        return OK, 'application/json', json.dumps(self.r.GetInfoDict())

    def ReceiverGetState(self, params):
        return OK, 'application/json', json.dumps(self.r.GetStateDict())

    def ReceiverStart(self, params):
        signalling.Get().SignalReceiverStart()
        stream_url = params['stream_url'][0]
        return self.r.Start(stream_url)

    def ReceiverStop(self, params):
        signalling.Get().SignalReceiverStop()
        self.r.Stop()
        return True

    def ReceiverSetFrequency(self, params):
        if 'hz' not in params:
            return False
        try:
            freq_hz = int(params['hz'][0])
        except ValueError:
            return False
        return self.r.SetHardwareTunerHz(freq_hz)

    def ReceiverWaterfallPNG(self, params):
        img = self.r.WaterfallImage()
        if img:
            data = cStringIO.StringIO()
            img.save(data, 'PNG')
            return OK, 'image/png', data.getvalue()
        else:
            return OK_NO_DATA, 'image/png', ''


class MotorHandlers(object):

    def __init__(self, motor):
        self.m = motor
    
    def RegisterHandlers(self, client):
        client.RegisterHandler('/MotorGetInfo', self.MotorGetInfo)
        client.RegisterHandler('/MotorGetState', self.MotorGetState)
        client.RegisterHandler('/MotorStart', self.MotorStart)
        client.RegisterHandler('/MotorStop', self.MotorStop)

    def MotorGetInfo(self, params):
        return OK, 'application/json', json.dumps(self.m.GetInfoDict())

    def MotorGetState(self, params):
        return OK, 'application/json', json.dumps(self.m.GetStateDict())

    def MotorStart(self, params):
        program = json.loads(params['program'][0])
        return self.m.Start(program)

    def MotorStop(self, params):
        return self.m.Stop()


class TNCHandlers(object):

    def __init__(self, tnc):
        self.t = tnc
    
    def RegisterHandlers(self, client):
        client.RegisterHandler('/TNCStart', self.TNCStart)
        client.RegisterHandler('/TNCStop', self.TNCStop)
        client.RegisterHandler('/TNCGetLatestFrames', self.TNCGetLatestFrames)
        client.RegisterHandler('/TNCGetState', self.TNCGetState)

    def TNCStart(self, params):
        api_host = params['api_host'][0]
        try:
            api_port = int(params['api_port'][0])
        except ValueError:
            return False
        satellite_id = params['satellite_id'][0]
        return self.t.Start(api_host, api_port, satellite_id)

    def TNCStop(self, params):
        return self.t.Stop()

    def TNCGetLatestFrames(self, params):
        # This is mostly for interactive console use.
        ok, frames = self.t.GetLatestFrames()
        if ok:
            b64_frames = [base64.b64encode(f) for f in frames]
            return OK, 'application/json', json.dumps(b64_frames)
        else:
            return False

    def TNCGetState(self, params):
        return OK, 'application/json', json.dumps(self.t.GetStateDict())
