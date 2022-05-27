package clihelper

import (
	"fmt"
	"io/fs"
	"math"
	"os/user"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"arhat.dev/pkg/fshelper"
	"arhat.dev/pkg/numhelper"
	"arhat.dev/pkg/regexphelper"
	"github.com/spf13/pflag"
)

// FindCli
func FindCli(osfs *fshelper.OSFS, startpath string, args ...string) (_ []string, err error) {
	var (
		pfs     pflag.FlagSet
		rawOpts FindCliOptions
	)

	InitFlagSet(&pfs, "find")
	rawOpts.AddFlags(&pfs)

	err = pfs.Parse(args)
	if err != nil {
		return
	}

	fopts, err := rawOpts.Resolve(time.Now().Unix())
	if err != nil {
		return
	}

	return osfs.Find(&fopts, startpath)
}

type FindCliOptions struct {
	flags uint32

	minDepth int32
	maxDepth int32
	depth    int32

	minSize    string
	maxSize    string
	size       string
	sizeApprox string

	typ string

	unknownOwner   bool // for -nouser
	followSymlinks bool

	createdAt, createdAfter, createdBefore, createdApprox     string
	updatedAt, updatedAfter, updatedBefore, updatedApprox     string // mtime
	accessedAt, accessedAfter, accessedBefore, accessedApprox string // atime
	changedAt, changedAfter, changedBefore, changedApprox     string // metadata changed: ctime

	user  string
	group string

	regexpr, regexprIgnoreCase string
	pathPtn, pathPtnLower      string
	namePtn, namePtnLower,
	namePtnFollowSymlink, namePtnLowerFollowSymlink string
}

func (opts *FindCliOptions) AddFlags(fs *pflag.FlagSet) {
	fs.Uint32Var(&opts.flags, "flags", 0, "")

	fs.StringVarP(&opts.typ, "type", "t", "", "")

	fs.Int32Var(&opts.depth, "depth", -1, "")
	fs.Int32Var(&opts.minDepth, "min-depth", -1, "")
	fs.Int32Var(&opts.maxDepth, "max-depth", -1, "")

	fs.StringVar(&opts.size, "size", "", "")
	fs.StringVar(&opts.minSize, "min-size", "", "")
	fs.StringVar(&opts.maxSize, "max-size", "", "")
	fs.StringVar(&opts.sizeApprox, "size-approx", "512k", "")

	fs.StringVar(&opts.user, "user", "", "")
	fs.StringVar(&opts.group, "group", "", "")

	fs.StringVar(&opts.namePtn, "name", "", "")
	fs.StringVar(&opts.namePtnFollowSymlink, "lname", "", "")
	fs.StringVar(&opts.namePtnLower, "iname", "", "")
	fs.StringVar(&opts.namePtnLowerFollowSymlink, "ilname", "", "")

	fs.StringVar(&opts.pathPtn, "path", "", "")
	fs.StringVar(&opts.pathPtnLower, "ipath", "", "")

	fs.StringVar(&opts.regexpr, "regex", "", "")
	fs.StringVar(&opts.regexprIgnoreCase, "iregex", "", "")

	fs.StringVar(&opts.createdAfter, "created-after", "", "")
	fs.StringVar(&opts.createdBefore, "created-before", "", "")
	fs.StringVar(&opts.createdAt, "created-at", "", "")
	fs.StringVar(&opts.createdApprox, "created-approx", "24h", "")

	fs.StringVar(&opts.accessedAfter, "accessed-after", "", "")
	fs.StringVar(&opts.accessedBefore, "accessed-before", "", "")
	fs.StringVar(&opts.accessedAt, "accessed-at", "", "")
	fs.StringVar(&opts.accessedApprox, "accessed-approx", "24h", "")

	fs.StringVar(&opts.updatedAfter, "updated-after", "", "")
	fs.StringVar(&opts.updatedBefore, "updated-before", "", "")
	fs.StringVar(&opts.updatedAt, "updated-at", "", "")
	fs.StringVar(&opts.updatedApprox, "updated-approx", "24h", "")

	fs.StringVar(&opts.changedAfter, "changed-after", "", "")
	fs.StringVar(&opts.changedBefore, "changed-before", "", "")
	fs.StringVar(&opts.changedAt, "changed-at", "", "")
	fs.StringVar(&opts.changedApprox, "changed-approx", "24h", "")

	fs.BoolVar(&opts.unknownOwner, "owner-invalid", false, "")
}

