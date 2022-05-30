// This may look like C code, but it is really -*- C++ -*-
//
// Copyright Bob Friesenhahn, 1999 - 2018
//
// Inclusion of GraphicsMagick headers (with namespace magic)

#ifndef Magick_Include_header
#define Magick_Include_header

#if !defined(_MAGICK_CONFIG_H)
# define _MAGICK_CONFIG_H
# if !defined(vms) && !defined(macintosh)
#  include "magick/magick_config.h"
# else
#  include "magick_config.h"
# endif
# undef inline // Remove possible definition from config.h
# undef class
#endif

#include <stdio.h>
#include <stdarg.h>
#include <stdlib.h>
#include <time.h>
#include <sys/types.h> /* POSIX 1990 header and declares size_t and ssize_t */

#if defined(__BORLANDC__)
# include <vcl.h> /* Borland C++ Builder 4.0 requirement */
#endif // defined(__BORLANDC__)

//
// Include GraphicsMagick headers into namespace "MagickLib". If
// MAGICK_IMPLEMENTATION is defined, include GraphicsMagick development
// headers.  This scheme minimizes the possibility of conflict with
// user code.
//
namespace MagickLib
{
#include <magick/api.h>
#undef inline // Remove possible definition from config.h

#undef class
}

//
// Provide appropriate DLL imports/exports for Visual C++,
// Borland C++Builder and MinGW builds.
//
#if (defined(WIN32) || defined(WIN64)) && !defined (__CYGWIN__) //&& !defined(__MINGW32__)
# define MagickCplusPlusDLLSupported
#endif
#if defined(MagickCplusPlusDLLSupported)
#  if defined(_MT) && defined(_DLL) && !defined(_LIB) && !defined(STATIC_MAGICK)
//
// In a native Windows build, the following defines are used:
//
//   _MT         = Multithreaded
//   _DLL        = Using code is part of a DLL
//   _LIB        = Using code is being built as a library.
//   _MAGICKMOD_ = Build uses loadable modules (Magick++ does not care about this)
//
// In the case where GraphicsMagick is built as a static library but the
// using code is dynamic, STATIC_MAGICK may be defined in the project to
// override triggering dynamic library behavior.
//
#    if defined(_VISUALC_)
#      define MagickDLLExplicitTemplate     /* Explicit template instantiation in DLLs */
#      pragma warning( disable: 4273 )      /* Disable the stupid dll linkage warnings */
#      pragma warning( disable: 4251 )
#    endif
#    if !defined(MAGICK_IMPLEMENTATION)
#      define MagickDLLDecl __declspec(dllimport)
#      define MagickDLLDeclExtern extern __declspec(dllimport)
#      if defined(_VISUALC_)
#        pragma message( "Magick++ lib DLL import" )
#      endif
#    else
#      if defined(__BORLANDC__) || defined(__MINGW32__)
#        define MagickDLLDecl __declspec(dllexport)
#        define MagickDLLDeclExtern __declspec(dllexport)
#        if defined(__BORLANDC__)
#          pragma message( "BCBMagick++ lib DLL export" )
#        endif
#      else
#        define MagickDLLDecl __declspec(dllexport)
#        define MagickDLLDeclExtern extern __declspec(dllexport)
#      endif
#      if defined(_VISUALC_)
#        pragma message( "Magick++ lib DLL export" )
#      endif
#    endif
#  else
#    define MagickDLLDecl
#    define MagickDLLDeclExtern
#    if defined(_VISUALC_)
#      pragma message( "Magick++ lib static interface" )
#    endif
#    if defined(_MSC_VER) && defined(STATIC_MAGICK) && !defined(NOAUTOLINK_MAGICK)
#      if defined(_DEBUG)
#        if defined(HasBZLIB)
#          pragma comment(lib, "CORE_DB_bzlib_.lib")
#        endif
#        pragma comment(lib, "CORE_DB_coders_.lib")
#        pragma comment(lib, "CORE_DB_filters_.lib")
#        if defined(HasJBIG)
#          pragma comment(lib, "CORE_DB_jbig_.lib")
#        endif
#        if defined(HasJP2)
#          pragma comment(lib, "CORE_DB_jp2_.lib")
#        endif
#        if defined(HasJPEG)
#          pragma comment(lib, "CORE_DB_jpeg_.lib")
#        endif
#        if defined(HasLCMS)
#          pragma comment(lib, "CORE_DB_lcms_.lib")
#        endif
#        if defined(HasXML)
#          pragma comment(lib, "CORE_DB_libxml_.lib")
#        endif
#        pragma comment(lib, "CORE_DB_magick_.lib")
#        pragma comment(lib, "CORE_DB_Magick++_.lib")
#        if defined(HasPNG)
#          pragma comment(lib, "CORE_DB_png_.lib")
#        endif
#        if defined(HasTIFF)
#          pragma comment(lib, "CORE_DB_tiff_.lib")
#        endif
#        if defined(HasTTF)
#          pragma comment(lib, "CORE_DB_ttf_.lib")
#        endif
#        pragma comment(lib, "CORE_DB_wand_.lib")
#        if defined(HasWEBP)
#          pragma comment(lib, "CORE_DB_webp_.lib")
#        endif
#        if defined(HasWMFlite)
#          pragma comment(lib, "CORE_DB_wmf_.lib")
#        endif
#        if defined(HasX11)
#          pragma comment(lib, "CORE_DB_xlib_.lib")
#        endif
#        if defined(HasZLIB)
#          pragma comment(lib, "CORE_DB_zlib_.lib")
#        endif
#      else
#        if defined(HasBZLIB)
#          pragma comment(lib, "CORE_RL_bzlib_.lib")
#        endif
#        pragma comment(lib, "CORE_RL_coders_.lib")
#        pragma comment(lib, "CORE_RL_filters_.lib")
#        if defined(HasJBIG)
#          pragma comment(lib, "CORE_RL_jbig_.lib")
#        endif
#        if defined(HasJP2)
#          pragma comment(lib, "CORE_RL_jp2_.lib")
#        endif
#        if defined(HasJPEG)
#          pragma comment(lib, "CORE_RL_jpeg_.lib")
#        endif
#        if defined(HasLCMS)
#          pragma comment(lib, "CORE_RL_lcms_.lib")
#        endif
#        if defined(HasXML)
#          pragma comment(lib, "CORE_RL_libxml_.lib")
#        endif
#        pragma comment(lib, "CORE_RL_magick_.lib")
#        pragma comment(lib, "CORE_RL_Magick++_.lib")
#        if defined(HasPNG)
#          pragma comment(lib, "CORE_RL_png_.lib")
#        endif
#        if defined(HasTIFF)
#          pragma comment(lib, "CORE_RL_tiff_.lib")
#        endif
#        if defined(HasTTF)
#          pragma comment(lib, "CORE_RL_ttf_.lib")
#        endif
#        pragma comment(lib, "CORE_RL_wand_.lib")
#        if defined(HasWEBP)
#          pragma comment(lib, "CORE_RL_webp_.lib")
#        endif
#        if defined(HasWMFlite)
#          pragma comment(lib, "CORE_RL_wmf_.lib")
#        endif
#        if defined(HasX11)
#          pragma comment(lib, "CORE_RL_xlib_.lib")
#        endif
#        if defined(HasZLIB)
#          pragma comment(lib, "CORE_RL_zlib_.lib")
#        endif
#      endif
#      if defined(_WIN32_WCE)
#        pragma comment(lib, "wsock32.lib")
#      else
#        pragma comment(lib, "ws2_32.lib")
#      endif
#    endif
#  endif
#else
#  define MagickDLLDecl
#  define MagickDLLDeclExtern
#endif

