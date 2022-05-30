/* Copyright (C) 2003-2009 GraphicsMagick Group */
/*
  ImageMagick Drawing Wand API.
*/
#ifndef _MAGICK_DRAWING_WAND_H
#define _MAGICK_DRAWING_WAND_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif

#include "wand/wand_symbols.h"
#include "wand/pixel_wand.h"

#undef CloneDrawingWand
#define CloneDrawingWand MagickCloneDrawingWand
#undef DestroyDrawingWand
#define DestroyDrawingWand MagickDestroyDrawingWand
#undef DrawAffine
#define DrawAffine MagickDrawAffine
#undef DrawAllocateWand
#define DrawAllocateWand MagickDrawAllocateWand
#undef DrawAnnotation
#define DrawAnnotation MagickDrawAnnotation
#undef DrawArc
#define DrawArc MagickDrawArc
#undef DrawBezier
#define DrawBezier MagickDrawBezier
#undef DrawClearException
#define DrawClearException MagickDrawClearException
#undef DrawCircle
#define DrawCircle MagickDrawCircle
#undef DrawColor
#define DrawColor MagickDrawColor
#undef DrawComment
#define DrawComment MagickDrawComment
#undef DrawComposite
#define DrawComposite MagickDrawComposite
#undef DrawEllipse
#define DrawEllipse MagickDrawEllipse
#undef DrawGetClipPath
#define DrawGetClipPath MagickDrawGetClipPath
#undef DrawGetClipRule
#define DrawGetClipRule MagickDrawGetClipRule
#undef DrawGetClipUnits
#define DrawGetClipUnits MagickDrawGetClipUnits
#undef DrawGetException
#define DrawGetException MagickDrawGetException
#undef DrawGetFillColor
#define DrawGetFillColor MagickDrawGetFillColor
#undef DrawGetFillOpacity
#define DrawGetFillOpacity MagickDrawGetFillOpacity
#undef DrawGetFillRule
#define DrawGetFillRule MagickDrawGetFillRule
#undef DrawGetFont
#define DrawGetFont MagickDrawGetFont
#undef DrawGetFontFamily
#define DrawGetFontFamily MagickDrawGetFontFamily
#undef DrawGetFontSize
#define DrawGetFontSize MagickDrawGetFontSize
#undef DrawGetFontStretch
#define DrawGetFontStretch MagickDrawGetFontStretch
#undef DrawGetFontStyle
#define DrawGetFontStyle MagickDrawGetFontStyle
#undef DrawGetFontWeight
#define DrawGetFontWeight MagickDrawGetFontWeight
#undef DrawGetGravity
#define DrawGetGravity MagickDrawGetGravity
#undef DrawGetStrokeAntialias
#define DrawGetStrokeAntialias MagickDrawGetStrokeAntialias
#undef DrawGetStrokeColor
#define DrawGetStrokeColor MagickDrawGetStrokeColor
#undef DrawGetStrokeDashArray
#define DrawGetStrokeDashArray MagickDrawGetStrokeDashArray
#undef DrawGetStrokeDashOffset
#define DrawGetStrokeDashOffset MagickDrawGetStrokeDashOffset
#undef DrawGetStrokeLineCap
#define DrawGetStrokeLineCap MagickDrawGetStrokeLineCap
#undef DrawGetStrokeLineJoin
#define DrawGetStrokeLineJoin MagickDrawGetStrokeLineJoin
#undef DrawGetStrokeMiterLimit
#define DrawGetStrokeMiterLimit MagickDrawGetStrokeMiterLimit
#undef DrawGetStrokeOpacity
#define DrawGetStrokeOpacity MagickDrawGetStrokeOpacity
#undef DrawGetStrokeWidth
#define DrawGetStrokeWidth MagickDrawGetStrokeWidth
#undef DrawGetTextAntialias
#define DrawGetTextAntialias MagickDrawGetTextAntialias
#undef DrawGetTextDecoration
#define DrawGetTextDecoration MagickDrawGetTextDecoration
#undef DrawGetTextEncoding
#define DrawGetTextEncoding MagickDrawGetTextEncoding
#undef DrawGetTextUnderColor
#define DrawGetTextUnderColor MagickDrawGetTextUnderColor
#undef DrawLine
#define DrawLine MagickDrawLine
#undef DrawMatte
#define DrawMatte MagickDrawMatte
#undef DrawPathClose
#define DrawPathClose MagickDrawPathClose
#undef DrawPathCurveToAbsolute
#define DrawPathCurveToAbsolute MagickDrawPathCurveToAbsolute
#undef DrawPathCurveToQuadraticBezierAbsolute
#define DrawPathCurveToQuadraticBezierAbsolute MagickDrawPathCurveToQuadraticBezierAbsolute
#undef DrawPathCurveToQuadraticBezierRelative
#define DrawPathCurveToQuadraticBezierRelative MagickDrawPathCurveToQuadraticBezierRelative
#undef DrawPathCurveToQuadraticBezierSmoothAbsolute
#define DrawPathCurveToQuadraticBezierSmoothAbsolute MagickDrawPathCurveToQuadraticBezierSmoothAbsolute
#undef DrawPathCurveToQuadraticBezierSmoothRelative
#define DrawPathCurveToQuadraticBezierSmoothRelative MagickDrawPathCurveToQuadraticBezierSmoothRelative
#undef DrawPathCurveToRelative
#define DrawPathCurveToRelative MagickDrawPathCurveToRelative
#undef DrawPathCurveToSmoothAbsolute
#define DrawPathCurveToSmoothAbsolute MagickDrawPathCurveToSmoothAbsolute
#undef DrawPathCurveToSmoothRelative
#define DrawPathCurveToSmoothRelative MagickDrawPathCurveToSmoothRelative
#undef DrawPathEllipticArcAbsolute
#define DrawPathEllipticArcAbsolute MagickDrawPathEllipticArcAbsolute
#undef DrawPathEllipticArcRelative
#define DrawPathEllipticArcRelative MagickDrawPathEllipticArcRelative
#undef DrawPathFinish
#define DrawPathFinish MagickDrawPathFinish
#undef DrawPathLineToAbsolute
#define DrawPathLineToAbsolute MagickDrawPathLineToAbsolute
#undef DrawPathLineToHorizontalAbsolute
#define DrawPathLineToHorizontalAbsolute MagickDrawPathLineToHorizontalAbsolute
#undef DrawPathLineToHorizontalRelative
#define DrawPathLineToHorizontalRelative MagickDrawPathLineToHorizontalRelative
#undef DrawPathLineToRelative
#define DrawPathLineToRelative MagickDrawPathLineToRelative
#undef DrawPathLineToVerticalAbsolute
#define DrawPathLineToVerticalAbsolute MagickDrawPathLineToVerticalAbsolute
#undef DrawPathLineToVerticalRelative
#define DrawPathLineToVerticalRelative MagickDrawPathLineToVerticalRelative
#undef DrawPathMoveToAbsolute
#define DrawPathMoveToAbsolute MagickDrawPathMoveToAbsolute
#undef DrawPathMoveToRelative
#define DrawPathMoveToRelative MagickDrawPathMoveToRelative
#undef DrawPathStart
#define DrawPathStart MagickDrawPathStart
#undef DrawPeekGraphicContext
#define DrawPeekGraphicContext MagickDrawPeekGraphicContext
#undef DrawPoint
#define DrawPoint MagickDrawPoint
#undef DrawPolygon
#define DrawPolygon MagickDrawPolygon
#undef DrawPolyline
#define DrawPolyline MagickDrawPolyline
#undef DrawPopClipPath
#define DrawPopClipPath MagickDrawPopClipPath
#undef DrawPopDefs
#define DrawPopDefs MagickDrawPopDefs
#undef DrawPopGraphicContext
#define DrawPopGraphicContext MagickDrawPopGraphicContext
#undef DrawPopPattern
#define DrawPopPattern MagickDrawPopPattern
#undef DrawPushClipPath
#define DrawPushClipPath MagickDrawPushClipPath
#undef DrawPushDefs
#define DrawPushDefs MagickDrawPushDefs
#undef DrawPushGraphicContext
#define DrawPushGraphicContext MagickDrawPushGraphicContext
#undef DrawPushPattern
#define DrawPushPattern MagickDrawPushPattern
#undef DrawRectangle
#define DrawRectangle MagickDrawRectangle
#undef DrawRender
#define DrawRender MagickDrawRender
#undef DrawRotate
#define DrawRotate MagickDrawRotate
#undef DrawRoundRectangle
#define DrawRoundRectangle MagickDrawRoundRectangle
#undef DrawScale
#define DrawScale MagickDrawScale
#undef DrawSetClipPath
#define DrawSetClipPath MagickDrawSetClipPath
#undef DrawSetClipRule
#define DrawSetClipRule MagickDrawSetClipRule
#undef DrawSetClipUnits
#define DrawSetClipUnits MagickDrawSetClipUnits
#undef DrawSetFillColor
#define DrawSetFillColor MagickDrawSetFillColor
#undef DrawSetFillOpacity
#define DrawSetFillOpacity MagickDrawSetFillOpacity
#undef DrawSetFillPatternURL
#define DrawSetFillPatternURL MagickDrawSetFillPatternURL
#undef DrawSetFillRule
#define DrawSetFillRule MagickDrawSetFillRule
#undef DrawSetFont
#define DrawSetFont MagickDrawSetFont
#undef DrawSetFontFamily
#define DrawSetFontFamily MagickDrawSetFontFamily
#undef DrawSetFontSize
#define DrawSetFontSize MagickDrawSetFontSize
#undef DrawSetFontStretch
#define DrawSetFontStretch MagickDrawSetFontStretch
#undef DrawSetFontStyle
#define DrawSetFontStyle MagickDrawSetFontStyle
#undef DrawSetFontWeight
#define DrawSetFontWeight MagickDrawSetFontWeight
#undef DrawSetGravity
#define DrawSetGravity MagickDrawSetGravity
#undef DrawSetStrokeAntialias
#define DrawSetStrokeAntialias MagickDrawSetStrokeAntialias
#undef DrawSetStrokeColor
#define DrawSetStrokeColor MagickDrawSetStrokeColor
#undef DrawSetStrokeDashArray
#define DrawSetStrokeDashArray MagickDrawSetStrokeDashArray
#undef DrawSetStrokeDashOffset
#define DrawSetStrokeDashOffset MagickDrawSetStrokeDashOffset
#undef DrawSetStrokeLineCap
#define DrawSetStrokeLineCap MagickDrawSetStrokeLineCap
#undef DrawSetStrokeLineJoin
#define DrawSetStrokeLineJoin MagickDrawSetStrokeLineJoin
#undef DrawSetStrokeMiterLimit
#define DrawSetStrokeMiterLimit MagickDrawSetStrokeMiterLimit
#undef DrawSetStrokeOpacity
#define DrawSetStrokeOpacity MagickDrawSetStrokeOpacity
#undef DrawSetStrokePatternURL
#define DrawSetStrokePatternURL MagickDrawSetStrokePatternURL
#undef DrawSetStrokeWidth
#define DrawSetStrokeWidth MagickDrawSetStrokeWidth
#undef DrawSetTextAntialias
#define DrawSetTextAntialias MagickDrawSetTextAntialias
#undef DrawSetTextDecoration
#define DrawSetTextDecoration MagickDrawSetTextDecoration
#undef DrawSetTextEncoding
#define DrawSetTextEncoding MagickDrawSetTextEncoding
#undef DrawSetTextUnderColor
#define DrawSetTextUnderColor MagickDrawSetTextUnderColor
#undef DrawSetViewbox
#define DrawSetViewbox MagickDrawSetViewbox
#undef DrawSkewX
#define DrawSkewX MagickDrawSkewX
#undef DrawSkewY
#define DrawSkewY MagickDrawSkewY
#undef DrawTranslate
#define DrawTranslate MagickDrawTranslate
#undef NewDrawingWand
#define NewDrawingWand MagickNewDrawingWand

