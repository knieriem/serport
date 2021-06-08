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

When using RS-485 transceivers, additional commands can
be used to configure Linux' serial_rs485 struct.
The commands have been assigned to a separate namespace;
an extended command syntax exists to call commands of a specific namespace:

	extended_command = [ "." [ namespace_id ] "." ] command_char arg

The namespace_id is "rs485", and it may be ommitted if specified previously.
Example:

	.rs485.s1  ..[0 ..]1 ..a0 ..e0

RS-485 specific commands:

	sn	Set logical level of RTS when sending.

	an	Set logical level of RTS after sending.

	[n	After adjusting RTS, delay send by n milliseconds.

	]n	After sending, delay RTS adjustment by n milliseconds.

	en	Receive during transmission (local echo).

	tn	Enable bus termination (if supported).

See https://www.kernel.org/doc/Documentation/serial/serial-rs485.txt
for details on the corresponding flags and fields of struct serial_rs485.

For remote access to serial ports, see sub-package serial9p.
*/
package serport
