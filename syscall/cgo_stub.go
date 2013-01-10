// +build !cgo

package syscall

func GetUserName() string {
	return "none"
}
