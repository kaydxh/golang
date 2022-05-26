/*
  Copyright (C) 2003 - 2018 GraphicsMagick Group
  Copyright (C) 2002 ImageMagick Studio

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  GraphicsMagick Pixel Cache Methods.
*/
#ifndef _MAGICK_CACHE_H
#define _MAGICK_CACHE_H

#include "magick/forward.h"

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif

  /*
    Enum declaractions.
  */

  typedef enum
    {
      UndefinedVirtualPixelMethod,
      ConstantVirtualPixelMethod,
      EdgeVirtualPixelMethod,
      MirrorVirtualPixelMethod,
      TileVirtualPixelMethod
    } VirtualPixelMethod;

  /*
    Typedef declaractions.
  */
  typedef _CacheInfoPtr_ Cache;

  /*****
   *
   * Default View interfaces
   *
   *****/

  /*
    Read only access to a rectangular pixel region.
  */
  MagickExport const PixelPacket
  *AcquireImagePixels(const Image *image,const long x,const long y,
                      const unsigned long columns,
                      const unsigned long rows,ExceptionInfo *exception);

  /*
    AccessImmutableIndexes() returns the read-only indexes
    associated with a rectangular pixel region already selected via
    AcquireImagePixels().
  */
  extern MagickExport const IndexPacket
  *AccessImmutableIndexes(const Image *image);

  /*
    Return one DirectClass pixel at the the specified (x,y) location.
    Similar function as GetOnePixel().  Note that the value returned
    by GetIndexes() may or may not be influenced by this function.
  */
  extern MagickExport PixelPacket
  AcquireOnePixel(const Image *image,const long x,const long y,
                  ExceptionInfo *exception);


  /*
    GetImagePixels() and GetImagePixelsEx() obtains a pixel region for
    read/write access.
  */
  extern MagickExport PixelPacket
  *GetImagePixels(Image *image,const long x,const long y,
                  const unsigned long columns,const unsigned long rows);
  extern MagickExport PixelPacket
  *GetImagePixelsEx(Image *image,const long x,const long y,
                    const unsigned long columns,const unsigned long rows,
                    ExceptionInfo *exception);

  /*
    GetImageVirtualPixelMethod() gets the "virtual pixels" method for
    the image.
  */
  extern MagickExport VirtualPixelMethod
  GetImageVirtualPixelMethod(const Image *image) MAGICK_FUNC_PURE;

  /*
    GetPixels() and AccessMutablePixels() return the pixels associated
    with the last call to SetImagePixels() or GetImagePixels().
  */
  extern MagickExport PixelPacket
  *GetPixels(const Image *image)
    MAGICK_FUNC_DEPRECATED; /* Prefer AccessMutablePixels instead */
  extern MagickExport PixelPacket
  *AccessMutablePixels(Image *image);

  /*
    GetIndexes() and AccessMutableIndexes() return the colormap
    indexes associated with the last call to SetImagePixels() or
    GetImagePixels().
  */
  extern MagickExport IndexPacket
  *GetIndexes(const Image *image)
    MAGICK_FUNC_DEPRECATED; /* Prefer AccessMutableIndexes() instead */
  extern MagickExport IndexPacket
  *AccessMutableIndexes(Image *image);

  /*
    GetOnePixel() returns a single DirectClass pixel at the specified
    (x,y) location.  Similar to AcquireOnePixel().  It is preferred to
    use AcquireOnePixel() since it allows reporting to a specified
    exception structure. Note that the value returned by GetIndexes()
    is not reliably influenced by this function.
  */
  extern MagickExport PixelPacket
  GetOnePixel(Image *image,const long x,const long y)
    MAGICK_FUNC_DEPRECATED; /* Prefer AcquireOnePixel() instead */

  /*
    GetPixelCacheArea() returns the area (width * height in pixels)
    consumed by the current pixel cache.
  */
  extern MagickExport magick_off_t
  GetPixelCacheArea(const Image *image);

  /*
    SetImagePixels() and SetImagePixelsEx() initialize a pixel region
    for write-only access.
  */
  extern MagickExport PixelPacket
  *SetImagePixels(Image *image,const long x,const long y,
                  const unsigned long columns,const unsigned long rows);
  extern MagickExport PixelPacket
  *SetImagePixelsEx(Image *image,const long x,const long y,
                    const unsigned long columns,const unsigned long rows,
                    ExceptionInfo *exception);

  /*
    SetImageVirtualPixelMethod() sets the "virtual pixels" method for
    the image.
  */
  extern MagickExport MagickPassFail
  SetImageVirtualPixelMethod(const Image *image,
                             const VirtualPixelMethod method);

  /*
    SyncImagePixels() and SyncImagePixelsEx() save the image pixels to
    the in-memory or disk cache.
  */
  extern MagickExport MagickPassFail
  SyncImagePixels(Image *image);
  extern MagickExport MagickPassFail
  SyncImagePixelsEx(Image *image,ExceptionInfo *exception);

  /****
   *
   * Cache view interfaces
   *
   ****/

  /*
    OpenCacheView() opens a cache view.
  */
  extern MagickExport ViewInfo
  *OpenCacheView(Image *image);

  /*
    CloseCacheView() closes a cache view.
  */
  extern MagickExport void
  CloseCacheView(ViewInfo *view);


  /*
    AccessCacheViewPixels() returns the pixels associated with the
    last request to select a view pixel region
    (i.e. AcquireCacheViewPixels() or GetCacheViewPixels()).
  */
  extern MagickExport PixelPacket
  *AccessCacheViewPixels(const ViewInfo *view) MAGICK_FUNC_PURE;

  /*
    AcquireCacheViewIndexes() returns read-only indexes associated
    with a cache view.
  */
  extern MagickExport const IndexPacket
  *AcquireCacheViewIndexes(const ViewInfo *view) MAGICK_FUNC_PURE;

  /*
    AcquireCacheViewPixels() obtains a pixel region from a cache view
    for read-only access.
  */
  extern MagickExport const PixelPacket
  *AcquireCacheViewPixels(ViewInfo *view,
                          const long x,const long y,
                          const unsigned long columns,
                          const unsigned long rows,
                          ExceptionInfo *exception);

  /*
    AcquireOneCacheViewPixel() returns one DirectClass pixel from a
    cache view. Note that the value returned by GetCacheViewIndexes()
    is not reliably influenced by this function.
  */
  extern MagickExport MagickPassFail
  AcquireOneCacheViewPixel(ViewInfo *view,PixelPacket *pixel,
                           const long x,const long y,ExceptionInfo *exception);

  /*
    GetCacheViewArea() returns the area (width * height in pixels)
    currently consumed by the pixel cache view.
  */
  extern MagickExport magick_off_t
  GetCacheViewArea(const ViewInfo *view) MAGICK_FUNC_PURE;

  /*
    GetCacheViewImage() obtains the image used to allocate the cache view.
  */
  extern Image *
  GetCacheViewImage(const ViewInfo *view) MAGICK_FUNC_PURE;

  /*
    GetCacheViewIndexes() returns the indexes associated with a cache view.
  */
  extern MagickExport IndexPacket
  *GetCacheViewIndexes(const ViewInfo *view) MAGICK_FUNC_PURE;

  /*
    GetCacheViewPixels() obtains a pixel region from a cache view for
    read/write access.
  */
  extern MagickExport PixelPacket
  *GetCacheViewPixels(ViewInfo *view,const long x,const long y,
                      const unsigned long columns,const unsigned long rows,
                      ExceptionInfo *exception);

  /*
    Obtain the offset and size of the selected region.
  */
  extern MagickExport RectangleInfo
  GetCacheViewRegion(const ViewInfo *view) MAGICK_FUNC_PURE;


  /*
    SetCacheViewPixels() gets blank writeable pixels from the pixel
    cache view.
  */
  extern MagickExport PixelPacket
  *SetCacheViewPixels(ViewInfo *view,const long x,const long y,
                      const unsigned long columns,const unsigned long rows,
                      ExceptionInfo *exception);

  /*
    SyncCacheViewPixels() saves any changes to the pixel cache view.
  */
  extern MagickExport MagickPassFail
  SyncCacheViewPixels(const ViewInfo *view,ExceptionInfo *exception);

#if defined(MAGICK_IMPLEMENTATION)

#include "magick/pixel_cache-private.h"

#endif /* defined(MAGICK_IMPLEMENTATION) */

#if defined(__cplusplus) || defined(c_plusplus)
}
#endif /* defined(__cplusplus) || defined(c_plusplus) */

#endif /* _MAGICK_CACHE_H */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 2
 * fill-column: 78
 * End:
 */
