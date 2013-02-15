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

// Returns NULL on error.
// FCDClose must be called later.
extern FCD* FCDOpen(int index);

extern void FCDClose(FCD* dev);

// Returns 0 on error.
extern long long FCDGetFreqHz(FCD* dev);
extern long long FCDSetFreqHz(FCD* dev, long long freq_hz);

// Returned pointer is owned by us and must not be freed.
// Returns NULL on error.
extern const char* FCDGetType(FCD* dev);

// Returned pointer is owned by caller and must be freed.
// Returns NULL on error.
extern const char* FCDGetFirmwareVersion(FCD* dev);


#ifdef __cplusplus
}  // extern "C"
#endif

#endif  // FCDCARPSD_H
