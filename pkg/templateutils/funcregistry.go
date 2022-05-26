package templateutils

import (
	"reflect"
)

type funcID uint32

type Funcs [funcID_COUNT]any

// var staticFuncs = Funcs{
// 	_unknown_template_func: nil,
//
// 	// funcID_ago:              dateAgo,
// 	// funcID_date:             date,
// 	// funcID_date_in_zone:     dateInZone,
// 	// funcID_date_modify:      dateModify,
// 	// funcID_dateInZone:       dateInZone,
// 	// funcID_dateModify:       dateModify,
// 	// funcID_duration:         duration,
// 	// funcID_durationRound:    durationRound,
// 	// funcID_htmlDate:         htmlDate,
// 	// funcID_htmlDateInZone:   htmlDateInZone,
// 	// funcID_must_date_modify: mustDateModify,
// 	// funcID_mustDateModify:   mustDateModify,
// 	// funcID_mustToDate:       mustToDate,
// 	// funcID_now:              time.Now,
// 	// funcID_toDate:    toDate,
// 	// funcID_unixEpoch: unixEpoch,
//
// 	// Strings
// 	funcID_randAlphaNum: randAlphaNumeric,
// 	funcID_randAlpha:    randAlpha,
// 	funcID_randAscii:    randAscii,
// 	funcID_randNumeric:  randNumeric,
// 	funcID_shuffle:      xstrings.Shuffle,
// 	funcID_wrap:         func(l int, s string) string { return util.Wrap(s, l) },
// 	funcID_wrapWith:     func(l int, sep, str string) string { return util.WrapCustom(str, l, sep, true) },
//
// 	// funcID_sha1sum:    sha1sum,
// 	// funcID_sha256sum:  sha256sum,
// 	// funcID_adler32sum: adler32sum,
// 	// funcID_toString: strval,
//
// 	// Wrap Atoi to stop errors.
// 	funcID_atoi: func(a string) int { i, _ := strconv.Atoi(a); return i },
// 	// funcID_int64:     toInt64,
// 	// funcID_int:       toInt,
// 	// funcID_float64:   toFloat64,
// 	// funcID_seq: seq,
// 	// funcID_toDecimal: toDecimal,
//
// 	// gt: func(a, b int) bool {return a > b},
// 	// gte: func(a, b int) bool {return a >= b},
// 	// lt: func(a, b int) bool {return a < b},
// 	// lte: func(a, b int) bool {return a <= b},
//
// 	// VERY basic arithmetic.
// 	// funcID_add1: func(i interface{}) int64 { return toInt64(i) + 1 },
// 	// funcID_add: func(i ...interface{}) int64 {
// 	// 	var a int64 = 0
// 	// 	for _, b := range i {
// 	// 		a += toInt64(b)
// 	// 	}
// 	// 	return a
// 	// },
// 	// funcID_sub: func(a, b interface{}) int64 { return toInt64(a) - toInt64(b) },
// 	// funcID_div: func(a, b interface{}) int64 { return toInt64(a) / toInt64(b) },
// 	// funcID_mod: func(a, b interface{}) int64 { return toInt64(a) % toInt64(b) },
// 	// funcID_mul: func(a interface{}, v ...interface{}) int64 {
// 	// 	val := toInt64(a)
// 	// 	for _, b := range v {
// 	// 		val = val * toInt64(b)
// 	// 	}
// 	// 	return val
// 	// },
// 	funcID_randInt: func(min, max int) int { return rand.Intn(max-min) + min },
// 	// funcID_biggest: max,
// 	// funcID_max:     max,
// 	// funcID_min:     min,
// 	// funcID_maxf:    maxf,
// 	// funcID_minf:    minf,
// 	// funcID_ceil:    ceil,
// 	// funcID_floor:   floor,
// 	// funcID_round:   round,
//
// 	// string slices. Note that we reverse the order b/c that's better
// 	// for template processing.
// 	// funcID_join:      join,
// 	// funcID_sortAlpha: sortAlpha,
//
// 	// Defaults
// 	// funcID_default:          dfault,
// 	// funcID_empty:            empty,
// 	// funcID_coalesce:         coalesce,
// 	// funcID_all:              allTrue,
// 	// funcID_any:              anyTrue,
// 	// funcID_fromJson:         fromJson,
// 	// funcID_toJson:           toJson,
// 	// funcID_toPrettyJson:     toPrettyJson,
// 	// funcID_toRawJson:        toRawJson,
// 	// funcID_mustFromJson:     mustFromJson,
// 	// funcID_mustToJson:       mustToJson,
// 	// funcID_mustToPrettyJson: mustToPrettyJson,
// 	// funcID_mustToRawJson:    mustToRawJson,
// 	funcID_ternary: ternary,
// 	// funcID_deepCopy:         deepCopy,
// 	// funcID_mustDeepCopy:     mustDeepCopy,
//
// 	// Reflection
//
// 	funcID_deepEqual: reflect.DeepEqual,
//
// 	// OS:
// 	funcID_env:       os.Getenv,
// 	funcID_expandenv: os.ExpandEnv,
//
// 	// Network:
//
// 	// Paths:
// 	funcID_base:  path.Base,
// 	funcID_dir:   path.Dir,
// 	funcID_clean: path.Clean,
// 	funcID_ext:   path.Ext,
// 	funcID_isAbs: path.IsAbs,
//
// 	// Filepaths:
// 	funcID_osBase:  filepath.Base,
// 	funcID_osClean: filepath.Clean,
// 	funcID_osDir:   filepath.Dir,
// 	funcID_osExt:   filepath.Ext,
// 	funcID_osIsAbs: filepath.IsAbs,
//
// 	// Encoding:
// 	// funcID_b64enc: base64encode,
// 	// funcID_b64dec: base64decode,
// 	// funcID_b32enc: base32encode,
// 	// funcID_b32dec: base32decode,
//
// 	// Data Structures:
// 	// funcID_tuple:              list, // FIXME: with the addition of append/prepend these are no longer immutable.
// 	// funcID_list:               list,
// 	// funcID_dict:               dict,
// 	// funcID_get:                get,
// 	funcID_set:   set,
// 	funcID_unset: unset,
// 	// funcID_hasKey:             hasKey,
// 	// funcID_pluck: pluck,
// 	// funcID_keys:               keys,
// 	// funcID_pick:               pick,
// 	// funcID_omit:               omit,
// 	// funcID_merge:              merge,
// 	// funcID_mergeOverwrite:     mergeOverwrite,
// 	// funcID_mustMerge:          mustMerge,
// 	// funcID_mustMergeOverwrite: mustMergeOverwrite,
// 	// funcID_values:             values,
//
// 	// funcID_append:      push,
// 	// funcID_push:        push,
// 	// funcID_mustAppend:  mustPush,
// 	// funcID_mustPush:    mustPush,
// 	// funcID_prepend:     prepend,
// 	// funcID_mustPrepend: mustPrepend,
// 	// funcID_first:       first,
// 	// funcID_mustFirst:   mustFirst,
// 	// funcID_rest:     rest,
// 	// funcID_mustRest: mustRest,
// 	// funcID_last:        last,
// 	// funcID_mustLast:    mustLast,
// 	// funcID_reverse:     reverse,
// 	// funcID_mustReverse: mustReverse,
// 	// funcID_uniq:        uniq,
// 	// funcID_mustUniq:    mustUniq,
// 	// funcID_without:     without,
// 	// funcID_mustWithout: mustWithout,
// 	// funcID_has:         has,
// 	// funcID_mustHas:     mustHas,
// 	// funcID_slice:     slice,
// 	// funcID_mustSlice: mustSlice,
// 	// funcID_dig:       dig,
// 	funcID_chunk:     chunk,
// 	funcID_mustChunk: mustChunk,
//
// 	// Crypto:
// 	// funcID_bcrypt:                   bcrypt,
// 	// funcID_htpasswd:                 htpasswd,
// 	funcID_genPrivateKey:            generatePrivateKey,
// 	funcID_derivePassword:           derivePassword,
// 	funcID_buildCustomCert:          buildCustomCertificate,
// 	funcID_genCA:                    generateCertificateAuthority,
// 	funcID_genCAWithKey:             generateCertificateAuthorityWithPEMKey,
// 	funcID_genSelfSignedCert:        generateSelfSignedCertificate,
// 	funcID_genSelfSignedCertWithKey: generateSelfSignedCertificateWithPEMKey,
// 	funcID_genSignedCert:            generateSignedCertificate,
// 	funcID_genSignedCertWithKey:     generateSignedCertificateWithPEMKey,
// 	funcID_encryptAES:               encryptAES,
// 	funcID_decryptAES:               decryptAES,
// 	funcID_randBytes:                randBytes,
//
// 	// UUIDs:
// 	// funcID_uuidv4: uuidv4,
//
// 	// SemVer:
// 	funcID_semver:        semver,
// 	funcID_semverCompare: semverCompare,
//
// 	// Flow Control:
// 	funcID_fail: func(msg string) (string, error) { return "", errors.New(msg) },
//
// 	// Regex
// 	// funcID_regexMatch:                 regexMatch,
// 	// funcID_mustRegexMatch:             mustRegexMatch,
// 	// funcID_regexFindAll:               regexFindAll,
// 	// funcID_mustRegexFindAll:           mustRegexFindAll,
// 	// funcID_regexFind:                  regexFind,
// 	// funcID_mustRegexFind:              mustRegexFind,
// 	// funcID_regexReplaceAll:            regexReplaceAll,
// 	// funcID_mustRegexReplaceAll:        mustRegexReplaceAll,
// 	// funcID_regexReplaceAllLiteral:     regexReplaceAllLiteral,
// 	// funcID_mustRegexReplaceAllLiteral: mustRegexReplaceAllLiteral,
// 	// funcID_regexSplit:                 regexSplit,
// 	// funcID_mustRegexSplit:             mustRegexSplit,
// 	// funcID_regexQuoteMeta:             regexQuoteMeta,
//
// 	// URLs:
// 	funcID_urlParse: urlParse,
// 	funcID_urlJoin:  urlJoin,
//
// 	funcID_abbrev:     stringsNS{}.Abbrev,
// 	funcID_abbrevboth: abbrevboth,
// 	// funcID_add
// 	// funcID_add1
// 	// funcID_add1f
// 	// funcID_addPrefix
// 	// funcID_addSuffix
// 	// funcID_addf
// 	// funcID_adler32sum
// 	// funcID_all
// 	// funcID_any
// 	// funcID_append
// 	funcID_archconv_AlpineArch:         archconvNS{}.AlpineArch,
// 	funcID_archconv_AlpineTripleName:   archconvNS{}.AlpineTripleName,
// 	funcID_archconv_DebianArch:         archconvNS{}.DebianArch,
// 	funcID_archconv_DebianTripleName:   archconvNS{}.DebianTripleName,
// 	funcID_archconv_DockerArch:         archconvNS{}.DockerArch,
// 	funcID_archconv_DockerArchVariant:  archconvNS{}.DockerArchVariant,
// 	funcID_archconv_DockerHubArch:      archconvNS{}.DockerHubArch,
// 	funcID_archconv_DockerOS:           archconvNS{}.DockerOS,
// 	funcID_archconv_DockerPlatformArch: archconvNS{}.DockerPlatformArch,
// 	funcID_archconv_GNUArch:            archconvNS{}.GNUArch,
// 	funcID_archconv_GNUTripleName:      archconvNS{}.GNUTripleName,
// 	funcID_archconv_GolangArch:         archconvNS{}.GolangArch,
// 	funcID_archconv_GolangOS:           archconvNS{}.GolangOS,
// 	funcID_archconv_HF:                 archconvNS{}.HardFloadArch,
// 	funcID_archconv_HardFloadArch:      archconvNS{}.HardFloadArch,
// 	funcID_archconv_LLVMArch:           archconvNS{}.LLVMArch,
// 	funcID_archconv_LLVMTripleName:     archconvNS{}.LLVMTripleName,
// 	funcID_archconv_OciArch:            archconvNS{}.OciArch,
// 	funcID_archconv_OciArchVariant:     archconvNS{}.OciArchVariant,
// 	funcID_archconv_OciOS:              archconvNS{}.OciOS,
// 	funcID_archconv_QemuArch:           archconvNS{}.QemuArch,
// 	funcID_archconv_SF:                 archconvNS{}.SoftFloadArch,
// 	funcID_archconv_SimpleArch:         archconvNS{}.SimpleArch,
// 	funcID_archconv_SoftFloadArch:      archconvNS{}.SoftFloadArch,
// 	// funcID_atoi
// 	funcID_b32dec: encNS{}.Base32,
// 	funcID_b64dec: encNS{}.Base64,
// 	// funcID_base
// 	// funcID_bcrypt
// 	// funcID_biggest
// 	// funcID_bool
// 	// funcID_buildCustomCert
// 	funcID_camelcase: stringsNS{}.CamelCase,
// 	// funcID_cat
// 	// funcID_ceil
// 	// funcID_chunk
// 	// funcID_clean
// 	// funcID_coalesce:
// 	funcID_coll_Append:  collNS{}.Append,
// 	funcID_coll_Dict:    collNS{}.MapStringAny,
// 	funcID_coll_Flatten: collNS{}.Flatten,
// 	funcID_coll_Has:     collNS{}.HasAll,
// 	funcID_coll_Keys:    collNS{}.Keys,
// 	funcID_coll_Merge:   collNS{}.Merge,
// 	funcID_coll_Omit:    collNS{}.Omit,
// 	funcID_coll_Pick:    collNS{}.Pick,
// 	funcID_coll_Prepend: collNS{}.Prepend,
// 	funcID_coll_Reverse: collNS{}.Reverse,
// 	funcID_coll_Slice:   collNS{}.List,
// 	funcID_coll_Sort:    collNS{}.Sort,
// 	funcID_coll_Uniq:    collNS{}.Unique,
// 	funcID_coll_Values:  collNS{}.Values,
// 	// funcID_compact
// 	// funcID_concat
// 	funcID_contains: stringsNS{}.Contains,
// 	// funcID_conv_Atoi
// 	// funcID_conv_Bool
// 	// funcID_conv_Default:
// 	// funcID_conv_Dict: conv.Dict,
// 	// funcID_conv_Has:  conv.Has,
// 	// funcID_conv_Join: conv.Join,
// 	// funcID_conv_ParseFloat
// 	// funcID_conv_ParseInt
// 	// funcID_conv_ParseUint
// 	// funcID_conv_Slice
// 	// funcID_conv_ToBool
// 	// funcID_conv_ToBools
// 	// funcID_conv_ToFloat64
// 	// funcID_conv_ToFloat64s
// 	// funcID_conv_ToInt
// 	// funcID_conv_ToInt64
// 	// funcID_conv_ToInt64s
// 	// funcID_conv_ToInts
// 	// funcID_conv_ToString
// 	// funcID_conv_ToStrings
// 	// funcID_conv_URL
// 	// funcID_crypto_Bcrypt
// 	// funcID_crypto_PBKDF2
// 	// funcID_crypto_RSADecrypt
// 	// funcID_crypto_RSADecryptBytes
// 	// funcID_crypto_RSADerivePublicKey
// 	// funcID_crypto_RSAEncrypt
// 	// funcID_crypto_RSAGenerateKey
// 	// funcID_crypto_SHA1
// 	// funcID_crypto_SHA224
// 	// funcID_crypto_SHA256
// 	// funcID_crypto_SHA384
// 	// funcID_crypto_SHA512
// 	// funcID_crypto_SHA512_224
// 	// funcID_crypto_SHA512_256
// 	// funcID_crypto_WPAPSK
// 	// funcID_date
// 	// funcID_dateInZone
// 	// funcID_dateModify
// 	// funcID_date_in_zone
// 	// funcID_date_modify
// 	// funcID_decryptAES
// 	// funcID_deepCopy
// 	// funcID_deepEqual
// 	// funcID_default
// 	// funcID_derivePassword
// 	// funcID_dict
// 	// funcID_dig
// 	// funcID_dir
// 	// funcID_div
// 	// funcID_divf
// 	// funcID_duration
// 	// funcID_durationRound
// 	// funcID_empty
// 	// funcID_encryptAES
// 	// funcID_env
// 	// funcID_expandenv
// 	// funcID_ext
// 	// funcID_fail
// 	// funcID_file_Exists
// 	// funcID_file_IsDir
// 	// funcID_file_Read
// 	// funcID_file_ReadDir
// 	// funcID_file_Stat
// 	// funcID_file_Walk
// 	// funcID_file_Write
// 	// funcID_first
// 	// funcID_flatten
// 	// funcID_float64
// 	// funcID_floor
// 	// funcID_fromJson
// 	// funcID_genCA
// 	// funcID_genCAWithKey
// 	// funcID_genPrivateKey
// 	// funcID_genSelfSignedCert
// 	// funcID_genSelfSignedCertWithKey
// 	// funcID_genSignedCert
// 	// funcID_genSignedCertWithKey
// 	// funcID_get
// 	// funcID_getHostByName
// 	// funcID_has
// 	// funcID_hasKey
// 	funcID_hasPrefix: stringsNS{}.HasPrefix,
// 	funcID_hasSuffix: stringsNS{}.HasSuffix,
// 	// funcID_htmlDate
// 	// funcID_htmlDateInZone
// 	// funcID_htpasswd
// 	funcID_indent: stringsNS{}.Indent,
// 	// funcID_initial:  initial,
// 	// funcID_int
// 	// funcID_int64
// 	// funcID_isAbs
// 	// funcID_join
// 	// funcID_jq:        textNS{}.JQ,
// 	// funcID_jqObj:     textNS{}.JQObj,
// 	funcID_kebabcase: stringsNS{}.KebabCase,
// 	// funcID_keys
// 	// funcID_kindIs
// 	// funcID_kindOf
// 	// funcID_last
// 	// funcID_list
// 	// funcID_lower: stringsNS{}.ToLower,
// 	// funcID_math_Abs
// 	// funcID_math_Add
// 	// funcID_math_Ceil
// 	// funcID_math_Div
// 	// funcID_math_Floor
// 	// funcID_math_IsFloat
// 	// funcID_math_IsInt
// 	// funcID_math_IsNum
// 	// funcID_math_Max
// 	// funcID_math_Min
// 	// funcID_math_Mul
// 	// funcID_math_Pow
// 	// funcID_math_Rem
// 	// funcID_math_Round
// 	// funcID_math_Seq
// 	// funcID_math_Sub
// 	// funcID_max
// 	// funcID_maxf
// 	// funcID_md5sum
// 	// funcID_merge
// 	// funcID_mergeOverwrite
// 	// funcID_min
// 	// funcID_minf
// 	// funcID_mod
// 	// funcID_mul
// 	// funcID_mulf
// 	// funcID_mustAppend
// 	// funcID_mustChunk
// 	// funcID_mustCompact
// 	// funcID_mustDateModify
// 	// funcID_mustDeepCopy
// 	// funcID_mustFirst
// 	// funcID_mustFromJson
// 	// funcID_mustHas
// 	// funcID_mustInitial
// 	// funcID_mustLast
// 	// funcID_mustMerge
// 	// funcID_mustMergeOverwrite
// 	// funcID_mustPrepend
// 	// funcID_mustPush
// 	// funcID_mustRegexFind
// 	// funcID_mustRegexFindAll
// 	// funcID_mustRegexMatch
// 	// funcID_mustRegexReplaceAll
// 	// funcID_mustRegexReplaceAllLiteral
// 	// funcID_mustRegexSplit
// 	// funcID_mustRest
// 	// funcID_mustReverse
// 	// funcID_mustSlice
// 	// funcID_mustToDate
// 	// funcID_mustToJson
// 	// funcID_mustToPrettyJson
// 	// funcID_mustToRawJson
// 	// funcID_mustUniq
// 	// funcID_mustWithout
// 	// funcID_must_date_modify
// 	// funcID_net_LookupCNAME
// 	// funcID_net_LookupIP
// 	// funcID_net_LookupIPs
// 	// funcID_net_LookupSRV
// 	// funcID_net_LookupSRVs
// 	// funcID_net_LookupTXT
// 	// funcID_nindent: nindent,
// 	funcID_nospace: util.DeleteWhiteSpace,
// 	// funcID_now_Add
// 	// funcID_now_AddDate
// 	// funcID_now_After
// 	// funcID_now_AppendFormat
// 	// funcID_now_Before
// 	// funcID_now_Clock
// 	// funcID_now_Date
// 	// funcID_now_Day
// 	// funcID_now_Equal
// 	// funcID_now_Format
// 	// funcID_now_GoString
// 	// funcID_now_GobEncode
// 	// funcID_now_Hour
// 	// funcID_now_ISOWeek
// 	// funcID_now_In
// 	// funcID_now_IsDST
// 	// funcID_now_IsZero
// 	// funcID_now_Local
// 	// funcID_now_Location
// 	// funcID_now_MarshalBinary
// 	// funcID_now_MarshalJSON
// 	// funcID_now_MarshalText
// 	// funcID_now_Minute
// 	// funcID_now_Month
// 	// funcID_now_Nanosecond
// 	// funcID_now_Round
// 	// funcID_now_Second
// 	// funcID_now_String
// 	// funcID_now_Sub
// 	// funcID_now_Truncate
// 	// funcID_now_UTC
// 	// funcID_now_Unix
// 	// funcID_now_UnixMicro
// 	// funcID_now_UnixMilli
// 	// funcID_now_UnixNano
// 	// funcID_now_Weekday
// 	// funcID_now_Year
// 	// funcID_now_YearDay
// 	// funcID_now_Zone
// 	// funcID_omit
// 	// funcID_osBase
// 	// funcID_osClean
// 	// funcID_osDir
// 	// funcID_osExt
// 	// funcID_osIsAbs
// 	// funcID_path_Base
// 	// funcID_path_Clean
// 	// funcID_path_Dir
// 	// funcID_path_Ext
// 	// funcID_path_IsAbs
// 	// funcID_path_Join
// 	// funcID_path_Match
// 	// funcID_path_Split
// 	// funcID_pick
// 	// funcID_pluck
// 	funcID_plural: plural,
// 	// funcID_pow
// 	// funcID_prepend
// 	// funcID_push
// 	funcID_quote: stringsNS{}.DoubleQuote,
// 	// funcID_randAlpha
// 	// funcID_randAlphaNum
// 	// funcID_randAscii
// 	// funcID_randBytes
// 	// funcID_randInt
// 	// funcID_randNumeric
// 	// funcID_random_ASCII
// 	// funcID_random_Alpha
// 	// funcID_random_AlphaNum
// 	// funcID_random_Float
// 	// funcID_random_Item
// 	// funcID_random_Number
// 	// funcID_random_String
// 	// funcID_regexFind
// 	// funcID_regexFindAll
// 	// funcID_regexMatch
// 	// funcID_regexQuoteMeta
// 	// funcID_regexReplaceAll
// 	// funcID_regexReplaceAllLiteral
// 	// funcID_regexSplit
// 	// funcID_regexp_Find
// 	// funcID_regexp_FindAll
// 	// funcID_regexp_Match
// 	// funcID_regexp_QuoteMeta
// 	// funcID_regexp_Replace
// 	// funcID_regexp_ReplaceLiteral
// 	// funcID_regexp_Split
// 	// funcID_rem
// 	// funcID_removePrefix
// 	// funcID_removeSuffix
// 	funcID_repeat: stringsNS{}.Repeat,
// 	// funcID_replace:    replace,
// 	funcID_replaceAll: stringsNS{}.ReplaceAll,
// 	// funcID_rest
// 	// funcID_reverse
// 	// funcID_round
// 	// funcID_semver
// 	// funcID_semverCompare
// 	// funcID_seq
// 	// funcID_set
// 	// funcID_sha1sum
// 	// funcID_sha256sum
// 	// funcID_shellQuote
// 	// funcID_shuffle
// 	// funcID_slice
// 	funcID_snakecase: stringsNS{}.SnakeCase,
// 	// funcID_sockaddr_Attr
// 	// funcID_sockaddr_Exclude
// 	// funcID_sockaddr_GetAllInterfaces
// 	// funcID_sockaddr_GetDefaultInterfaces
// 	// funcID_sockaddr_GetInterfaceIP
// 	// funcID_sockaddr_GetInterfaceIPs
// 	// funcID_sockaddr_GetPrivateIP
// 	// funcID_sockaddr_GetPrivateIPs
// 	// funcID_sockaddr_GetPrivateInterfaces
// 	// funcID_sockaddr_GetPublicIP
// 	// funcID_sockaddr_GetPublicIPs
// 	// funcID_sockaddr_GetPublicInterfaces
// 	// funcID_sockaddr_Include
// 	// funcID_sockaddr_Join
// 	// funcID_sockaddr_Limit
// 	// funcID_sockaddr_Math
// 	// funcID_sockaddr_Offset
// 	// funcID_sockaddr_Sort
// 	// funcID_sockaddr_Unique
// 	// funcID_sort
// 	// funcID_sortAlpha
// 	funcID_split: stringsNS{}.Split,
// 	// funcID_splitList
// 	funcID_splitN: stringsNS{}.SplitN,
// 	// funcID_splitn
// 	funcID_squote: stringsNS{}.SingleQuote,
// 	// funcID_strconv_Unquote
// 	funcID_strings_Abbrev:    stringsNS{}.Abbrev,
// 	funcID_strings_AddPrefix: stringsNS{}.AddPrefix,
// 	funcID_strings_AddSuffix: stringsNS{}.AddSuffix,
// 	funcID_strings_CamelCase: stringsNS{}.CamelCase,
// 	funcID_strings_Contains:  stringsNS{}.Contains,
// 	funcID_strings_HasPrefix: stringsNS{}.HasPrefix,
// 	funcID_strings_HasSuffix: stringsNS{}.HasSuffix,
// 	funcID_strings_Indent:    stringsNS{}.Indent,
// 	// funcID_strings_JQ:           textNS{}.JQ,
// 	// funcID_strings_JQObj:        textNS{}.JQObj,
// 	funcID_strings_KebabCase:    stringsNS{}.KebabCase,
// 	funcID_strings_Quote:        stringsNS{}.DoubleQuote,
// 	funcID_strings_RemovePrefix: stringsNS{}.RemovePrefix,
// 	funcID_strings_RemoveSuffix: stringsNS{}.RemoveSuffix,
// 	funcID_strings_Repeat:       stringsNS{}.Repeat,
// 	funcID_strings_ReplaceAll:   stringsNS{}.ReplaceAll,
// 	funcID_strings_RuneCount:    stringsNS{}.RuneCount,
// 	funcID_strings_ShellQuote:   stringsNS{}.ShellQuote,
// 	funcID_strings_Slug:         stringsNS{}.Slug,
// 	funcID_strings_SnakeCase:    stringsNS{}.SnakeCase,
// 	funcID_strings_Split:        stringsNS{}.Split,
// 	funcID_strings_SplitN:       stringsNS{}.SplitN,
// 	funcID_strings_Squote:       stringsNS{}.SingleQuote,
// 	funcID_strings_Title:        stringsNS{}.Title,
// 	// funcID_strings_ToLower:      stringsNS{}.ToLower,
// 	// funcID_strings_ToUpper:      stringsNS{}.ToUpper,
// 	funcID_strings_Trim:       stringsNS{}.Trim,
// 	funcID_strings_TrimPrefix: stringsNS{}.TrimPrefix,
// 	funcID_strings_TrimSpace:  stringsNS{}.TrimSpace,
// 	funcID_strings_TrimSuffix: stringsNS{}.TrimSuffix,
// 	// funcID_strings_Trunc:        stringsNS{}.Trunc,
// 	funcID_strings_WordWrap: stringsNS{}.WordWrap,
// 	// funcID_strings_YQ:           textNS{}.YQ,
// 	// funcID_sub
// 	// funcID_subf
// 	// funcID_substr:   substring,
// 	funcID_swapcase: util.SwapCase,
// 	// funcID_ternary
// 	funcID_time_Hour:        timeNS{}.Hour,
// 	funcID_time_Microsecond: timeNS{}.Microsecond,
// 	funcID_time_Millisecond: timeNS{}.Millisecond,
// 	funcID_time_Minute:      timeNS{}.Minute,
// 	funcID_time_Nanosecond:  timeNS{}.Nanosecond,
// 	funcID_time_Now:         timeNS{}.Now,
// 	funcID_time_Parse:       timeNS{}.Parse,
// 	// funcID_time_ParseDuration
// 	// funcID_time_ParseInLocation
// 	// funcID_time_ParseLocal
// 	// funcID_time_Second
// 	// funcID_time_Since
// 	// funcID_time_Unix
// 	// funcID_time_Until
// 	// funcID_time_ZoneName
// 	// funcID_time_ZoneOffset
// 	funcID_title: stringsNS{}.Title,
// 	// funcID_toBytes
// 	// funcID_toDate
// 	// funcID_toDecimal
// 	// funcID_toJson
// 	// funcID_toLower
// 	// funcID_toPrettyJson
// 	// funcID_toRawJson
// 	// funcID_toString
// 	// funcID_toStrings: strslice,
// 	// funcID_toUpper: stringsNS{}.ToUpper,
// 	// funcID_toYaml
// 	funcID_totp:       credentialNS{}.Totp,
// 	funcID_trim:       stringsNS{}.Trim,
// 	funcID_trimAll:    stringsNS{}.Trim,
// 	funcID_trimPrefix: stringsNS{}.TrimPrefix,
// 	funcID_trimSpace:  stringsNS{}.TrimSpace,
// 	funcID_trimSuffix: stringsNS{}.TrimSuffix,
// 	// funcID_trunc:      stringsNS{}.Trunc,
// 	// funcID_tuple
// 	// funcID_typeIs
// 	// funcID_typeIsLike
// 	// funcID_typeOf
// 	// funcID_uniq
// 	// funcID_unixEpoch
// 	// funcID_unset
// 	// funcID_until
// 	// funcID_untilStep
// 	// funcID_upper: stringsNS{}.ToUpper,
// 	// funcID_urlJoin
// 	// funcID_urlParse
// 	// funcID_uuid_IsValid
// 	// funcID_uuid_Nil
// 	// funcID_uuid_Parse
// 	// funcID_uuid_V1
// 	// funcID_uuid_V4
// 	// funcID_without
// 	// funcID_wrap
// 	// funcID_wrapWith
// 	// funcID_yq: textNS{}.YQ,
// }

type ExecFuncs [funcID_COUNT]reflect.Value

var allExecFuncs ExecFuncs

type AllTemplateFuncs struct {
	staticFuncs *Funcs
}
