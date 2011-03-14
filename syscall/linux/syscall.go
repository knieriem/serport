package syscall

//sys	Fcntl(fd int, cmd int, arg int) (val int, errno int)

//sys	IoctlTermios(fd int, action int, t *Termios) (errno int) = SYS_IOCTL
//sys	IoctlModem(fd int, action int, flags *Int) (errno int) = SYS_IOCTL

func (t *Termios) SetInSpeed(s int) {
//	t.Iflag = t.Iflag&^CBAUD | uint32(s)&CBAUD
}
func (t *Termios) SetOutSpeed(s int) {
	t.Cflag = t.Cflag&^CBAUD | uint32(s)&CBAUD
}
