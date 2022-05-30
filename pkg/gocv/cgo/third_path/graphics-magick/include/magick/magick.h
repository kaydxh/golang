/*
  Copyright (C) 2003 - 2020 GraphicsMagick Group
  Copyright (C) 2002 ImageMagick Studio
  Copyright 1991-1999 E. I. du Pont de Nemours and Company

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  GraphicsMagick Application Programming Interface declarations.
*/
#ifndef _MAGICK_MAGICK_H
#define _MAGICK_MAGICK_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif

/*
  Flags to form options passed to InitializeMagickEx
*/
#define MAGICK_OPT_NO_SIGNAL_HANDER 0x0001 /* Don't register ANSI/POSIX signal handlers */

typedef Image
  *(*DecoderHandler)(const ImageInfo *,ExceptionInfo *);

typedef unsigned int
  (*EncoderHandler)(const ImageInfo *,Image *),
  (*MagickHandler)(const unsigned char *,const size_t);

/*
  Stability and usefulness of the coder.
*/
typedef enum
{
  BrokenCoderClass = -1,   /* Known to sometimes/often not work properly or
                              might not be useful at all */
  UnstableCoderClass = 0,  /* Weak implementation, poorly designed file
                              format, and/or hardly ever used */
  StableCoderClass,        /* Well maintained, but not often used */
  PrimaryCoderClass        /* Well maintained and commonly used */
} CoderClass;

/*
  How the file extension should be treated (e.g. in SetImageInfo()).
*/
typedef enum
{
  HintExtensionTreatment = 0, /* Extension is a useful hint to indicate format */
  ObeyExtensionTreatment,     /* Extension must be obeyed as format indicator */
  IgnoreExtensionTreatment    /* Format has no associated extension */
} ExtensionTreatment;

typedef struct _MagickInfo
{
  struct _MagickInfo
    *next,              /* private, next member in list */
    *previous;          /* private, previous member in list */

  const char
    *name;              /* format ID ("magick") */

  const char
    *description,       /* format description */
    *note,              /* usage note for user */
    *version,           /* support library version */
    *module;            /* name of loadable module */

  DecoderHandler
    decoder;            /* function vector to decoding routine */

  EncoderHandler
    encoder;            /* function vector to encoding routine */

  MagickHandler
    magick;             /* function vector to format test routine */

  void
    *client_data;       /* arbitrary user supplied data */

  MagickBool
    adjoin,             /* coder may read/write multiple frames (default True) */
    raw,                /* coder requires that size be set (default False) */
    stealth,            /* coder should not appear in formats listing (default MagickFalse) */
    seekable_stream,    /* coder requires BLOB "seek" and "tell" APIs (default MagickFalse)
                         *   Note that SetImageInfo() currently always copies input
                         *   from a pipe, .gz, or .bz2 file, to a temporary file so
                         *   that it can retrieve a bit of the file header in order to
                         *   support the file header magic logic.
                         */
    blob_support,       /* coder uses BLOB APIs (default True) */
    thread_support;     /* coder is thread safe (default True) */

  CoderClass
    coder_class;        /* Coder usefulness/stability level */

  ExtensionTreatment
    extension_treatment; /* How much faith should be placed on file extension? */

  unsigned long
    signature;          /* private, structure validator */

} MagickInfo;

/*
  Magick method declaractions.
*/
extern MagickExport char
  *MagickToMime(const char *magick);

extern MagickExport const char
  *GetImageMagick(const unsigned char *magick,const size_t length);

extern MagickExport MagickBool
  IsMagickConflict(const char *magick) MAGICK_FUNC_PURE;

extern MagickExport MagickPassFail
  ListModuleMap(FILE *file,ExceptionInfo *exception),
  ListMagickInfo(FILE *file,ExceptionInfo *exception),
  InitializeMagickEx(const char *path, unsigned int options,
                     ExceptionInfo *exception),
  UnregisterMagickInfo(const char *name);

extern MagickExport void
  DestroyMagick(void),
  InitializeMagick(const char *path),
  PanicDestroyMagick(void);

extern MagickExport const MagickInfo
  *GetMagickInfo(const char *name,ExceptionInfo *exception);

extern MagickExport MagickInfo
  **GetMagickInfoArray(ExceptionInfo *exception);

extern MagickExport MagickInfo
  *RegisterMagickInfo(MagickInfo *magick_info),
  *SetMagickInfo(const char *name);

#if defined(MAGICK_IMPLEMENTATION)
#  include "magick/magick-private.h"
#endif /* defined(MAGICK_IMPLEMENTATION) */


#if defined(__cplusplus) || defined(c_plusplus)
}
#endif

#endif /* _MAGICK_MAGICK_H */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 2
 * fill-column: 78
 * End:
 */
