/*
  Copyright (C) 2003 - 2020 GraphicsMagick Group
  Copyright (C) 2002 ImageMagick Studio

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  Drawing methods.
*/
#ifndef _MAGICK_TYPE_H
#define _MAGICK_TYPE_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif

/*
  Enum declarations.
*/
typedef enum
{
  NormalStretch,
  UltraCondensedStretch,
  ExtraCondensedStretch,
  CondensedStretch,
  SemiCondensedStretch,
  SemiExpandedStretch,
  ExpandedStretch,
  ExtraExpandedStretch,
  UltraExpandedStretch,
  AnyStretch
} StretchType;

typedef enum
{
  NormalStyle,
  ItalicStyle,
  ObliqueStyle,
  AnyStyle
} StyleType;

/*
  Typedef declarations.
*/
typedef struct _TypeInfo
{
  char
    *path,
    *name,
    *description,
    *family;

  StyleType
    style;

  StretchType
    stretch;

  unsigned long
    weight;

  char
    *encoding,
    *foundry,
    *format,
    *metrics,
    *glyphs;

  unsigned int
    stealth;

  unsigned long
    signature;

  struct _TypeInfo
    *previous,
    *next;
} TypeInfo;

/*
  Method declarations.
*/
extern MagickExport char
  **GetTypeList(const char *,unsigned long *);

extern MagickExport MagickPassFail
  ListTypeInfo(FILE *,ExceptionInfo *);

extern MagickExport const TypeInfo
  *GetTypeInfo(const char *,ExceptionInfo *),
  *GetTypeInfoByFamily(const char *,const StyleType,const StretchType,
    const unsigned long,ExceptionInfo *);

#if defined(MAGICK_IMPLEMENTATION)
#  include "magick/type-private.h"
#endif /* MAGICK_IMPLEMENTATION */

#if defined(__cplusplus) || defined(c_plusplus)
}
#endif

#endif /* _MAGICK_TYPE_H */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 2
 * fill-column: 78
 * End:
 */
