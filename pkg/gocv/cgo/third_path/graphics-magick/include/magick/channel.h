/*
  Copyright (C) 2004 GraphicsMagick Group

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

*/
#ifndef _MAGICK_CHANNEL_H
#define _MAGICK_CHANNEL_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif


extern MagickExport Image
  *ExportImageChannel(const Image *image,
                      const ChannelType channel,
                      ExceptionInfo *exception);

extern MagickExport unsigned int
  GetImageChannelDepth(const Image *image,
                       const ChannelType channel,
                       ExceptionInfo *exception);

extern MagickExport MagickPassFail
  ChannelImage(Image *image,const ChannelType channel),
  ImportImageChannel(const Image *src_image,
                     Image *dst_image,
                     const ChannelType channel),
  ImportImageChannelsMasked(const Image *source_image,
                            Image *update_image,
                            const ChannelType channels),
  SetImageChannelDepth(Image *image,
                       const ChannelType channel,
                       const unsigned int depth);


#if defined(__cplusplus) || defined(c_plusplus)
}
#endif

#endif /* _MAGICK_CHANNEL_H */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 2
 * fill-column: 78
 * End:
 */
