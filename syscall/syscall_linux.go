package syscall

//sys	Fcntl(fd uintptr, cmd int, arg int) (val int, err error)

//sys	IoctlTermios(fd uintptr, action int, t *Termios) (err error) = SYS_IOCTL
//sys	IoctlModem(fd uintptr, action int, flags *Int) (err error) = SYS_IOCTL

func (t *Termios) SetInSpeed(s int) {
	//	t.Iflag = t.Iflag&^CBAUD | uint32(s)&CBAUD
}
func (t *Termios) SetOutSpeed(s int) {
	t.Cflag = t.Cflag&^CBAUD | uint32(s)&CBAUD
}
