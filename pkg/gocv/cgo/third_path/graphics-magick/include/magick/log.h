/*
  Copyright (C) 2003 - 2020 GraphicsMagick Group
  Copyright (C) 2002 ImageMagick Studio

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  Log methods.
*/
#ifndef _MAGICK_LOG_H
#define _MAGICK_LOG_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif

  /*
    Define declarations.
  */
#define MagickLogFilename  "log.mgk"

  /*
    Obtain the current C function name (if possible)
  */
#  if !defined(GetCurrentFunction)
#    if (((defined(__cplusplus) || defined(c_plusplus)) && defined(HAS_CPP__func__)) || \
        (!(defined(__cplusplus) || defined(c_plusplus)) && defined(HAS_C__func__)))
#      define GetCurrentFunction() (__func__)
#    elif defined(_VISUALC_) && defined(__FUNCTION__)
#      define GetCurrentFunction() (__FUNCTION__)
#    else
#      define GetCurrentFunction() ("unknown")
#    endif
#  endif

  /*
    Obtain current source file, function name, and source file line,
    in a form acceptable for use with LogMagickEvent.
  */
#  if !defined(GetMagickModule)
#    define GetMagickModule()  __FILE__,GetCurrentFunction(),__LINE__
#  endif


/* NOTE: any changes to this effect PerlMagick */
typedef enum
{
  UndefinedEventMask     = 0x00000000,
  NoEventsMask           = 0x00000000,
  ConfigureEventMask     = 0x00000001,
  AnnotateEventMask      = 0x00000002,
  RenderEventMask        = 0x00000004,
  TransformEventMask     = 0x00000008,
  LocaleEventMask        = 0x00000010,
  CoderEventMask         = 0x00000020,
  X11EventMask           = 0x00000040,
  CacheEventMask         = 0x00000080,
  BlobEventMask          = 0x00000100,
  DeprecateEventMask     = 0x00000200,
  UserEventMask          = 0x00000400,
  ResourceEventMask      = 0x00000800,
  TemporaryFileEventMask = 0x00001000,
  /* ExceptionEventMask = WarningEventMask | ErrorEventMask |  FatalErrorEventMask */
  ExceptionEventMask     = 0x00070000,
  OptionEventMask        = 0x00004000,
  InformationEventMask   = 0x00008000,
  WarningEventMask       = 0x00010000,
  ErrorEventMask         = 0x00020000,
  FatalErrorEventMask    = 0x00040000,
  AllEventsMask          = 0x7FFFFFFF
} LogEventType;

typedef void
  (*LogMethod)(const ExceptionType,const char *);

/*
  Method declarations.
*/
extern MagickExport MagickBool
  IsEventLogging(void) MAGICK_FUNC_PURE,
  LogMagickEvent(const ExceptionType type,
    const char *module,const char *function,const unsigned long line,
    const char *format,...) MAGICK_ATTRIBUTE((__format__ (__printf__,5,6))),
  LogMagickEventList(const ExceptionType type,
    const char *module,const char *function,const unsigned long line,
    const char *format,va_list operands) MAGICK_ATTRIBUTE((__format__ (__printf__,5,0)));

extern MagickExport unsigned long
  SetLogEventMask(const char *events);

extern MagickExport void
  SetLogFormat(const char *format),
  SetLogMethod(LogMethod);

#if defined(MAGICK_IMPLEMENTATION)
#  include "magick/log-private.h"
#endif /* MAGICK_IMPLEMENTATION */

#if defined(__cplusplus) || defined(c_plusplus)
}
#endif

#endif /* _MAGICK_LOG_H */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 2
 * fill-column: 78
 * End:
 */
