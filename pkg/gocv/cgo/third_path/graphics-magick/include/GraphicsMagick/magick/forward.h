/*
  Copyright (C) 2007 GraphicsMagick Group

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  Forward declarations for types used in public structures.

*/
#ifndef _MAGICK_FORWARD_H
#define _MAGICK_FORWARD_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif

typedef struct _Image *ImagePtr;

typedef struct _Ascii85Info* _Ascii85InfoPtr_;

typedef struct _BlobInfo* _BlobInfoPtr_;

typedef struct _CacheInfo* _CacheInfoPtr_;

typedef struct _ImageAttribute* _ImageAttributePtr_;

typedef struct _SemaphoreInfo* _SemaphoreInfoPtr_;

typedef struct _ThreadViewSet* _ThreadViewSetPtr_;

typedef void *ViewInfo;

#if defined(__cplusplus) || defined(c_plusplus)
}
#endif

#endif /* _MAGICK_FORWARD_H */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 2
 * fill-column: 78
 * End:
 */