#if (defined(WIN32) || defined(WIN64)) && defined(_VISUALC_)
#  pragma warning(disable : 4996) /* function deprecation warnings */
#endif

#if defined(MAGICK_IMPLEMENTATION)
namespace MagickLib
{
#  include "magick/enum_strings.h"
}
#endif

//
// Import GraphicsMagick symbols and types which are used as part of the
// Magick++ API definition into namespace "Magick".
//
namespace Magick
{
  // The datatype for an RGB component
  using MagickLib::Quantum;

  // Image class types
  using MagickLib::ClassType;
  using MagickLib::UndefinedClass;
  using MagickLib::DirectClass;
  using MagickLib::PseudoClass;

  // Channel types
  using MagickLib::ChannelType;
  using MagickLib::UndefinedChannel;
  using MagickLib::RedChannel;
  using MagickLib::CyanChannel;
  using MagickLib::GreenChannel;
  using MagickLib::MagentaChannel;
  using MagickLib::BlueChannel;
  using MagickLib::YellowChannel;
  using MagickLib::OpacityChannel;
  using MagickLib::BlackChannel;
  using MagickLib::MatteChannel;
  using MagickLib::AllChannels;
  using MagickLib::GrayChannel;

  // Color-space types
  using MagickLib::ColorspaceType;
  using MagickLib::UndefinedColorspace;
  using MagickLib::RGBColorspace;
  using MagickLib::GRAYColorspace;
  using MagickLib::TransparentColorspace;
  using MagickLib::OHTAColorspace;
  using MagickLib::XYZColorspace;
  using MagickLib::YCbCrColorspace;
  using MagickLib::YCCColorspace;
  using MagickLib::YIQColorspace;
  using MagickLib::YPbPrColorspace;
  using MagickLib::YUVColorspace;
  using MagickLib::CMYKColorspace;
  using MagickLib::sRGBColorspace;
  using MagickLib::HSLColorspace;
  using MagickLib::HWBColorspace;
  using MagickLib::LABColorspace;
  using MagickLib::CineonLogRGBColorspace;
  using MagickLib::Rec601LumaColorspace;
  using MagickLib::Rec709LumaColorspace;
  using MagickLib::Rec709YCbCrColorspace;

  // Composition operations
  using MagickLib::AddCompositeOp;
  using MagickLib::AtopCompositeOp;
  using MagickLib::BumpmapCompositeOp;
  using MagickLib::ClearCompositeOp;
  using MagickLib::ColorizeCompositeOp;
  using MagickLib::CompositeOperator;
  using MagickLib::CopyBlueCompositeOp;
  using MagickLib::CopyCompositeOp;
  using MagickLib::CopyGreenCompositeOp;
  using MagickLib::CopyOpacityCompositeOp;
  using MagickLib::CopyRedCompositeOp;
  using MagickLib::DarkenCompositeOp;
  using MagickLib::DifferenceCompositeOp;
  using MagickLib::DisplaceCompositeOp;
  using MagickLib::DissolveCompositeOp;
  using MagickLib::HueCompositeOp;
  using MagickLib::InCompositeOp;
  using MagickLib::LightenCompositeOp;
  using MagickLib::LuminizeCompositeOp;
  using MagickLib::MinusCompositeOp;
  using MagickLib::ModulateCompositeOp;
  using MagickLib::MultiplyCompositeOp;
  using MagickLib::NoCompositeOp;
  using MagickLib::OutCompositeOp;
  using MagickLib::OverCompositeOp;
  using MagickLib::OverlayCompositeOp;
  using MagickLib::PlusCompositeOp;
  using MagickLib::SaturateCompositeOp;
  using MagickLib::ScreenCompositeOp;
  using MagickLib::SubtractCompositeOp;
  using MagickLib::ThresholdCompositeOp;
  using MagickLib::UndefinedCompositeOp;
  using MagickLib::XorCompositeOp;
  using MagickLib::CopyCyanCompositeOp;
  using MagickLib::CopyMagentaCompositeOp;
  using MagickLib::CopyYellowCompositeOp;
  using MagickLib::CopyBlackCompositeOp;
  using MagickLib::DivideCompositeOp;
  using MagickLib::HardLightCompositeOp;
  using MagickLib::ExclusionCompositeOp;
  using MagickLib::ColorDodgeCompositeOp;
  using MagickLib::ColorBurnCompositeOp;
  using MagickLib::SoftLightCompositeOp;
  using MagickLib::LinearBurnCompositeOp;
  using MagickLib::LinearDodgeCompositeOp;
  using MagickLib::LinearLightCompositeOp;
  using MagickLib::VividLightCompositeOp;
  using MagickLib::PinLightCompositeOp;
  using MagickLib::HardMixCompositeOp;

