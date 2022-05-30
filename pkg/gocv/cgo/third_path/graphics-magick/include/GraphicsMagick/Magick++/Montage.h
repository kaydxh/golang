// This may look like C code, but it is really -*- C++ -*-
//
// Copyright Bob Friesenhahn, 1999, 2000, 2001, 2002, 2003
//
// Definition of Montage class used to specify montage options.
//

#if !defined(Magick_Montage_header)
#define Magick_Montage_header

#include "Magick++/Include.h"
#include <string>
#include "Magick++/Color.h"
#include "Magick++/Geometry.h"

//
// Basic (Un-framed) Montage
//
namespace Magick
{
  class MagickDLLDecl Montage
  {
  public:
    Montage( void );
    virtual ~Montage( void );

    // Specifies the background color that thumbnails are imaged upon.
    void              backgroundColor ( const Color &backgroundColor_ );
    Color             backgroundColor ( void ) const;

    // Specifies the image composition algorithm for thumbnails. This
    // controls the algorithm by which the thumbnail image is placed
    // on the background. Use of OverCompositeOp is recommended for
    // use with images that have transparency. This option may have
    // negative side-effects for images without transparency.
    void              compose ( CompositeOperator compose_ );
    CompositeOperator compose ( void ) const;

    // Specifies the image filename to be used for the generated
    // montage images. To handle the case were multiple montage images
    // are generated, a printf-style format may be embedded within the
    // filename. For example, a filename specification of
    // image%02d.miff names the montage images as image00.miff,
    // image01.miff, etc.
    void              fileName( const std::string &fileName_ );
    std::string       fileName( void ) const;

    // Specifies the fill color to use for the label text.
    void              fillColor ( const Color &fill_ );
    Color             fillColor ( void ) const;

    // Specifies the thumbnail label font.
    void              font ( const std::string &font_ );
    std::string       font ( void ) const;

    // Specifies the size of the generated thumbnail.
    void              geometry ( const Geometry &geometry_ );
    Geometry          geometry ( void ) const;

    // Specifies the thumbnail positioning within the specified
    // geometry area. If the thumbnail is smaller in any dimension
    // than the geometry, then it is placed according to this
    // specification
    void              gravity ( GravityType gravity_ );
    GravityType       gravity ( void ) const;

    // Specifies the format used for the image label. Special format
    // characters may be embedded in the format string to include
    // information about the image.
    void              label( const std::string &label_ );
    std::string       label( void ) const;

    // Specifies the pen color to use for the label text (same as fill).
    void              penColor ( const Color &pen_ );
    Color             penColor ( void ) const;

    // Specifies the thumbnail label font size.
    void              pointSize ( unsigned int pointSize_ );
    unsigned int      pointSize ( void ) const;

    // Enable/disable drop-shadow on thumbnails.
    void              shadow ( bool shadow_ );
    bool              shadow ( void ) const;

    // Specifies the stroke color to use for the label text .
    void              strokeColor ( const Color &stroke_ );
    Color             strokeColor ( void ) const;

    // Specifies a texture image to use as montage background. The
    // built-in textures "granite:" and "plasma:" are available. A
    // texture is the same as a background image.
    void              texture ( const std::string &texture_ );
    std::string       texture ( void ) const;

    // Specifies the maximum number of montage columns and rows in the
    // montage. The montage is built by filling out all cells in a row
    // before advancing to the next row. Once the montage has reached
    // the maximum number of columns and rows, a new montage image is
    // started.
    void              tile ( const Geometry &tile_ );
    Geometry          tile ( void ) const;

    // Specifies the montage title
    void              title ( const std::string &title_ );
    std::string       title ( void ) const;

    // Specifies a montage color to set transparent. This option can
    // be set the same as the background color in order for the
    // thumbnails to appear without a background when rendered on an
    // HTML page. For best effect, ensure that the transparent color
    // selected does not occur in the rendered thumbnail colors.
    void              transparentColor ( const Color &transparentColor_ );
    Color             transparentColor ( void ) const;

    //
    // Implementation methods/members
    //

    // Update elements in existing MontageInfo structure
    virtual void      updateMontageInfo ( MagickLib::MontageInfo &montageInfo_ ) const;

  protected:

  private:
    Color             _backgroundColor;   // Color that thumbnails are composed on
    CompositeOperator _compose;           // Composition algorithm to use (e.g. ReplaceCompositeOp)
    std::string       _fileName;          // Filename to save montages to
    Color             _fill;              // Fill color
    std::string       _font;              // Label font
    Geometry          _geometry;          // Thumbnail width & height plus border width & height
    GravityType       _gravity;           // Thumbnail position (e.g. SouthWestGravity)
    std::string       _label;             // Thumbnail label (applied to image prior to montage)
    unsigned int      _pointSize;         // Font point size
    bool              _shadow;            // Enable drop-shadows on thumbnails
    Color             _stroke;            // Outline color
    std::string       _texture;           // Background texture image
    Geometry          _tile;              // Thumbnail rows and colmns
    std::string       _title;             // Montage title
    Color             _transparentColor;  // Transparent color
  };

  //
  // Montage With Frames (Extends Basic Montage)
  //
  class MagickDLLDecl MontageFramed : public Montage
  {
  public:
    MontageFramed ( void );
    /* virtual */ ~MontageFramed ( void );

    // Specifies the background color within the thumbnail frame.
    void           borderColor ( const Color &borderColor_ );
    Color          borderColor ( void ) const;

