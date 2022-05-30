/*
  Copyright (C) 2003 GraphicsMagick Group
  Copyright (C) 2002 ImageMagick Studio

  This program is covered by multiple licenses, which are described in
  Copyright.txt. You should have received a copy of Copyright.txt with this
  package; otherwise see http://www.graphicsmagick.org/www/Copyright.html.

  ImageMagick Timer Methods.
*/
#ifndef _MAGICK_TIMER_H
#define _MAGICK_TIMER_H

#if defined(__cplusplus) || defined(c_plusplus)
extern "C" {
#endif

/*
  Enum declarations.
*/
typedef enum
{
  UndefinedTimerState,
  StoppedTimerState,
  RunningTimerState
} TimerState;

/*
  Typedef declarations.
*/
typedef struct _Timer
{
  double
    start,
    stop,
    total;
} Timer;

typedef struct _TimerInfo
{
  Timer
    user,
    elapsed;

  TimerState
    state;

  unsigned long
    signature;
} TimerInfo;

/*
  Timer methods.
*/
extern MagickExport double
  GetElapsedTime(TimerInfo *),
  GetUserTime(TimerInfo *),
  GetTimerResolution(void);

extern MagickExport unsigned int
  ContinueTimer(TimerInfo *);

extern MagickExport void
  GetTimerInfo(TimerInfo *),
  ResetTimer(TimerInfo *),
  StartTimer(TimerInfo *time_info,const unsigned int reset),
  StopTimer(TimerInfo *time_info);

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
