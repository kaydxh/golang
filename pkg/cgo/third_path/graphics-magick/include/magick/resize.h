/*
  Copyright (C) 2003 - 2012 GraphicsMagick Group
  Copyright (C) 2002 ImageMagick Studio

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  ImageMagick Image Resize Methods.
*/
#ifndef _MAGICK_RESIZE_H
#define _MAGICK_RESIZE_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif /* defined(__cplusplus) || defined(c_plusplus) */

#define DefaultResizeFilter LanczosFilter
#define DefaultThumbnailFilter BoxFilter

extern MagickExport Image
  *MagnifyImage(const Image *,ExceptionInfo *),
  *MinifyImage(const Image *,ExceptionInfo *),
  *ResizeImage(const Image *,const unsigned long,const unsigned long,
     const FilterTypes,const double,ExceptionInfo *),
  *SampleImage(const Image *,const unsigned long,const unsigned long,
   ExceptionInfo *),
  *ScaleImage(const Image *,const unsigned long,const unsigned long,
     ExceptionInfo *),
  *ThumbnailImage(const Image *,const unsigned long,const unsigned long,
   ExceptionInfo *),
  *ZoomImage(const Image *,const unsigned long,const unsigned long,
     ExceptionInfo *);

#if defined(__cplusplus) || defined(c_plusplus)
}
#endif /* defined(__cplusplus) || defined(c_plusplus) */

#endif /* _MAGICK_RESIZE_H */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 2
 * fill-column: 78
 * End:
 */
