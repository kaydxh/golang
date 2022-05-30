/*
  Copyright (C) 2003 - 2009 GraphicsMagick Group
  Copyright (C) 2002 ImageMagick Studio
  Copyright 1991-1999 E. I. du Pont de Nemours and Company

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  GraphicsMagick Image Averaging Methods.
*/
#ifndef _MAGICK_AVERAGE_H
#define _MAGICK_AVERAGE_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif

/*
  Include declarations.
*/
#include "magick/image.h"
#include "magick/error.h"

extern MagickExport Image
  *AverageImages(const Image *,ExceptionInfo *);

#if defined(__cplusplus) || defined(c_plusplus)
}
#endif

#endif /* _MAGICK_AVERAGE_H */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 2
 * fill-column: 78
 * End:
 */
