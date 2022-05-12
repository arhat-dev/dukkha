package constant

// GetGolangOS get GOOS value by kernel value, return true if kernel value is known to this package
func GetGolangOS(kernel string) (string, bool) {
	if kid := kernel_id_of(kernel); kid != _unknown_kernel {
		return kid.String(), true
	}

	return kernel, false
}

// GetDockerOS is currently an alias of GetGolangOS
func GetDockerOS(kernel string) (string, bool) { return GetGolangOS(kernel) }

// GetOciOS is currently an alias of GetGolangOS
func GetOciOS(kernel string) (string, bool) { return GetGolangOS(kernel) }
