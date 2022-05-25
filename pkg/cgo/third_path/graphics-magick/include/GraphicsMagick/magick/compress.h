/*
  Copyright (C) 2003 - 2009 GraphicsMagick Group
  Copyright (C) 2002 ImageMagick Studio
  Copyright 1991-1999 E. I. du Pont de Nemours and Company

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  Image Compression/Decompression Methods.
*/
#ifndef _MAGICK_COMPRESS_H
#define _MAGICK_COMPRESS_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif

/*
  Typedef declaration.
*/
typedef struct _Ascii85Info
{
  long
    offset,
    line_break;

  magick_uint8_t
    buffer[10];
} Ascii85Info;

/*
  TODO: Clean up the interface between BLOB write functions,
  compression functions, and encoding functions so they
  may be hooked into/stacked on top of each other. Most are
  (or can be changed to be) stream based.
*/
typedef unsigned int
  (*WriteByteHook)(Image *, const magick_uint8_t, void *info);

/*
  Commonly used byte writer hooks.
*/
extern MagickExport unsigned int
  Ascii85WriteByteHook(Image *image, const magick_uint8_t code, void *info),
  BlobWriteByteHook(Image *image, const magick_uint8_t code, void *info);

/*
  Compress methods.
*/
extern MagickExport MagickPassFail
  HuffmanDecodeImage(Image *image),
  HuffmanEncodeImage(const ImageInfo *image_info,Image *image),
  HuffmanEncode2Image(const ImageInfo *image_info,Image *image,WriteByteHook write_byte,void *info),
  LZWEncodeImage(Image *image,const size_t length,magick_uint8_t *pixels),
  LZWEncode2Image(Image *image,const size_t length,magick_uint8_t *pixels,WriteByteHook write_byte,void *info),
  PackbitsEncodeImage(Image *image,const size_t length,magick_uint8_t *pixels),
  PackbitsEncode2Image(Image *image,const size_t length,magick_uint8_t *pixels,WriteByteHook write_byte,void *info);

extern MagickExport unsigned char
  *ImageToHuffman2DBlob(const Image *image,const ImageInfo *image_info,
     size_t *length,ExceptionInfo *exception),
  *ImageToJPEGBlob(const Image *image,const ImageInfo *image_info,
     size_t *length,ExceptionInfo *exception);

extern MagickExport void
  Ascii85Encode(Image *image,const magick_uint8_t code),
  Ascii85Flush(Image *image),
  Ascii85Initialize(Image *image);

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