func (opts *FindCliOptions) Resolve(now int64) (ret fshelper.FindOptions, err error) {
	if len(opts.user) != 0 {
		ret.Ops |= fshelper.FindOp_CheckUser

		var u *user.User

		switch runtime.GOOS {
		case "windows":
			u, err = user.Lookup(opts.user)
			if err != nil {
				return
			}

			ret.WindowsOrPlan9User = u.Uid
		case "plan9":
			ret.WindowsOrPlan9User = opts.user
		default:
			u, err = user.Lookup(opts.user)
			if err != nil {
				return
			}

			var uid uint64
			uid, err = strconv.ParseUint(u.Uid, 10, 32)
			if err != nil {
				return
			}

			ret.UnixUID = uint32(uid)
		}
	}

	if len(opts.group) != 0 {
		ret.Ops |= fshelper.FindOp_CheckGroup

		var g *user.Group
		g, err = user.LookupGroup(opts.group)

		if runtime.GOOS == "windows" {
			ret.WindowsOrPlan9User = g.Gid
		} else {
			var gid uint64
			gid, err = strconv.ParseUint(g.Gid, 10, 32)
			if err != nil {
				return
			}
			ret.UnixUID = uint32(gid)
		}
	}

	if len(opts.namePtn) != 0 {
		ret.Ops |= fshelper.FindOp_CheckName
		ret.NamePattern = opts.namePtn
	}

	if len(opts.namePtnLower) != 0 {
		ret.Ops |= fshelper.FindOp_CheckNameIgnoreCase
		ret.NamePatternLower = strings.ToLower(opts.namePtnLower)
	}

	if len(opts.namePtnFollowSymlink) != 0 {
		ret.Ops |= fshelper.FindOp_CheckNameFollowSymlink
		ret.NamePatternFollowSymlink = opts.namePtnFollowSymlink
	}

	if len(opts.namePtnLowerFollowSymlink) != 0 {
		ret.Ops |= fshelper.FindOp_CheckNameIgnoreCaseFollowSymlink
		ret.NamePatternLowerFollowSymlink = strings.ToLower(opts.namePtnLowerFollowSymlink)
	}

	if len(opts.pathPtn) != 0 {
		ret.Ops |= fshelper.FindOp_CheckPath
		ret.PathPattern = opts.pathPtn
	}

	if len(opts.pathPtnLower) != 0 {
		ret.Ops |= fshelper.FindOp_CheckPathIgnoreCase
		ret.PathPatternLower = strings.ToLower(opts.pathPtnLower)
	}

	if len(opts.regexpr) != 0 {
		ret.Ops |= fshelper.FindOp_CheckRegex
		ret.Regexpr, err = regexp.Compile(opts.regexpr)
		if err != nil {
			return
		}
	}

	if len(opts.regexprIgnoreCase) != 0 {
		ret.Ops |= fshelper.FindOp_CheckRegexIgnoreCase
		ret.RegexprIgnoreCase, err = regexp.Compile(regexphelper.Options{CaseInsensitive: true}.Wrap(opts.regexprIgnoreCase))
		if err != nil {
			return
		}
	}

	if len(opts.typ) != 0 {
		ret.Ops |= fshelper.FindOp_CheckTypeNotFile

		switch opts.typ {
		case "b", "block":
			// TODO: FIXME
			ret.FileType = fs.ModeDevice
		case "c", "char":
			ret.FileType = fs.ModeCharDevice
		case "d", "dir":
			ret.FileType = fs.ModeDir
		case "f", "file":
			ret.Ops &= ^fshelper.FindOp_CheckTypeNotFile
			ret.Ops |= fshelper.FindOp_CheckTypeIsFile
		case "l", "symlink":
			ret.FileType = fs.ModeSymlink
		case "p", "fifo", "namedpipe":
			ret.FileType = fs.ModeNamedPipe
		case "s", "socket":
			ret.FileType = fs.ModeSocket
		default:
			err = fmt.Errorf("unsupported file type %q", opts.typ)
			return
		}
	}

	if opts.depth > -1 {
		ret.Ops |= fshelper.FindOp_CheckDepth

		ret.MinDepth, ret.MaxDepth = opts.depth, opts.depth
	} else if opts.minDepth > -1 || opts.maxDepth > -1 {
		ret.Ops |= fshelper.FindOp_CheckDepth

		ret.MinDepth, ret.MaxDepth = opts.minDepth, opts.maxDepth
		if ret.MaxDepth < 0 {
			ret.MaxDepth = math.MaxInt32
		}

		if ret.MinDepth > ret.MaxDepth {
			err = fmt.Errorf("invalid pair of depth min (%d) > max (%d)", ret.MinDepth, ret.MaxDepth)
			return
		}
	}

	err = resolveFindOptionsMinMax(
		opts.size, opts.sizeApprox, 0, 512*1024,
		parseFindFileSize, parseFindFileSize,
		opts.minSize, opts.maxSize, 0, math.MaxInt64,
		parseFindFileSize,
		func(min, max int64) error {
			ret.Ops |= fshelper.FindOp_CheckSize
			if min > max {
				return fmt.Errorf("invalid pair of size min (%d) > max (%d)", min, max)
			}

			ret.MinSize, ret.MaxSize = min, max
			return nil
		},
	)
	if err != nil {
		return
	}

	err = resolveFindOptionsMinMax(
		opts.createdAt, opts.createdApprox, now, 24*60*60,
		parseFindFileTime, parseFindFileDuration,
		opts.createdAfter, opts.createdBefore, math.MinInt64, math.MaxInt64,
		parseFindFileTime,
		func(min, max int64) error {
			ret.Ops |= fshelper.FindOp_CheckCreationTime
			if min > max {
				return fmt.Errorf("invalid pair of creation time min (%d) > max (%d)", min, max)
			}

			ret.MinCreationTime, ret.MaxCreationTime = min, max
			return nil
		},
	)
	if err != nil {
		return
	}

	err = resolveFindOptionsMinMax(
		opts.accessedAt, opts.accessedApprox, now, 24*60*60,
		parseFindFileTime, parseFindFileDuration,
		opts.accessedAfter, opts.accessedBefore, math.MinInt64, math.MaxInt64,
		parseFindFileTime,
		func(min, max int64) error {
			ret.Ops |= fshelper.FindOp_CheckLastAccessTime
			if min > max {
				return fmt.Errorf("invalid pair of atime min (%d) > max (%d)", min, max)
			}

			ret.MinAtime, ret.MaxAtime = min, max
			return nil
		},
	)
	if err != nil {
		return
	}

	err = resolveFindOptionsMinMax(
		opts.updatedAt, opts.updatedApprox, now, 24*60*60,
		parseFindFileTime, parseFindFileDuration,
		opts.updatedAfter, opts.updatedBefore, math.MinInt64, math.MaxInt64,
		parseFindFileTime,
		func(min, max int64) error {
			ret.Ops |= fshelper.FindOp_CheckLastContentUpdatedTime
			if min > max {
				return fmt.Errorf("invalid pair of mtime min (%d) > max (%d)", min, max)
			}

			ret.MinMtime, ret.MaxMtime = min, max
			return nil
		},
	)
	if err != nil {
		return
	}

	err = resolveFindOptionsMinMax(
		opts.changedAt, opts.changedApprox, now, 24*60*60,
		parseFindFileTime, parseFindFileDuration,
		opts.changedAfter, opts.changedBefore, math.MinInt64, math.MaxInt64,
		parseFindFileTime,
		func(min, max int64) error {
			ret.Ops |= fshelper.FindOp_CheckLastMetadataChangeTime
			if min > max {
				return fmt.Errorf("invalid pair of ctime min (%d) > max (%d)", min, max)
			}

			ret.MinCtime, ret.MaxCtime = min, max
			return nil
		},
	)
	if err != nil {
		return
	}

	return
}

