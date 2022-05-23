/*
  Copyright (C) 2003 - 2020 GraphicsMagick Group
  Copyright (C) 2002 ImageMagick Studio

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  GraphicsMagick Modules Methods.
*/
#ifndef _MAGICK_MODULE_H
#define _MAGICK_MODULE_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif

/*
  Modules declarations.
*/

extern MagickExport MagickPassFail
  ExecuteModuleProcess(const char *tag,Image **image,
                       const int argc,char **argv);


#if defined(MAGICK_IMPLEMENTATION)
#include "magick/module-private.h"
#endif /* defined(MAGICK_IMPLEMENTATION) */


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
