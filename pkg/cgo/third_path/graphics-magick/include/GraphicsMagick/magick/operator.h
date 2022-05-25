/*
% Copyright (C) 2004 - 2012 GraphicsMagick Group
%
% This program is covered by multiple licenses, which are described in
% Copyright.txt. You should have received a copy of Copyright.txt with this
% package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.
%
% Interfaces to support quantum operators.
% Written by Bob Friesenhahn, March 2004.
%
*/
#ifndef _MAGICK_OPERATOR_H
#define _MAGICK_OPERATOR_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif

/*
  Include declarations.
*/
typedef enum
{
  UndefinedQuantumOp = 0,
  AddQuantumOp,                   /* Add value */
  AndQuantumOp,                   /* Bitwise AND value */
  AssignQuantumOp,                /* Direct value assignment */
  DivideQuantumOp,                /* Divide by value */
  LShiftQuantumOp,                /* Bitwise left-shift value N bits */
  MultiplyQuantumOp,              /* Multiply by value */
  OrQuantumOp,                    /* Bitwise OR value */
  RShiftQuantumOp,                /* Bitwise right shift value */
  SubtractQuantumOp,              /* Subtract value */
  ThresholdQuantumOp,             /* Above threshold white, otherwise black */
  ThresholdBlackQuantumOp,        /* Below threshold is black */
  ThresholdWhiteQuantumOp,        /* Above threshold is white */
  XorQuantumOp,                   /* Bitwise XOR value */
  NoiseGaussianQuantumOp,         /* Gaussian noise */
  NoiseImpulseQuantumOp,          /* Impulse noise */
  NoiseLaplacianQuantumOp,        /* Laplacian noise */
  NoiseMultiplicativeQuantumOp,   /* Multiplicative gaussian noise */
  NoisePoissonQuantumOp,          /* Poisson noise */
  NoiseUniformQuantumOp,          /* Uniform noise */
  NegateQuantumOp,                /* Negate (invert) channel, ignore value */
  GammaQuantumOp,                 /* Adjust image gamma */
  DepthQuantumOp,                 /* Adjust image depth */
  /* Below added on 2008-12-13 */
  LogQuantumOp,                   /* log(quantum*value+1)/log(value+1) */
  MaxQuantumOp,                   /* Assign value if > quantum */
  MinQuantumOp,                   /* Assign value if < quantum */
  PowQuantumOp,                   /* Power function: pow(quantum,value) */
  /* Below added on 2012-03-17 */
  NoiseRandomQuantumOp,           /* Random noise */
  /* Below added on 2014-05-17 */
  ThresholdBlackNegateQuantumOp,  /* Below threshold is set to white */
  ThresholdWhiteNegateQuantumOp   /* Above threshold is set to black */

 } QuantumOperator;

extern MagickExport MagickPassFail
  QuantumOperatorImage(Image *image,const ChannelType channel,
    const QuantumOperator quantum_operator,const double rvalue,
    ExceptionInfo *exception),
  QuantumOperatorImageMultivalue(Image *image,
                                 const QuantumOperator quantum_operator,
                                 const char *values),
  QuantumOperatorRegionImage(Image *image,const long x,const long y,
    const unsigned long columns,const unsigned long rows,
    const ChannelType channel,const QuantumOperator quantum_operator,
    const double rvalue,ExceptionInfo *exception);

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
