/*
  Copyright (C) 2003 - 2020 GraphicsMagick Group
  Copyright (C) 2002 ImageMagick Studio

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  GraphicsMagick Constitute Methods.
*/
#ifndef _MAGICK_CONSTITUTE_H
#define _MAGICK_CONSTITUTE_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif /* defined(__cplusplus) || defined(c_plusplus) */

/*
  Quantum import/export types as used by ImportImagePixelArea() and
  ExportImagePixelArea(). Values are imported or exported in network
  byte order ("big endian") by default, but little endian may be
  selected via the 'endian' option in ExportPixelAreaOptions and
  ImportPixelAreaOptions.
*/
typedef enum
{
  UndefinedQuantum,  /* Not specified */
  IndexQuantum,      /* Colormap indexes */
  GrayQuantum,       /* Grayscale values (minimum value is black) */
  IndexAlphaQuantum, /* Colormap indexes with transparency */
  GrayAlphaQuantum,  /* Grayscale values with transparency */
  RedQuantum,        /* Red values only (RGBA) */
  CyanQuantum,       /* Cyan values only (CMYKA) */
  GreenQuantum,      /* Green values only (RGBA) */
  YellowQuantum,     /* Yellow values only (CMYKA) */
  BlueQuantum,       /* Blue values only (RGBA) */
  MagentaQuantum,    /* Magenta values only (CMYKA) */
  AlphaQuantum,      /* Transparency values (RGBA or CMYKA) */
  BlackQuantum,      /* Black values only (CMYKA) */
  RGBQuantum,        /* Red, green, and blue values */
  RGBAQuantum,       /* Red, green, blue, and transparency values */
  CMYKQuantum,       /* Cyan, magenta, yellow, and black values */
  CMYKAQuantum,      /* Cyan, magenta, yellow, black, and transparency values */
  CIEYQuantum,       /* CIE Y values, based on CCIR-709 primaries */
  CIEXYZQuantum      /* CIE XYZ values, based on CCIR-709 primaries */
} QuantumType;

/*
  Quantum sample type for when exporting/importing a pixel area.
*/
typedef enum
{
  UndefinedQuantumSampleType, /* Not specified */
  UnsignedQuantumSampleType,  /* Unsigned integral type (1-32 or 64 bits) */
  FloatQuantumSampleType      /* Floating point type (16, 24, 32, or 64 bit) */
} QuantumSampleType;

/*
  Quantum size types as used by ConstituteImage() and DispatchImage()/
*/
typedef enum
{
  CharPixel,         /* Unsigned 8 bit 'unsigned char' */
  ShortPixel,        /* Unsigned 16 bit 'unsigned short int' */
  IntegerPixel,      /* Unsigned 32 bit 'unsigned int' */
  LongPixel,         /* Unsigned 32 or 64 bit (CPU dependent) 'unsigned long' */
  FloatPixel,        /* Floating point 32-bit 'float' */
  DoublePixel        /* Floating point 64-bit 'double' */
} StorageType;

/*
  Additional options for ExportImagePixelArea()
*/
typedef struct _ExportPixelAreaOptions
{
  QuantumSampleType
    sample_type;          /* Quantum sample type */

  double
    double_minvalue,      /* Minimum value (default 0.0) for linear floating point samples */
    double_maxvalue;      /* Maximum value (default 1.0) for linear floating point samples */

  MagickBool
    grayscale_miniswhite; /* Grayscale minimum value is white rather than black */

  unsigned long
    pad_bytes;            /* Number of pad bytes to output after pixel data */

  unsigned char
    pad_value;            /* Value to use when padding end of pixel data */

  EndianType
    endian;               /* Endian orientation for 16/32/64 bit types (default MSBEndian) */

  unsigned long
    signature;
} ExportPixelAreaOptions;

/*
  Optional results info for ExportImagePixelArea()
*/
typedef struct _ExportPixelAreaInfo
{
  size_t
    bytes_exported;       /* Number of bytes which were exported */

} ExportPixelAreaInfo;

