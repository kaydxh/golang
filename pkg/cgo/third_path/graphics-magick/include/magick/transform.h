/*
  Copyright (C) 2003 - 2010 GraphicsMagick Group
  Copyright (C) 2002 ImageMagick Studio

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  ImageMagick Image Transform Methods.
*/
#ifndef _MAGICK_TRANSFORM_H
#define _MAGICK_TRANSFORM_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif /* defined(__cplusplus) || defined(c_plusplus) */

extern MagickExport Image
  *ChopImage(const Image *image,const RectangleInfo *chop_info,ExceptionInfo *exception),
  *CoalesceImages(const Image *image,ExceptionInfo *exception),
  *CropImage(const Image *image,const RectangleInfo *geometry,ExceptionInfo *exception),
  *DeconstructImages(const Image *image,ExceptionInfo *exception),
  *ExtentImage(const Image *image,const RectangleInfo *geometry,ExceptionInfo *exception),
  *FlattenImages(const Image *image,ExceptionInfo *exception),
  *FlipImage(const Image *image,ExceptionInfo *exception),
  *FlopImage(const Image *image,ExceptionInfo *exception),
  *MosaicImages(const Image *image,ExceptionInfo *exception),
  *RollImage(const Image *image,const long x_offset,const long y_offset,ExceptionInfo *exception),
  *ShaveImage(const Image *image,const RectangleInfo *shave_info,ExceptionInfo *exception);

extern MagickExport MagickPassFail
  TransformImage(Image **,const char *,const char *);

#if defined(__cplusplus) || defined(c_plusplus)
}
#endif /* defined(__cplusplus) || defined(c_plusplus) */

#endif /* _MAGICK_TRANSFORM_H */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 2
 * fill-column: 78
 * End:
 */
