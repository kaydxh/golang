/*
  Copyright (C) 2003-2020 GraphicsMagick Group
  Copyright (C) 2002 ImageMagick Studio

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  Magick registry methods.
*/
#ifndef _MAGICK_REGISTRY_H
#define _MAGICK_REGISTRY_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif

/*
  Enum declarations.
*/
typedef enum
{
  UndefinedRegistryType,
  ImageRegistryType,
  ImageInfoRegistryType
} RegistryType;

/*
  Magick registry methods.
*/
extern MagickExport Image
  *GetImageFromMagickRegistry(const char *name,long *id,
     ExceptionInfo *exception);

extern MagickExport long
  SetMagickRegistry(const RegistryType type,const void *blob,
    const size_t length,ExceptionInfo *exception);

extern MagickExport MagickPassFail
  DeleteMagickRegistry(const long id);

extern MagickExport void
  *GetMagickRegistry(const long id,RegistryType *type,size_t *length,
     ExceptionInfo *exception);

#if defined(MAGICK_IMPLEMENTATION)
#  include "magick/registry-private.h"
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
