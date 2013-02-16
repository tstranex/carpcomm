// carpsd-fcd: command-line interface for the FCD.
// Copyright (C) 2013 Timothy Stranex <tstranex@carpcomm.com>

#include "fcdcarpsd.h"
#include <iostream>
#include <string>
#include <cstdlib>

using std::cout;
using std::cerr;
using std::endl;
using std::string;

static int PrintDeviceInfo(int device_index) {
  FCD* dev = FCDOpen(device_index);
  if (!dev) {
    cerr << "Error opening device " << device_index << endl;
    return 1;
  }

  cout << "Device: " << device_index << endl;
  cout << "Type: " << FCDGetType(dev) << endl;
  cout << "Firmware Version: " << FCDGetFirmwareVersion(dev) << endl;
  cout << "Frequency [Hz]: " << FCDGetFreqHz(dev) << endl;

  FCDClose(dev);

  return 0;
}

static int PrintInfo(int device_index) {
  if (device_index >= 0) {
    return PrintDeviceInfo(device_index);
  }

  int num = FCDCountDevices();
  for (int i = 0; i < num; i++) {
    if (!PrintDeviceInfo(i)) {
      continue;
    }
    cout << endl;
  }
  return 0;
}

static int SetFreqHz(int device_index, const char* freq_hz_str) {
  if (device_index < 0) {
    cerr << "--device is missing." << endl;
    return 1;
  }

  long long freq_hz = atoll(freq_hz_str);

  FCD* dev = FCDOpen(device_index);
  if (!dev) {
    cerr << "Error opening device." << endl;
    return 1;
  }

  long long new_freq_hz = FCDSetFreqHz(dev, freq_hz);
  FCDClose(dev);
  if (!new_freq_hz) {
    cerr << "Error setting frequency." << endl;
    return 1;
  }

  cout << "Frequency [Hz]: " << new_freq_hz << endl;
  return 0;
}

static void PrintHelp() {
  cout << "carpsd-fcd: Command-line FUNcube Dongle controller" << endl;
  cout << "Examples:" << endl;
  cout << "$ carpsd-fcd [--device 1]" << endl;
  cout << "    Prints FCD status info." << endl;
  cout << "$ carpsd-fcd --device 1 --set_freq_hz 437505000" << endl;
  cout << "    Tunes to the given frequency." << endl;
}

int main(int argc, char* argv[]) {
  int device_index = -1;

  for (int i = 1; i < argc; i++) {
    string arg = argv[i];

    if (arg == "--help") {
      PrintHelp();
      return 0;

    } else if (arg == "--device") {
      if (i + 1 < argc) {
	device_index = atoi(argv[i+1]);
	i++;
      } else {
	cerr << "Missing argument to --device." << endl;
	return 1;
      }

    } else if (arg == "--set_freq_hz") {
      if (i + 1 < argc) {
	return SetFreqHz(device_index, argv[i+1]);
      } else {
	cerr << "Missing argument to --set_freq_hz." << endl;
	return 1;
      }

    } else {
      cerr << "Unrecognised option '" << arg << "'." << endl;
      return 1;
    }
  }

  return PrintInfo(device_index);
}
