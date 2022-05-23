/*
  Copyright (C) 2003-2018 GraphicsMagick Group

  GraphicsMagick Wand API Methods
*/
#ifndef _MAGICK_WAND_API_H
#define _MAGICK_WAND_API_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif

#include <magick/api.h>
#include <wand/wand_symbols.h>

#if defined(_VISUALC_)

/**
 * Under VISUALC we have single threaded static libraries, or
 * mutli-threaded DLLs using the multithreaded runtime DLLs.
 **/
#  if defined(_MT) && defined(_DLL) && !defined(_LIB)
#    pragma warning( disable: 4273 )    /* Disable the stupid dll linkage warnings */
#    if !defined(_WANDLIB_)
#      define WandExport __declspec(dllimport)
#    else
#     define WandExport __declspec(dllexport)
#    endif
#  else
#    define WandExport
#  endif

#  pragma warning(disable : 4018)
#  pragma warning(disable : 4244)
#  pragma warning(disable : 4142)
#else
#  define WandExport
#endif

#include <wand/drawing_wand.h>
#include <wand/magick_wand.h>
#include <wand/pixel_wand.h>

#if defined(__cplusplus) || defined(c_plusplus)
}
#endif /* defined(__cplusplus) || defined(c_plusplus) */

#endif /* _MAGICK_WAND_API_H */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 2
 * fill-column: 78
 * End:
 */
