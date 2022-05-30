/*
  Copyright (C) 2003 - 2009 GraphicsMagick Group
  Copyright (C) 2002 ImageMagick Studio
  Copyright 1991-1999 E. I. du Pont de Nemours and Company

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  GraphicsMagick Image Statistics Methods.
*/
#ifndef _MAGICK_STATISTICS_H
#define _MAGICK_STATISTICS_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif

/*
  Include declarations.
*/
#include "magick/image.h"
#include "magick/error.h"

typedef struct _ImageChannelStatistics
 {
   /* Minimum value observed */
   double maximum;
   /* Maximum value observed */
   double minimum;
   /* Average (mean) value observed */
   double mean;
   /* Standard deviation, sqrt(variance) */
   double standard_deviation;
   /* Variance */
   double variance;
 } ImageChannelStatistics;

typedef struct _ImageStatistics
 {
   ImageChannelStatistics red;
   ImageChannelStatistics green;
   ImageChannelStatistics blue;
   ImageChannelStatistics opacity;
 } ImageStatistics;

extern MagickExport MagickPassFail
  GetImageStatistics(const Image *image,ImageStatistics *statistics,
                     ExceptionInfo *exception);

#if defined(__cplusplus) || defined(c_plusplus)
}
#endif

#endif /* _MAGICK_STATISTICS_H */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 2
 * fill-column: 78
 * End:
 */
