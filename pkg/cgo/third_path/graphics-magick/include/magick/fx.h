/*
  Copyright (C) 2003-2009 GraphicsMagick Group
  Copyright (C) 2002 ImageMagick Studio

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  GraphicsMagick Image FX Methods.
*/
#ifndef _MAGICK_FX_H
#define _MAGICK_FX_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif  /* defined(__cplusplus) || defined(c_plusplus) */

extern MagickExport Image
  *CharcoalImage(const Image *,const double,const double,ExceptionInfo *),
  *ColorizeImage(const Image *,const char *,const PixelPacket,ExceptionInfo *),
  *ImplodeImage(const Image *,const double,ExceptionInfo *),
  *MorphImages(const Image *,const unsigned long,ExceptionInfo *),
  *OilPaintImage(const Image *,const double,ExceptionInfo *),
  *SteganoImage(const Image *,const Image *,ExceptionInfo *),
  *StereoImage(const Image *,const Image *,ExceptionInfo *),
  *SwirlImage(const Image *,double,ExceptionInfo *),
  *WaveImage(const Image *,const double,const double,ExceptionInfo *);

extern MagickExport MagickPassFail
  ColorMatrixImage(Image *image,const unsigned int order,const double *matrix),
  SolarizeImage(Image *,const double);

#if defined(__cplusplus) || defined(c_plusplus)
}
#endif /* defined(__cplusplus) || defined(c_plusplus) */

#endif /* _MAGICK_FX_H */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 2
 * fill-column: 78
 * End:
 */
