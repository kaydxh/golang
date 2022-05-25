// This may look like C code, but it is really -*- C++ -*-
//
// Copyright Bob Friesenhahn, 1999 - 2014
//
// Geometry Definition
//
// Representation of an ImageMagick geometry specification
// X11 geometry specification plus hints

#if !defined (Magick_Geometry_header)
#define Magick_Geometry_header

#include "Magick++/Include.h"
#include <string>

namespace Magick
{

  class MagickDLLDecl Geometry;

  // Compare two Geometry objects regardless of LHS/RHS
  int MagickDLLDecl operator == ( const Magick::Geometry& left_, const Magick::Geometry& right_ );
  int MagickDLLDecl operator != ( const Magick::Geometry& left_, const Magick::Geometry& right_ );
  int MagickDLLDecl operator >  ( const Magick::Geometry& left_, const Magick::Geometry& right_ );
  int MagickDLLDecl operator <  ( const Magick::Geometry& left_, const Magick::Geometry& right_ );
  int MagickDLLDecl operator >= ( const Magick::Geometry& left_, const Magick::Geometry& right_ );
  int MagickDLLDecl operator <= ( const Magick::Geometry& left_, const Magick::Geometry& right_ );

  class MagickDLLDecl Geometry
  {
  public:

    Geometry ( unsigned int width_,
               unsigned int height_,
               unsigned int xOff_ = 0,
               unsigned int yOff_ = 0,
               bool xNegative_ = false,
               bool yNegative_ = false );
    Geometry ( const std::string &geometry_ );
    Geometry ( const char * geometry_ );
    Geometry ( const Geometry &geometry_ );
    Geometry ( );
    ~Geometry ( void );

    // Width
    void          width ( unsigned int width_ );
    unsigned int  width ( void ) const;

    // Height
    void          height ( unsigned int height_ );
    unsigned int  height ( void ) const;

    // X offset from origin
    void          xOff ( unsigned int xOff_ );
    unsigned int  xOff ( void ) const;

    // Y offset from origin
    void          yOff ( unsigned int yOff_ );
    unsigned int  yOff ( void ) const;

    // Sign of X offset negative? (X origin at right)
    void          xNegative ( bool xNegative_ );
    bool          xNegative ( void ) const;

    // Sign of Y offset negative? (Y origin at bottom)
    void          yNegative ( bool yNegative_ );
    bool          yNegative ( void ) const;

    // Width and height are expressed as percentages
    void          percent ( bool percent_ );
    bool          percent ( void ) const;

    // Resize without preserving aspect ratio (!)
    void          aspect ( bool aspect_ );
    bool          aspect ( void ) const;

    // Resize if image is greater than size (>)
    void          greater ( bool greater_ );
    bool          greater ( void ) const;

    // Resize if image is less than size (<)
    void          less ( bool less_ );
    bool          less ( void ) const;

    // Resize image to fit total pixel area specified by dimensions (@).
    void          limitPixels ( bool limitPixels_ );
    bool          limitPixels ( void ) const;

    // Dimensions are treated as minimum rather than maximum values (^)
    void          fillArea ( bool fillArea_ );
    bool          fillArea ( void ) const;

    // Does object contain valid geometry?
    void          isValid ( bool isValid_ );
    bool          isValid ( void ) const;

    // Set via geometry string
    const Geometry& operator = ( const std::string &geometry_ );
    const Geometry& operator = ( const char * geometry_ );

    // Assignment operator
    Geometry& operator= ( const Geometry& Geometry_ );

    // Return geometry string
    operator std::string() const;

    //
    // Public methods below this point are for Magick++ use only.
    //

    // Construct from RectangleInfo
    Geometry ( const MagickLib::RectangleInfo &rectangle_ );

    // Return an ImageMagick RectangleInfo struct
    operator MagickLib::RectangleInfo() const;

  private:
    unsigned int  _width;
    unsigned int  _height;
    unsigned int  _xOff;
    unsigned int  _yOff;
    union
    {
      struct
      {
        // Bit-field for compact boolean storage
        bool          _xNegative : 1;
        bool          _yNegative : 1;
        bool          _isValid : 1;
        bool          _percent : 1;    // Interpret width & height as percentages (%)
        bool          _aspect : 1;     // Force exact size (!)
        bool          _greater : 1;    // Re-size only if larger than geometry (>)
        bool          _less : 1;       // Re-size only if smaller than geometry (<)
        bool          _limitPixels : 1;// Resize image to fit total pixel area (@).
        bool          _fillArea : 1;   // Dimensions are treated as
                                       // minimum rather than maximum
                                       // values (^)
      } _b;
      struct
      {
        // Padding for future use.
        unsigned int pad[2];
      } _padding;
    } _flags; // union
  }; // class Geometry;
} // namespace Magick

//
// Inlines
//


#endif // Magick_Geometry_header