func parseFindFileSize(sz string, def int64) (ret int64, err error) {
	const (
		B  = 1
		KB = 1024 * B
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
		PB = 1024 * TB
	)

	n := len(sz)
	if n == 0 {
		return def, nil
	}

	switch c := sz[n-1]; c {
	default:
		if c < '0' || c > '9' {
			return -1, fmt.Errorf("invalid size format %q: expecting n[kmgtp]/[KMGTP]", sz)
		}
		n++
		fallthrough
	case 'B', 'b':
		ret = B
	case 'K', 'k':
		ret = KB
	case 'M', 'm':
		ret = MB
	case 'G', 'g':
		ret = GB
	case 'T', 't':
		ret = TB
	case 'P', 'p':
		ret = PB
	}

	iv, isFloat, err := numhelper.ParseNumber(sz[:n-1])
	if err != nil {
		return
	}

	if isFloat {
		return int64(math.Float64frombits(iv) * float64(ret)), nil
	}

	return int64(iv) * ret, nil
}

func parseFindFileDuration(dur string, def int64) (secs int64, err error) {
	const (
		SECOND = 1
		MINUTE = 60 * SECOND
		HOUR   = 60 * MINUTE
		DAY    = 24 * HOUR
		WEEK   = 7 * DAY
	)

	n := len(dur)
	if n == 0 {
		return def, nil
	}

	switch c := dur[n-1]; c {
	default:
		if c < '0' || c > '9' {
			return -1, fmt.Errorf("invalid time format %q: expecting n[wdhms]", dur)
		}
		n++
		fallthrough
	case 'd':
		secs = DAY
	case 'w':
		secs = WEEK
	case 'h':
		secs = HOUR
	case 'm':
		secs = MINUTE
	case 's':
		secs = SECOND
	}

	iv, isFloat, err := numhelper.ParseNumber(dur[:n-1])
	if err != nil {
		return
	}

	if isFloat {
		return int64(math.Float64frombits(iv) * float64(secs)), nil
	}

	return int64(iv) * secs, nil
}

