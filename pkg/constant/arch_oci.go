package constant

func GetOCIArch(mArch string) string {
	return GetGolangArch(mArch)
}

func GetOCIArchVariant(mArch string) string {
	return GetDockerArchVariant(mArch)
}