  // Compression algorithms
  using MagickLib::CompressionType;
  using MagickLib::UndefinedCompression;
  using MagickLib::NoCompression;
  using MagickLib::BZipCompression;
  using MagickLib::FaxCompression;
  using MagickLib::Group3Compression;
  using MagickLib::Group4Compression;
  using MagickLib::JPEGCompression;
  using MagickLib::LZWCompression;
  using MagickLib::RLECompression;
  using MagickLib::ZipCompression;
  using MagickLib::LZMACompression;
  using MagickLib::JPEG2000Compression;
  using MagickLib::JBIG1Compression;
  using MagickLib::JBIG2Compression;
  using MagickLib::ZSTDCompression;
  using MagickLib::WebPCompression;

  using MagickLib::DisposeType;
  using MagickLib::UndefinedDispose;
  using MagickLib::NoneDispose;
  using MagickLib::BackgroundDispose;
  using MagickLib::PreviousDispose;

  // Endian options
  using MagickLib::EndianType;
  using MagickLib::UndefinedEndian;
  using MagickLib::LSBEndian;
  using MagickLib::MSBEndian;
  using MagickLib::NativeEndian;

  // Exception types
  using MagickLib::ExceptionType;
  using MagickLib::UndefinedException;
  using MagickLib::EventException;
  using MagickLib::ExceptionEvent;
  using MagickLib::ResourceEvent;
  using MagickLib::ResourceLimitEvent;
  using MagickLib::TypeEvent;
  using MagickLib::AnnotateEvent;
  using MagickLib::OptionEvent;
  using MagickLib::DelegateEvent;
  using MagickLib::MissingDelegateEvent;
  using MagickLib::CorruptImageEvent;
  using MagickLib::FileOpenEvent;
  using MagickLib::BlobEvent;
  using MagickLib::StreamEvent;
  using MagickLib::CacheEvent;
  using MagickLib::CoderEvent;
  using MagickLib::ModuleEvent;
  using MagickLib::DrawEvent;
  using MagickLib::RenderEvent;
  using MagickLib::ImageEvent;
  using MagickLib::WandEvent;
  using MagickLib::TemporaryFileEvent;
  using MagickLib::TransformEvent;
  using MagickLib::XServerEvent;
  using MagickLib::X11Event;
  using MagickLib::UserEvent;
  using MagickLib::MonitorEvent;
  using MagickLib::LocaleEvent;
  using MagickLib::DeprecateEvent;
  using MagickLib::RegistryEvent;
  using MagickLib::ConfigureEvent;
  using MagickLib::WarningException;
  using MagickLib::ExceptionWarning;
  using MagickLib::ResourceWarning;
  using MagickLib::ResourceLimitWarning;
  using MagickLib::TypeWarning;
  using MagickLib::AnnotateWarning;
  using MagickLib::OptionWarning;
  using MagickLib::DelegateWarning;
  using MagickLib::MissingDelegateWarning;
  using MagickLib::CorruptImageWarning;
  using MagickLib::FileOpenWarning;
  using MagickLib::BlobWarning;
  using MagickLib::StreamWarning;
  using MagickLib::CacheWarning;
  using MagickLib::CoderWarning;
  using MagickLib::ModuleWarning;
  using MagickLib::DrawWarning;
  using MagickLib::RenderWarning;
  using MagickLib::ImageWarning;
  using MagickLib::WandWarning;
  using MagickLib::TemporaryFileWarning;
  using MagickLib::TransformWarning;
  using MagickLib::XServerWarning;
  using MagickLib::X11Warning;
  using MagickLib::UserWarning;
  using MagickLib::MonitorWarning;
  using MagickLib::LocaleWarning;
  using MagickLib::DeprecateWarning;
  using MagickLib::RegistryWarning;
  using MagickLib::ConfigureWarning;
  using MagickLib::ErrorException;
  using MagickLib::ExceptionError;
  using MagickLib::ResourceError;
  using MagickLib::ResourceLimitError;
  using MagickLib::TypeError;
  using MagickLib::AnnotateError;
  using MagickLib::OptionError;
  using MagickLib::DelegateError;
  using MagickLib::MissingDelegateError;
  using MagickLib::CorruptImageError;
  using MagickLib::FileOpenError;
  using MagickLib::BlobError;
  using MagickLib::StreamError;
  using MagickLib::CacheError;
  using MagickLib::CoderError;
  using MagickLib::ModuleError;
  using MagickLib::DrawError;
  using MagickLib::RenderError;
  using MagickLib::ImageError;
  using MagickLib::WandError;
  using MagickLib::TemporaryFileError;
  using MagickLib::TransformError;
  using MagickLib::XServerError;
  using MagickLib::X11Error;
  using MagickLib::UserError;
  using MagickLib::MonitorError;
  using MagickLib::LocaleError;
  using MagickLib::DeprecateError;
  using MagickLib::RegistryError;
  using MagickLib::ConfigureError;
  using MagickLib::FatalErrorException;
  using MagickLib::ExceptionFatalError;
  using MagickLib::ResourceFatalError;
  using MagickLib::ResourceLimitFatalError;
  using MagickLib::TypeFatalError;
  using MagickLib::AnnotateFatalError;
  using MagickLib::OptionFatalError;
  using MagickLib::DelegateFatalError;
  using MagickLib::MissingDelegateFatalError;
  using MagickLib::CorruptImageFatalError;
  using MagickLib::FileOpenFatalError;
  using MagickLib::BlobFatalError;
  using MagickLib::StreamFatalError;
  using MagickLib::CacheFatalError;
  using MagickLib::CoderFatalError;
  using MagickLib::ModuleFatalError;
  using MagickLib::DrawFatalError;
  using MagickLib::RenderFatalError;
  using MagickLib::ImageFatalError;
  using MagickLib::WandFatalError;
  using MagickLib::TemporaryFileFatalError;
  using MagickLib::TransformFatalError;
  using MagickLib::XServerFatalError;
  using MagickLib::X11FatalError;
  using MagickLib::UserFatalError;
  using MagickLib::MonitorFatalError;
  using MagickLib::LocaleFatalError;
  using MagickLib::DeprecateFatalError;
  using MagickLib::RegistryFatalError;
  using MagickLib::ConfigureFatalError;

