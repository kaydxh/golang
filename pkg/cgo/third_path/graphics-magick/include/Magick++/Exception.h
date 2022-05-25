// This may look like C code, but it is really -*- C++ -*-
//
// Copyright Bob Friesenhahn, 1999, 2000, 2001, 2002, 2003
//
// Definition of Magick::Exception and derived classes
// Magick::Warning* and Magick::Error*.  Derived from C++ STD
// 'exception' class for convenience.
//
// These classes form part of the Magick++ user interface.
//

#if !defined(Magick_Exception_header)
#define Magick_Exception_header

#include "Magick++/Include.h"
#include <string>
#include <exception>

namespace Magick
{
  class MagickDLLDecl Exception : public std::exception
  {
  public:
    Exception( const std::string& what_ );
    Exception( const Exception& original_ );
    Exception& operator= (const Exception& original_ );
    virtual const char* what () const throw();
    virtual ~Exception ( ) throw ();

  private:
    std::string _what;
  };

  //
  // Warnings
  //

  class MagickDLLDecl Warning : public Exception
  {
  public:
    explicit Warning ( const std::string& what_ );
    ~Warning() throw ();
  };

  class MagickDLLDecl WarningUndefined : public Warning
  {
  public:
    explicit WarningUndefined ( const std::string& what_ );
    ~WarningUndefined() throw ();
  };

  class MagickDLLDecl WarningBlob: public Warning
  {
  public:
    explicit WarningBlob ( const std::string& what_ );
    ~WarningBlob() throw ();
  };

  class MagickDLLDecl WarningCache: public Warning
  {
  public:
    explicit WarningCache ( const std::string& what_ );
    ~WarningCache() throw ();
  };

  class MagickDLLDecl WarningCoder: public Warning
  {
  public:
    explicit WarningCoder ( const std::string& what_ );
    ~WarningCoder() throw ();
  };

  class MagickDLLDecl WarningConfigure: public Warning
  {
  public:
    explicit WarningConfigure ( const std::string& what_ );
    ~WarningConfigure() throw ();
  };

  class MagickDLLDecl WarningCorruptImage: public Warning
  {
  public:
    explicit WarningCorruptImage ( const std::string& what_ );
    ~WarningCorruptImage() throw ();
  };

  class MagickDLLDecl WarningDelegate : public Warning
  {
  public:
    explicit WarningDelegate ( const std::string& what_ );
    ~WarningDelegate() throw ();
  };

  class MagickDLLDecl WarningDraw : public Warning
  {
  public:
    explicit WarningDraw ( const std::string& what_ );
    ~WarningDraw() throw ();
  };

  class MagickDLLDecl WarningFileOpen: public Warning
  {
  public:
    explicit WarningFileOpen ( const std::string& what_ );
    ~WarningFileOpen() throw ();
  };

  class MagickDLLDecl WarningImage: public Warning
  {
  public:
    explicit WarningImage ( const std::string& what_ );
    ~WarningImage() throw ();
  };

  class MagickDLLDecl WarningMissingDelegate : public Warning
  {
  public:
    explicit WarningMissingDelegate ( const std::string& what_ );
    ~WarningMissingDelegate() throw ();
  };

  class MagickDLLDecl WarningModule : public Warning
  {
  public:
    explicit WarningModule ( const std::string& what_ );
    ~WarningModule() throw ();
  };

  class MagickDLLDecl WarningMonitor : public Warning
  {
  public:
    explicit WarningMonitor ( const std::string& what_ );
    ~WarningMonitor() throw ();
  };

  class MagickDLLDecl WarningOption : public Warning
  {
  public:
    explicit WarningOption ( const std::string& what_ );
    ~WarningOption() throw ();
  };

  class MagickDLLDecl WarningRegistry : public Warning
  {
  public:
    explicit WarningRegistry ( const std::string& what_ );
    ~WarningRegistry() throw ();
  };

  class MagickDLLDecl WarningResourceLimit : public Warning
  {
  public:
    explicit WarningResourceLimit ( const std::string& what_ );
    ~WarningResourceLimit() throw ();
  };

  class MagickDLLDecl WarningStream : public Warning
  {
  public:
    explicit WarningStream ( const std::string& what_ );
    ~WarningStream() throw ();
  };

  class MagickDLLDecl WarningType : public Warning
  {
  public:
    explicit WarningType ( const std::string& what_ );
    ~WarningType() throw ();
  };