typedef struct _DrawingWand
  DrawingWand;

extern WandExport char
  *DrawGetClipPath(const DrawingWand *),
  *DrawGetException(const DrawingWand *,ExceptionType *),
  *DrawGetFont(const DrawingWand *),
  *DrawGetFontFamily(const DrawingWand *),
  *DrawGetTextEncoding(const DrawingWand *);

extern WandExport ClipPathUnits
  DrawGetClipUnits(const DrawingWand *) MAGICK_FUNC_PURE;

extern WandExport DecorationType
  DrawGetTextDecoration(const DrawingWand *) MAGICK_FUNC_PURE;

extern WandExport double
  DrawGetFillOpacity(const DrawingWand *) MAGICK_FUNC_PURE,
  DrawGetFontSize(const DrawingWand *) MAGICK_FUNC_PURE,
  *DrawGetStrokeDashArray(const DrawingWand *,unsigned long *),
  DrawGetStrokeDashOffset(const DrawingWand *) MAGICK_FUNC_PURE,
  DrawGetStrokeOpacity(const DrawingWand *) MAGICK_FUNC_PURE,
  DrawGetStrokeWidth(const DrawingWand *) MAGICK_FUNC_PURE;

extern WandExport DrawInfo
  *DrawPeekGraphicContext(const DrawingWand *);

