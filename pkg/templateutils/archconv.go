package templateutils

import (
	"arhat.dev/dukkha/pkg/constant"
)

type archconvNS struct{}

func (archconvNS) SimpleArch(arch String) string { return constant.SimpleArch(toString(arch)) }

// HF is an alias of HardFloadArch
func (archconvNS) HF(arch String) string            { return constant.HardFloadArch(toString(arch)) }
func (archconvNS) HardFloadArch(arch String) string { return constant.HardFloadArch(toString(arch)) }

// SF is an alias of SoftFloadArch
func (archconvNS) SF(arch String) string            { return constant.SoftFloadArch(toString(arch)) }
func (archconvNS) SoftFloadArch(arch String) string { return constant.SoftFloadArch(toString(arch)) }

func (archconvNS) AlpineArch(arch String) string {
	v, _ := constant.GetAlpineArch(toString(arch))
	return v
}

func (archconvNS) AlpineTripleName(arch String) string {
	v, _ := constant.GetAlpineTripleName(toString(arch))
	return v
}

func (archconvNS) DebianArch(arch String) string {
	v, _ := constant.GetDebianArch(toString(arch))
	return v
}

func (archconvNS) DebianTripleName(arch String, other ...String) string {
	var targetKernel, targetLibc string
	if len(other) > 0 {
		targetKernel = toString(other[0])
	}
	if len(other) > 1 {
		targetLibc = toString(other[1])
	}

	v, _ := constant.GetDebianTripleName(toString(arch), targetKernel, targetLibc)
	return v
}

func (archconvNS) GNUArch(arch String) string {
	v, _ := constant.GetGNUArch(toString(arch))
	return v
}

func (archconvNS) GNUTripleName(arch String, other ...String) string {
	var targetKernel, targetLibc string
	if len(other) > 0 {
		targetKernel = toString(other[0])
	}
	if len(other) > 1 {
		targetLibc = toString(other[1])
	}

	v, _ := constant.GetGNUTripleName(toString(arch), targetKernel, targetLibc)
	return v
}

func (archconvNS) QemuArch(arch String) string {
	v, _ := constant.GetQemuArch(toString(arch))
	return v
}

func (archconvNS) OciOS(mKernel string) string {
	v, _ := constant.GetOciOS(mKernel)
	return v
}

func (archconvNS) OciArch(arch String) string {
	v, _ := constant.GetOciArch(toString(arch))
	return v
}

func (archconvNS) OciArchVariant(arch String) string {
	v, _ := constant.GetOciArchVariant(toString(arch))
	return v
}

func (archconvNS) DockerOS(mKernel string) string {
	v, _ := constant.GetDockerOS(mKernel)
	return v
}

func (archconvNS) DockerArch(arch String) string {
	v, _ := constant.GetDockerArch(toString(arch))
	return v
}

func (archconvNS) DockerArchVariant(arch String) string {
	v, _ := constant.GetDockerArchVariant(toString(arch))
	return v
}

func (archconvNS) DockerHubArch(arch String, other ...String) string {
	mKernel := ""
	if len(other) > 0 {
		mKernel = toString(other[0])
	}

	v, _ := constant.GetDockerHubArch(toString(arch), mKernel)
	return v
}

func (archconvNS) DockerPlatformArch(arch String) string {
	mArch, ok := constant.GetDockerArch(toString(arch))
	if !ok {
		return ""
	}

	variant, _ := constant.GetDockerArchVariant(toString(arch))
	if len(variant) != 0 {
		return mArch + "/" + variant
	}

	return mArch
}

func (archconvNS) GolangOS(kernel String) string {
	v, _ := constant.GetGolangOS(toString(kernel))
	return v
}

func (archconvNS) GolangArch(arch String) string {
	v, _ := constant.GetGolangArch(toString(arch))
	return v
}

func (archconvNS) LLVMArch(arch String) string {
	v, _ := constant.GetLLVMArch(toString(arch))
	return v
}

func (archconvNS) LLVMTripleName(arch String, other ...String) string {
	targetKernel, targetLibc := "", ""
	if len(other) > 0 {
		targetKernel = toString(other[0])
	}
	if len(other) > 1 {
		targetLibc = toString(other[1])
	}

	v, _ := constant.GetLLVMTripleName(toString(arch), targetKernel, targetLibc)
	return v
}
