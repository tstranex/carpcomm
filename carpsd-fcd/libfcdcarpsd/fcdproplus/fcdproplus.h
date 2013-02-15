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

#ifndef FCD_H
#define FCD_H 1


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
#include "fcdproplushidcmd.h"


/** \brief FCD mode enumeration. */
typedef enum {
    FCD_MODE_NONE,  /*!< No FCD detected. */
    FCD_MODE_BL,    /*!< FCD present in bootloader mode. */
    FCD_MODE_APP    /*!< FCD present in application mode. */
} FCD_MODE_ENUM; // The current mode of the FCD: no FCD, in bootloader mode or in normal application mode

#ifdef __cplusplus
extern "C" {
#endif

/* Application functions */
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusGetMode(void);
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusGetFwVerStr(char *str);
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusAppReset(void);

EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusAppSetFreqKhz(int nFreq);
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusAppSetFreq(unsigned int uFreq, unsigned int *rFreq);
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusAppGetFreq(unsigned int *rFreq);

EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusAppSetLna(char enabled);
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusAppGetLna(char *enabled);

EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusAppSetRfFilter(tuner_rf_filter_t filter);
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusAppGetRfFilter(tuner_rf_filter_t *filter);

EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusAppSetMixerGain(char enabled);
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusAppGetMixerGain(char *enabled);

EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusAppSetIfGain(unsigned char gain);
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusAppGetIfGain(unsigned char *gain);

EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusAppSetIfFilter(tuner_if_filter_t filter);
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusAppGetIfFilter(tuner_if_filter_t *filter);

EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusAppSetBiasTee(char enabled);
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProPlusAppGetBiasTee(char *enabled);

#ifdef __cplusplus
}
#endif

#endif // FCD_H
