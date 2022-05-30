/*
  Copyright (C) 2003-2009 GraphicsMagick Group
  Copyright (C) 2002 ImageMagick Studio

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  GraphicsMagick Image Effect Methods.
*/
#ifndef _MAGICK_EFFECT_H
#define _MAGICK_EFFECT_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif /* defined(__cplusplus) || defined(c_plusplus) */

extern MagickExport Image
  *AdaptiveThresholdImage(const Image *,const unsigned long,const unsigned long,
     const double,ExceptionInfo *),
  *AddNoiseImage(const Image *,const NoiseType,ExceptionInfo *),
  *AddNoiseImageChannel(const Image *image,const ChannelType channel,
     const NoiseType noise_type,ExceptionInfo *exception),
  *BlurImage(const Image *,const double,const double,ExceptionInfo *),
  *BlurImageChannel(const Image *image,const ChannelType channel,
     const double radius,const double sigma,ExceptionInfo *exception),
  *ConvolveImage(const Image *,const unsigned int,const double *,
     ExceptionInfo *),
  *DespeckleImage(const Image *,ExceptionInfo *),
  *EdgeImage(const Image *,const double,ExceptionInfo *),
  *EmbossImage(const Image *,const double,const double,ExceptionInfo *),
  *EnhanceImage(const Image *,ExceptionInfo *),
  *GaussianBlurImage(const Image *,const double,const double,ExceptionInfo *),
  *GaussianBlurImageChannel(const Image *image,
     const ChannelType channel,const double radius,const double sigma,
   ExceptionInfo *exception),
  *MedianFilterImage(const Image *,const double,ExceptionInfo *),
  *MotionBlurImage(const Image *,const double,const double,const double,
     ExceptionInfo *),
  *ReduceNoiseImage(const Image *,const double,ExceptionInfo *),
  *ShadeImage(const Image *,const unsigned int,double,double,ExceptionInfo *),
  *SharpenImage(const Image *,const double,const double,ExceptionInfo *),
  *SharpenImageChannel(const Image *image,const ChannelType channel,
     const double radius,const double sigma,ExceptionInfo *exception),
  *SpreadImage(const Image *,const unsigned int,ExceptionInfo *),
  *UnsharpMaskImage(const Image *,const double,const double,const double,
                    const double,ExceptionInfo *),
  *UnsharpMaskImageChannel(const Image *image,
     const ChannelType channel,const double radius,const double sigma,
     const double amount,const double threshold,
     ExceptionInfo *exception);

extern MagickExport MagickPassFail
  BlackThresholdImage(Image *image,const char *thresholds),
  ChannelThresholdImage(Image *,const char *),
  RandomChannelThresholdImage(Image *,const char *,const char *,
      ExceptionInfo *exception),
  ThresholdImage(Image *,const double),
  WhiteThresholdImage(Image *image,const char *thresholds);

#if defined(__cplusplus) || defined(c_plusplus)
}
#endif /* defined(__cplusplus) || defined(c_plusplus) */

#endif /* _MAGICK_EFFECT_H */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 2
 * fill-column: 78
 * End:
 */