extern WandExport DrawingWand
  *CloneDrawingWand(const DrawingWand *drawing_wand),
  *DrawAllocateWand(const DrawInfo *,Image *) MAGICK_ATTRIBUTE ((deprecated)),
  *NewDrawingWand(void);

extern WandExport FillRule
  DrawGetClipRule(const DrawingWand *) MAGICK_FUNC_PURE,
  DrawGetFillRule(const DrawingWand *) MAGICK_FUNC_PURE;

extern WandExport GravityType
  DrawGetGravity(const DrawingWand *) MAGICK_FUNC_PURE;

extern WandExport LineCap
  DrawGetStrokeLineCap(const DrawingWand *) MAGICK_FUNC_PURE;

extern WandExport LineJoin
  DrawGetStrokeLineJoin(const DrawingWand *) MAGICK_FUNC_PURE;

extern WandExport StretchType
  DrawGetFontStretch(const DrawingWand *) MAGICK_FUNC_PURE;

extern WandExport StyleType
  DrawGetFontStyle(const DrawingWand *) MAGICK_FUNC_PURE;

extern WandExport unsigned int
  DrawClearException(DrawingWand *),
  DrawGetStrokeAntialias(const DrawingWand *) MAGICK_FUNC_PURE,
  DrawGetTextAntialias(const DrawingWand *) MAGICK_FUNC_PURE,
  DrawRender(const DrawingWand *) MAGICK_ATTRIBUTE ((deprecated)); /* Use MagickDrawImage() instead */