/*
  Additional options for ImportImagePixelArea()
*/
typedef struct _ImportPixelAreaOptions
{
  QuantumSampleType
    sample_type;          /* Quantum sample type */

  double
    double_minvalue,      /* Minimum value (default 0.0) for linear floating point samples */
    double_maxvalue;      /* Maximum value (default 1.0) for linear floating point samples */

  MagickBool
    grayscale_miniswhite; /* Grayscale minimum value is white rather than black */

  EndianType
    endian;               /* Endian orientation for 16/32/64 bit types (default MSBEndian) */

  unsigned long
    signature;
} ImportPixelAreaOptions;

/*
  Optional results info for ImportImagePixelArea()
*/
typedef struct _ImportPixelAreaInfo
{
  size_t
    bytes_imported;       /* Number of bytes which were imported */

} ImportPixelAreaInfo;

extern MagickExport const char
  *StorageTypeToString(const StorageType storage_type),
  *QuantumSampleTypeToString(const QuantumSampleType sample_type),
  *QuantumTypeToString(const QuantumType quantum_type);

extern MagickExport Image
  *ConstituteImage(const unsigned long width,const unsigned long height,
     const char *map,const StorageType type,const void *pixels,
     ExceptionInfo *exception),
  *ConstituteTextureImage(const unsigned long columns,const unsigned long rows,
     const Image *texture,ExceptionInfo *exception),
  *PingImage(const ImageInfo *image_info,ExceptionInfo *exception),
  *ReadImage(const ImageInfo *image_info,ExceptionInfo *exception),
  *ReadInlineImage(const ImageInfo *image_info,const char *content,
     ExceptionInfo *exception);

extern MagickExport MagickPassFail
  DispatchImage(const Image *image,const long x_offset,const long y_offset,
    const unsigned long columns,const unsigned long rows,const char *map,
    const StorageType type,void *pixels,ExceptionInfo *exception),
  ExportImagePixelArea(const Image *image,const QuantumType quantum_type,
    const unsigned int quantum_size,unsigned char *destination,
    const ExportPixelAreaOptions *options,ExportPixelAreaInfo *export_info),
  ExportViewPixelArea(const ViewInfo *view,const QuantumType quantum_type,
    const unsigned int quantum_size,unsigned char *destination,
    const ExportPixelAreaOptions *options,ExportPixelAreaInfo *export_info);

extern MagickExport MagickPassFail
  ImportImagePixelArea(Image *image,const QuantumType quantum_type,
    const unsigned int quantum_size,const unsigned char *source,
    const ImportPixelAreaOptions *options,ImportPixelAreaInfo *import_info),
  ImportViewPixelArea(ViewInfo *view,const QuantumType quantum_type,
    const unsigned int quantum_size,const unsigned char *source,
    const ImportPixelAreaOptions *options,ImportPixelAreaInfo *import_info),
  WriteImage(const ImageInfo *image_info,Image *image),
  WriteImages(const ImageInfo *image_info,Image *image,const char *filename,
    ExceptionInfo *exception),
  WriteImagesFile(const ImageInfo *image_info,Image *image,FILE * file,
    ExceptionInfo *exception);

extern MagickExport void
  ExportPixelAreaOptionsInit(ExportPixelAreaOptions *options),
  ImportPixelAreaOptionsInit(ImportPixelAreaOptions *options);

extern MagickExport MagickPassFail
  MagickFindRawImageMinMax(Image *image, EndianType endian,
    unsigned long width, unsigned long height,StorageType type,
    unsigned scanline_octets, void *scanline_buffer,
    double *min, double *max);

extern MagickExport unsigned int
  MagickGetQuantumSamplesPerPixel(const QuantumType quantum_type) MAGICK_FUNC_CONST;

#if defined(MAGICK_IMPLEMENTATION)
#  include "magick/constitute-private.h"
#endif /* defined(MAGICK_IMPLEMENTATION) */

#if defined(__cplusplus) || defined(c_plusplus)
}
#endif /* defined(__cplusplus) || defined(c_plusplus) */

#endif /* _MAGICK_CONSTITUTE_H */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 2
 * fill-column: 78
 * End:
 */
