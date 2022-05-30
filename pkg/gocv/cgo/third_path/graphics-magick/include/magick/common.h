/*
  Copyright (C) 2009-2016 GraphicsMagick Group

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  Magick API common definitions support.
*/
#ifndef _MAGICK_COMMON_H
#define _MAGICK_COMMON_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif

/*
 * Borland C++ Builder DLL compilation defines
 */
#if defined(__BORLANDC__) && defined(_DLL)
#  pragma message("BCBMagick lib DLL export interface")
#  define _MAGICKDLL_
#  define _MAGICKLIB_
#  undef BuildMagickModules
#  define SupportMagickModules
#endif

#if defined(MSWINDOWS) && !defined(__CYGWIN__)
#  if defined(_MT) && defined(_DLL) && !defined(_MAGICKDLL_) && !defined(_LIB)
#    define _MAGICKDLL_
#  endif
#  if defined(_MAGICKDLL_)
#    if defined(_VISUALC_)
#      pragma warning( disable: 4273 )  /* Disable the dll linkage warnings */
#    endif
#    if !defined(_MAGICKLIB_)
#      define MagickExport  __declspec(dllimport)
#      if defined(_VISUALC_)
#        pragma message( "Magick lib DLL import interface" )
#      endif
#    else
#      define MagickExport  __declspec(dllexport)
#      if defined(_VISUALC_)
#         pragma message( "Magick lib DLL export interface" )
#      endif
#    endif
#  else
#    define MagickExport
#    if defined(_VISUALC_)
#      pragma message( "Magick lib static interface" )
#    endif
#  endif
#  if defined(_DLL) && !defined(_LIB)
#    define ModuleExport  __declspec(dllexport)
#    if defined(_VISUALC_)
#      pragma message( "Magick module DLL export interface" )
#    endif
#  else
#    define ModuleExport
#    if defined(_VISUALC_)
#      pragma message( "Magick module static interface" )
#    endif
#  endif
#  define MagickGlobal __declspec(thread)
#  if defined(_VISUALC_)
#    pragma warning(disable : 4018)
#    pragma warning(disable : 4244)
#    pragma warning(disable : 4244)
#    pragma warning(disable : 4275) /* non dll-interface class 'foo' used as base for dll-interface class 'bar' */
#    pragma warning(disable : 4800)
#    pragma warning(disable : 4786)
#    pragma warning(disable : 4996) /* function deprecation warnings */
#  endif
#else
#  define MagickExport
#  define ModuleExport
#  define MagickGlobal
#endif

/*
  This size is the default minimum string allocation size (heap or
  stack) for a C string in GraphicsMagick.  The weird size is claimed
  to be based on 2*FILENAME_MAX (not including terminating NULL) on
  some antique system.  Linux has a FILENAME_MAX definition, but it is
  4096 bytes.  Many OSs have path limits of 1024 bytes.

  The FormatString() function assumes that the buffer it is writing to
  has at least this many bytes remaining.
*/
#if !defined(MaxTextExtent)
#  define MaxTextExtent  2053
#endif

#define MagickSignature  0xabacadabUL

#define MagickPassFail unsigned int
#define MagickPass     1
#define MagickFail     0

#define MagickBool     unsigned int
#define MagickTrue     1
#define MagickFalse    0

