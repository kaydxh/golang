/*
  Copyright (C) 2004 - 2009 GraphicsMagick Group

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  GraphicsMagick Methods For Manipulating Embedded Image Profiles.
*/
#ifndef _MAGICK_PROFILE_H
#define _MAGICK_PROFILE_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif /* defined(__cplusplus) || defined(c_plusplus) */

/*
  Retrieve a profile from the image by name.
*/
extern MagickExport const unsigned char
  *GetImageProfile(const Image* image, const char *name, size_t *length);

/*
  Remove a profile from the image by name.
*/
extern MagickExport MagickPassFail
  DeleteImageProfile(Image *image,const char *name);

/*
  Apply (or add) a profile to the image.
*/
extern MagickExport MagickPassFail
  ProfileImage(Image *image,const char *name,unsigned char *profile,
               const size_t length,MagickBool clone);

/*
  Add (or replace) profile to the image by name.
*/
extern MagickExport MagickPassFail
  SetImageProfile(Image *image,const char *name,const unsigned char *profile,
                  const size_t length);

/*
  Add (or append) profile to the image by name.
 */
  extern MagickExport MagickPassFail
  AppendImageProfile(Image *image,const char *name,
                     const unsigned char *profile_chunk,
                     const size_t chunk_length);

/*
  Generic iterator for traversing profiles.
*/
typedef void *ImageProfileIterator;

/*
  Allocate an image profile iterator which points to one before the
  list so NextImageProfile() must be used to advance to first entry.
*/
extern MagickExport ImageProfileIterator
  AllocateImageProfileIterator(const Image *image);

/*
  Advance to next image profile.  Name, profile, and length are
  updated with information on current profile. MagickFail is returned
  when there are no more entries.
*/
extern MagickExport MagickPassFail
  NextImageProfile(ImageProfileIterator profile_iterator,const char **name,
                   const unsigned char **profile,size_t *length);

/*
  Deallocate profile iterator.
*/
extern MagickExport void
  DeallocateImageProfileIterator(ImageProfileIterator profile_iterator);

#if defined(__cplusplus) || defined(c_plusplus)
}
#endif /* defined(__cplusplus) || defined(c_plusplus) */

#endif /* _MAGICK_PROFILE_H */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 2
 * fill-column: 78
 * End:
 */
