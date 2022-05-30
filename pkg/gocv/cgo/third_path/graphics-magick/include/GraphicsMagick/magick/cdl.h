/*
 * Copyright (C) 2009 GraphicsMagick Group
 *
 * American Society of Cinematographers Color Decision List (ASC-CDL)
 * implementation.
 *
 * Original implementation by Clément Follet.
 */

#ifndef _MAGICK_CDL_H
#define _MAGICK_CDL_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif  /* defined(__cplusplus) || defined(c_plusplus) */

extern MagickExport MagickPassFail
  CdlImage(Image *image,const char *cdl);

#if defined(__cplusplus) || defined(c_plusplus)
}
#endif /* defined(__cplusplus) || defined(c_plusplus) */

#endif /* _MAGICK_CDL_H */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 2
 * fill-column: 78
 * End:
 */