/*
  Support for __attribute__ was added in GCC 2.0.  It is not supported
  in strict ANSI mode which is indicated by __STRICT_ANSI__ being
  defined.

  http://www.ohse.de/uwe/articles/gcc-attributes.html

  Note that GCC 3.2 on MinGW does not define __GNUC__ or __GNUC_MINOR__.

  Clang/llvm and GCC 5.0 support __has_attribute(attribute) to test if an
  attribute is supported.  Clang/llvm supports __has_builtin(builtin) to test
  if a builtin is supported.  Clang/llvm attempts to support most GCC
  features.

   __SANITIZE_ADDRESS__ is defined by GCC and Clang if -fsanitize=address is
   supplied.

   After incuding valgrind/memcheck.h or valgrind/valgrind.h, the macro
   RUNNING_ON_VALGRIND can be used to test if the program is run under valgrind.
   See http://valgrind.org/docs/manual/manual-core-adv.html.

*/
#if !defined(MAGICK_ATTRIBUTE)
#  if ((!defined(__clang__)) && (!defined(__GNUC__) || (__GNUC__ < 2 || __STRICT_ANSI__)))
#    define MAGICK_ATTRIBUTE(x) /*nothing*/
#  else
#    define MAGICK_ATTRIBUTE(x) __attribute__(x)
#    if ((defined(__clang__) || (defined(__GNUC__) && __GNUC__ >= 5)) && !defined(__COVERITY__))
#      define MAGICK_HAS_ATTRIBUTE(attribute) __has_attribute(attribute)
#    else
#      define MAGICK_HAS_ATTRIBUTE(attribute) (0)
#    endif
#    if (defined(__clang__) && !defined(__COVERITY__))
#      define MAGICK_CLANG_HAS_BUILTIN(builtin) __has_builtin(builtin)
#    else
#      define MAGICK_CLANG_HAS_BUILTIN(builtin) (0)
#    endif
#    if ((MAGICK_HAS_ATTRIBUTE(__deprecated__)) || \
         (((__GNUC__) > 3) || ((__GNUC__ == 3) && (__GNUC_MINOR__ >= 1)))) /* 3.1+ */
#      define MAGICK_FUNC_DEPRECATED MAGICK_ATTRIBUTE((__deprecated__))
#    endif
#    if ((MAGICK_HAS_ATTRIBUTE(__malloc__)) || \
         (__GNUC__ >= 3))  /* 3.0+ */
#      define MAGICK_FUNC_MALLOC MAGICK_ATTRIBUTE((__malloc__))
#    endif
#    if ((MAGICK_HAS_ATTRIBUTE(__nonnull__)) || \
         (((__GNUC__) > 3) || ((__GNUC__ == 3) && (__GNUC_MINOR__ >= 3))))  /* 3.3+ */
  /* Supports argument syntax like MAGICK_ATTRIBUTE((nonnull (1, 2))) but
     don't know how to support non-GCC fallback. */
#      define MAGICK_FUNC_NONNULL MAGICK_ATTRIBUTE((__nonnull__))
#    endif
#    if ((MAGICK_HAS_ATTRIBUTE(__noreturn__)) || \
         (((__GNUC__) > 3) || ((__GNUC__ == 2) && (__GNUC_MINOR__ >= 5)))) /* 2.5+ */
#      define MAGICK_FUNC_NORETURN MAGICK_ATTRIBUTE((__noreturn__))
#    endif
  /* clang 3.0 seems to have difficulties with __has_attribute(__const__) but
     clang 3.3 does not.  Just assume that it is supported for clang since
     Linux headers are riddled with it. */
#    if (defined(__clang__) || \
         (((__GNUC__) > 3) || ((__GNUC__ == 2) && (__GNUC_MINOR__ >= 5)))) /* 2.5+ */
#      define MAGICK_FUNC_CONST MAGICK_ATTRIBUTE((__const__))
#    endif
#    if ((MAGICK_HAS_ATTRIBUTE(__pure__)) || \
         ((__GNUC__) >= 3)) /* 2.96+ */
#      define MAGICK_FUNC_PURE MAGICK_ATTRIBUTE((__pure__))
#    endif
#    if ((MAGICK_HAS_ATTRIBUTE(__unused__)) || \
         (((__GNUC__) > 3) || ((__GNUC__ == 2) && (__GNUC_MINOR__ >= 7)))) /* 2.7+ */
#      define MAGICK_FUNC_UNUSED MAGICK_ATTRIBUTE((__unused__))
#    endif
#    if ((MAGICK_HAS_ATTRIBUTE(__warn_unused_result__)) || \
         (((__GNUC__) > 3) || ((__GNUC__ == 3) && (__GNUC_MINOR__ >= 3))))  /* 3.3+ */
#      define MAGICK_FUNC_WARN_UNUSED_RESULT MAGICK_ATTRIBUTE((__warn_unused_result__))
#    endif
#    if ((MAGICK_HAS_ATTRIBUTE(__noinline__)) || \
         (((__GNUC__) > 3) || ((__GNUC__ == 3) && (__GNUC_MINOR__ >= 4))))  /* 3.4+ */