  // Fill rules
  using MagickLib::FillRule;
  using MagickLib::UndefinedRule;
  using MagickLib::EvenOddRule;
  using MagickLib::NonZeroRule;

  // Filter types
  using MagickLib::FilterTypes;
  using MagickLib::UndefinedFilter;
  using MagickLib::PointFilter;
  using MagickLib::BoxFilter;
  using MagickLib::TriangleFilter;
  using MagickLib::HermiteFilter;
  using MagickLib::HanningFilter;
  using MagickLib::HammingFilter;
  using MagickLib::BlackmanFilter;
  using MagickLib::GaussianFilter;
  using MagickLib::QuadraticFilter;
  using MagickLib::CubicFilter;
  using MagickLib::CatromFilter;
  using MagickLib::MitchellFilter;
  using MagickLib::LanczosFilter;
  using MagickLib::BesselFilter;
  using MagickLib::SincFilter;

  // Bit gravity
  using MagickLib::GravityType;
  using MagickLib::ForgetGravity;
  using MagickLib::NorthWestGravity;
  using MagickLib::NorthGravity;
  using MagickLib::NorthEastGravity;
  using MagickLib::WestGravity;
  using MagickLib::CenterGravity;
  using MagickLib::EastGravity;
  using MagickLib::SouthWestGravity;
  using MagickLib::SouthGravity;
  using MagickLib::SouthEastGravity;
  using MagickLib::StaticGravity;

  // Image types
  using MagickLib::ImageType;
  using MagickLib::UndefinedType;
  using MagickLib::BilevelType;
  using MagickLib::GrayscaleType;
  using MagickLib::GrayscaleMatteType;
  using MagickLib::PaletteType;
  using MagickLib::PaletteMatteType;
  using MagickLib::TrueColorType;
  using MagickLib::TrueColorMatteType;
  using MagickLib::ColorSeparationType;
  using MagickLib::ColorSeparationMatteType;
  using MagickLib::OptimizeType;

  // Interlace types
  using MagickLib::InterlaceType;
  using MagickLib::UndefinedInterlace;
  using MagickLib::NoInterlace;
  using MagickLib::LineInterlace;
  using MagickLib::PlaneInterlace;
  using MagickLib::PartitionInterlace;

  // Line cap types
  using MagickLib::LineCap;
  using MagickLib::UndefinedCap;
  using MagickLib::ButtCap;
  using MagickLib::RoundCap;
  using MagickLib::SquareCap;

  // Line join types
  using MagickLib::LineJoin;
  using MagickLib::UndefinedJoin;
  using MagickLib::MiterJoin;
  using MagickLib::RoundJoin;
  using MagickLib::BevelJoin;

  // Noise types
  using MagickLib::NoiseType;
  using MagickLib::UniformNoise;
  using MagickLib::GaussianNoise;
  using MagickLib::MultiplicativeGaussianNoise;
  using MagickLib::ImpulseNoise;
  using MagickLib::LaplacianNoise;
  using MagickLib::PoissonNoise;
  using MagickLib::RandomNoise;

  // Orientation types
  using MagickLib::OrientationType;
  using MagickLib::UndefinedOrientation;
  using MagickLib::TopLeftOrientation;
  using MagickLib::TopRightOrientation;
  using MagickLib::BottomRightOrientation;
  using MagickLib::BottomLeftOrientation;
  using MagickLib::LeftTopOrientation;
  using MagickLib::RightTopOrientation;
  using MagickLib::RightBottomOrientation;
  using MagickLib::LeftBottomOrientation;

  // Paint methods
  using MagickLib::PaintMethod;
  using MagickLib::PointMethod;
  using MagickLib::ReplaceMethod;
  using MagickLib::FloodfillMethod;
  using MagickLib::FillToBorderMethod;
  using MagickLib::ResetMethod;

