/*
Package serport provides access to local, and remote (9P served)
serial ports on Linux and Windows. The Port interface is still
incomplete -- basic access to serial ports works, but some
functions might be missing.

Besides a number of SetXY() methods that allow modification
of single parameters, a method named Ctl() is available that
processes a string containing one or more commands as described
in http://plan9.bell-labs.com/magic/man2html/3/uart.  Only a
subset is supported:

	bn	Set the baud rate to n.

 	dn	Set DTR if n is non-zero; else clear it.

	//f	Flush output queue.

	ln	Set number of bits per byte to n. Legal values are 5,
		6, 7, or 8.

	mn	Obey modem CTS signal if n is non-zero; else clear it.

	pc	Set parity to odd if c is o, to even if c is e; else
		set no parity.

	rn	Set RTS if n is non-zero; else clear it.

	sn	Set number of stop bits to n. Legal values are 1 or 2.

  Additional commands:

	Dn	Delay exection for n milli-seconds

	Wn	Write a byte with value n

	{	Postpone execution of commands until '}' is sent.

	}	Execute pending commands

*/
package serport
