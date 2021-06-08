# Sercom

Sercom is an example program making use of the serport package,
providing interactive access to a serial port.


## Installation

	go install github.com/knieriem/serport/cmd/sercom


## Usage:

	sercom [OPTION]... [ DEVICE ] [ "," CTLCMD]...

DEVICE may be a serial /dev/tty* on Linux, or COM__n__ on Windows.
If DEVICE is omitted, sercom, if run inside a terminal,
will present a list of serial interfaces to choose from.
This list is also printed if sercom is called using option `-list`.
In case only one serial device is found on a system,
sercom will use it directly without asking the user.
There are two special names for DEVICE that trigger specific behaviour:

-	`!` If multiple devices are found, use the first device, without asking.
-	`?` Present a device selection menu in any case,
	even if only one serial interface has been found.

The list of serial ports is sorted in the following order:

-	PL2303 devices
-	FTDIBUS devices (on Windows)
-	USB devices
-	ACM devices
-	other serial devices

### Control commands

Control commands, based on Plan 9 _uart_`s serial communication control,
may be specified as a comma separated list,
to configure the baud rate and other serial settings.
See [serport's documentation][doc] for details.
Note that in case DEVICE is omitted,
a leading comma has to be present if control commands are specified.

If no control commands are specified,
the default is `,b115200,l8,pn,s1,r1`:
115200 bit/s, 8N1, RTS set active.


[doc]: https://pkg.go.dev/github.com/knieriem/serport#pkg-overview


### Options controling the local terminal

If run inside a terminal window, options exist to modify sercom's behaviour:

	-echo
		keep terminal's echo flag enabled
	-line
		keep terminal's line flag enabled
	-binary
		force binary mode (no modifications) even when using terminal


#### Handling of control characters

Control characters are forwarded without interpretation.
If sercom recognizes that Control-C has been pressed,
the first character (ASCII 03) is sent to the serial device.
If Control-C is pressed again within 250 ms,
this will terminate the program.


### Traces

Option `-trace` enables printing sent and received bytes in hexadecimal notation in _stderr_.