func parseFindFileTime(t string, def int64) (stamp int64, err error) {
	if len(t) == 0 {
		return def, nil
	}

	// relative: [+-] duration
	if t[0] == '-' || t[0] == '+' {
		dur, err := time.ParseDuration(t)
		if err != nil {
			return time.Now().Add(dur).Unix(), nil
		}
	}

	// clock: hour:min
	ti, err := time.ParseInLocation("15:04", t, time.Local)
	if err != nil {
		y, m, d := time.Now().Date()
		return ti.AddDate(y, int(m), d).Unix(), nil
	}

	// date without time
	ti, err = time.ParseInLocation("2006-01-02", t, time.Local)
	if err == nil {
		return ti.Unix(), nil
	}

	// date full with timezone
	ti, err = time.ParseInLocation(time.RFC3339, t, time.Local)
	if err == nil {
		return ti.Unix(), nil
	}

	// date full without timezone
	ti, err = time.ParseInLocation("2006-01-02T15:04:05", t, time.Local)
	if err == nil {
		return ti.Unix(), nil
	}

	// clock: hour:min:sec
	ti, err = time.ParseInLocation("15:04:05", t, time.Local)
	if err == nil {
		y, m, d := time.Now().Date()
		return ti.AddDate(y, int(m), d).Unix(), nil
	}

	// clock: hour
	ti, err = time.ParseInLocation("15", t, time.Local)
	if err != nil {
		y, m, d := time.Now().Date()
		return ti.AddDate(y, int(m), d).Unix(), nil
	}

	err = fmt.Errorf("unrecognized time format")
	return
}

func resolveFindOptionsMinMax(
	exact, approx string,
	defaultExact, defaultApprox int64,
	parseExact, parseApprox func(string, int64) (int64, error),
	min, max string,
	defaultMin, defaultMax int64,
	parseMinMax func(string, int64) (int64, error),
	setMinMax func(min, max int64) error,
) (err error) {
	var a, v int64
	if len(exact) != 0 {
		a, err = parseApprox(approx, defaultApprox)
		if err != nil {
			return
		}

		if a < 0 {
			return fmt.Errorf("invalid approx %q < 0", approx)
		}

		v, err = parseExact(exact, 0)
		if err != nil {
			return
		}

		return setMinMax(v-a, v+a)
	} else if len(min) != 0 || len(max) != 0 {
		var minV, maxV int64
		minV, err = parseMinMax(min, defaultMin)
		if err != nil {
			return
		}

		maxV, err = parseMinMax(max, defaultMax)
		if err != nil {
			return
		}

		return setMinMax(minV, maxV)
	}

	return nil
}
