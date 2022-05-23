/*
  Copyright (C) 2003 GraphicsMagick Group
  Copyright (C) 2002 ImageMagick Studio

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  Methods to Get/Set/Destroy Image Text Attributes.
*/
#ifndef _MAGICK_ATTRIBUTE_H
#define _MAGICK_ATTRIBUTE_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif

#include "magick/image.h"

typedef struct _ImageAttribute
{
  char
    *key,           /* identifying key */
    *value;         /* value string */

  size_t
    length;         /* value string length */

  struct _ImageAttribute
    *previous,
    *next;
} ImageAttribute;

/*
  MagickExported text attribute methods.
*/
extern MagickExport const ImageAttribute
  *GetImageAttribute(const Image *image,const char *key),
  *GetImageClippingPathAttribute(const Image *image),
  *GetImageInfoAttribute(const ImageInfo *image_info,const Image *image,const char *key);

extern MagickExport MagickPassFail
  CloneImageAttributes(Image* clone_image, const Image* original_image),
  SetImageAttribute(Image *image,const char *key,const char *value);

extern MagickExport void
  DestroyImageAttributes(Image *image);

#if defined(MAGICK_IMPLEMENTATION)
#  include "magick/attribute-private.h"
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
