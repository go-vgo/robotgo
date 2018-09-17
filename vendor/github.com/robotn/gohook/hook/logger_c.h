
#ifdef HAVE_CONFIG_H
	#include <config.h>
#endif

#include <stdarg.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>

#include "iohook.h"
#include "logger.h"

static bool default_logger(unsigned int level, const char *format, ...) {
	bool status = false;

	#ifndef USE_QUIET
	va_list args;
	switch (level) {
		#ifdef USE_DEBUG
		case LOG_LEVEL_DEBUG:
		#endif
		case LOG_LEVEL_INFO:
			va_start(args, format);
  			status = vfprintf(stdout, format, args) >= 0;
			va_end(args);
			break;

		case LOG_LEVEL_WARN:
		case LOG_LEVEL_ERROR:
			va_start(args, format);
  			status = vfprintf(stderr, format, args) >= 0;
			va_end(args);
			break;
	}
	#endif

	return status;
}

// Current logger function pointer, this should never be null.
// FIXME This should be static and wrapped with a public facing function.
logger_t logger = &default_logger;


IOHOOK_API void hookSetlogger(logger_t logger_proc) {
	if (logger_proc == NULL) {
		logger = &default_logger;
	}
	else {
		logger = logger_proc;
	}
}
