/*
  Copyright (C) 2003 GraphicsMagick Group
  Copyright (C) 2002 ImageMagick Studio

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  Digital signature methods.
*/
#ifndef _MAGICK_SIGNATURE_H
#define _MAGICK_SIGNATURE_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif

/*
  Define declarations.
*/
#define SignatureSize  64

/*
  Typedef declarations.
*/
typedef struct _SignatureInfo
{
  unsigned long
    digest[8],
    low_order,
    high_order;

  long
    offset;

  unsigned char
    message[SignatureSize];
} SignatureInfo;

/*
  Method declarations.
*/
extern MagickExport unsigned int
  SignatureImage(Image *);

extern MagickExport void
  FinalizeSignature(SignatureInfo *),
  GetSignatureInfo(SignatureInfo *),
  TransformSignature(SignatureInfo *),
  UpdateSignature(SignatureInfo *,const unsigned char *,const size_t);

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
