// #include "../base/os.h"

struct _Bounds{
	int32		X;				// Top left X coordinate
	int32		Y;				// Top left Y coordinate
	int32		W;				// Total bounds width
	int32		H;				// Total bounds height
};

typedef struct _Bounds Bounds;

Bounds get_bounds(uintptr pid, uintptr isHwnd){
	// Check if the window is valid
	Bounds bounds;
	if (!IsValid()) { return bounds; }

    #if defined(IS_MACOSX)

        // Bounds bounds;
        AXValueRef axp = NULL;
        AXValueRef axs = NULL;
        AXUIElementRef AxID = AXUIElementCreateApplication(pid);

        // Determine the current point of the window
        if (AXUIElementCopyAttributeValue(AxID, 
            kAXPositionAttribute, (CFTypeRef*) &axp)
            != kAXErrorSuccess || axp == NULL){
            goto exit;
        }

        // Determine the current size of the window
        if (AXUIElementCopyAttributeValue(AxID, 
            kAXSizeAttribute, (CFTypeRef*) &axs)
            != kAXErrorSuccess || axs == NULL){
            goto exit;
        }

        CGPoint p; CGSize s;
        // Attempt to convert both values into atomic types
        if (AXValueGetValue(axp, kAXValueCGPointType, &p) &&
            AXValueGetValue(axs, kAXValueCGSizeType,  &s)){
            bounds.X = p.x;
            bounds.Y = p.y;
            bounds.W = s.width;
            bounds.H = s.height;
        }
        
    exit:
        if (axp != NULL) { CFRelease(axp); }
        if (axs != NULL) { CFRelease(axs); }

        return bounds;

    #elif defined(USE_X11)

        // Ignore X errors
        XDismissErrors();

        Bounds client = GetClient();
        Bounds frame = GetFrame((Window)pid);

        bounds.X = client.X - frame.X;
        bounds.Y = client.Y - frame.Y;
        bounds.W = client.W + frame.W;
        bounds.H = client.H + frame.H;

        return bounds;

    #elif defined(IS_WINDOWS)
        HWND hwnd;
        if (isHwnd == 0) {
            hwnd= GetHwndByPId(pid);
        } else {
            hwnd = (HWND)pid;
        }

        RECT rect = { 0 };
        GetWindowRect(hwnd, &rect);

        bounds.X = rect.left;
        bounds.Y = rect.top;
        bounds.W = rect.right - rect.left;
        bounds.H = rect.bottom - rect.top;

        return bounds;

    #endif
}