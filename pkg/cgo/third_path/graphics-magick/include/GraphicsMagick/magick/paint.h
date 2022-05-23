/*
  Copyright (C) 2003 GraphicsMagick Group
  Copyright (C) 2002 ImageMagick Studio
  Copyright 1991-1999 E. I. du Pont de Nemours and Company

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  ImageMagick Image Paint Methods.
*/
#ifndef _MAGICK_PAINT_H
#define _MAGICK_PAINT_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif

#include "magick/render.h"

extern MagickExport unsigned int
  ColorFloodfillImage(Image *,const DrawInfo *,const PixelPacket,const long,
    const long,const PaintMethod),
  MatteFloodfillImage(Image *,const PixelPacket,const unsigned int,const long,
    const long,const PaintMethod);

extern MagickExport unsigned int
  OpaqueImage(Image *,const PixelPacket,const PixelPacket),
  TransparentImage(Image *,const PixelPacket,const unsigned int);

#if defined(__cplusplus) || defined(c_plusplus)
}
#endif

#endif

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 2
 * fill-column: 78
 * End:
 */
