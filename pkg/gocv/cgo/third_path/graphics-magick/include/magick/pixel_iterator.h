/*
  Copyright (C) 2004-2016 GraphicsMagick Group

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  Interfaces to support simple iterative pixel read/update access within an
  image or between two images.  These interfaces exist in order to eliminate
  large amounts of redundant code and to allow changing the underlying
  implementation without changing the using code. These interfaces
  intentionally omit any pixel position information in order to not constrain
  the implementation and to improve performance.

  User-provided callbacks must be thread-safe (preferably re-entrant) since
  they may be invoked by multiple threads.

  These interfaces have proven to be future safe (since implemented) and may
  be safely used by other applications/libraries.

  Written by Bob Friesenhahn, March 2004, Updated for regions 2008.

*/
#ifndef _PIXEL_ROW_ITERATOR_H
#define _PIXEL_ROW_ITERATOR_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif

  /*
    Pixel iterator options.
  */
  typedef struct _PixelIteratorOptions
  {
    int           max_threads; /* Desired number of threads */
    unsigned long signature;
  } PixelIteratorOptions;


  /*
    Initialize pixel iterator options with defaults.
  */
  extern MagickExport void
  InitializePixelIteratorOptions(PixelIteratorOptions *options,
                                 ExceptionInfo *exception);

  /*
    Read-only access across pixel region.
  */

  typedef MagickPassFail (*PixelIteratorMonoReadCallback)
    (
     void *mutable_data,                   /* User provided mutable data */
     const void *immutable_data,       /* User provided immutable data */
     const Image *const_image,          /* Input image */
     const PixelPacket *pixels,         /* Pixel row */
     const IndexPacket *indexes,        /* Pixel indexes */
     const long npixels,                /* Number of pixels in row */
     ExceptionInfo *exception           /* Exception report */
     );

  extern MagickExport MagickPassFail
  PixelIterateMonoRead(PixelIteratorMonoReadCallback call_back,
                       const PixelIteratorOptions *options,
                       const char *description,
                       void *mutable_data,
                       const void *immutable_data,
                       const long x,
                       const long y,
                       const unsigned long columns,
                       const unsigned long rows,
                       const Image *image,
                       ExceptionInfo *exception);


  typedef MagickPassFail (*PixelIteratorMonoModifyCallback)
    (
     void *mutable_data,                /* User provided mutable data */
     const void *immutable_data,        /* User provided immutable data */
     Image *image,                      /* Modify image */
     PixelPacket *pixels,               /* Pixel row */
     IndexPacket *indexes,              /* Pixel row indexes */
     const long npixels,                /* Number of pixels in row */
     ExceptionInfo *exception           /* Exception report */
     );

  /*
    Write access across pixel region.
  */
  extern MagickExport MagickPassFail
  PixelIterateMonoSet(PixelIteratorMonoModifyCallback call_back,
                      const PixelIteratorOptions *options,
                      const char *description,
                      void *mutable_data,
                      const void *immutable_data,
                      const long x,
                      const long y,
                      const unsigned long columns,
                      const unsigned long rows,
                      Image *image,
                      ExceptionInfo *exception);

  /*
    Read-write access across pixel region.
  */
  extern MagickExport MagickPassFail
  PixelIterateMonoModify(PixelIteratorMonoModifyCallback call_back,
                         const PixelIteratorOptions *options,
                         const char *description,
                         void *mutable_data,
                         const void *immutable_data,
                         const long x,
                         const long y,
                         const unsigned long columns,
                         const unsigned long rows,
                         Image *image,
                         ExceptionInfo *exception);

  /*
    Read-only access across pixel regions of two images.
  */

  typedef MagickPassFail (*PixelIteratorDualReadCallback)
    (
     void *mutable_data,                /* User provided mutable data */
     const void *immutable_data,        /* User provided immutable data */
     const Image *first_image,          /* First Input image */
     const PixelPacket *first_pixels,   /* Pixel row in first image */
     const IndexPacket *first_indexes,  /* Pixel row indexes in first image */
     const Image *second_image,         /* Second Input image */
     const PixelPacket *second_pixels,  /* Pixel row in second image */
     const IndexPacket *second_indexes, /* Pixel row indexes in second image */
     const long npixels,                /* Number of pixels in row */
     ExceptionInfo *exception           /* Exception report */
     );

  extern MagickExport MagickPassFail
  PixelIterateDualRead(PixelIteratorDualReadCallback call_back,
                       const PixelIteratorOptions *options,
                       const char *description,
                       void *mutable_data,
                       const void *immutable_data,
                       const unsigned long columns,
                       const unsigned long rows,
                       const Image *first_image,
                       const long first_x,
                       const long first_y,
                       const Image *second_image,
                       const long second_x,
                       const long second_y,
                       ExceptionInfo *exception);

  /*
    Read-write access across pixel regions of two images. The first
    (source) image is accessed read-only while the second (update)
    image is accessed as read-write.
  */

  typedef MagickPassFail (*PixelIteratorDualModifyCallback)
    (
     void *mutable_data,                /* User provided mutable data */
     const void *immutable_data,        /* User provided immutable data */
     const Image *source_image,         /* Source image */
     const PixelPacket *source_pixels,  /* Pixel row in source image */
     const IndexPacket *source_indexes, /* Pixel row indexes in source image */
     Image *update_image,               /* Update image */
     PixelPacket *update_pixels,        /* Pixel row in update image */
     IndexPacket *update_indexes,       /* Pixel row indexes in update image */
     const long npixels,                /* Number of pixels in row */
     ExceptionInfo *exception           /* Exception report */
     );

  extern MagickExport MagickPassFail
  PixelIterateDualModify(PixelIteratorDualModifyCallback call_back,
                         const PixelIteratorOptions *options,
                         const char *description,
                         void *mutable_data,
                         const void *immutable_data,
                         const unsigned long columns,
                         const unsigned long rows,
                         const Image *source_image,
                         const long source_x,
                         const long source_y,
                         Image *update_image,
                         const long update_x,
                         const long update_y,
                         ExceptionInfo *exception);

  /*
    Read-write access across pixel regions of two images. The first
    (source) image is accessed read-only while the second (new)
    image is accessed for write (uninitialized pixels).
  */
  typedef PixelIteratorDualModifyCallback PixelIteratorDualNewCallback;

  extern MagickExport MagickPassFail
  PixelIterateDualNew(PixelIteratorDualNewCallback call_back,
                      const PixelIteratorOptions *options,
                      const char *description,
                      void *mutable_data,
                      const void *immutable_data,
                      const unsigned long columns,
                      const unsigned long rows,
                      const Image *source_image,
                      const long source_x,
                      const long source_y,
                      Image *new_image,
                      const long new_x,
                      const long new_y,
                      ExceptionInfo *exception);

  /*
    Read-read-write access across pixel regions of three images. The
    first two images are accessed read-only while the third is
    accessed as read-write.
  */

  typedef MagickPassFail (*PixelIteratorTripleModifyCallback)
    (
     void *mutable_data,                 /* User provided mutable data */
     const void *immutable_data,         /* User provided immutable data */
     const Image *source1_image,         /* Source 1 image */
     const PixelPacket *source1_pixels,  /* Pixel row in source 1 image */
     const IndexPacket *source1_indexes, /* Pixel row indexes in source 1 image */
     const Image *source2_image,         /* Source 2 image */
     const PixelPacket *source2_pixels,  /* Pixel row in source 2 image */
     const IndexPacket *source2_indexes, /* Pixel row indexes in source 2 image */
     Image *update_image,                /* Update image */
     PixelPacket *update_pixels,         /* Pixel row in update image */
     IndexPacket *update_indexes,        /* Pixel row indexes in update image */
     const long npixels,                 /* Number of pixels in row */
     ExceptionInfo *exception            /* Exception report */
     );

  extern MagickExport MagickPassFail
  PixelIterateTripleModify(PixelIteratorTripleModifyCallback call_back,
                           const PixelIteratorOptions *options,
                           const char *description,
                           void *mutable_data,
                           const void *immutable_data,
                           const unsigned long columns,
                           const unsigned long rows,
                           const Image *source1_image,
                           const Image *source2_image,
                           const long source_x,
                           const long source_y,
                           Image *update_image,
                           const long update_x,
                           const long update_y,
                           ExceptionInfo *exception);

  /*
    Read-write access across pixel regions of two images. The first
    (source) image is accessed read-only while the second (new)
    image is accessed for write (uninitialized pixels).
  */
  typedef PixelIteratorTripleModifyCallback PixelIteratorTripleNewCallback;

  extern MagickExport MagickPassFail
  PixelIterateTripleNew(PixelIteratorTripleNewCallback call_back,
                        const PixelIteratorOptions *options,
                        const char *description,
                        void *mutable_data,
                        const void *immutable_data,
                        const unsigned long columns,
                        const unsigned long rows,
                        const Image *source1_image,
                        const Image *source2_image,
                        const long source_x,
                        const long source_y,
                        Image *new_image,
                        const long new_x,
                        const long new_y,
                        ExceptionInfo *exception);

#if defined(__cplusplus) || defined(c_plusplus)
}
#endif

#endif /* _PIXEL_ROW_ITERATOR_H */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 2
 * fill-column: 78
 * End:
 */
