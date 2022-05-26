/*
  Copyright (C) 2003 GraphicsMagick Group
  Copyright (C) 2002 ImageMagick Studio

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  ImageMagick Decorate Methods.
*/
#ifndef _MAGICK_DECORATE_H
#define _MAGICK_DECORATE_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif /* defined(__cplusplus) || defined(c_plusplus) */

extern MagickExport Image
  *BorderImage(const Image *,const RectangleInfo *,ExceptionInfo *),
  *FrameImage(const Image *,const FrameInfo *,ExceptionInfo *);

MagickExport unsigned int
  RaiseImage(Image *,const RectangleInfo *,const int);

#if defined(__cplusplus) || defined(c_plusplus)
}
#endif /* defined(__cplusplus) || defined(c_plusplus) */

#endif /* _MAGICK_DECORATE_H */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 2
 * fill-column: 78
 * End:
 */
