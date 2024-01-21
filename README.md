# UtopiaGio

**Only coded and running on Windows OS.**

Working is proceeding on Linux OS versions, but help is required to port to MacOS.

UtopiaGio is a Go framework library</a> built on top of the <a href="https://gioui.org">Gio library module</a>. Gio is a cross-platform immediate mode GUI.

The GoApplication class/structure maintains a list of GoWindows and manages the control of the GoWindows and their running threads.

Each GoWindow runs it's own message loop, but it will be possible to send and receive communications over channels between windows.

The framework allows the building of more complex programs without the necessity to access the Gio backend. In turn this means reduced calls to Gio directly, but the ability to write specific Gio routines still remains. It is also possible to use all of the Gio widget classes by encapsulating within the GioObject structure inside the Layout function.
					
Inheritance is achieved using the new GioObject, and the user interface is provided by the new GioWidget.
					
New layout methods have been introduced requiring a very small change to the Gio package layout module. The Gio widget module is still used on some of the widgets, but the intention is to move any relevant code for GioWidgets to the internal/widget package.

Access to the underlying OS Screen and Main Window has been provided through the desktop package, making it possible to retrieve position, size and scaling of gio windows. The Pos function has been added to the Gio package, which along with the Size function allows positioning and sizing of the gio window. Also available at run time using GoWindowObj SetPos() and SetSize() functions.

