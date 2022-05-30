/*
  Copyright (C) 2003 GraphicsMagick Group
  Copyright (C) 2002 ImageMagick Studio
  Copyright 1991-1999 E. I. du Pont de Nemours and Company

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  Methods to Reduce the Number of Unique Colors in an Image.
*/
#ifndef _MAGICK_QUANTIZE_H
#define _MAGICK_QUANTIZE_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif

/*
  Define declarations.
*/
#define MaxTreeDepth  8
#define NodesInAList  1536

/*
  Typedef declarations.
*/
typedef struct _QuantizeInfo
{
  unsigned long
    number_colors;

  unsigned int
    tree_depth,
    dither;

  ColorspaceType
    colorspace;

  unsigned int
    measure_error;

  unsigned long
    signature;
} QuantizeInfo;

/*
  Quantization utilities methods.
*/
extern MagickExport QuantizeInfo
  *CloneQuantizeInfo(const QuantizeInfo *);

extern MagickExport unsigned int
  GetImageQuantizeError(Image *),
  MapImage(Image *,const Image *,const unsigned int),
  MapImages(Image *,const Image *,const unsigned int),
  OrderedDitherImage(Image *),
  QuantizeImage(const QuantizeInfo *,Image *),
  QuantizeImages(const QuantizeInfo *,Image *),
  SegmentImage(Image *,const ColorspaceType,const unsigned int,const double,
    const double);

extern MagickExport void
  CompressImageColormap(Image *),
  DestroyQuantizeInfo(QuantizeInfo *),
  GetQuantizeInfo(QuantizeInfo *),
  GrayscalePseudoClassImage(Image *,unsigned int);

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