  // Arithmetic and bitwise operators
  using MagickLib::UndefinedQuantumOp;
  using MagickLib::AddQuantumOp;
  using MagickLib::AndQuantumOp;
  using MagickLib::AssignQuantumOp;
  using MagickLib::DivideQuantumOp;
  using MagickLib::LShiftQuantumOp;
  using MagickLib::MultiplyQuantumOp;
  using MagickLib::OrQuantumOp;
  using MagickLib::RShiftQuantumOp;
  using MagickLib::SubtractQuantumOp;
  using MagickLib::ThresholdQuantumOp;
  using MagickLib::ThresholdBlackQuantumOp;
  using MagickLib::ThresholdWhiteQuantumOp;
  using MagickLib::ThresholdBlackNegateQuantumOp;
  using MagickLib::ThresholdWhiteNegateQuantumOp;
  using MagickLib::XorQuantumOp;
  using MagickLib::NoiseGaussianQuantumOp;
  using MagickLib::NoiseImpulseQuantumOp;
  using MagickLib::NoiseLaplacianQuantumOp;
  using MagickLib::NoiseMultiplicativeQuantumOp;
  using MagickLib::NoisePoissonQuantumOp;
  using MagickLib::NoiseUniformQuantumOp;
  using MagickLib::NegateQuantumOp;
  using MagickLib::GammaQuantumOp;
  using MagickLib::DepthQuantumOp;
  using MagickLib::LogQuantumOp;
  using MagickLib::MaxQuantumOp;
  using MagickLib::MinQuantumOp;
  using MagickLib::PowQuantumOp;
  using MagickLib::QuantumOperator;

  // Preview types.  Not currently used by Magick++
  using MagickLib::PreviewType;
  using MagickLib::UndefinedPreview;
  using MagickLib::RotatePreview;
  using MagickLib::ShearPreview;
  using MagickLib::RollPreview;
  using MagickLib::HuePreview;
  using MagickLib::SaturationPreview;
  using MagickLib::BrightnessPreview;
  using MagickLib::GammaPreview;
  using MagickLib::SpiffPreview;
  using MagickLib::DullPreview;
  using MagickLib::GrayscalePreview;
  using MagickLib::QuantizePreview;
  using MagickLib::DespecklePreview;
  using MagickLib::ReduceNoisePreview;
  using MagickLib::AddNoisePreview;
  using MagickLib::SharpenPreview;
  using MagickLib::BlurPreview;
  using MagickLib::ThresholdPreview;
  using MagickLib::EdgeDetectPreview;
  using MagickLib::SpreadPreview;
  using MagickLib::SolarizePreview;
  using MagickLib::ShadePreview;
  using MagickLib::RaisePreview;
  using MagickLib::SegmentPreview;
  using MagickLib::SwirlPreview;
  using MagickLib::ImplodePreview;
  using MagickLib::WavePreview;
  using MagickLib::OilPaintPreview;
  using MagickLib::CharcoalDrawingPreview;
  using MagickLib::JPEGPreview;

  // Quantum types
  using MagickLib::QuantumType;
  using MagickLib::IndexQuantum;
  using MagickLib::GrayQuantum;
  using MagickLib::IndexAlphaQuantum;
  using MagickLib::GrayAlphaQuantum;
  using MagickLib::RedQuantum;
  using MagickLib::CyanQuantum;
  using MagickLib::GreenQuantum;
  using MagickLib::YellowQuantum;
  using MagickLib::BlueQuantum;
  using MagickLib::MagentaQuantum;
  using MagickLib::AlphaQuantum;
  using MagickLib::BlackQuantum;
  using MagickLib::RGBQuantum;
  using MagickLib::RGBAQuantum;
  using MagickLib::CMYKQuantum;
  using MagickLib::CIEYQuantum;
  using MagickLib::CIEXYZQuantum;

  // Quantum sample types
  using MagickLib::QuantumSampleType;
  using MagickLib::UndefinedQuantumSampleType;
  using MagickLib::UnsignedQuantumSampleType;
  using MagickLib::FloatQuantumSampleType;

  // Rendering intents
  using MagickLib::RenderingIntent;
  using MagickLib::UndefinedIntent;
  using MagickLib::SaturationIntent;
  using MagickLib::PerceptualIntent;
  using MagickLib::AbsoluteIntent;
  using MagickLib::RelativeIntent;

  // Resolution units
  using MagickLib::ResolutionType;
  using MagickLib::UndefinedResolution;
  using MagickLib::PixelsPerInchResolution;
  using MagickLib::PixelsPerCentimeterResolution;

  // PixelPacket structure
  using MagickLib::PixelPacket;

  // IndexPacket type
  using MagickLib::IndexPacket;

  // ImageStatistics type
  using MagickLib::ImageStatistics;

  // StorageType type
  using MagickLib::StorageType;
  using MagickLib::CharPixel;
  using MagickLib::ShortPixel;
  using MagickLib::IntegerPixel;
  using MagickLib::LongPixel;
  using MagickLib::FloatPixel;
  using MagickLib::DoublePixel;

  // StretchType type
  using MagickLib::StretchType;
  using MagickLib::NormalStretch;
  using MagickLib::UltraCondensedStretch;
  using MagickLib::ExtraCondensedStretch;
  using MagickLib::CondensedStretch;
  using MagickLib::SemiCondensedStretch;
  using MagickLib::SemiExpandedStretch;
  using MagickLib::ExpandedStretch;
  using MagickLib::ExtraExpandedStretch;
  using MagickLib::UltraExpandedStretch;
  using MagickLib::AnyStretch;

  // StyleType type
  using MagickLib::StyleType;
  using MagickLib::NormalStyle;
  using MagickLib::ItalicStyle;
  using MagickLib::ObliqueStyle;
  using MagickLib::AnyStyle;

  // Decoration types
  using MagickLib::DecorationType;
  using MagickLib::NoDecoration;
  using MagickLib::UnderlineDecoration;
  using MagickLib::OverlineDecoration;
  using MagickLib::LineThroughDecoration;

