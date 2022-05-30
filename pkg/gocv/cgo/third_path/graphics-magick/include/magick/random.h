/*
  Copyright (C) 2009, 2014 GraphicsMagick Group

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  Random number generator (semi-public interfaces).

  Currently based on George Marsaglia's multiply-with-carry generator.
  This is a k=2 generator with a period >2^60.
*/

#ifndef _MAGICK_RANDOM_H
#define _MAGICK_RANDOM_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif /* defined(__cplusplus) || defined(c_plusplus) */

typedef struct _MagickRandomKernel
{
  magick_uint32_t z;
  magick_uint32_t w;
} MagickRandomKernel;

#define MAGICK_RANDOM_MAX 4294967295

  /*
    Generate a random integer value (0 - MAGICK_RANDOM_MAX)
  */
  MagickExport magick_uint32_t MagickRandomInteger(void);

  /*
    Generate a random double value (0.0 - 1.0)
  */
  MagickExport double MagickRandomReal(void);

#if defined(MAGICK_IMPLEMENTATION)
#include "magick/random-private.h"
#endif /* defined(MAGICK_IMPLEMENTATION) */

#if defined(__cplusplus) || defined(c_plusplus)
}
#endif /* defined(__cplusplus) || defined(c_plusplus) */

#endif /* ifndef _MAGICK_RANDOM_H */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 2
 * fill-column: 78
 * End:
 */
