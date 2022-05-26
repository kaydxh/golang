/*
  Copyright (C) 2003 - 2009 GraphicsMagick Group
  Copyright (C) 2002 ImageMagick Studio
  Copyright 1991-1999 E. I. du Pont de Nemours and Company

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  GraphicsMagick Plasma Methods.
*/
#ifndef _MAGICK_PLASMA_H
#define _MAGICK_PLASMA_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif

/*
  Include declarations.
*/
#include "magick/image.h"

extern MagickExport MagickPassFail
PlasmaImage(Image *image,const SegmentInfo *segment,
            unsigned long attenuate,unsigned long depth);

#if defined(__cplusplus) || defined(c_plusplus)
}
#endif

#endif /* _MAGICK_PLASMA_H */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 2
 * fill-column: 78
 * End:
 */
