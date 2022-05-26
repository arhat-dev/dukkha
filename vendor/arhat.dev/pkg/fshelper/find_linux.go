//go:build linux

package fshelper

import (
	"fmt"
	"syscall"
)

func (ofs *OSFS) matchFileSysinfo(opts *FindOptions, path string, f any) (_ bool, err error) {
	const (
		needUserGroup = FindOp_CheckUser | FindOp_CheckGroup
		needTime      = FindOp_CheckCreationTime | FindOp_CheckLastAccessTime | FindOp_CheckLastMetadataChangeTime | FindOp_CheckLastContentUpdatedTime
	)

	info, ok := f.(*syscall.Stat_t)
	if !ok {
		return false, fmt.Errorf("unexpected not *syscall.Stat_t: %T", f)
	}

	ops := opts.Ops

	if ops&FindOp_CheckUserInvalid != 0 {
		// TODO: currently os/user.LookupID is based on reading /etc/passwd when bulit without cgo
		//       doesn't work properly nor efficiently
	}

	if ops&needUserGroup != 0 {
		if ops&FindOp_CheckUser != 0 && info.Uid != opts.UnixUID {
			return false, nil
		}

		if ops&FindOp_CheckGroup != 0 && info.Gid != opts.UnixGID {
			return false, nil
		}
	}

	if ops&needTime != 0 {
		if ops&FindOp_CheckCreationTime != 0 {
			// TODO: check creation time

			var xinfo statx_t
			const (
				_AT_FDCWD = -0x64

				STATX_BTIME = 0x00000800
			)

			_, path, err = ofs.getRealPath(path)
			if err != nil {
				return
			}

			err = fstat_statx(_AT_FDCWD, path, &xinfo, STATX_BTIME)
			if err != nil {
				return
			}

			if xinfo.stx_btime.Sec < opts.MinCreationTime || xinfo.stx_btime.Sec > opts.MaxCreationTime {
				return false, nil
			}
		}

		if ops&FindOp_CheckLastAccessTime != 0 &&
			(info.Atim.Sec < opts.MinAtime ||
				info.Atim.Sec > opts.MaxAtime) {
			return false, nil
		}

		if ops&FindOp_CheckLastMetadataChangeTime != 0 &&
			(info.Ctim.Sec < opts.MinCtime ||
				info.Ctim.Sec > opts.MaxCtime) {
			return false, nil
		}

		if ops&FindOp_CheckLastContentUpdatedTime != 0 &&
			(info.Mtim.Sec < opts.MinMtime ||
				info.Mtim.Sec > opts.MaxMtime) {
			return false, nil
		}
	}

	return true, nil
}
