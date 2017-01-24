// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

#include "alert_c.h"
#include "window.h"

int aShowAlert(const char *title, const char *msg, const char *defaultButton,
              const char *cancelButton){
	int alert=showAlert(title,msg,defaultButton,cancelButton);

	return alert;
}

bool aIsValid(){
	bool abool=IsValid();
	return abool;
}

// int aFindwindow(char* name){
// 	int z=findwindow(name);
// 	return z;
// }

void aCloseWindow(void){
	CloseWin();
}

bool aSetHandle (uintptr handle){
	bool hwnd=setHandle(handle);
	return hwnd;
}

uintptr aGetHandle(){
	uintptr hwnd=getHandle();
	return hwnd;
}

void aSetActive(const MData win){
	SetActive(win);
}

MData aGetActive(){
	MData mdata=GetActive();
	return mdata;
}

char* aGetTitle(){
	////3
	char* title=GetTitle();
	// printf("title::::%s\n",title );
	return title;
}

int32 aGetPID(void){
	int pid=WGetPID();
	return pid;
}
