/*
  Copyright (C) 2003 GraphicsMagick Group
  Copyright (C) 2002 ImageMagick Studio

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  ImageMagick Image Enhancement Methods.
*/
#ifndef _MAGICK_ENHANCE_H
#define _MAGICK_ENHANCE_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif  /* defined(__cplusplus) || defined(c_plusplus) */

extern MagickExport MagickPassFail
  ContrastImage(Image *,const unsigned int),
  EqualizeImage(Image *),
  GammaImage(Image *,const char *),
  LevelImage(Image *,const char *),
  LevelImageChannel(Image *,const ChannelType,const double,const double,
    const double),
  ModulateImage(Image *,const char *),
  NegateImage(Image *,const unsigned int),
  NormalizeImage(Image *);

#if defined(__cplusplus) || defined(c_plusplus)
}
#endif /* defined(__cplusplus) || defined(c_plusplus) */

#endif /* _MAGICK_ENHANCE_H */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 2
 * fill-column: 78
 * End:
 */