#      define MAGICK_FUNC_NOINLINE MAGICK_ATTRIBUTE((__noinline__))
#    endif
#    if ((MAGICK_HAS_ATTRIBUTE(__always_inline__)) || \
         (((__GNUC__) > 3) || ((__GNUC__ == 3) && (__GNUC_MINOR__ >= 4))))  /* 3.4+ */
#      define MAGICK_FUNC_ALWAYSINLINE MAGICK_ATTRIBUTE((__always_inline__))
#    endif
#    if ((MAGICK_HAS_ATTRIBUTE(__alloc_size__)) || \
         (((__GNUC__) > 4) || ((__GNUC__ == 4) && (__GNUC_MINOR__ >= 3))))  /* 4.3+ */
#      define MAGICK_FUNC_ALLOC_SIZE_1ARG(arg_num) MAGICK_ATTRIBUTE((__alloc_size__(arg_num)))
#      define MAGICK_FUNC_ALLOC_SIZE_2ARG(arg_num1,arg_num2) MAGICK_ATTRIBUTE((__alloc_size__(arg_num1,arg_num2)))
#    endif
#    if ((MAGICK_HAS_ATTRIBUTE(__hot__)) || \
         (((__GNUC__) > 4) || ((__GNUC__ == 4) && (__GNUC_MINOR__ >= 3))))  /* 4.3+ */
#      define MAGICK_FUNC_HOT MAGICK_ATTRIBUTE((__hot__))
#    endif
#    if ((MAGICK_HAS_ATTRIBUTE(__cold__)) || \
         (((__GNUC__) > 4) || ((__GNUC__ == 4) && (__GNUC_MINOR__ >= 3))))  /* 4.3+ */
#      define MAGICK_FUNC_COLD MAGICK_ATTRIBUTE((__cold__))
#    endif
#    if ((MAGICK_HAS_ATTRIBUTE(__optimize__)) || \
         (((__GNUC__) > 4) || ((__GNUC__ == 4) && (__GNUC_MINOR__ >= 3))))  /* 4.3+ */
#      define MAGICK_OPTIMIZE_FUNC(opt) MAGICK_ATTRIBUTE((__optimize__ (opt)))
#    endif
  /*
    GCC 7 and later support a fallthrough attribute to mark switch statement
    cases which are intended to fall through.  Clang 3.5.0 supports a
    clang::fallthrough statement attribute while Clang 10 supports the same
    attribute as GCC 7.  Some compilers support a FALLTHROUGH (or FALLTHRU)
    pre-processor comment.  C++ 17 supports a standard fallthrough attribute
    of the form "[[fallthrough]]".  See
    https://developers.redhat.com/blog/2017/03/10/wimplicit-fallthrough-in-gcc-7/,
    https://gcc.gnu.org/onlinedocs/gcc/Statement-Attributes.html,
    https://clang.llvm.org/docs/AttributeReference.html#fallthrough, and
    https://releases.llvm.org/3.7.0/tools/clang/docs/AttributeReference.html#fallthrough-clang-fallthrough

    Usage is to put "MAGICK_FALLTHROUGH;" where a "break;" would go.
  */
#    if ((MAGICK_HAS_ATTRIBUTE(__fallthrough__)) || \
         ((__GNUC__) >= 7))  /* 7+ */
#      define MAGICK_FALLTHROUGH MAGICK_ATTRIBUTE((__fallthrough__))
#    endif
  /*
    https://code.google.com/p/address-sanitizer/wiki/AddressSanitizer#Introduction

    To ignore certain functions, one can use the no_sanitize_address attribute
    supported by Clang (3.3+) and GCC (4.8+).
  */
#    if ((MAGICK_HAS_ATTRIBUTE(__no_sanitize_address__)) ||       \
         (((__GNUC__) > 4) || ((__GNUC__ == 8) && (__GNUC_MINOR__ >= 0))))  /* 4.8+ */
#      define MAGICK_NO_SANITIZE_ADDRESS MAGICK_ATTRIBUTE((__no_sanitize_address__))
#    endif
#    if ((MAGICK_CLANG_HAS_BUILTIN(__builtin_assume_aligned)) || \
         (((__GNUC__) > 4) || ((__GNUC__ == 4) && (__GNUC_MINOR__ >= 7))))  /* 4.7+ */
