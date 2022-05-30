/*
  Copyright (C) 2003-2012 GraphicsMagick Group
  Copyright (C) 2002 ImageMagick Studio

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  ImageMagick Image Shear Methods.
*/
#ifndef _MAGICK_SHEAR_H
#define _MAGICK_SHEAR_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif /* defined(__cplusplus) || defined(c_plusplus) */

extern MagickExport Image
  *AffineTransformImage(const Image *,const AffineMatrix *,ExceptionInfo *),
  *AutoOrientImage(const Image *image,const OrientationType current_orientation,
                   ExceptionInfo *exception),
  *RotateImage(const Image *,const double,ExceptionInfo *),
  *ShearImage(const Image *,const double,const double,ExceptionInfo *);

#if defined(__cplusplus) || defined(c_plusplus)
}
#endif /* defined(__cplusplus) || defined(c_plusplus) */

#endif /* _MAGICK_SHEAR_H */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 2
 * fill-column: 78
 * End:
 */
