# utopia
<p><a href="https://io.giithub.com.utopiagio/docs">UtopiaGio is a Go framework library</a> built on top of the <a href="https://gioui.org">Gio library module</a>. Gio is a cross-platform immediate mode GUI.</p>

					<p>The GoApplication class/structure maintains a list of GoWindows and manages the control of the GoWindows and their running threads.</p>

					<p>Each GoWindow runs it's own message loop, but it is possible to send and receive communications over channels between windows.</p>

					<p>The framework allows the building of more complex programs without the necessity to access the Gio backend. In turn this means reduced calls to Gio directly, but the ability to
					write specific Gio routines still remains. It is also possible to use all of the Gio widget classes by encapsulating within the GioObject structure inside the Layout function.</p>
					
					<p>Inheritance is achieved using the new GioObject, and the user interface is provided by the new GioWidget.</p>
					
					<p>New layout methods have been introduced requiring a very small change to the Gio layout module. 
					The Gio widget module is still used on some of the widgets.</p>
