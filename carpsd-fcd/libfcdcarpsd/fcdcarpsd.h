/* libfcdcarpsd - FUNcube Dongle controller library.
 *
 * This library provides a uniform way to access both FCD Pro and FCD Pro+
 * devices.
 *
 * Copyright 2013 Timothy Stranex <tstranex@carpcomm.com>
 *
 * libfcdcarpsd is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * libfcdcarpsd is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with libfcdcarpsd.  If not, see <http://www.gnu.org/licenses/>.
 */

#ifndef FCDCARPSD_H
#define FCDCARPSD_H

struct FCD_;
typedef struct FCD_ FCD;

#ifdef __cplusplus
extern "C" {
#endif


// Returns the number of connected devices.
extern int FCDCountDevices();

// Open the FCD device with the given index (in the range 0 to FCDCountDevices).
// Returns NULL on error.
// FCDClose must be called to destroy the FCD* object.
extern FCD* FCDOpen(int index);

// Close the device. The pointer shouldn't be used afterwards.
extern void FCDClose(FCD* dev);

// Get the frequency in Hz.
// Returns 0 on error.
extern long long FCDGetFreqHz(FCD* dev);

// Set the frequency in Hz and return the frequency that was actually set.
// Returns 0 on error.
extern long long FCDSetFreqHz(FCD* dev, long long freq_hz);

// Get the device type e.g. FCD Pro or FCD Pro+.
// The returned pointer should not be freed.
// Returns NULL on error.
extern const char* FCDGetType(FCD* dev);

// Get the device firmware version string.
// The returned pointer must be freed by the caller.
// Returns NULL on error.
extern const char* FCDGetFirmwareVersion(FCD* dev);


#ifdef __cplusplus
}  // extern "C"
#endif

#endif  // FCDCARPSD_H
