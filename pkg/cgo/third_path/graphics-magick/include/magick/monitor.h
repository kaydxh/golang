/*
  Copyright (C) 2003 - 2020 GraphicsMagick Group
  Copyright (C) 2002 ImageMagick Studio
  Copyright 1991-1999 E. I. du Pont de Nemours and Company

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  GraphicsMagick Progress Monitor Methods.
*/
#ifndef _MAGICK_MONITOR_H
#define _MAGICK_MONITOR_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif

  /*
    Monitor typedef declarations.
  */
  typedef MagickPassFail
  (*MonitorHandler)(const char *text,const magick_int64_t quantum,
                    const magick_uint64_t span,ExceptionInfo *exception);

  /*
    Monitor declarations.
  */
  extern MagickExport MonitorHandler
  SetMonitorHandler(MonitorHandler handler);

  extern MagickExport MagickPassFail
  MagickMonitor(const char *text,
                const magick_int64_t quantum,const magick_uint64_t span,
                ExceptionInfo *exception) MAGICK_FUNC_DEPRECATED;

  extern MagickExport MagickPassFail
  MagickMonitorFormatted(const magick_int64_t quantum,
                         const magick_uint64_t span,
                         ExceptionInfo *exception,
                         const char *format,...) MAGICK_ATTRIBUTE((__format__ (__printf__,4,5)));

#if defined(MAGICK_IMPLEMENTATION)
#  include "magick/monitor-private.h"
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