extern WandExport unsigned long
  DrawGetFontWeight(const DrawingWand *) MAGICK_FUNC_PURE,
  DrawGetStrokeMiterLimit(const DrawingWand *) MAGICK_FUNC_PURE;

extern WandExport void
  DrawAffine(DrawingWand *,const AffineMatrix *),
  DrawAnnotation(DrawingWand *,const double,const double,const unsigned char *),
  DrawArc(DrawingWand *,const double,const double,const double,const double,
    const double,const double),
  DrawBezier(DrawingWand *,const unsigned long,const PointInfo *),
  DrawCircle(DrawingWand *,const double,const double,const double,const double),
  DrawColor(DrawingWand *,const double,const double,const PaintMethod),
  DrawComment(DrawingWand *,const char *),
  DestroyDrawingWand(DrawingWand *),
  DrawEllipse(DrawingWand *,const double,const double,const double,const double,
    const double,const double),
  DrawComposite(DrawingWand *,const CompositeOperator,const double,const double,
    const double,const double,const Image *),
  DrawGetFillColor(const DrawingWand *,PixelWand *),
  DrawGetStrokeColor(const DrawingWand *,PixelWand *),
  DrawGetTextUnderColor(const DrawingWand *,PixelWand *),
  DrawLine(DrawingWand *,const double, const double,const double,const double),
  DrawMatte(DrawingWand *,const double,const double,const PaintMethod),
  DrawPathClose(DrawingWand *),
  DrawPathCurveToAbsolute(DrawingWand *,const double,const double,const double,
    const double,const double,const double),
  DrawPathCurveToRelative(DrawingWand *,const double,const double,const double,
    const double,const double, const double),
  DrawPathCurveToQuadraticBezierAbsolute(DrawingWand *,const double,
    const double,const double,const double),
  DrawPathCurveToQuadraticBezierRelative(DrawingWand *,const double,
    const double,const double,const double),
  DrawPathCurveToQuadraticBezierSmoothAbsolute(DrawingWand *,const double,
    const double),
  DrawPathCurveToQuadraticBezierSmoothRelative(DrawingWand *,const double,
    const double),
  DrawPathCurveToSmoothAbsolute(DrawingWand *,const double,const double,
    const double,const double),
  DrawPathCurveToSmoothRelative(DrawingWand *,const double,const double,
    const double,const double),
  DrawPathEllipticArcAbsolute(DrawingWand *,const double,const double,
    const double,unsigned int,unsigned int,const double,const double),
  DrawPathEllipticArcRelative(DrawingWand *,const double,const double,
    const double,unsigned int,unsigned int,const double,const double),
  DrawPathFinish(DrawingWand *),
  DrawPathLineToAbsolute(DrawingWand *,const double,const double),
  DrawPathLineToRelative(DrawingWand *,const double,const double),
  DrawPathLineToHorizontalAbsolute(DrawingWand *,const double),
  DrawPathLineToHorizontalRelative(DrawingWand *,const double),
  DrawPathLineToVerticalAbsolute(DrawingWand *,const double),
  DrawPathLineToVerticalRelative(DrawingWand *,const double),
  DrawPathMoveToAbsolute(DrawingWand *,const double,const double),
  DrawPathMoveToRelative(DrawingWand *,const double,const double),
  DrawPathStart(DrawingWand *),
  DrawPoint(DrawingWand *,const double,const double),
  DrawPolygon(DrawingWand *,const unsigned long,const PointInfo *),
  DrawPolyline(DrawingWand *,const unsigned long,const PointInfo *),
  DrawPopClipPath(DrawingWand *),
  DrawPopDefs(DrawingWand *),
  DrawPopGraphicContext(DrawingWand *),
  DrawPopPattern(DrawingWand *),
  DrawPushClipPath(DrawingWand *,const char *),
  DrawPushDefs(DrawingWand *),
  DrawPushGraphicContext(DrawingWand *),
  DrawPushPattern(DrawingWand *,const char *,const double,const double,
    const double,const double),
  DrawRectangle(DrawingWand *,const double,const double,const double,
    const double),
  DrawRotate(DrawingWand *,const double),
  DrawRoundRectangle(DrawingWand *,double,double,double,double,double,double),
  DrawScale(DrawingWand *,const double,const double),
  DrawSetClipPath(DrawingWand *,const char *),
  DrawSetClipRule(DrawingWand *,const FillRule),
  DrawSetClipUnits(DrawingWand *,const ClipPathUnits),
  DrawSetFillColor(DrawingWand *,const PixelWand *),
  DrawSetFillOpacity(DrawingWand *,const double),
  DrawSetFillRule(DrawingWand *,const FillRule),
  DrawSetFillPatternURL(DrawingWand *,const char *),
  DrawSetFont(DrawingWand *,const char *),
  DrawSetFontFamily(DrawingWand *,const char *),
  DrawSetFontSize(DrawingWand *,const double),
  DrawSetFontStretch(DrawingWand *,const StretchType),
  DrawSetFontStyle(DrawingWand *,const StyleType),
  DrawSetFontWeight(DrawingWand *,const unsigned long),
  DrawSetGravity(DrawingWand *,const GravityType),
  DrawSkewX(DrawingWand *,const double),
  DrawSkewY(DrawingWand *,const double),
  DrawSetStrokeAntialias(DrawingWand *,const unsigned int),
  DrawSetStrokeColor(DrawingWand *,const PixelWand *),
  DrawSetStrokeDashArray(DrawingWand *,const unsigned long,const double *),
  DrawSetStrokeDashOffset(DrawingWand *,const double dashoffset),
  DrawSetStrokeLineCap(DrawingWand *,const LineCap),
  DrawSetStrokeLineJoin(DrawingWand *,const LineJoin),
  DrawSetStrokeMiterLimit(DrawingWand *,const unsigned long),
  DrawSetStrokeOpacity(DrawingWand *, const double),
  DrawSetStrokePatternURL(DrawingWand *,const char *),
  DrawSetStrokeWidth(DrawingWand *,const double),
  DrawSetTextAntialias(DrawingWand *,const unsigned int),
  DrawSetTextDecoration(DrawingWand *,const DecorationType),
  DrawSetTextEncoding(DrawingWand *,const char *),
  DrawSetTextUnderColor(DrawingWand *,const PixelWand *),
  DrawSetViewbox(DrawingWand *,unsigned long,unsigned long,unsigned long,
    unsigned long),
  DrawTranslate(DrawingWand *,const double,const double);

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
