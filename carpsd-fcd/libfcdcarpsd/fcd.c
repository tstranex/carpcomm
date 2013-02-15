#include "fcdcarpsd.h"
#include "fcdpro/fcd_pro.h"
#include "fcdproplus/fcd_pro_plus.h"

#include <stdlib.h>
#include <string.h>

struct FCD_ {
  enum {
    DRIVER_FCD_PRO = 0,
    DRIVER_FCD_PRO_PLUS = 1,
  } driver;

  hid_device* dev;
};

int FCDCountDevices() {
  return fcdProCountDevices() + fcdProPlusCountDevices();
}

FCD* FCDOpen(int index) {
  int num_pro = fcdProCountDevices();
  int num_pro_plus = fcdProPlusCountDevices();

  if (index < num_pro) {
    hid_device* dev = fcdProOpen(index);
    if (!dev) {
      return NULL;
    }
    FCD* fcd = malloc(sizeof(FCD));
    fcd->driver = DRIVER_FCD_PRO;
    fcd->dev = dev;
    return fcd;

  } else if (index < num_pro + num_pro_plus) {

    hid_device* dev = fcdProPlusOpen(index - num_pro);
    if (!dev) {
      return NULL;
    }
    FCD* fcd = malloc(sizeof(FCD));
    fcd->driver = DRIVER_FCD_PRO_PLUS;
    fcd->dev = dev;
    return fcd;

  } else {
    return NULL;
  }
}

void FCDClose(FCD* dev) {
  if (dev->driver == DRIVER_FCD_PRO) {
    fcdProClose(dev->dev);
  } else if (dev->driver == DRIVER_FCD_PRO_PLUS) {
    fcdProPlusClose(dev->dev);
  }
}

long long FCDGetFreqHz(FCD* dev) {
  if (dev->driver == DRIVER_FCD_PRO) {
    unsigned char buf[4];
    unsigned int freq = 0;
    fcdProAppGetParam(dev->dev, FCD_CMD_APP_GET_FREQ_HZ, buf, 4);
    freq += buf[0];
    freq += buf[1] << 8;
    freq += buf[2] << 16;
    freq += buf[3] << 24;
    return freq;

  } else if (dev->driver == DRIVER_FCD_PRO_PLUS) {
    unsigned int freq;
    fcdProPlusAppGetFreq(dev->dev, &freq);
    return freq;

  } else {
    return 0;
  }
}

long long FCDSetFreqHz(FCD* dev, long long freq_hz) {
  if (dev->driver == DRIVER_FCD_PRO) {
    int freq_khz = (int)(freq_hz / 1000);
    if (fcdProAppSetFreqkHz(dev->dev, freq_khz) != FCD_MODE_APP) {
      return 0;
    }
    return FCDGetFreqHz(dev);

  } else if (dev->driver == DRIVER_FCD_PRO_PLUS) {
    unsigned int actual_freq_hz;
    if (fcdProPlusAppSetFreq(dev->dev, freq_hz, &actual_freq_hz)
	!= FCD_MODE_APP) {
      return 0;
    }
    return actual_freq_hz;

  } else {
    return 0;
  }
}

const char* FCDGetType(FCD* dev) {
  if (dev->driver == DRIVER_FCD_PRO) {
    return "FUNcube Dongle Pro";
  } else if (dev->driver == DRIVER_FCD_PRO_PLUS) {
    return "FUNcube Dongle Pro+";
  } else {
    return NULL;
  }
}

const char* FCDGetFirmwareVersion(FCD* dev) {
  char version[6];
  if (dev->driver == DRIVER_FCD_PRO) {
    fcdProGetFwVerStr(dev->dev, version);
    return strdup(version);
  } else if (dev->driver == DRIVER_FCD_PRO_PLUS) {
    fcdProPlusGetFwVerStr(dev->dev, version);
    return strdup(version);
  } else {
    return NULL;
  }
}
