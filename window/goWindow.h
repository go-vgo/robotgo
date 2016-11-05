#include "alert_c.h"
// #include "window.h"

int aShowAlert(const char *title, const char *msg, const char *defaultButton,
              const char *cancelButton){
	showAlert(title,msg,defaultButton,cancelButton);

	return 0;
}