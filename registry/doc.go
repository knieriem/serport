/*
A package allowing access to Windows' Registry database.
At the moment read-only access is implemented only.

Example: The following code prints a list of serial devices
currently present:

	key, err := registry.KeyLocalMachine.Subkey("HARDWARE", "DEVICEMAP", "SERIALCOMM")
	if err != nil {
		return
	}
	for _, v := range key.Values() {
		if s := v.String(); s != "" {
			fmt.Println(s)
		}
	}
*/
package registry
