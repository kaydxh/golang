/*
  Copyright (C) 2003 - 2020 GraphicsMagick Group
  Copyright (C) 2002 ImageMagick Studio

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  GraphicsMagick Image Analysis Methods.
*/
#ifndef _MAGICK_ANALYZE_H
#define _MAGICK_ANALYZE_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif /* defined(__cplusplus) || defined(c_plusplus) */

#include "magick/image.h"

typedef struct _ImageCharacteristics
{
  MagickBool
    cmyk,               /* CMYK(A) image */
    grayscale,          /* Grayscale image */
    monochrome,         /* Black/white image */
    opaque,             /* Opaque image */
    palette;            /* Colormapped image */
} ImageCharacteristics;

/* Functions which return unsigned int to indicate operation pass/fail */
extern MagickExport MagickPassFail
  GetImageCharacteristics(const Image *image,ImageCharacteristics *characteristics,
                          const MagickBool optimize,ExceptionInfo *exception);

extern MagickExport unsigned long
  GetImageDepth(const Image *,ExceptionInfo *);

extern MagickExport MagickBool
  IsGrayImage(const Image *image,ExceptionInfo *exception),
  IsMonochromeImage(const Image *image,ExceptionInfo *exception),
  IsOpaqueImage(const Image *image,ExceptionInfo *exception);

extern MagickExport ImageType
  GetImageType(const Image *,ExceptionInfo *);

extern MagickExport RectangleInfo
  GetImageBoundingBox(const Image *,ExceptionInfo *exception);

#if defined(__cplusplus) || defined(c_plusplus)
}
#endif /* defined(__cplusplus) || defined(c_plusplus) */

#endif /* _MAGICK_ANALYZE_H */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 2
 * fill-column: 78
 * End:
 */
