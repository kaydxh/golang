/*
  Copyright (C) 2003 - 2010 GraphicsMagick Group
  Copyright (C) 2002 ImageMagick Studio

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  GraphicsMagick Image Composite Methods.
*/
#ifndef _MAGICK_COMPOSITE_H
#define _MAGICK_COMPOSITE_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif /* defined(__cplusplus) || defined(c_plusplus) */

/*
  Special options required by some composition operators.
*/
typedef struct _CompositeOptions_t
{
  /* ModulateComposite */
  double            percent_brightness;

  /* ThresholdComposite */
  double            amount;
  double            threshold;
} CompositeOptions_t;

extern MagickExport MagickPassFail
  CompositeImage(Image *canvas_image,const CompositeOperator compose,
                 const Image *update_image,
                 const long x_offset,const long y_offset),
  CompositeImageRegion(const CompositeOperator compose,
                       const CompositeOptions_t *options,
                       const unsigned long columns,
                       const unsigned long rows,
                       const Image *update_image,
                       const long update_x,
                       const long update_y,
                       Image *canvas_image,
                       const long canvas_x,
                       const long canvas_y,
                       ExceptionInfo *exception),
  MagickCompositeImageUnderColor(Image *image,const PixelPacket *undercolor,
                                 ExceptionInfo *exception);

#if defined(__cplusplus) || defined(c_plusplus)
}
#endif /* defined(__cplusplus) || defined(c_plusplus) */

#endif /* _MAGICK_COMPOSITE_H */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 2
 * fill-column: 78
 * End:
 */
