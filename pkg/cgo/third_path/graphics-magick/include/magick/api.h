/*
  Copyright (C) 2003 - 2012 GraphicsMagick Group
  Copyright (C) 2002 ImageMagick Studio
  Copyright 1991-1999 E. I. du Pont de Nemours and Company

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  GraphicsMagick Application Programming Interface declarations.

*/

#if !defined(_MAGICK_API_H)
#define _MAGICK_API_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif

#include "magick/magick_config.h"
#if defined(__cplusplus) || defined(c_plusplus)
#  undef inline
#endif

#include <stdio.h>
#include <stdarg.h>
#include <stdlib.h>
#include <sys/types.h> /* POSIX 1990 header and declares size_t and ssize_t */

/*
  Note that the WIN32 and WIN64 definitions are provided by the build
  configuration rather than the compiler.  Definitions available from
  the Windows compiler are _WIN32 and _WIN64.
*/
#if defined(WIN32) || defined(WIN64)
#  define MSWINDOWS
#endif /* defined(WIN32) || defined(WIN64) */

#if defined(MAGICK_IMPLEMENTATION)
#  if defined(MSWINDOWS)
  /* Use Visual C++ C inline method extension to improve performance */
#    if !defined(inline) && !defined(__cplusplus) && !defined(c_plusplus)
#      define inline __inline
#    endif
#  endif
#endif

#if defined(PREFIX_MAGICK_SYMBOLS)
#  include "magick/symbols.h"
#endif /* defined(PREFIX_MAGICK_SYMBOLS) */

#include "magick/common.h"
#include "magick/magick_types.h"
#include "magick/analyze.h"
#include "magick/attribute.h"
#include "magick/average.h"
#include "magick/blob.h"
#include "magick/cdl.h"
#include "magick/channel.h"
#include "magick/color.h"
#include "magick/color_lookup.h"
#include "magick/colormap.h"
#include "magick/command.h"
#include "magick/compare.h"
#include "magick/composite.h"
#include "magick/compress.h"
#include "magick/confirm_access.h"
#include "magick/constitute.h"
#include "magick/decorate.h"
#include "magick/delegate.h"
#include "magick/deprecate.h"
#include "magick/describe.h"
#include "magick/draw.h"
#include "magick/effect.h"
#include "magick/enhance.h"
#include "magick/error.h"
#include "magick/fx.h"
#include "magick/gem.h"
#include "magick/gradient.h"
#include "magick/hclut.h"
#include "magick/image.h"
#include "magick/list.h"
#include "magick/log.h"
#include "magick/magic.h"
#include "magick/magick.h"
#include "magick/memory.h"
#include "magick/module.h"
#include "magick/monitor.h"
#include "magick/montage.h"
#include "magick/operator.h"
#include "magick/paint.h"
#include "magick/pixel_cache.h"
#include "magick/pixel_iterator.h"
#include "magick/plasma.h"
#include "magick/profile.h"
#include "magick/quantize.h"
  /*#include "magick/random.h"*/
#include "magick/registry.h"
#include "magick/render.h"
#include "magick/resize.h"
#include "magick/resource.h"
#include "magick/shear.h"
#include "magick/signature.h"
#include "magick/statistics.h"
#include "magick/texture.h"
#include "magick/timer.h"
#include "magick/transform.h"
#include "magick/type.h"
#include "magick/utility.h"
#include "magick/version.h"

#if defined(__cplusplus) || defined(c_plusplus)
}
#endif /* defined(__cplusplus) || defined(c_plusplus) */

#endif /* _MAGICK_API_H */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 2
 * fill-column: 78
 * End:
 */