    // Specifies the border (in pixels) to place between a thumbnail
    // and its surrounding frame. This option only takes effect if
    // thumbnail frames are enabled (via frameGeometry) and the
    // thumbnail geometry specification doesn't also specify the
    // thumbnail border width.
    void           borderWidth ( unsigned int borderWidth_ );
    unsigned int   borderWidth ( void ) const;

    // Specifies the geometry specification for frame to place around
    // thumbnail. If this parameter is not specified, then the montage
    // is unframed.
    void           frameGeometry ( const Geometry &frame_ );
    Geometry       frameGeometry ( void ) const;

    // Specifies the thumbnail frame color.
    void           matteColor ( const Color &matteColor_ );
    Color          matteColor ( void ) const;

    //
    // Implementation methods/members
    //

    // Update elements in existing MontageInfo structure
    /* virtual */ void updateMontageInfo ( MagickLib::MontageInfo &montageInfo_ ) const;

  protected:

  private:

    Color          _borderColor;        // Frame border color
    unsigned int   _borderWidth;        // Pixels between thumbnail and surrounding frame
    Geometry       _frame;              // Frame geometry (width & height frame thickness)
    Color          _matteColor;         // Frame foreground color
  };
} // namespace Magick

//
// Inlines
//

//
// Implementation of Montage
//

inline void Magick::Montage::backgroundColor ( const Magick::Color &backgroundColor_ )
{
  _backgroundColor = backgroundColor_;
}
inline Magick::Color Magick::Montage::backgroundColor ( void ) const
{
  return _backgroundColor;
}

inline void Magick::Montage::compose ( Magick::CompositeOperator compose_ )
{
  _compose = compose_;
}
inline Magick::CompositeOperator Magick::Montage::compose ( void ) const
{
  return _compose;
}

inline void Magick::Montage::fileName( const std::string &fileName_ )
{
  _fileName = fileName_;
}
inline std::string Magick::Montage::fileName( void ) const
{
  return _fileName;
}

inline void Magick::Montage::fillColor ( const Color &fill_ )
{
  _fill=fill_;
}
inline Magick::Color Magick::Montage::fillColor ( void ) const
{
  return _fill;
}

inline void Magick::Montage::font ( const std::string &font_ )
{
  _font = font_;
}
inline std::string Magick::Montage::font ( void ) const
{
  return _font;
}

inline void Magick::Montage::geometry ( const Magick::Geometry &geometry_ )
{
  _geometry = geometry_;
}
inline Magick::Geometry Magick::Montage::geometry ( void ) const
{
  return _geometry;
}

inline void Magick::Montage::gravity ( Magick::GravityType gravity_ )
{
  _gravity = gravity_;
}
inline Magick::GravityType Magick::Montage::gravity ( void ) const
{
  return _gravity;
}

// Apply as attribute to all images before doing montage
inline void Magick::Montage::label( const std::string &label_ )
{
  _label = label_;
}
inline std::string Magick::Montage::label( void ) const
{
  return _label;
}

inline void Magick::Montage::penColor ( const Color &pen_ )
{
  _fill=pen_;
  _stroke=Color("none");
}
inline Magick::Color Magick::Montage::penColor ( void ) const
{
  return _fill;
}

inline void Magick::Montage::pointSize ( unsigned int pointSize_ )
{
  _pointSize = pointSize_;
}
inline unsigned int Magick::Montage::pointSize ( void ) const
{
  return _pointSize;
}

inline void Magick::Montage::shadow ( bool shadow_ )
{
  _shadow = shadow_;
}
inline bool Magick::Montage::shadow ( void ) const
{
  return _shadow;
}

inline void Magick::Montage::strokeColor ( const Color &stroke_ )
{
  _stroke=stroke_;
}
inline Magick::Color Magick::Montage::strokeColor ( void ) const
{
  return _stroke;
}

inline void Magick::Montage::texture ( const std::string &texture_ )
{
  _texture = texture_;
}
inline std::string Magick::Montage::texture ( void ) const
{
  return _texture;
}

inline void Magick::Montage::tile ( const Geometry &tile_ )
{
  _tile = tile_;
}
inline Magick::Geometry Magick::Montage::tile ( void ) const
{
  return _tile;
}

inline void Magick::Montage::title ( const std::string &title_ )
{
  _title = title_;
}
inline std::string Magick::Montage::title ( void ) const
{
  return _title;
}

// Applied after the fact to montage with TransparentImage()
inline void Magick::Montage::transparentColor ( const Magick::Color &transparentColor_ )
{
  _transparentColor = transparentColor_;
}
inline Magick::Color Magick::Montage::transparentColor ( void ) const
{
  return _transparentColor;
}

//
// Implementation of MontageFramed
//

inline void Magick::MontageFramed::borderColor ( const Magick::Color &borderColor_ )
{
  _borderColor = borderColor_;
}
inline Magick::Color Magick::MontageFramed::borderColor ( void ) const
{
  return _borderColor;
}

inline void Magick::MontageFramed::borderWidth ( unsigned int borderWidth_ )
{
  _borderWidth = borderWidth_;
}
inline unsigned int Magick::MontageFramed::borderWidth ( void ) const
{
  return _borderWidth;
}

inline void Magick::MontageFramed::frameGeometry ( const Magick::Geometry &frame_ )
{
  _frame = frame_;
}
inline Magick::Geometry Magick::MontageFramed::frameGeometry ( void ) const
{
  return _frame;
}

inline void Magick::MontageFramed::matteColor ( const Magick::Color &matteColor_ )
{
  _matteColor = matteColor_;
}
inline Magick::Color Magick::MontageFramed::matteColor ( void ) const
{
  return _matteColor;
}

#endif // Magick_Montage_header
