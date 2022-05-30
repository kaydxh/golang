/*
  Copyright (C) 2003 GraphicsMagick Group
  Copyright (C) 2002 ImageMagick Studio

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  ImageMagick Image Montage Methods.
*/
#ifndef _MAGICK_MONTAGE_H
#define _MAGICK_MONTAGE_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif /* defined(__cplusplus) || defined(c_plusplus) */

extern MagickExport Image
  *MontageImages(const Image *,const MontageInfo *,ExceptionInfo *);

extern MagickExport MontageInfo
  *CloneMontageInfo(const ImageInfo *,const MontageInfo *);

extern MagickExport void
  DestroyMontageInfo(MontageInfo *),
  GetMontageInfo(const ImageInfo *,MontageInfo *);

#if defined(__cplusplus) || defined(c_plusplus)
}
#endif /* defined(__cplusplus) || defined(c_plusplus) */

#endif /* _MAGICK_MONTAGE_H */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 2
 * fill-column: 78
 * End:
 */
