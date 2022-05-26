package templateutils

import (
	"context"
	"net"

	"arhat.dev/pkg/clihelper"
	"github.com/spf13/pflag"
)

type dnsNS struct{}

// IP lookup A and AAA records of host
func (dnsNS) IP(args ...String) (ret []string, err error) {
	err = handleDNSTemplateFunc_HOST(args, func(ctx context.Context, r *net.Resolver, host string) error {
		v, err := r.LookupNetIP(ctx, "ip", host)
		if err != nil {
			return err
		}

		ret, err = toStrings(v)
		return err
	})

	return
}

// CNAME lookup CNAME record of host
func (dnsNS) CNAME(args ...String) (ret string, err error) {
	err = handleDNSTemplateFunc_HOST(args, func(ctx context.Context, r *net.Resolver, host string) error {
		ret, err = r.LookupCNAME(ctx, host)
		return err
	})

	return
}

func (dnsNS) HOST(args ...String) (ret []string, err error) {
	err = handleDNSTemplateFunc_HOST(args, func(ctx context.Context, r *net.Resolver, s string) error {
		ret, err = r.LookupHost(ctx, s)
		return err
	})

	return
}

// SRV lookup SRV records
func (dnsNS) SRV(args ...String) (ret []*net.SRV, err error) {
	// TODO: support _service, _proto

	err = handleDNSTemplateFunc_HOST(args, func(ctx context.Context, r *net.Resolver, host string) error {
		_, ret, err = r.LookupSRV(ctx, "", "", host)
		return err
	})

	return
}

// TXT lookup TXT records
func (dnsNS) TXT(args ...String) (ret []string, err error) {
	err = handleDNSTemplateFunc_HOST(args, func(ctx context.Context, r *net.Resolver, host string) error {
		ret, err = r.LookupTXT(ctx, host)
		return err
	})

	return
}

func handleDNSTemplateFunc_HOST(args []String, do func(context.Context, *net.Resolver, string) error) error {
	n := len(args)
	if n == 0 {
		return errAtLeastOneArgGotZero
	}

	if n == 1 {
		return do(context.TODO(), net.DefaultResolver, must(toString(args[0])))
	}

	flags, err := toStrings(args[:n-1])
	if err != nil {
		return err
	}

	opts, err := parseDNSOptions(flags)
	if err != nil {
		return err
	}

	return opts.lookup(func(r *net.Resolver) error {
		return do(context.TODO(), r, must(toString(args[n-1])))
	})
}

type dnsOptions struct {
	nameservers []string
}

func (opts dnsOptions) lookup(do func(*net.Resolver) error) error {
	if len(opts.nameservers) != 0 {
		// TODO: support DoT,
		resolver := net.Resolver{
			PreferGo: true,
		}

		return do(&resolver)
	} else {
		return do(net.DefaultResolver)
	}
}

func parseDNSOptions(flags []string) (ret dnsOptions, err error) {
	var (
		fs pflag.FlagSet
	)

	clihelper.InitFlagSet(&fs, "dns")

	fs.StringSliceVarP(&ret.nameservers, "nameserver", "n", nil, "")

	err = fs.Parse(flags)
	if err != nil {
		return
	}

	return
}
