/*
  Copyright (C) 2008 GraphicsMagick Group

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  GraphicsMagick Image Comparison Methods.
*/
#ifndef _MAGICK_COMPARE_H
#define _MAGICK_COMPARE_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif /* defined(__cplusplus) || defined(c_plusplus) */

/*
  Pixel differencing algorithms.
*/
typedef enum
{
  UndefinedHighlightStyle,
  AssignHighlightStyle,
  ThresholdHighlightStyle,
  TintHighlightStyle,
  XorHighlightStyle
} HighlightStyle;

typedef struct _DifferenceImageOptions
{
  ChannelType              channel; /* Channel(s) to difference */
  HighlightStyle           highlight_style; /* Pixel annotation style */
  PixelPacket              highlight_color; /* Changed pixel highlight color */
} DifferenceImageOptions;

extern MagickExport void
  InitializeDifferenceImageOptions(DifferenceImageOptions *options,
                                   ExceptionInfo *exception);

extern MagickExport Image
  *DifferenceImage(const Image *reference_image,const Image *compare_image,
                   const DifferenceImageOptions *difference_options,
                   ExceptionInfo *exception);

/*
  Pixel error metrics.
*/
typedef enum
{
  UndefinedMetric,
  MeanAbsoluteErrorMetric,
  MeanSquaredErrorMetric,
  PeakAbsoluteErrorMetric,
  PeakSignalToNoiseRatioMetric,
  RootMeanSquaredErrorMetric
} MetricType;

/*
  Pixel difference statistics.
*/
typedef struct _DifferenceStatistics
{
  double
    red,
    green,
    blue,
    opacity,
    combined;
} DifferenceStatistics;

extern MagickExport void
  InitializeDifferenceStatistics(DifferenceStatistics *difference_statistics,
                                 ExceptionInfo *exception);

extern MagickExport MagickPassFail
  GetImageChannelDifference(const Image *reference_image,
                            const Image *compare_image,
                            const MetricType metric,
                            DifferenceStatistics *statistics,
                            ExceptionInfo *exception),
  GetImageChannelDistortion(const Image *reference_image,
                            const Image *compare_image,
                            const ChannelType channel,
                            const MetricType metric,
                            double *distortion,
                            ExceptionInfo *exception),
  GetImageDistortion(const Image *reference_image,
                     const Image *compare_image,
                     const MetricType metric,
                     double *distortion,
                     ExceptionInfo *exception);

extern MagickExport MagickBool
  IsImagesEqual(Image *,const Image *);

#if defined(__cplusplus) || defined(c_plusplus)
}
#endif /* defined(__cplusplus) || defined(c_plusplus) */

#endif /* _MAGICK_COMPARE_H */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 2
 * fill-column: 78
 * End:
 */
