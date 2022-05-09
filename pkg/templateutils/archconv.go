package templateutils

import (
	"arhat.dev/dukkha/pkg/constant"
	"arhat.dev/pkg/archconst"
)

type archconvNS struct{}

func (ns archconvNS) SimpleArch(arch string) string { return constant.SimpleArch(arch) }

// HF is an alias of HardFloadArch
func (ns archconvNS) HF(arch string) string { return ns.HardFloadArch(arch) }

// HardFloadArch returns hard-float version of arch
func (ns archconvNS) HardFloadArch(arch string) string {
	spec, ok := archconst.Split(arch)
	if !ok {
		return string(arch)
	}

	spec.SoftFloat = false
	return spec.String()
}

func (ns archconvNS) AlpineArch(arch string) string {
	v, _ := constant.GetAlpineArch(arch)
	return v
}

func (ns archconvNS) AlpineTripleName(arch string) string {
	v, _ := constant.GetAlpineTripleName(arch)
	return v
}

func (ns archconvNS) DebianArch(mArch string) string {
	v, _ := constant.GetDebianArch(mArch)
	return v
}

func (ns archconvNS) DebianTripleName(mArch string, other ...string) string {
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

func (ns archconvNS) GNUArch(mArch string) string {
	v, _ := constant.GetGNUArch(mArch)
	return v
}

func (ns archconvNS) GNUTripleName(mArch string, other ...string) string {
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

func (ns archconvNS) QemuArch(mArch string) string {
	v, _ := constant.GetQemuArch(mArch)
	return v
}

func (ns archconvNS) OciOS(mKernel string) string {
	v, _ := constant.GetOciOS(mKernel)
	return v
}

func (ns archconvNS) OciArch(mArch string) string {
	v, _ := constant.GetOciArch(mArch)
	return v
}

func (ns archconvNS) OciArchVariant(mArch string) string {
	v, _ := constant.GetOciArchVariant(mArch)
	return v
}

func (ns archconvNS) DockerOS(mKernel string) string {
	v, _ := constant.GetDockerOS(mKernel)
	return v
}

func (ns archconvNS) DockerArch(mArch string) string {
	v, _ := constant.GetDockerArch(mArch)
	return v
}

func (ns archconvNS) DockerArchVariant(mArch string) string {
	v, _ := constant.GetDockerArchVariant(mArch)
	return v
}

func (ns archconvNS) DockerHubArch(mArch string, other ...string) string {
	mKernel := ""
	if len(other) > 0 {
		mKernel = other[0]
	}

	v, _ := constant.GetDockerHubArch(mArch, mKernel)
	return v
}

func (ns archconvNS) DockerPlatformArch(mArch string) string {
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

func (ns archconvNS) GolangOS(mKernel string) string {
	v, _ := constant.GetGolangOS(mKernel)
	return v
}

func (ns archconvNS) GolangArch(mArch string) string {
	v, _ := constant.GetGolangArch(mArch)
	return v
}

func (ns archconvNS) LLVMArch(mArch string) string {
	v, _ := constant.GetLLVMArch(mArch)
	return v
}

func (ns archconvNS) LLVMTripleName(mArch string, other ...string) string {
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
