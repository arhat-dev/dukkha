package templateutils

import (
	"arhat.dev/dukkha/pkg/dukkha"
	"arhat.dev/pkg/clihelper"
	"github.com/spf13/pflag"
)

func createTagNS(rc dukkha.RenderingContext) tagNS {
	return tagNS{rc: rc}
}

type tagNS struct {
	rc dukkha.RenderingContext
}

type tagOptions struct {
	keepKernelInfo bool
}

func (opts *tagOptions) addFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&opts.keepKernelInfo, "keep-kernel-info", false, "")
}

func handleImageTagging(args []String, do func(opts *tagOptions, name string) (string, error)) (_ string, err error) {
	n := len(args)
	if n == 0 {
		err = errAtLeastOneArgGotZero
		return
	}

	flags, err := toStrings(args)
	if err != nil {
		return
	}

	var opts tagOptions

	if n > 1 {
		var fs pflag.FlagSet
		clihelper.InitFlagSet(&fs, "tag")
		opts.addFlags(&fs)

		err = fs.Parse(flags[:n-1])
		if err != nil {
			return
		}
	}

	return do(&opts, flags[n-1])
}

// ImageName generates a tagged image name for the last argument based on git, matrix info
func (ns tagNS) ImageName(args ...String) (string, error) {
	return handleImageTagging(args, func(opts *tagOptions, name string) (string, error) {
		return GetFullImageName_UseDefault_IfIfNoTagSet(ns.rc, name, opts.keepKernelInfo), nil
	})
}

// ManifestName generates a tagged image name for the last argument based on git, matrix info
func (ns tagNS) ManifestName(args ...String) (string, error) {
	return handleImageTagging(args, func(opts *tagOptions, name string) (string, error) {
		return GetFullManifestName_UseDefault_IfNoTagSet(ns.rc, name), nil
	})
}

// ImageTag generates a tag for container image name (the last argument) based on git, matrix info
func (ns tagNS) ImageTag(args ...String) (string, error) {
	return handleImageTagging(args, func(opts *tagOptions, name string) (string, error) {
		return GetImageTag(ns.rc, name, opts.keepKernelInfo), nil
	})
}

// ManifestTag generates a tag for container manifest name (the last arugment)
func (ns tagNS) ManifestTag(args ...String) (_ string, err error) {
	return handleImageTagging(args, func(opts *tagOptions, name string) (string, error) {
		return GetManifestTag(ns.rc, name), nil
	})
}
