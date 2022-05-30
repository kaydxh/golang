/*
  Copyright (C) 2009 GraphicsMagick Group

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  Hald CLUT implementation

*/
#ifndef _MAGICK_HCLUT_H
#define _MAGICK_HCLUT_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif  /* defined(__cplusplus) || defined(c_plusplus) */

extern MagickExport MagickPassFail
  HaldClutImage(Image *,const Image * clut);

#if defined(__cplusplus) || defined(c_plusplus)
}
#endif /* defined(__cplusplus) || defined(c_plusplus) */

#endif /* _MAGICK_HCLUT_H */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 2
 * fill-column: 78
 * End:
 */
