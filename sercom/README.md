* Linux

	When setting device parameters (struct termios) via
	an ioctl call in `updateCtl()`, they should be read
	back and compared to the struct that has been written.
