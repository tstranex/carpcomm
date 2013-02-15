/***************************************************************************
 *  This file is part of Qthid.
 *
 *  Copyright (C) 2010       Howard Long, G6LVB
 *  Copyright (C) 2011       Mario Lorenz, DL5MLO
 *  Copyright (C) 2011-2012  Alexandru Csete, OZ9AEC
 *
 *  Qthid is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  Qthid is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with Qthid.  If not, see <http://www.gnu.org/licenses/>.
 *
 ***************************************************************************/

#ifndef FCD_PRO_PLUS_H
#define FCD_PRO_PLUS_H 1


#ifdef FCD
#define EXTERN
#define ASSIGN(x) =x
#else
#define EXTERN extern
#define ASSIGN(x)
#endif

#ifdef _WIN32
#define FCD_API_EXPORT __declspec(dllexport)
#define FCD_API_CALL  _stdcall
#else
#define FCD_API_EXPORT
#define FCD_API_CALL
#endif

#include <inttypes.h>
#include "fcd_pro_plus_hidcmd.h"


/** \brief FCD mode enumeration. */
typedef enum {
    FCD_MODE_NONE,  /*!< No FCD detected. */
    FCD_MODE_BL,    /*!< FCD present in bootloader mode. */
    FCD_MODE_APP    /*!< FCD present in application mode. */
} FCD_MODE_ENUM; // The current mode of the FCD: no FCD, in bootloader mode or in normal application mode

#ifdef __cplusplus
extern "C" {
#endif

struct hid_device;

hid_device* fcdProPlusOpen();
void fcdProPlusClose(hid_device* dev);

/* Application functions */
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusGetMode(hid_device* phd);
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusGetFwVerStr(hid_device* phd, char *str);
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusAppReset(hid_device* phd);

EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusAppSetFreqKhz(hid_device* phd, int nFreq);
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusAppSetFreq(hid_device* phd, unsigned int uFreq, unsigned int *rFreq);
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusAppGetFreq(hid_device* phd, unsigned int *rFreq);

EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusAppSetLna(hid_device* phd, char enabled);
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusAppGetLna(hid_device* phd, char *enabled);

EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusAppSetRfFilter(hid_device* phd, tuner_rf_filter_t filter);
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusAppGetRfFilter(hid_device* phd, tuner_rf_filter_t *filter);

EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusAppSetMixerGain(hid_device* phd, char enabled);
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusAppGetMixerGain(hid_device* phd, char *enabled);

EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusAppSetIfGain(hid_device* phd, unsigned char gain);
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusAppGetIfGain(hid_device* phd, unsigned char *gain);

EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusAppSetIfFilter(hid_device* phd, tuner_if_filter_t filter);
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusAppGetIfFilter(hid_device* phd, tuner_if_filter_t *filter);

EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusAppSetBiasTee(hid_device* phd, char enabled);
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusAppGetBiasTee(hid_device* phd, char *enabled);

#ifdef __cplusplus
}
#endif

#endif // FCD_PRO_PLUS_H
