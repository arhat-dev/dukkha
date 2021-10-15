package templateutils

import "arhat.dev/dukkha/pkg/constant"

var archconvNS = &_archconvNS{}

type _archconvNS struct{}

func (ns *_archconvNS) AlpineArch(arch string) string {
	v, _ := constant.GetAlpineArch(arch)
	return v
}

func (ns *_archconvNS) AlpineTripleName(arch string) string {
	v, _ := constant.GetAlpineTripleName(arch)
	return v
}

func (ns *_archconvNS) DebianArch(mArch string) string {
	v, _ := constant.GetDebianArch(mArch)
	return v
}

func (ns *_archconvNS) DebianTripleName(mArch string, other ...string) string {
	targetKernel, targetLibc := "", ""
	if len(other) > 0 {
		targetKernel = other[0]
	}
	if len(other) > 1 {
		targetLibc = other[1]
	}

	v, _ := constant.GetDebianTripleName(mArch, targetKernel, targetLibc)
	return v
}

func (ns *_archconvNS) GNUArch(mArch string) string {
	v, _ := constant.GetGNUArch(mArch)
	return v
}

func (ns *_archconvNS) GNUTripleName(mArch string, other ...string) string {
	targetKernel, targetLibc := "", ""
	if len(other) > 0 {
		targetKernel = other[0]
	}
	if len(other) > 1 {
		targetLibc = other[1]
	}

	v, _ := constant.GetGNUTripleName(mArch, targetKernel, targetLibc)
	return v
}

func (ns *_archconvNS) QemuArch(mArch string) string {
	v, _ := constant.GetQemuArch(mArch)
	return v
}

func (ns *_archconvNS) OciOS(mKernel string) string {
	v, _ := constant.GetOciOS(mKernel)
	return v
}
func (ns *_archconvNS) OciArch(mArch string) string {
	v, _ := constant.GetOciArch(mArch)
	return v
}

func (ns *_archconvNS) OciArchVariant(mArch string) string {
	v, _ := constant.GetOciArchVariant(mArch)
	return v
}

func (ns *_archconvNS) DockerOS(mKernel string) string {
	v, _ := constant.GetDockerOS(mKernel)
	return v
}

func (ns *_archconvNS) DockerArch(mArch string) string {
	v, _ := constant.GetDockerArch(mArch)
	return v
}

func (ns *_archconvNS) DockerArchVariant(mArch string) string {
	v, _ := constant.GetDockerArchVariant(mArch)
	return v
}

func (ns *_archconvNS) DockerHubArch(mArch string, other ...string) string {
	mKernel := ""
	if len(other) > 0 {
		mKernel = other[0]
	}

	v, _ := constant.GetDockerHubArch(mArch, mKernel)
	return v
}

func (ns *_archconvNS) DockerPlatformArch(mArch string) string {
	arch, ok := constant.GetDockerArch(mArch)
	if !ok {
		return ""
	}

	variant, _ := constant.GetDockerArchVariant(mArch)
	if len(variant) != 0 {
		return arch + "/" + variant
	}

	return arch
}

func (ns *_archconvNS) GolangOS(mKernel string) string {
	v, _ := constant.GetGolangOS(mKernel)
	return v
}

func (ns *_archconvNS) GolangArch(mArch string) string {
	v, _ := constant.GetGolangArch(mArch)
	return v
}

func (ns *_archconvNS) LLVMArch(mArch string) string {
	v, _ := constant.GetLLVMArch(mArch)
	return v
}

func (ns *_archconvNS) LLVMTripleName(mArch string, other ...string) string {
	targetKernel, targetLibc := "", ""
	if len(other) > 0 {
		targetKernel = other[0]
	}
	if len(other) > 1 {
		targetLibc = other[1]
	}

	v, _ := constant.GetLLVMTripleName(mArch, targetKernel, targetLibc)
	return v
}
