package templateutils

import (
	"arhat.dev/dukkha/pkg/constant"
)

type archconvNS struct{}

func (archconvNS) SimpleArch(arch String) string { return constant.SimpleArch(must(toString(arch))) }

// HF is an alias of HardFloadArch
func (archconvNS) HF(arch String) string { return constant.HardFloadArch(must(toString(arch))) }
func (archconvNS) HardFloadArch(arch String) string {
	return constant.HardFloadArch(must(toString(arch)))
}

// SF is an alias of SoftFloadArch
func (archconvNS) SF(arch String) string { return constant.SoftFloadArch(must(toString(arch))) }
func (archconvNS) SoftFloadArch(arch String) string {
	return constant.SoftFloadArch(must(toString(arch)))
}

func (archconvNS) AlpineArch(arch String) string {
	v, _ := constant.GetAlpineArch(must(toString(arch)))
	return v
}

func (archconvNS) AlpineTripleName(arch String) string {
	v, _ := constant.GetAlpineTripleName(must(toString(arch)))
	return v
}

func (archconvNS) DebianArch(arch String) string {
	v, _ := constant.GetDebianArch(must(toString(arch)))
	return v
}

func (archconvNS) DebianTripleName(arch String, other ...String) string {
	var targetKernel, targetLibc string
	if len(other) > 0 {
		targetKernel = must(toString(other[0]))
	}
	if len(other) > 1 {
		targetLibc = must(toString(other[1]))
	}

	v, _ := constant.GetDebianTripleName(must(toString(arch)), targetKernel, targetLibc)
	return v
}

func (archconvNS) GNUArch(arch String) string {
	v, _ := constant.GetGNUArch(must(toString(arch)))
	return v
}

func (archconvNS) GNUTripleName(arch String, other ...String) string {
	var targetKernel, targetLibc string
	if len(other) > 0 {
		targetKernel = must(toString(other[0]))
	}
	if len(other) > 1 {
		targetLibc = must(toString(other[1]))
	}

	v, _ := constant.GetGNUTripleName(must(toString(arch)), targetKernel, targetLibc)
	return v
}

func (archconvNS) QemuArch(arch String) string {
	v, _ := constant.GetQemuArch(must(toString(arch)))
	return v
}

func (archconvNS) OciOS(mKernel string) string {
	v, _ := constant.GetOciOS(mKernel)
	return v
}

func (archconvNS) OciArch(arch String) string {
	v, _ := constant.GetOciArch(must(toString(arch)))
	return v
}

func (archconvNS) OciArchVariant(arch String) string {
	v, _ := constant.GetOciArchVariant(must(toString(arch)))
	return v
}

func (archconvNS) DockerOS(mKernel string) string {
	v, _ := constant.GetDockerOS(mKernel)
	return v
}

func (archconvNS) DockerArch(arch String) string {
	v, _ := constant.GetDockerArch(must(toString(arch)))
	return v
}

func (archconvNS) DockerArchVariant(arch String) string {
	v, _ := constant.GetDockerArchVariant(must(toString(arch)))
	return v
}

func (archconvNS) DockerHubArch(arch String, other ...String) string {
	mKernel := ""
	if len(other) > 0 {
		mKernel = must(toString(other[0]))
	}

	v, _ := constant.GetDockerHubArch(must(toString(arch)), mKernel)
	return v
}

func (archconvNS) DockerPlatformArch(arch String) string {
	mArch, ok := constant.GetDockerArch(must(toString(arch)))
	if !ok {
		return ""
	}

	variant, _ := constant.GetDockerArchVariant(must(toString(arch)))
	if len(variant) != 0 {
		return mArch + "/" + variant
	}

	return mArch
}

func (archconvNS) GolangOS(kernel String) string {
	v, _ := constant.GetGolangOS(must(toString(kernel)))
	return v
}

func (archconvNS) GolangArch(arch String) string {
	v, _ := constant.GetGolangArch(must(toString(arch)))
	return v
}

func (archconvNS) LLVMArch(arch String) string {
	v, _ := constant.GetLLVMArch(must(toString(arch)))
	return v
}

func (archconvNS) LLVMTripleName(arch String, other ...String) string {
	targetKernel, targetLibc := "", ""
	if len(other) > 0 {
		targetKernel = must(toString(other[0]))
	}
	if len(other) > 1 {
		targetLibc = must(toString(other[1]))
	}

	v, _ := constant.GetLLVMTripleName(must(toString(arch)), targetKernel, targetLibc)
	return v
}
