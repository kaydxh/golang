/*
  Copyright (C) 2003 GraphicsMagick Group
  Copyright (C) 2002 ImageMagick Studio

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  Methods to Read/Write/Invoke Delegates.
*/
#ifndef _MAGICK_DELEGATE_H
#define _MAGICK_DELEGATE_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif

/*
  Delegate structure definitions.
*/
typedef struct _DelegateInfo
{
  char
    *path,    /* Path to delegate configuation file */
    *decode,  /* Decode from format */
    *encode;  /* Transcode to format */

  char
    *commands; /* Commands to execute */

  int mode;    /* <0 = encoder, >0 = decoder */

  MagickBool
    stealth;   /* Don't list this delegate */

  unsigned long
    signature;

  struct _DelegateInfo
    *previous,
    *next;
} DelegateInfo;

/*
  Magick delegate methods.
*/
extern MagickExport char
  *GetDelegateCommand(const ImageInfo *image_info,Image *image,
                      const char *decode,const char *encode,
                      ExceptionInfo *exception);

extern MagickExport const DelegateInfo
  *GetDelegateInfo(const char *decode,const char *encode,
                   ExceptionInfo *exception),
  *GetPostscriptDelegateInfo(const ImageInfo *image_info,
                   unsigned int *antialias, ExceptionInfo *exception);

extern MagickExport DelegateInfo
  *SetDelegateInfo(DelegateInfo *);

extern MagickExport MagickPassFail
  InvokePostscriptDelegate(const unsigned int verbose,const char *command,
                           ExceptionInfo *exception),
  InvokeDelegate(ImageInfo *image_info,Image *image,const char *decode,
                 const char *encode,ExceptionInfo *exception),
  ListDelegateInfo(FILE *file,ExceptionInfo *exception);

#if defined(MAGICK_IMPLEMENTATION)

#if defined(HasGS)
#include "ghostscript/iapi.h"
#endif

#ifndef gs_main_instance_DEFINED
# define gs_main_instance_DEFINED
typedef struct gs_main_instance_s gs_main_instance;
#endif

#if !defined(MagickDLLCall)
#  if defined(MSWINDOWS)
#    define MagickDLLCall __stdcall
#  else
#    define MagickDLLCall
#  endif
#endif

/*
  Define a vector of Ghostscript library callback functions so that
  DLL/shared and static Ghostscript libbraries may be handled identically.
  These definitions must be compatible with those in the Ghostscript API
  headers (which we don't require).

  http://pages.cs.wisc.edu/~ghost/doc/cvs/API.htm
  */
typedef struct _GhostscriptVectors
{
  /* Exit the interpreter (gsapi_exit)*/
  int  (MagickDLLCall *exit)(gs_main_instance *instance);

  /* Destroy instance of Ghostscript.  Call exit first! (gsapi_delete_instance) */
  void (MagickDLLCall *delete_instance)(gs_main_instance *instance);

  /* Initialize the Ghostscript interpreter (gsapi_init_with_args) */
  int  (MagickDLLCall *init_with_args)(gs_main_instance *instance,int argc,
                                       char **argv);

  /* Create a new instance of the Ghostscript interpreter (gsapi_new_instance) */
  int  (MagickDLLCall *new_instance)(gs_main_instance **pinstance,
                                     void *caller_handle);

  /* Execute string command in Ghostscript interpreter (gsapi_run_string) */
  int  (MagickDLLCall *run_string)(gs_main_instance *instance,const char *str,
                                   int user_errors,int *pexit_code);
} GhostscriptVectors;

extern MagickExport void
  DestroyDelegateInfo(void);

extern MagickPassFail
  InitializeDelegateInfo(void);

#endif /* MAGICK_IMPLEMENTATION */

#if defined(__cplusplus) || defined(c_plusplus)
}
#endif

#endif /* _MAGICK_DELEGATE_H */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 2
 * fill-column: 78
 * End:
 */
