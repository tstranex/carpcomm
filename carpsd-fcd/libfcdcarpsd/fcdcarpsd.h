/* 2013 Timothy Stranex <tstranex@carpcomm.com>
 */

#ifndef FCDCARPSD_H
#define FCDCARPSD_H

struct FCD_;
typedef struct FCD_ FCD;

#ifdef __cplusplus
extern "C" {
#endif

extern int FCDCountDevices();
extern FCD* FCDOpen(int index);
extern void FCDClose(FCD* dev);
extern int FCDGetFreqHz(FCD* dev);

#ifdef __cplusplus
}  // extern "C"
#endif

#endif  // FCDCARPSD_H