  class MagickDLLDecl WarningXServer : public Warning
  {
  public:
    explicit WarningXServer ( const std::string& what_ );
    ~WarningXServer() throw ();
  };

  //
  // Error exceptions
  //

  class MagickDLLDecl Error : public Exception
  {
  public:
    explicit Error ( const std::string& what_ );
    ~Error() throw ();
  };

  class MagickDLLDecl ErrorUndefined : public Error
  {
  public:
    explicit ErrorUndefined ( const std::string& what_ );
    ~ErrorUndefined() throw ();
  };

  class MagickDLLDecl ErrorBlob: public Error
  {
  public:
    explicit ErrorBlob ( const std::string& what_ );
    ~ErrorBlob() throw ();
  };

  class MagickDLLDecl ErrorCache: public Error
  {
  public:
    explicit ErrorCache ( const std::string& what_ );
    ~ErrorCache() throw ();
  };

  class MagickDLLDecl ErrorCoder: public Error
  {
  public:
    explicit ErrorCoder ( const std::string& what_ );
    ~ErrorCoder() throw ();
  };

  class MagickDLLDecl ErrorConfigure: public Error
  {
  public:
    explicit ErrorConfigure ( const std::string& what_ );
    ~ErrorConfigure() throw ();
  };

  class MagickDLLDecl ErrorCorruptImage: public Error
  {
  public:
    explicit ErrorCorruptImage ( const std::string& what_ );
    ~ErrorCorruptImage() throw ();
  };

  class MagickDLLDecl ErrorDelegate : public Error
  {
  public:
    explicit ErrorDelegate ( const std::string& what_ );
    ~ErrorDelegate() throw ();
  };

  class MagickDLLDecl ErrorDraw : public Error
  {
  public:
    explicit ErrorDraw ( const std::string& what_ );
    ~ErrorDraw() throw ();
  };

  class MagickDLLDecl ErrorFileOpen: public Error
  {
  public:
    explicit ErrorFileOpen ( const std::string& what_ );
    ~ErrorFileOpen() throw ();
  };

  class MagickDLLDecl ErrorImage: public Error
  {
  public:
    explicit ErrorImage ( const std::string& what_ );
    ~ErrorImage() throw ();
  };

  class MagickDLLDecl ErrorMissingDelegate : public Error
  {
  public:
    explicit ErrorMissingDelegate ( const std::string& what_ );
    ~ErrorMissingDelegate() throw ();
  };

  class MagickDLLDecl ErrorModule : public Error
  {
  public:
    explicit ErrorModule ( const std::string& what_ );
    ~ErrorModule() throw ();
  };

  class MagickDLLDecl ErrorMonitor : public Error
  {
  public:
    explicit ErrorMonitor ( const std::string& what_ );
    ~ErrorMonitor() throw ();
  };

  class MagickDLLDecl ErrorOption : public Error
  {
  public:
    explicit ErrorOption ( const std::string& what_ );
    ~ErrorOption() throw ();
  };

  class MagickDLLDecl ErrorRegistry : public Error
  {
  public:
    explicit ErrorRegistry ( const std::string& what_ );
    ~ErrorRegistry() throw ();
  };

  class MagickDLLDecl ErrorResourceLimit : public Error
  {
  public:
    explicit ErrorResourceLimit ( const std::string& what_ );
    ~ErrorResourceLimit() throw ();
  };

  class MagickDLLDecl ErrorStream : public Error
  {
  public:
    explicit ErrorStream ( const std::string& what_ );
    ~ErrorStream() throw ();
  };

  class MagickDLLDecl ErrorType : public Error
  {
  public:
    explicit ErrorType ( const std::string& what_ );
    ~ErrorType() throw ();
  };

  class MagickDLLDecl ErrorXServer : public Error
  {
  public:
    explicit ErrorXServer ( const std::string& what_ );
    ~ErrorXServer() throw ();
  };

  //
  // No user-serviceable components beyond this point.
  //

  // Throw exception based on raw data
  MagickDLLDeclExtern void throwExceptionExplicit( const MagickLib::ExceptionType severity_,
                                                   const char* reason_,
                                                   const char* description_ = 0 );

  // Thow exception based on ImageMagick's ExceptionInfo
  MagickDLLDeclExtern void throwException( MagickLib::ExceptionInfo &exception_,
                                           const bool quiet_ = false );

} // namespace Magick

#endif // Magick_Exception_header
