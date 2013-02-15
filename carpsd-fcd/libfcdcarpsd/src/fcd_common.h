#ifndef FCD_COMMON_H
#define FCD_COMMON_H

/** \brief FCD mode enumeration. */
typedef enum {
    FCD_MODE_NONE,  /*!< No FCD detected. */
    FCD_MODE_BL,    /*!< FCD present in bootloader mode. */
    FCD_MODE_APP    /*!< FCD present in application mode. */
} FCD_MODE_ENUM; // The current mode of the FCD: none inserted, in bootloader mode or in normal application mode

#endif  // FCD_COMMON_H
