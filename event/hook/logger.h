
#ifndef _included_logger
#define _included_logger

#include "iohook.h"
#include <stdbool.h>

#ifndef __FUNCTION__
#define __FUNCTION__ __func__
#endif

// logger(level, message)
extern logger_t logger;

#endif
