/*
  Copyright (C) 2009 GraphicsMagick Group

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  GraphicsMagick Access Confirmation Methods.
*/
#ifndef _MAGICK_CONFIRM_ACCESS_H
#define _MAGICK_CONFIRM_ACCESS_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif

  typedef enum
    {
      UndefinedConfirmAccessMode,
      FileExecuteConfirmAccessMode, /* Path is to be opened for execution */
      FileReadConfirmAccessMode,    /* Path is to be opened for read */
      FileWriteConfirmAccessMode,   /* Path is to be opened for write */
      URLGetFTPConfirmAccessMode,   /* ftp:// URL get */
      URLGetFileConfirmAccessMode,  /* file:// URL get */
      URLGetHTTPConfirmAccessMode   /* http:// URL get */
    } ConfirmAccessMode;

  typedef MagickPassFail
  (*ConfirmAccessHandler)(const ConfirmAccessMode mode,
                          const char *path,
                          ExceptionInfo *exception);

  extern MagickExport MagickPassFail
  MagickConfirmAccess(const ConfirmAccessMode mode,
                      const char *path,
                      ExceptionInfo *exception);

  extern MagickExport ConfirmAccessHandler
  MagickSetConfirmAccessHandler(ConfirmAccessHandler handler);

#if defined(__cplusplus) || defined(c_plusplus)
}
#endif

#endif /* _MAGICK_CONFIRM_ACCESS_H */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 2
 * fill-column: 78
 * End:
 */