  // Resource types
  using MagickLib::ResourceType;
  using MagickLib::DiskResource;
  using MagickLib::FileResource;
  using MagickLib::MapResource;
  using MagickLib::MemoryResource;
  using MagickLib::PixelsResource;
  using MagickLib::ThreadsResource;
  using MagickLib::WidthResource;
  using MagickLib::HeightResource;

  // Virtual pixel methods
  using MagickLib::VirtualPixelMethod;
  using MagickLib::UndefinedVirtualPixelMethod;
  using MagickLib::ConstantVirtualPixelMethod;
  using MagickLib::EdgeVirtualPixelMethod;
  using MagickLib::MirrorVirtualPixelMethod;
  using MagickLib::TileVirtualPixelMethod;

#if defined(MAGICK_IMPLEMENTATION)
  //
  // GraphicsMagick symbols used in implementation code
  //
  using MagickLib::AccessDefinition;
  using MagickLib::AccessImmutableIndexes;
  using MagickLib::AccessMutableIndexes;
  using MagickLib::AcquireCacheViewPixels;
  using MagickLib::AcquireImagePixels;
  using MagickLib::AdaptiveThresholdImage;
  using MagickLib::AddDefinition;
  using MagickLib::AddDefinitions;
  using MagickLib::AddNoiseImage;
  using MagickLib::AddNoiseImageChannel;
  using MagickLib::AffineMatrix;
  using MagickLib::AffineTransformImage;
  using MagickLib::AllocateImage;
  using MagickLib::AnnotateImage;
  using MagickLib::AreaValue;
  using MagickLib::AspectValue;
  using MagickLib::Base64Decode;
  using MagickLib::Base64Encode;
  using MagickLib::BlobError;
  using MagickLib::BlobFatalError;
  using MagickLib::BlobToImage;
  using MagickLib::BlobWarning;
  using MagickLib::BlurImage;
  using MagickLib::BlurImageChannel;
  using MagickLib::BorderImage;
  using MagickLib::CacheError;
  using MagickLib::CacheFatalError;
  using MagickLib::CacheWarning;
  using MagickLib::CdlImage;
  using MagickLib::ChannelImage;
  using MagickLib::CharcoalImage;
  using MagickLib::ChopImage;
  using MagickLib::CloneDrawInfo;
  using MagickLib::CloneImage;
  using MagickLib::CloneImageInfo;
  using MagickLib::CloneQuantizeInfo;
  using MagickLib::CloseCacheView;
  using MagickLib::CoderError;
  using MagickLib::CoderFatalError;
  using MagickLib::CoderWarning;
  using MagickLib::ColorFloodfillImage;
  using MagickLib::ColorizeImage;
  using MagickLib::ColorMatrixImage;
  using MagickLib::CompositeImage;
  using MagickLib::ConfigureError;
  using MagickLib::ConfigureFatalError;
  using MagickLib::ConfigureWarning;
  using MagickLib::ConstituteImage;
  using MagickLib::ContrastImage;
  using MagickLib::ConvolveImage;
  using MagickLib::CopyException;
  using MagickLib::CorruptImageError;
  using MagickLib::CorruptImageFatalError;
  using MagickLib::CorruptImageWarning;
  using MagickLib::CropImage;
  using MagickLib::CycleColormapImage;
  using MagickLib::DelegateError;
  using MagickLib::DelegateFatalError;
  using MagickLib::DelegateWarning;
  using MagickLib::DeleteMagickRegistry;
  using MagickLib::DespeckleImage;
  using MagickLib::DestroyDrawInfo;
  using MagickLib::DestroyExceptionInfo;
  using MagickLib::DestroyImageInfo;
  using MagickLib::DestroyImageList;
  using MagickLib::DestroyMagick;
  using MagickLib::DestroyQuantizeInfo;
  using MagickLib::DispatchImage;
  using MagickLib::DisplayImages;
  using MagickLib::DrawAffine;
  using MagickLib::DrawAllocateContext;
  using MagickLib::DrawAnnotation;
  using MagickLib::DrawArc;
  using MagickLib::DrawBezier;
  using MagickLib::DrawCircle;
  using MagickLib::DrawColor;
  using MagickLib::DrawComment;
  using MagickLib::DrawComposite;
  using MagickLib::DrawContext;
  using MagickLib::DrawDestroyContext;
  using MagickLib::DrawEllipse;
  using MagickLib::DrawError;
  using MagickLib::DrawFatalError;
  using MagickLib::DrawImage;
  using MagickLib::DrawInfo;
  using MagickLib::DrawLine;
  using MagickLib::DrawMatte;
  using MagickLib::DrawPathClose;
  using MagickLib::DrawPathCurveToAbsolute;
  using MagickLib::DrawPathCurveToQuadraticBezierAbsolute;
  using MagickLib::DrawPathCurveToQuadraticBezierRelative;
  using MagickLib::DrawPathCurveToQuadraticBezierSmoothAbsolute;
  using MagickLib::DrawPathCurveToQuadraticBezierSmoothRelative;
  using MagickLib::DrawPathCurveToRelative;
  using MagickLib::DrawPathCurveToSmoothAbsolute;
  using MagickLib::DrawPathCurveToSmoothRelative;
  using MagickLib::DrawPathEllipticArcAbsolute;
  using MagickLib::DrawPathEllipticArcRelative;
  using MagickLib::DrawPathFinish;
  using MagickLib::DrawPathLineToAbsolute;
  using MagickLib::DrawPathLineToHorizontalAbsolute;
  using MagickLib::DrawPathLineToHorizontalRelative;
  using MagickLib::DrawPathLineToRelative;
  using MagickLib::DrawPathLineToVerticalAbsolute;
  using MagickLib::DrawPathLineToVerticalRelative;
  using MagickLib::DrawPathMoveToAbsolute;
  using MagickLib::DrawPathMoveToRelative;
  using MagickLib::DrawPathStart;
  using MagickLib::DrawPoint;
  using MagickLib::DrawPolygon;
  using MagickLib::DrawPolyline;
  using MagickLib::DrawPopClipPath;
  using MagickLib::DrawPopDefs;
  using MagickLib::DrawPopGraphicContext;
  using MagickLib::DrawPopPattern;
  using MagickLib::DrawPushClipPath;
  using MagickLib::DrawPushDefs;
  using MagickLib::DrawPushGraphicContext;
  using MagickLib::DrawPushPattern;
  using MagickLib::DrawRectangle;
  using MagickLib::DrawRender;
  using MagickLib::DrawRotate;
  using MagickLib::DrawRoundRectangle;
  using MagickLib::DrawScale;
  using MagickLib::DrawSetClipPath;
  using MagickLib::DrawSetClipRule;
  using MagickLib::DrawSetClipUnits;
  using MagickLib::DrawSetFillColor;
  using MagickLib::DrawSetFillColorString;
  using MagickLib::DrawSetFillOpacity;
  using MagickLib::DrawSetFillPatternURL;
  using MagickLib::DrawSetFillRule;
  using MagickLib::DrawSetFont;
  using MagickLib::DrawSetFontFamily;
  using MagickLib::DrawSetFontSize;
  using MagickLib::DrawSetFontStretch;
  using MagickLib::DrawSetFontStyle;
  using MagickLib::DrawSetFontWeight;
  using MagickLib::DrawSetGravity;
  using MagickLib::DrawSetStrokeAntialias;
  using MagickLib::DrawSetStrokeColor;
  using MagickLib::DrawSetStrokeColorString;
  using MagickLib::DrawSetStrokeDashArray;
  using MagickLib::DrawSetStrokeDashOffset;
  using MagickLib::DrawSetStrokeLineCap;
  using MagickLib::DrawSetStrokeLineJoin;
  using MagickLib::DrawSetStrokeMiterLimit;
  using MagickLib::DrawSetStrokeOpacity;
  using MagickLib::DrawSetStrokePatternURL;
  using MagickLib::DrawSetStrokeWidth;
  using MagickLib::DrawSetTextAntialias;
  using MagickLib::DrawSetTextDecoration;
  using MagickLib::DrawSetTextEncoding;
  using MagickLib::DrawSetTextUnderColor;
  using MagickLib::DrawSetTextUnderColorString;
  using MagickLib::DrawSetViewbox;
  using MagickLib::DrawSkewX;
  using MagickLib::DrawSkewY;
  using MagickLib::DrawTranslate;
  using MagickLib::DrawWarning;
  using MagickLib::EdgeImage;
  using MagickLib::EmbossImage;
  using MagickLib::EnhanceImage;
  using MagickLib::EqualizeImage;
  using MagickLib::ExceptionInfo;
  using MagickLib::ExecuteModuleProcess;
  using MagickLib::ExportImagePixelArea;
  using MagickLib::ExtentImage;
  using MagickLib::FileOpenError;
  using MagickLib::FileOpenFatalError;
  using MagickLib::FileOpenWarning;
  using MagickLib::FlattenImages;
  using MagickLib::FlipImage;
  using MagickLib::FlopImage;
  using MagickLib::FormatString;
  using MagickLib::FrameImage;
  using MagickLib::FrameInfo;
  using MagickLib::GammaImage;
  using MagickLib::GammaImage;
  using MagickLib::GaussianBlurImage;
  using MagickLib::GaussianBlurImageChannel;
  using MagickLib::GetBlobSize;
  using MagickLib::GetCacheViewIndexes;
  using MagickLib::GetCacheViewPixels;
  using MagickLib::GetColorTuple;
  using MagickLib::GetDrawInfo;
  using MagickLib::GetExceptionInfo;
  using MagickLib::GetGeometry;
  using MagickLib::GetImageAttribute;
  using MagickLib::GetImageBoundingBox;
  using MagickLib::GetImageChannelDepth;
  using MagickLib::GetImageClipMask;
  using MagickLib::GetImageDepth;
  using MagickLib::GetImageGeometry;
  using MagickLib::GetImageInfo;
  using MagickLib::GetImagePixels;
  using MagickLib::GetImageProfile;
  using MagickLib::GetImageQuantizeError;
  using MagickLib::GetImageStatistics;
  using MagickLib::GetImageType;
  using MagickLib::GetMagickGeometry;
  using MagickLib::GetMagickInfo;
  using MagickLib::GetMagickInfoArray;
  using MagickLib::GetMagickRegistry;
  using MagickLib::GetNumberColors;
  using MagickLib::GetPageGeometry;
  using MagickLib::GetQuantizeInfo;
  using MagickLib::GetTypeMetrics;
  using MagickLib::GlobExpression;
  using MagickLib::GreaterValue;
  using MagickLib::HaldClutImage;
  using MagickLib::HSLTransform;
  using MagickLib::HeightValue;
  using MagickLib::IdentityAffine;
  using MagickLib::ImageAttribute;
  using MagickLib::ImageError;
  using MagickLib::ImageFatalError;
  using MagickLib::ImageInfo;
  using MagickLib::ImageInfoRegistryType;
  using MagickLib::ImageRegistryType;
  using MagickLib::ImageToBlob;
  using MagickLib::ImageWarning;
  using MagickLib::ImplodeImage;
  using MagickLib::ImportImagePixelArea;
  using MagickLib::IsEventLogging;
  using MagickLib::IsGeometry;
  using MagickLib::IsImagesEqual;
  using MagickLib::IsSubimage;
  using MagickLib::LessValue;
  using MagickLib::LevelImage;
  using MagickLib::LevelImageChannel;
  using MagickLib::LocaleCompare;
  using MagickLib::LogMagickEvent;
  using MagickLib::MagickFree;
  using MagickLib::MagickInfo;
  using MagickLib::MagickMalloc;
  using MagickLib::MagickRealloc;
  using MagickLib::MagickStrlCpy;
  using MagickLib::MagickToMime;
  using MagickLib::MagnifyImage;
  using MagickLib::MapImage;
  using MagickLib::MatteFloodfillImage;
  using MagickLib::MedianFilterImage;
  using MagickLib::MinifyImage;
  using MagickLib::MinimumValue;
  using MagickLib::MissingDelegateError;
  using MagickLib::MissingDelegateFatalError;
  using MagickLib::MissingDelegateWarning;
  using MagickLib::ModulateImage;
  using MagickLib::ModuleError;
  using MagickLib::ModuleFatalError;
  using MagickLib::ModuleWarning;
  using MagickLib::MonitorError;
  using MagickLib::MonitorFatalError;
  using MagickLib::MonitorWarning;
  using MagickLib::MontageInfo;
  using MagickLib::MotionBlurImage;
  using MagickLib::NegateImage;
  using MagickLib::NoValue;
  using MagickLib::NoiseType;
  using MagickLib::NormalizeImage;
  using MagickLib::OilPaintImage;
  using MagickLib::OpaqueImage;
  using MagickLib::OpenCacheView;
  using MagickLib::OptionError;
  using MagickLib::OptionFatalError;
  using MagickLib::OptionWarning;
  using MagickLib::PercentValue;
  using MagickLib::PingBlob;
  using MagickLib::PingImage;
  using MagickLib::PointInfo;
  using MagickLib::PopImagePixels;
  using MagickLib::ProfileImage;
  using MagickLib::ProfileInfo;
  using MagickLib::PushImagePixels;
  using MagickLib::QuantizeImage;
  using MagickLib::QuantizeInfo;
  using MagickLib::QuantumOperatorImage;
  using MagickLib::QuantumOperatorRegionImage;
  using MagickLib::QueryColorDatabase;
  using MagickLib::RGBTransformImage;
  using MagickLib::RaiseImage;
  using MagickLib::RandomChannelThresholdImage;
  using MagickLib::ReadImage;
  using MagickLib::RectangleInfo;
  using MagickLib::RectangleInfo;
  using MagickLib::ReduceNoiseImage;
  using MagickLib::RegisterMagickInfo;
  using MagickLib::RegistryError;
  using MagickLib::RegistryFatalError;
  using MagickLib::RegistryType;
  using MagickLib::RegistryWarning;
  using MagickLib::RemoveDefinitions;
  using MagickLib::ResizeImage;
  using MagickLib::ResourceLimitError;
  using MagickLib::ResourceLimitFatalError;
  using MagickLib::ResourceLimitWarning;
  using MagickLib::RollImage;
  using MagickLib::RotateImage;
  using MagickLib::SampleImage;
  using MagickLib::ScaleImage;
  using MagickLib::SegmentImage;
  using MagickLib::SetCacheViewPixels;
  using MagickLib::SetClientName;
  using MagickLib::SetImage;
  using MagickLib::SetImageAttribute;
  using MagickLib::SetImageChannelDepth;
  using MagickLib::SetImageClipMask;
  using MagickLib::SetImageDepth;
  using MagickLib::SetImageInfo;
  using MagickLib::SetImageOpacity;
  using MagickLib::SetImagePixels;
  using MagickLib::SetImageProfile;
  using MagickLib::SetImageType;
  using MagickLib::SetLogEventMask;
  using MagickLib::SetMagickInfo;
  using MagickLib::SetMagickRegistry;
  using MagickLib::SetMagickResourceLimit;
  using MagickLib::SetMagickResourceLimit;
  using MagickLib::ShadeImage;
  using MagickLib::SharpenImage;
  using MagickLib::SharpenImageChannel;
  using MagickLib::ShaveImage;
  using MagickLib::ShearImage;
  using MagickLib::SignatureImage;
  using MagickLib::SolarizeImage;
  using MagickLib::SpreadImage;
  using MagickLib::SteganoImage;
  using MagickLib::StereoImage;
  using MagickLib::StreamError;
  using MagickLib::StreamFatalError;
  using MagickLib::StreamWarning;
  using MagickLib::SwirlImage;
  using MagickLib::SyncCacheViewPixels;
  using MagickLib::SyncImage;
  using MagickLib::SyncImagePixels;
  using MagickLib::TextureImage;
  using MagickLib::ThresholdImage;
  using MagickLib::ThrowException;
  using MagickLib::ThrowLoggedException;
  using MagickLib::ThumbnailImage;
  using MagickLib::TransformHSL;
  using MagickLib::TransformImage;
  using MagickLib::TransformRGBImage;
  using MagickLib::TransparentImage;
  using MagickLib::TypeError;
  using MagickLib::TypeFatalError;
  using MagickLib::TypeWarning;
  using MagickLib::UndefinedException;
  using MagickLib::UndefinedRegistryType;
  using MagickLib::UnregisterMagickInfo;
  using MagickLib::UnsharpMaskImage;
  using MagickLib::UnsharpMaskImageChannel;
  using MagickLib::ViewInfo;
  using MagickLib::WaveImage;
  using MagickLib::WidthValue;
  using MagickLib::WriteImage;
  using MagickLib::XNegative;
  using MagickLib::XServerError;
  using MagickLib::XServerFatalError;
  using MagickLib::XServerWarning;
  using MagickLib::XValue;
  using MagickLib::YNegative;
  using MagickLib::YValue;
  using MagickLib::ZoomImage;


#endif // MAGICK_IMPLEMENTATION

}

#endif // Magick_Include_header
