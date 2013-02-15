#include "fcdcarpsd.h"
#include "fcdpro/fcd_pro.h"
#include "fcdproplus/fcd_pro_plus.h"

#include <stdlib.h>

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

int FCDGetFreqHz(FCD* dev) {
  if (dev->driver == DRIVER_FCD_PRO) {
    unsigned char buf[4];
    int freq = 0;
    fcdProAppGetParam(dev->dev, FCD_CMD_APP_GET_FREQ_HZ, buf, 4);
    freq += buf[0];
    freq += buf[1] << 8;
    freq += buf[2] << 16;
    freq += buf[3] << 24;
    return freq;

  } else if (dev->driver == DRIVER_FCD_PRO_PLUS) {
    unsigned int freq_hz;
    fcdProPlusAppGetFreq(dev->dev, &freq_hz);
    return freq_hz;

  } else {
    return -1;
  }
}