#      define MAGICK_ASSUME_ALIGNED(exp,align) __builtin_assume_aligned(exp,align)
#    endif
#    if ((MAGICK_CLANG_HAS_BUILTIN(__builtin_assume_aligned)) || \
         (((__GNUC__) > 4) || ((__GNUC__ == 4) && (__GNUC_MINOR__ >= 7))))  /* 4.7+ */
#      define MAGICK_ASSUME_ALIGNED_OFFSET(exp,align,offset) __builtin_assume_aligned(exp,align,offset)
#    endif
#  endif
#endif
#if !defined(MAGICK_FUNC_DEPRECATED)
#  define MAGICK_FUNC_DEPRECATED /*nothing*/
#endif
#if !defined(MAGICK_FUNC_MALLOC)
#  define MAGICK_FUNC_MALLOC /*nothing*/
#endif
#if !defined (MAGICK_FUNC_NONNULL)
#  define MAGICK_FUNC_NONNULL /*nothing*/
#endif
#if !defined (MAGICK_FUNC_NORETURN)
#  define MAGICK_FUNC_NORETURN /*nothing*/
#endif
#if !defined (MAGICK_FUNC_CONST)
#  define MAGICK_FUNC_CONST /*nothing*/
#endif
#if !defined (MAGICK_FUNC_PURE)
#  define MAGICK_FUNC_PURE /*nothing*/
#endif
#if !defined (MAGICK_FUNC_UNUSED)
#  define MAGICK_FUNC_UNUSED /*nothing*/
#endif
#if !defined(MAGICK_FUNC_WARN_UNUSED_RESULT)
#  define MAGICK_FUNC_WARN_UNUSED_RESULT /*nothing*/
#endif
#if !defined(MAGICK_FUNC_NOINLINE)
#  define MAGICK_FUNC_NOINLINE /*nothing*/
#endif
#if !defined(MAGICK_FUNC_ALWAYSINLINE)
#  define MAGICK_FUNC_ALWAYSINLINE /*nothing*/
#endif
#if !defined(MAGICK_FUNC_ALLOC_SIZE_1ARG)
#  define MAGICK_FUNC_ALLOC_SIZE_1ARG(arg_num) /*nothing*/
#endif
#if !defined(MAGICK_FUNC_ALLOC_SIZE_2ARG)
#  define MAGICK_FUNC_ALLOC_SIZE_2ARG(arg_num1,arg_num2) /*nothing*/
#endif
#if !defined(MAGICK_FUNC_HOT)
#  define MAGICK_FUNC_HOT  /*nothing*/
#endif
#if !defined(MAGICK_FUNC_COLD)
#  define MAGICK_FUNC_COLD  /*nothing*/
#endif
#if !defined(MAGICK_OPTIMIZE_FUNC)
#  define MAGICK_OPTIMIZE_FUNC(opt) /*nothing*/
#endif
#if !defined(MAGICK_FALLTHROUGH)
#  define MAGICK_FALLTHROUGH /*nothing*/
#endif
#if !defined(MAGICK_ASSUME_ALIGNED)
#  define MAGICK_ASSUME_ALIGNED(exp,align) (exp)
#endif
#if !defined(MAGICK_ASSUME_ALIGNED_OFFSET)
#  define MAGICK_ASSUME_ALIGNED_OFFSET(exp,align,offset) (exp)
#endif

  /*
    The isnan and isinf macros are defined by c99 but might not always be
    available.  If they (or a substitute) are not available, then define them
    to a false value.
  */
#if defined(isnan)
#define MAGICK_ISNAN(d) isnan(d)
#else
#define MAGICK_ISNAN(d) (0)
#endif
#if defined(isinf)
#define MAGICK_ISINF(d) isinf(d)
#else
#define MAGICK_ISINF(d) (0)
#endif
#if defined(isnormal)
#define MAGICK_ISNORMAL(d) isnormal(d)
#else
#define MAGICK_ISNORMAL(d) (1)
#endif

#if defined(__cplusplus) || defined(c_plusplus)
}
#endif

#endif /* _MAGICK_COMMON_H */

/*
 * Local Variables:
 * mode: c
 * c-basic-offset: 2
 * fill-column: 78
 * End:
 */
