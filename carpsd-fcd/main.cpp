// Command-line interface for the FCD.
// Copyright (C) 2013 Timothy Stranex <tstranex@carpcomm.com>

#include "fcdcarpsd.h"
#include <iostream>
#include <string>
#include <cstdlib>

using std::cout;
using std::cerr;
using std::endl;
using std::string;

static int PrintInfo() {
  int num = FCDCountDevices();
  cout << "Number of devices: " << num << endl << endl;
  for (int i = 0; i < num; i++) {
    cout << "Device " << i << ":" << endl;

    FCD* dev = FCDOpen(i);
    if (!dev) {
      cerr << "Error opening device." << endl;
      continue;
    }

    cout << "Type: " << FCDGetType(dev) << endl;
    cout << "Firmware version: " << FCDGetFirmwareVersion(dev) << endl;
    cout << "Frequency: " << FCDGetFreqHz(dev) << endl;

    FCDClose(dev);
  }

  return 0;
}

static int SetFreqHz(int device_index, const char* freq_hz_str) {
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

  cout << "Frequency set to: " << new_freq_hz << endl;
  return 0;
}

static void PrintHelp() {
  cout << "carpsd-fcd: Command-line FUNcube Dongle controller" << endl;
  cout << "Examples:" << endl;
  cout << "$ carpsd-fcd" << endl;
  cout << "    Prints FCD status info." << endl;
  cout << "$ carpsd-fcd --device 1 --set_freq_hz 437505000" << endl;
  cout << "    Tunes to the given frequency." << endl;
}

int main(int argc, char* argv[]) {
  int device_index = 0;

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

  return PrintInfo();
}
