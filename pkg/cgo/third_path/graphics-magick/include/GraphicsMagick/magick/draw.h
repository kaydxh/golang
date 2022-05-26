/*
  Copyright (C) 2003-2018 GraphicsMagick Group
  Copyright (C) 2002 ImageMagick Studio

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  ImageMagick Drawing API.

  Usage synopsis:

  DrawContext context;
  context = DrawAllocateContext((DrawInfo*)NULL, image);
    [ any number of drawing commands ]
    DrawSetStrokeColorString(context,"black");
    DrawSetFillColorString(context,"#ff00ff");
    DrawSetStrokeWidth(context,4);
    DrawRectangle(context,72,72,144,144);
  DrawRender(context);
  DrawDestroyContext(context);

*/
#ifndef _MAGICK_DRAW_H
#define _MAGICK_DRAW_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif

#include "magick/render.h"


typedef struct _DrawContext *DrawContext;

extern MagickExport ClipPathUnits
  DrawGetClipUnits(DrawContext context) MAGICK_FUNC_PURE;

extern MagickExport DrawInfo
  *DrawPeekGraphicContext(const DrawContext context);

extern MagickExport DecorationType
  DrawGetTextDecoration(DrawContext context) MAGICK_FUNC_PURE;

extern MagickExport DrawContext
  DrawAllocateContext(const DrawInfo *draw_info, Image *image);

extern MagickExport FillRule
  DrawGetClipRule(DrawContext context) MAGICK_FUNC_PURE,
  DrawGetFillRule(DrawContext context) MAGICK_FUNC_PURE;

extern MagickExport GravityType
  DrawGetGravity(DrawContext context) MAGICK_FUNC_PURE;

extern MagickExport LineCap
  DrawGetStrokeLineCap(DrawContext context) MAGICK_FUNC_PURE;

extern MagickExport LineJoin
  DrawGetStrokeLineJoin(DrawContext context) MAGICK_FUNC_PURE;

extern MagickExport PixelPacket
  DrawGetFillColor(DrawContext context) MAGICK_FUNC_PURE,
  DrawGetStrokeColor(DrawContext context) MAGICK_FUNC_PURE,
  DrawGetTextUnderColor(DrawContext context) MAGICK_FUNC_PURE;

extern MagickExport StretchType
  DrawGetFontStretch(DrawContext context) MAGICK_FUNC_PURE;

extern MagickExport StyleType
  DrawGetFontStyle(DrawContext context) MAGICK_FUNC_PURE;

extern MagickExport char
  *DrawGetClipPath(DrawContext context),
  *DrawGetFont(DrawContext context),
  *DrawGetFontFamily(DrawContext context),
  *DrawGetTextEncoding(DrawContext context);

extern MagickExport int
  DrawRender(const DrawContext context);

extern MagickExport unsigned int
  DrawGetStrokeAntialias(DrawContext context) MAGICK_FUNC_PURE,
  DrawGetTextAntialias(DrawContext context) MAGICK_FUNC_PURE;

extern MagickExport unsigned long
  DrawGetFontWeight(DrawContext context) MAGICK_FUNC_PURE,
  DrawGetStrokeMiterLimit(DrawContext context) MAGICK_FUNC_PURE;

extern MagickExport double
  DrawGetFillOpacity(DrawContext context) MAGICK_FUNC_PURE,
  DrawGetFontSize(DrawContext context) MAGICK_FUNC_PURE,
  *DrawGetStrokeDashArray(DrawContext context, unsigned long *num_elems),
  DrawGetStrokeDashOffset(DrawContext context) MAGICK_FUNC_PURE,
  DrawGetStrokeOpacity(DrawContext context) MAGICK_FUNC_PURE,
  DrawGetStrokeWidth(DrawContext context) MAGICK_FUNC_PURE;

