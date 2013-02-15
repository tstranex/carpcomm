/***************************************************************************
 *  This file is part of Qthid.
 *
 *  Copyright (C) 2010  Howard Long, G6LVB
 *  CopyRight (C) 2011  Alexandru Csete, OZ9AEC
 *                      Mario Lorenz, DL5MLO
 *  Copyright (C) 2013  Timothy Stranex, HB9FFH
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

#ifndef FCD_PRO_H
#define FCD_PRO_H 1


#ifdef FCD
#define EXTERN
#define ASSIGN (x) =x
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
#include "hidapi.h"
#include "fcd_common.h"


/** \brief FCD capabilities that depend on both hardware and firmware. */
typedef struct {
    unsigned char hasBiasT;     /*!< Whether FCD has hardware bias tee (1=yes, 0=no) */
    unsigned char hasCellBlock; /*!< Whether FCD has cellular blocking. */
} FCD_CAPS_STRUCT;

#ifdef __cplusplus
extern "C" {
#endif

int fcdProCountDevices();
hid_device* fcdProOpen(int index);
void fcdProClose(hid_device* phd);

/* Application functions */
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProGetMode(hid_device* phd);
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProGetFwVerStr(hid_device* phd, char *str);
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProGetCaps(hid_device* phd, FCD_CAPS_STRUCT *fcd_caps);
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProGetCapsStr(hid_device* phd, char *caps_str);
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProAppReset(hid_device* phd);
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProAppSetFreqkHz(hid_device* phd, int nFreq);

EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProAppSetParam(hid_device* phd, uint8_t u8Cmd, uint8_t *pu8Data, uint8_t u8len);
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProAppGetParam(hid_device* phd, uint8_t u8Cmd, uint8_t *pu8Data, uint8_t u8len);


/* Bootloader functions */
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProBlReset(hid_device* phd);
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProBlErase(hid_device* phd);
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProBlWriteFirmware(hid_device* phd, char *pc, int64_t n64Size);
EXTERN FCD_API_EXPORT FCD_API_CALL FCD_MODE_ENUM fcdProBlVerifyFirmware(hid_device* phd, char *pc, int64_t n64Size);


#ifdef __cplusplus
}
#endif

#endif // FCD_PRO_H
