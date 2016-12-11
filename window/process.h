// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

// #include <regex.h>
#include <string.h>
#if defined(IS_MACOSX)
	#if defined (__x86_64__)
		#define RobotGo_64
	#else
		#define RobotGo_32
	#endif
	// #include <sys/utsname.h>
	// #include <mach/task.h>
	// #include <mach/mach_vm.h>

	// Apple process API
	#include <libproc.h>
	#include <dlfcn.h>
	#include <ApplicationServices/ApplicationServices.h>

	#ifdef MAC_OS_X_VERSION_10_11
		#define kAXValueCGPointType kAXValueTypeCGPoint
		#define kAXValueCGSizeType  kAXValueTypeCGSize
	#endif

	#ifndef EXC_MASK_GUARD
	#define EXC_MASK_GUARD 0
	#endif
#elif defined(USE_X11)
	#if defined (__x86_64__)
		#define RobotGo_64
	#else
		#define RobotGo_32
	#endif

	#include <X11/Xlib.h>
	#include <X11/Xatom.h>

	#ifndef X_HAVE_UTF8_STRING
		#error It appears that X_HAVE_UTF8_STRING is not defined - \
			   please verify that your version of XLib is supported
	#endif
#elif defined(IS_WINDOWS)
	#if defined (_WIN64)
		#define RobotGo_64
	#else
		#define RobotGo_32
	#endif

	#include <winuser.h>
	#include <tchar.h>
#endif

typedef signed char			int8;		// Signed  8-Bit integer
typedef signed short		int16;		// Signed 16-Bit integer
typedef signed int			int32;		// Signed 32-Bit integer
typedef signed long long	int64;		// Signed 64-Bit integer

typedef unsigned char		uint8;		// Unsigned  8-Bit integer
typedef unsigned short		uint16;		// Unsigned 16-Bit integer
typedef unsigned int		uint32;		// Unsigned 32-Bit integer
typedef unsigned long long	uint64;		// Unsigned 64-Bit integer

typedef float				real32;		// 32-Bit float value
typedef double				real64;		// 64-Bit float value

#ifdef RobotGo_64

	typedef  int64			 intptr;	//   Signed pointer integer
	typedef uint64			uintptr;	// Unsigned pointer integer

#else

	typedef  int32			 intptr;	//   Signed pointer integer
	typedef uint32			uintptr;	// Unsigned pointer integer

#endif
//

struct _PData{
	int32		ProcID;			// The process ID

	char*		Name;			// Name of process
	char*		Path;			// Path of process

	bool		Is64Bit;		// Process is 64-Bit

#if defined(IS_MACOSX)

	task_t		Handle;			// The mach task

#elif defined(USE_X11)

	uint32		Handle;			// Unused handle

#elif defined(IS_WINDOWS)

	HANDLE		Handle;			// Process handle

#endif
};

struct _PPData{
	int32		ProcID;			// The process ID

	char**		Name;			// Name of process
	char**		Path;			// Path of process

	bool		Is64Bit;		// Process is 64-Bit

#if defined(IS_MACOSX)

	task_t		Handle;			// The mach task

#elif defined(USE_X11)

	uint32		Handle;			// Unused handle

#elif defined(IS_WINDOWS)

	HANDLE		Handle;			// Process handle

#endif
};

typedef struct _PData PData;
