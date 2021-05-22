/*
Package serport provides access to local serial ports on
Linux and Windows.

Besides a number of SetXY() methods that allow modification
of single parameters, a method named Ctl() is available that
processes a string containing one or more commands as described in
https://plan9.io/magic/man2html/3/uart.  The following subset
is supported:

	bn	Set the baud rate to n.

 	dn	Set DTR if n is non-zero; else clear it.

	kn	Send a break lasting n milliseconds.

	ln	Set number of bits per byte to n. Legal values are 5,
		6, 7, or 8.

	mn	Obey modem CTS signal if n is non-zero; else clear it.

	pc	Set parity to odd if c is o, to even if c is e; else
		set no parity.

	rn	Set RTS if n is non-zero; else clear it.

	sn	Set number of stop bits to n. Legal values are 1 or 2.

  Additional commands:

	Dn	Delay execution for n milli-seconds

	Wn	Write a byte with value n

	{	Postpone execution of commands until '}' is sent.

	}	Execute pending commands

For remote access to serial ports, see sub-package serial9p.
*/
package serport
