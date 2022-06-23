package on_machine

func OnTermux() bool {
	res, _, _ := Cache.Memoize("on-termux", func() (interface{}, error) {
		return commandExists("termux-setup-storage") && DirExists("/system/app/") && DirExists("/system/priv-app/"), nil
	})
	if castRes, ok := res.(bool); ok {
		return castRes
	} else {
		return false
	}
}
