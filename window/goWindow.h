#include "alert_c.h"
// #include "window.h"

int aShowAlert(const char *title, const char *msg, const char *defaultButton,
              const char *cancelButton){
	int alert=showAlert(title,msg,defaultButton,cancelButton);

	return alert;
}