extern MagickExport void
  DrawAffine(DrawContext context, const AffineMatrix *affine),
  DrawAnnotation(DrawContext context,
    const double x, const double y,
    const unsigned char *text),
  DrawArc(DrawContext context,
          const double sx, const double sy,
          const double ex, const double ey,
          const double sd, const double ed),
  DrawBezier(DrawContext context,
             const unsigned long num_coords, const PointInfo * coordinates),
  DrawCircle(DrawContext context,
             const double ox, const double oy,
             const double px, const double py),
  DrawColor(DrawContext context,
            const double x, const double y,
            const PaintMethod paintMethod),
  DrawComment(DrawContext context,const char* comment),
  DrawDestroyContext(DrawContext context),
  DrawEllipse(DrawContext context,
              const double ox, const double oy,
              const double rx, const double ry,
              const double start, const double end),
  DrawComposite(DrawContext context,
                const CompositeOperator composite_operator,
                const double x, const double y,
                const double width, const double height,
                const Image * image ),
  DrawLine(DrawContext context,
           const double sx, const double sy,
           const double ex, const double ey),
  DrawMatte(DrawContext context,
            const double x, const double y,
            const PaintMethod paint_method),
  DrawPathClose(DrawContext context),
  DrawPathCurveToAbsolute(DrawContext context,
                          const double x1, const double y1,
                          const double x2, const double y2,
                          const double x, const double y),
  DrawPathCurveToRelative(DrawContext context,
                          const double x1, const double y1,
                          const double x2, const double y2,
                          const double x, const double y),
  DrawPathCurveToQuadraticBezierAbsolute(DrawContext context,
                                         const double x1, const double y1,
                                         const double x, const double y),
  DrawPathCurveToQuadraticBezierRelative(DrawContext context,
                                         const double x1, const double y1,
                                         const double x, const double y),
  DrawPathCurveToQuadraticBezierSmoothAbsolute(DrawContext context,
                                               const double x, const double y),
  DrawPathCurveToQuadraticBezierSmoothRelative(DrawContext context,
                                               const double x, const double y),
  DrawPathCurveToSmoothAbsolute(DrawContext context,
                                const double x2, const double y2,
                                const double x, const double y),
  DrawPathCurveToSmoothRelative(DrawContext context,
                                const double x2, const double y2,
                                const double x, const double y),
  DrawPathEllipticArcAbsolute(DrawContext context,
                              const double rx, const double ry,
                              const double x_axis_rotation,
                              unsigned int large_arc_flag,
                              unsigned int sweep_flag,
                              const double x, const double y),
  DrawPathEllipticArcRelative(DrawContext context,
                              const double rx, const double ry,
                              const double x_axis_rotation,
                              unsigned int large_arc_flag,
                              unsigned int sweep_flag,
                              const double x, const double y),
  DrawPathFinish(DrawContext context),
  DrawPathLineToAbsolute(DrawContext context,
                         const double x, const double y),
  DrawPathLineToRelative(DrawContext context,
                         const double x, const double y),
  DrawPathLineToHorizontalAbsolute(DrawContext context, const double x),
  DrawPathLineToHorizontalRelative(DrawContext context, const double x),
  DrawPathLineToVerticalAbsolute(DrawContext context, const double y),
  DrawPathLineToVerticalRelative(DrawContext context, const double y),
  DrawPathMoveToAbsolute(DrawContext context,
                         const double x, const double y),
  DrawPathMoveToRelative(DrawContext context,
                         const double x, const double y),
  DrawPathStart(DrawContext context),
  DrawPoint(DrawContext context, const double x, const double y),
  DrawPolygon(DrawContext context,
              const unsigned long num_coords, const PointInfo * coordinates),
  DrawPolyline(DrawContext context,
               const unsigned long num_coords, const PointInfo * coordinates),
  DrawPopClipPath(DrawContext context),
  DrawPopDefs(DrawContext context),
  DrawPopGraphicContext(DrawContext context),
  DrawPopPattern(DrawContext context),
  DrawPushClipPath(DrawContext context, const char *clip_path_id),
  DrawPushDefs(DrawContext context),
  DrawPushGraphicContext(DrawContext context),
  DrawPushPattern(DrawContext context,
                  const char *pattern_id,
                  const double x, const double y,
                  const double width, const double height),
  DrawRectangle(DrawContext context,
                const double x1, const double y1,
                const double x2, const double y2),
  DrawRoundRectangle(DrawContext context,
                     double x1, double y1,
                     double x2, double y2,
                     double rx, double ry),
  DrawScale(DrawContext context, const double x, const double y),
  DrawSetClipPath(DrawContext context, const char *clip_path),
  DrawSetClipRule(DrawContext context, const FillRule fill_rule),
  DrawSetClipUnits(DrawContext context, const ClipPathUnits clip_units),
  DrawSetFillColor(DrawContext context, const PixelPacket *fill_color),
  DrawSetFillColorString(DrawContext context, const char *fill_color),
  DrawSetFillOpacity(DrawContext context, const double fill_opacity),
  DrawSetFillRule(DrawContext context, const FillRule fill_rule),
  DrawSetFillPatternURL(DrawContext context, const char *fill_url),
  DrawSetFont(DrawContext context, const char *font_name),
  DrawSetFontFamily(DrawContext context, const char *font_family),
  DrawSetFontSize(DrawContext context, const double font_pointsize),
  DrawSetFontStretch(DrawContext context, const StretchType font_stretch),
  DrawSetFontStyle(DrawContext context, const StyleType font_style),
  DrawSetFontWeight(DrawContext context, const unsigned long font_weight),
  DrawSetGravity(DrawContext context, const GravityType gravity),
  DrawRotate(DrawContext context, const double degrees),
  DrawSkewX(DrawContext context, const double degrees),
  DrawSkewY(DrawContext context, const double degrees),
  /*
   DrawSetStopColor(DrawContext context, const PixelPacket * color,
                    const double offset),
  */
  DrawSetStrokeAntialias(DrawContext context, const unsigned int true_false),
  DrawSetStrokeColor(DrawContext context, const PixelPacket *stroke_color),
  DrawSetStrokeColorString(DrawContext context, const char *stroke_color),
  DrawSetStrokeDashArray(DrawContext context, const unsigned long num_elems,
                         const double *dasharray),
  DrawSetStrokeDashOffset(DrawContext context,const double dashoffset),
  DrawSetStrokeLineCap(DrawContext context, const LineCap linecap),
  DrawSetStrokeLineJoin(DrawContext context, const LineJoin linejoin),
  DrawSetStrokeMiterLimit(DrawContext context,const unsigned long miterlimit),
  DrawSetStrokeOpacity(DrawContext context, const double opacity),
  DrawSetStrokePatternURL(DrawContext context, const char* stroke_url),
  DrawSetStrokeWidth(DrawContext context, const double width),
  DrawSetTextAntialias(DrawContext context, const unsigned int true_false),
  DrawSetTextDecoration(DrawContext context, const DecorationType decoration),
  DrawSetTextEncoding(DrawContext context, const char *encoding),
  DrawSetTextUnderColor(DrawContext context, const PixelPacket *color),
  DrawSetTextUnderColorString(DrawContext context, const char *under_color),
  DrawSetViewbox(DrawContext context,
                 unsigned long x1, unsigned long y1,
                 unsigned long x2, unsigned long y2),
  DrawTranslate(DrawContext context, const double x, const double y);

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
