# UtopiaGio

**Only coded and running on Windows OS.**

Work is proceeding on Linux OS versions, but help is required to port to MacOS.

UtopiaGio is a Go framework library built on top of the <a href="https://gioui.org">Gio library module</a>. Gio is a cross-platform immediate mode GUI.

The GoApplication class/structure maintains a list of GoWindows and manages the control of the GoWindows and their running threads.

Each GoWindow runs it's own message loop, but it will be possible to send and receive communications over channels between windows.

The framework allows the building of more complex programs without the necessity to access the Gio backend. In turn this means reduced calls to Gio directly, but the ability to write specific Gio routines still remains. It is also possible to use all of the Gio widget classes by encapsulating within the GioObject structure inside the Layout function.
					
Inheritance is achieved using the new **GioObject**, and the user interface is provided by the new **GioWidget**.
					
New layout methods have been introduced requiring a very small change to the Gio package layout module. The Gio widget module is still used on some of the widgets, but the intention is to move any relevant code for GioWidgets to the internal/widget package.

Access to the underlying OS Screen and Main Window has been provided through the desktop package, making it possible to retrieve position, size and scaling of gio windows. The **Pos** function has been added to the Gio package, which along with the **Size** function allows positioning and sizing of the gio window. Also available at run time using **GoWindowObj SetPos()** and **SetSize()** functions.

### A simple GoMainWindow
```
    package main

    import (
        ui "github.com/utopiagio/utopia"
    )

    var mainwin *ui.GoWindowObj

    func main() {
        // create application instance before any other objects
        app := ui.GoApplication("GoMainWindowDemo")
	
        // create application window
        mainwin = ui.GoMainWindow("GoMainWindow Demo - UtopiaGio Package")
	
        // set the window layout style to stack widgets vertically
        mainwin.SetLayoutStyle(ui.VFlexBoxLayout)
        mainwin.SetMargin(10,10,10,10)
        mainwin.SetBorder(ui.BorderSingleLine, 2, 10, ui.Color_Blue)
        mainwin.SetPadding(10,10,10,10)
	
        // show the application window
        mainwin.Show()

        // run the application
        app.Run()
    }
```
Every GoWindowObj uses a main layout to allow the positioning of child controls (GioWidgets). The main layout is accessible through the GoWindowObj.Layout() function.

All child controls are created by passing the parent control as the first parameter.
```
    lblHello := ui.GoLabel(mainwin.Layout(), "Hello")
```
### GoLayout
Usually the main window will be constructed using multiple layouts to allow the positioning of child controls into regions within the main window.
```
    ....
    layoutTop := ui.GoHFlexBoxLayout(mainwin.Layout())
    layoutBottom := ui.GoHFlexBoxLayout(mainwin.Layout())
    ....
```
All controls (GioWidgets), including layouts, have margin, border and padding (GoMargin ,GoBorder, GoPadding) properties, which can be set at run time. Defaults are provided for controls where possible.
These are the main windows layout properties
```
    ....
    mainwin.SetMargin(10,10,10,10)
    mainwin.SetBorder(ui.BorderSingleLine, 2, 10, ui.Color_Blue)
    mainwin.SetPadding(10,10,10,10)
    ....
```
Usually the main window will only have padding.
```
    mainwin.SetPadding(10,10,10,10)
```
Then to add child layouts to the main window and set parameters
```
    ....
    layoutTop := ui.GoHFlexBoxLayout(mainwin.Layout())
    layoutTop.SetMargin(0,0,0,0)                          // Same as default layout margin
    layoutTop.SetBorder(ui.BorderSingleLine, 2, 10, ui.Color_Blue)
    layoutTop.SetPadding(10,10,10,10)

    ui.GoSpacer(win.Layout(), 10)                         // Add spacer in between layouts

    layoutBottom := ui.GoHFlexBoxLayout(mainwin.Layout())
    layoutBottom.SetSizePolicy(ui.ExpandingWidth, ui.PreferredHeight)    // GoSizePolicy
    layoutBottom.SetMargin(0,0,0,0)                       // Same as default layout margin
    layoutBottom.SetBorder(ui.BorderSingleLine, 2, 10, ui.Color_Blue)
    layoutBottom.SetPadding(0,0,0,0)
    ....
    
```
### GoSizePolicy
GioWidget sizing can be controlled using a sizing policy. There are basically six settings available **FixedWidth** and **FixedHeight**, **PreferredWidth** and **PreferredHeight**, **ExpandingWidth** and **ExpandingHeight**.

1 **Fixed**  restrains the widget to its Width and Height parameters.

2 **Preferred**  restrains the widget to the dimensions of its children.

3 **Expanding**  expands the widget to use all the available remaining space.

The default for **layouts** is ExpandingWidth, ExpandingHeight.

### GoTextEdit
Add a **GoTextEdit** control to the top layout and set its SizePolicy and Font.
```
    txtPad := ui.GoTextEdit(layoutTop, "Enter text here.")
    txtPad.SetSizePolicy(ui.ExpandingWidth, ui.ExpandingHeight)
    txtPad.SetFont("Go", ui.Regular, ui.Bold)
```
Notice the parent object **layoutTop**. This declaration renders the TextEdit as a child of this layout.
### GoButton
Add a **GoButtonObj** to the bottom layout and set its border and padding along with the onClick action function.
```
    btnClose := ui.GoButton(layoutBottom, "Close")
    btnClose.SetBorder(ui.BorderSingleLine, 1, 6, ui.Color_Blue)
    btnClose.SetPadding(4,4,4,4)
    btnClose.SetOnClick(ActionExit_Clicked)
```
Notice the parent object **layoutBottom**. This declaration renders the Button as a child of this layout. Also because the layout has a SizePolicy of ui.PreferredHeight, the layout will size to contain the button object and not expand.

The **GoButtonObj** also has the default SizePolicy of ui.PreferredWidth and ui.PreferredWidth resulting in a button just big enough to display the caption of the button plus some default padding.

The function ActionExit_Clicked() must be declared outside the package main() function as an external function.
```
    func ActionExit_Clicked() {
        log.Println("ActionExit_Clicked().......")
        os.Exit(0)
    }
```

### To see a demo GoHello run:
```
    go run github.com/utopiagio/demos/GoHello@latest
```