package templateutils

func funcNameToFuncID(name string) funcID {
	switch name {

	// start of static funcs
	case funcName_abbrev:
		return funcID_abbrev
	case funcName_abbrevboth:
		return funcID_abbrevboth
	case funcName_add:
		return funcID_add
	case funcName_add1:
		return funcID_add1
	case funcName_add1f:
		return funcID_add1f
	case funcName_addPrefix:
		return funcID_addPrefix
	case funcName_addSuffix:
		return funcID_addSuffix
	case funcName_addf:
		return funcID_addf
	case funcName_adler32sum:
		return funcID_adler32sum
	case funcName_ago:
		return funcID_ago
	case funcName_all:
		return funcID_all
	case funcName_any:
		return funcID_any
	case funcName_append:
		return funcID_append
	case funcName_archconv_AlpineArch:
		return funcID_archconv_AlpineArch
	case funcName_archconv_AlpineTripleName:
		return funcID_archconv_AlpineTripleName
	case funcName_archconv_DebianArch:
		return funcID_archconv_DebianArch
	case funcName_archconv_DebianTripleName:
		return funcID_archconv_DebianTripleName
	case funcName_archconv_DockerArch:
		return funcID_archconv_DockerArch
	case funcName_archconv_DockerArchVariant:
		return funcID_archconv_DockerArchVariant
	case funcName_archconv_DockerHubArch:
		return funcID_archconv_DockerHubArch
	case funcName_archconv_DockerOS:
		return funcID_archconv_DockerOS
	case funcName_archconv_DockerPlatformArch:
		return funcID_archconv_DockerPlatformArch
	case funcName_archconv_GNUArch:
		return funcID_archconv_GNUArch
	case funcName_archconv_GNUTripleName:
		return funcID_archconv_GNUTripleName
	case funcName_archconv_GolangArch:
		return funcID_archconv_GolangArch
	case funcName_archconv_GolangOS:
		return funcID_archconv_GolangOS
	case funcName_archconv_HF:
		return funcID_archconv_HF
	case funcName_archconv_HardFloadArch:
		return funcID_archconv_HardFloadArch
	case funcName_archconv_LLVMArch:
		return funcID_archconv_LLVMArch
	case funcName_archconv_LLVMTripleName:
		return funcID_archconv_LLVMTripleName
	case funcName_archconv_OciArch:
		return funcID_archconv_OciArch
	case funcName_archconv_OciArchVariant:
		return funcID_archconv_OciArchVariant
	case funcName_archconv_OciOS:
		return funcID_archconv_OciOS
	case funcName_archconv_QemuArch:
		return funcID_archconv_QemuArch
	case funcName_archconv_SF:
		return funcID_archconv_SF
	case funcName_archconv_SimpleArch:
		return funcID_archconv_SimpleArch
	case funcName_archconv_SoftFloadArch:
		return funcID_archconv_SoftFloadArch
	case funcName_atoi:
		return funcID_atoi
	case funcName_b32dec:
		return funcID_b32dec
	case funcName_b32enc:
		return funcID_b32enc
	case funcName_b64dec:
		return funcID_b64dec
	case funcName_b64enc:
		return funcID_b64enc
	case funcName_base:
		return funcID_base
	case funcName_bcrypt:
		return funcID_bcrypt
	case funcName_biggest:
		return funcID_biggest
	case funcName_bool:
		return funcID_bool
	case funcName_buildCustomCert:
		return funcID_buildCustomCert
	case funcName_camelcase:
		return funcID_camelcase
	case funcName_cat:
		return funcID_cat
	case funcName_ceil:
		return funcID_ceil
	case funcName_chunk:
		return funcID_chunk
	case funcName_clean:
		return funcID_clean
	case funcName_coalesce:
		return funcID_coalesce
	case funcName_coll_Append:
		return funcID_coll_Append
	case funcName_coll_Dict:
		return funcID_coll_Dict
	case funcName_coll_Flatten:
		return funcID_coll_Flatten
	case funcName_coll_Has:
		return funcID_coll_Has
	case funcName_coll_Keys:
		return funcID_coll_Keys
	case funcName_coll_Merge:
		return funcID_coll_Merge
	case funcName_coll_Omit:
		return funcID_coll_Omit
	case funcName_coll_Pick:
		return funcID_coll_Pick
	case funcName_coll_Prepend:
		return funcID_coll_Prepend
	case funcName_coll_Reverse:
		return funcID_coll_Reverse
	case funcName_coll_Slice:
		return funcID_coll_Slice
	case funcName_coll_Sort:
		return funcID_coll_Sort
	case funcName_coll_Uniq:
		return funcID_coll_Uniq
	case funcName_coll_Values:
		return funcID_coll_Values
	case funcName_compact:
		return funcID_compact
	case funcName_concat:
		return funcID_concat
	case funcName_contains:
		return funcID_contains
	case funcName_conv_Atoi:
		return funcID_conv_Atoi
	case funcName_conv_Bool:
		return funcID_conv_Bool
	case funcName_conv_Default:
		return funcID_conv_Default
	case funcName_conv_Dict:
		return funcID_conv_Dict
	case funcName_conv_Has:
		return funcID_conv_Has
	case funcName_conv_Join:
		return funcID_conv_Join
	case funcName_conv_ParseFloat:
		return funcID_conv_ParseFloat
	case funcName_conv_ParseInt:
		return funcID_conv_ParseInt
	case funcName_conv_ParseUint:
		return funcID_conv_ParseUint
	case funcName_conv_Slice:
		return funcID_conv_Slice
	case funcName_conv_ToBool:
		return funcID_conv_ToBool
	case funcName_conv_ToBools:
		return funcID_conv_ToBools
	case funcName_conv_ToFloat64:
		return funcID_conv_ToFloat64
	case funcName_conv_ToFloat64s:
		return funcID_conv_ToFloat64s
	case funcName_conv_ToInt:
		return funcID_conv_ToInt
	case funcName_conv_ToInt64:
		return funcID_conv_ToInt64
	case funcName_conv_ToInt64s:
		return funcID_conv_ToInt64s
	case funcName_conv_ToInts:
		return funcID_conv_ToInts
	case funcName_conv_ToString:
		return funcID_conv_ToString
	case funcName_conv_ToStrings:
		return funcID_conv_ToStrings
	case funcName_conv_URL:
		return funcID_conv_URL
	case funcName_crypto_Bcrypt:
		return funcID_crypto_Bcrypt
	case funcName_crypto_PBKDF2:
		return funcID_crypto_PBKDF2
	case funcName_crypto_RSADecrypt:
		return funcID_crypto_RSADecrypt
	case funcName_crypto_RSADecryptBytes:
		return funcID_crypto_RSADecryptBytes
	case funcName_crypto_RSADerivePublicKey:
		return funcID_crypto_RSADerivePublicKey
	case funcName_crypto_RSAEncrypt:
		return funcID_crypto_RSAEncrypt
	case funcName_crypto_RSAGenerateKey:
		return funcID_crypto_RSAGenerateKey
	case funcName_crypto_SHA1:
		return funcID_crypto_SHA1
	case funcName_crypto_SHA224:
		return funcID_crypto_SHA224
	case funcName_crypto_SHA256:
		return funcID_crypto_SHA256
	case funcName_crypto_SHA384:
		return funcID_crypto_SHA384
	case funcName_crypto_SHA512:
		return funcID_crypto_SHA512
	case funcName_crypto_SHA512_224:
		return funcID_crypto_SHA512_224
	case funcName_crypto_SHA512_256:
		return funcID_crypto_SHA512_256
	case funcName_crypto_WPAPSK:
		return funcID_crypto_WPAPSK
	case funcName_date:
		return funcID_date
	case funcName_dateInZone:
		return funcID_dateInZone
	case funcName_dateModify:
		return funcID_dateModify
	case funcName_date_in_zone:
		return funcID_date_in_zone
	case funcName_date_modify:
		return funcID_date_modify
	case funcName_decryptAES:
		return funcID_decryptAES
	case funcName_deepCopy:
		return funcID_deepCopy
	case funcName_deepEqual:
		return funcID_deepEqual
	case funcName_default:
		return funcID_default
	case funcName_derivePassword:
		return funcID_derivePassword
	case funcName_dict:
		return funcID_dict
	case funcName_dig:
		return funcID_dig
	case funcName_dir:
		return funcID_dir
	case funcName_div:
		return funcID_div
	case funcName_divf:
		return funcID_divf
	case funcName_duration:
		return funcID_duration
	case funcName_durationRound:
		return funcID_durationRound
	case funcName_empty:
		return funcID_empty
	case funcName_encryptAES:
		return funcID_encryptAES
	case funcName_env:
		return funcID_env
	case funcName_expandenv:
		return funcID_expandenv
	case funcName_ext:
		return funcID_ext
	case funcName_fail:
		return funcID_fail
	case funcName_file_Exists:
		return funcID_file_Exists
	case funcName_file_IsDir:
		return funcID_file_IsDir
	case funcName_file_Read:
		return funcID_file_Read
	case funcName_file_ReadDir:
		return funcID_file_ReadDir
	case funcName_file_Stat:
		return funcID_file_Stat
	case funcName_file_Walk:
		return funcID_file_Walk
	case funcName_file_Write:
		return funcID_file_Write
	case funcName_first:
		return funcID_first
	case funcName_flatten:
		return funcID_flatten
	case funcName_float64:
		return funcID_float64
	case funcName_floor:
		return funcID_floor
	case funcName_fromJson:
		return funcID_fromJson
	case funcName_genCA:
		return funcID_genCA
	case funcName_genCAWithKey:
		return funcID_genCAWithKey
	case funcName_genPrivateKey:
		return funcID_genPrivateKey
	case funcName_genSelfSignedCert:
		return funcID_genSelfSignedCert
	case funcName_genSelfSignedCertWithKey:
		return funcID_genSelfSignedCertWithKey
	case funcName_genSignedCert:
		return funcID_genSignedCert
	case funcName_genSignedCertWithKey:
		return funcID_genSignedCertWithKey
	case funcName_get:
		return funcID_get
	case funcName_getHostByName:
		return funcID_getHostByName
	case funcName_has:
		return funcID_has
	case funcName_hasKey:
		return funcID_hasKey
	case funcName_hasPrefix:
		return funcID_hasPrefix
	case funcName_hasSuffix:
		return funcID_hasSuffix
	case funcName_htmlDate:
		return funcID_htmlDate
	case funcName_htmlDateInZone:
		return funcID_htmlDateInZone
	case funcName_htpasswd:
		return funcID_htpasswd
	case funcName_indent:
		return funcID_indent
	case funcName_initial:
		return funcID_initial
	case funcName_initials:
		return funcID_initials
	case funcName_int:
		return funcID_int
	case funcName_int64:
		return funcID_int64
	case funcName_isAbs:
		return funcID_isAbs
	case funcName_join:
		return funcID_join
	case funcName_jq:
		return funcID_jq
	case funcName_jqObj:
		return funcID_jqObj
	case funcName_kebabcase:
		return funcID_kebabcase
	case funcName_keys:
		return funcID_keys
	case funcName_kindIs:
		return funcID_kindIs
	case funcName_kindOf:
		return funcID_kindOf
	case funcName_last:
		return funcID_last
	case funcName_list:
		return funcID_list
	case funcName_lower:
		return funcID_lower
	case funcName_math_Abs:
		return funcID_math_Abs
	case funcName_math_Add:
		return funcID_math_Add
	case funcName_math_Ceil:
		return funcID_math_Ceil
	case funcName_math_Div:
		return funcID_math_Div
	case funcName_math_Floor:
		return funcID_math_Floor
	case funcName_math_IsFloat:
		return funcID_math_IsFloat
	case funcName_math_IsInt:
		return funcID_math_IsInt
	case funcName_math_IsNum:
		return funcID_math_IsNum
	case funcName_math_Max:
		return funcID_math_Max
	case funcName_math_Min:
		return funcID_math_Min
	case funcName_math_Mul:
		return funcID_math_Mul
	case funcName_math_Pow:
		return funcID_math_Pow
	case funcName_math_Rem:
		return funcID_math_Rem
	case funcName_math_Round:
		return funcID_math_Round
	case funcName_math_Seq:
		return funcID_math_Seq
	case funcName_math_Sub:
		return funcID_math_Sub
	case funcName_max:
		return funcID_max
	case funcName_maxf:
		return funcID_maxf
	case funcName_md5sum:
		return funcID_md5sum
	case funcName_merge:
		return funcID_merge
	case funcName_mergeOverwrite:
		return funcID_mergeOverwrite
	case funcName_min:
		return funcID_min
	case funcName_minf:
		return funcID_minf
	case funcName_mod:
		return funcID_mod
	case funcName_mul:
		return funcID_mul
	case funcName_mulf:
		return funcID_mulf
	case funcName_mustAppend:
		return funcID_mustAppend
	case funcName_mustChunk:
		return funcID_mustChunk
	case funcName_mustCompact:
		return funcID_mustCompact
	case funcName_mustDateModify:
		return funcID_mustDateModify
	case funcName_mustDeepCopy:
		return funcID_mustDeepCopy
	case funcName_mustFirst:
		return funcID_mustFirst
	case funcName_mustFromJson:
		return funcID_mustFromJson
	case funcName_mustHas:
		return funcID_mustHas
	case funcName_mustInitial:
		return funcID_mustInitial
	case funcName_mustLast:
		return funcID_mustLast
	case funcName_mustMerge:
		return funcID_mustMerge
	case funcName_mustMergeOverwrite:
		return funcID_mustMergeOverwrite
	case funcName_mustPrepend:
		return funcID_mustPrepend
	case funcName_mustPush:
		return funcID_mustPush
	case funcName_mustRegexFind:
		return funcID_mustRegexFind
	case funcName_mustRegexFindAll:
		return funcID_mustRegexFindAll
	case funcName_mustRegexMatch:
		return funcID_mustRegexMatch
	case funcName_mustRegexReplaceAll:
		return funcID_mustRegexReplaceAll
	case funcName_mustRegexReplaceAllLiteral:
		return funcID_mustRegexReplaceAllLiteral
	case funcName_mustRegexSplit:
		return funcID_mustRegexSplit
	case funcName_mustRest:
		return funcID_mustRest
	case funcName_mustReverse:
		return funcID_mustReverse
	case funcName_mustSlice:
		return funcID_mustSlice
	case funcName_mustToDate:
		return funcID_mustToDate
	case funcName_mustToJson:
		return funcID_mustToJson
	case funcName_mustToPrettyJson:
		return funcID_mustToPrettyJson
	case funcName_mustToRawJson:
		return funcID_mustToRawJson
	case funcName_mustUniq:
		return funcID_mustUniq
	case funcName_mustWithout:
		return funcID_mustWithout
	case funcName_must_date_modify:
		return funcID_must_date_modify
	case funcName_net_LookupCNAME:
		return funcID_net_LookupCNAME
	case funcName_net_LookupIP:
		return funcID_net_LookupIP
	case funcName_net_LookupIPs:
		return funcID_net_LookupIPs
	case funcName_net_LookupSRV:
		return funcID_net_LookupSRV
	case funcName_net_LookupSRVs:
		return funcID_net_LookupSRVs
	case funcName_net_LookupTXT:
		return funcID_net_LookupTXT
	case funcName_nindent:
		return funcID_nindent
	case funcName_nospace:
		return funcID_nospace
	case funcName_now_Add:
		return funcID_now_Add
	case funcName_now_AddDate:
		return funcID_now_AddDate
	case funcName_now_After:
		return funcID_now_After
	case funcName_now_AppendFormat:
		return funcID_now_AppendFormat
	case funcName_now_Before:
		return funcID_now_Before
	case funcName_now_Clock:
		return funcID_now_Clock
	case funcName_now_Date:
		return funcID_now_Date
	case funcName_now_Day:
		return funcID_now_Day
	case funcName_now_Equal:
		return funcID_now_Equal
	case funcName_now_Format:
		return funcID_now_Format
	case funcName_now_GoString:
		return funcID_now_GoString
	case funcName_now_GobEncode:
		return funcID_now_GobEncode
	case funcName_now_Hour:
		return funcID_now_Hour
	case funcName_now_ISOWeek:
		return funcID_now_ISOWeek
	case funcName_now_In:
		return funcID_now_In
	case funcName_now_IsDST:
		return funcID_now_IsDST
	case funcName_now_IsZero:
		return funcID_now_IsZero
	case funcName_now_Local:
		return funcID_now_Local
	case funcName_now_Location:
		return funcID_now_Location
	case funcName_now_MarshalBinary:
		return funcID_now_MarshalBinary
	case funcName_now_MarshalJSON:
		return funcID_now_MarshalJSON
	case funcName_now_MarshalText:
		return funcID_now_MarshalText
	case funcName_now_Minute:
		return funcID_now_Minute
	case funcName_now_Month:
		return funcID_now_Month
	case funcName_now_Nanosecond:
		return funcID_now_Nanosecond
	case funcName_now_Round:
		return funcID_now_Round
	case funcName_now_Second:
		return funcID_now_Second
	case funcName_now_String:
		return funcID_now_String
	case funcName_now_Sub:
		return funcID_now_Sub
	case funcName_now_Truncate:
		return funcID_now_Truncate
	case funcName_now_UTC:
		return funcID_now_UTC
	case funcName_now_Unix:
		return funcID_now_Unix
	case funcName_now_UnixMicro:
		return funcID_now_UnixMicro
	case funcName_now_UnixMilli:
		return funcID_now_UnixMilli
	case funcName_now_UnixNano:
		return funcID_now_UnixNano
	case funcName_now_Weekday:
		return funcID_now_Weekday
	case funcName_now_Year:
		return funcID_now_Year
	case funcName_now_YearDay:
		return funcID_now_YearDay
	case funcName_now_Zone:
		return funcID_now_Zone
	case funcName_omit:
		return funcID_omit
	case funcName_osBase:
		return funcID_osBase
	case funcName_osClean:
		return funcID_osClean
	case funcName_osDir:
		return funcID_osDir
	case funcName_osExt:
		return funcID_osExt
	case funcName_osIsAbs:
		return funcID_osIsAbs
	case funcName_path_Base:
		return funcID_path_Base
	case funcName_path_Clean:
		return funcID_path_Clean
	case funcName_path_Dir:
		return funcID_path_Dir
	case funcName_path_Ext:
		return funcID_path_Ext
	case funcName_path_IsAbs:
		return funcID_path_IsAbs
	case funcName_path_Join:
		return funcID_path_Join
	case funcName_path_Match:
		return funcID_path_Match
	case funcName_path_Split:
		return funcID_path_Split
	case funcName_pick:
		return funcID_pick
	case funcName_pluck:
		return funcID_pluck
	case funcName_plural:
		return funcID_plural
	case funcName_pow:
		return funcID_pow
	case funcName_prepend:
		return funcID_prepend
	case funcName_push:
		return funcID_push
	case funcName_quote:
		return funcID_quote
	case funcName_randAlpha:
		return funcID_randAlpha
	case funcName_randAlphaNum:
		return funcID_randAlphaNum
	case funcName_randAscii:
		return funcID_randAscii
	case funcName_randBytes:
		return funcID_randBytes
	case funcName_randInt:
		return funcID_randInt
	case funcName_randNumeric:
		return funcID_randNumeric
	case funcName_random_ASCII:
		return funcID_random_ASCII
	case funcName_random_Alpha:
		return funcID_random_Alpha
	case funcName_random_AlphaNum:
		return funcID_random_AlphaNum
	case funcName_random_Float:
		return funcID_random_Float
	case funcName_random_Item:
		return funcID_random_Item
	case funcName_random_Number:
		return funcID_random_Number
	case funcName_random_String:
		return funcID_random_String
	case funcName_regexFind:
		return funcID_regexFind
	case funcName_regexFindAll:
		return funcID_regexFindAll
	case funcName_regexMatch:
		return funcID_regexMatch
	case funcName_regexQuoteMeta:
		return funcID_regexQuoteMeta
	case funcName_regexReplaceAll:
		return funcID_regexReplaceAll
	case funcName_regexReplaceAllLiteral:
		return funcID_regexReplaceAllLiteral
	case funcName_regexSplit:
		return funcID_regexSplit
	case funcName_regexp_Find:
		return funcID_regexp_Find
	case funcName_regexp_FindAll:
		return funcID_regexp_FindAll
	case funcName_regexp_Match:
		return funcID_regexp_Match
	case funcName_regexp_QuoteMeta:
		return funcID_regexp_QuoteMeta
	case funcName_regexp_Replace:
		return funcID_regexp_Replace
	case funcName_regexp_ReplaceLiteral:
		return funcID_regexp_ReplaceLiteral
	case funcName_regexp_Split:
		return funcID_regexp_Split
	case funcName_rem:
		return funcID_rem
	case funcName_removePrefix:
		return funcID_removePrefix
	case funcName_removeSuffix:
		return funcID_removeSuffix
	case funcName_repeat:
		return funcID_repeat
	case funcName_replace:
		return funcID_replace
	case funcName_replaceAll:
		return funcID_replaceAll
	case funcName_rest:
		return funcID_rest
	case funcName_reverse:
		return funcID_reverse
	case funcName_round:
		return funcID_round
	case funcName_semver:
		return funcID_semver
	case funcName_semverCompare:
		return funcID_semverCompare
	case funcName_seq:
		return funcID_seq
	case funcName_set:
		return funcID_set
	case funcName_sha1sum:
		return funcID_sha1sum
	case funcName_sha256sum:
		return funcID_sha256sum
	case funcName_shellQuote:
		return funcID_shellQuote
	case funcName_shuffle:
		return funcID_shuffle
	case funcName_slice:
		return funcID_slice
	case funcName_snakecase:
		return funcID_snakecase
	case funcName_sockaddr_Attr:
		return funcID_sockaddr_Attr
	case funcName_sockaddr_Exclude:
		return funcID_sockaddr_Exclude
	case funcName_sockaddr_GetAllInterfaces:
		return funcID_sockaddr_GetAllInterfaces
	case funcName_sockaddr_GetDefaultInterfaces:
		return funcID_sockaddr_GetDefaultInterfaces
	case funcName_sockaddr_GetInterfaceIP:
		return funcID_sockaddr_GetInterfaceIP
	case funcName_sockaddr_GetInterfaceIPs:
		return funcID_sockaddr_GetInterfaceIPs
	case funcName_sockaddr_GetPrivateIP:
		return funcID_sockaddr_GetPrivateIP
	case funcName_sockaddr_GetPrivateIPs:
		return funcID_sockaddr_GetPrivateIPs
	case funcName_sockaddr_GetPrivateInterfaces:
		return funcID_sockaddr_GetPrivateInterfaces
	case funcName_sockaddr_GetPublicIP:
		return funcID_sockaddr_GetPublicIP
	case funcName_sockaddr_GetPublicIPs:
		return funcID_sockaddr_GetPublicIPs
	case funcName_sockaddr_GetPublicInterfaces:
		return funcID_sockaddr_GetPublicInterfaces
	case funcName_sockaddr_Include:
		return funcID_sockaddr_Include
	case funcName_sockaddr_Join:
		return funcID_sockaddr_Join
	case funcName_sockaddr_Limit:
		return funcID_sockaddr_Limit
	case funcName_sockaddr_Math:
		return funcID_sockaddr_Math
	case funcName_sockaddr_Offset:
		return funcID_sockaddr_Offset
	case funcName_sockaddr_Sort:
		return funcID_sockaddr_Sort
	case funcName_sockaddr_Unique:
		return funcID_sockaddr_Unique
	case funcName_sort:
		return funcID_sort
	case funcName_sortAlpha:
		return funcID_sortAlpha
	case funcName_split:
		return funcID_split
	case funcName_splitList:
		return funcID_splitList
	case funcName_splitN:
		return funcID_splitN
	case funcName_splitn:
		return funcID_splitn
	case funcName_squote:
		return funcID_squote
	case funcName_strconv_Unquote:
		return funcID_strconv_Unquote
	case funcName_strings_Abbrev:
		return funcID_strings_Abbrev
	case funcName_strings_AddPrefix:
		return funcID_strings_AddPrefix
	case funcName_strings_AddSuffix:
		return funcID_strings_AddSuffix
	case funcName_strings_CamelCase:
		return funcID_strings_CamelCase
	case funcName_strings_Contains:
		return funcID_strings_Contains
	case funcName_strings_HasPrefix:
		return funcID_strings_HasPrefix
	case funcName_strings_HasSuffix:
		return funcID_strings_HasSuffix
	case funcName_strings_Indent:
		return funcID_strings_Indent
	case funcName_strings_JQ:
		return funcID_strings_JQ
	case funcName_strings_JQObj:
		return funcID_strings_JQObj
	case funcName_strings_KebabCase:
		return funcID_strings_KebabCase
	case funcName_strings_Quote:
		return funcID_strings_Quote
	case funcName_strings_RemovePrefix:
		return funcID_strings_RemovePrefix
	case funcName_strings_RemoveSuffix:
		return funcID_strings_RemoveSuffix
	case funcName_strings_Repeat:
		return funcID_strings_Repeat
	case funcName_strings_ReplaceAll:
		return funcID_strings_ReplaceAll
	case funcName_strings_RuneCount:
		return funcID_strings_RuneCount
	case funcName_strings_ShellQuote:
		return funcID_strings_ShellQuote
	case funcName_strings_Slug:
		return funcID_strings_Slug
	case funcName_strings_SnakeCase:
		return funcID_strings_SnakeCase
	case funcName_strings_Split:
		return funcID_strings_Split
	case funcName_strings_SplitN:
		return funcID_strings_SplitN
	case funcName_strings_Squote:
		return funcID_strings_Squote
	case funcName_strings_Title:
		return funcID_strings_Title
	case funcName_strings_ToLower:
		return funcID_strings_ToLower
	case funcName_strings_ToUpper:
		return funcID_strings_ToUpper
	case funcName_strings_Trim:
		return funcID_strings_Trim
	case funcName_strings_TrimPrefix:
		return funcID_strings_TrimPrefix
	case funcName_strings_TrimSpace:
		return funcID_strings_TrimSpace
	case funcName_strings_TrimSuffix:
		return funcID_strings_TrimSuffix
	case funcName_strings_Trunc:
		return funcID_strings_Trunc
	case funcName_strings_WordWrap:
		return funcID_strings_WordWrap
	case funcName_strings_YQ:
		return funcID_strings_YQ
	case funcName_sub:
		return funcID_sub
	case funcName_subf:
		return funcID_subf
	case funcName_substr:
		return funcID_substr
	case funcName_swapcase:
		return funcID_swapcase
	case funcName_ternary:
		return funcID_ternary
	case funcName_time_Add:
		return funcID_time_Add
	case funcName_time_Ceil:
		return funcID_time_Ceil
	case funcName_time_CeilDuration:
		return funcID_time_CeilDuration
	case funcName_time_Day:
		return funcID_time_Day
	case funcName_time_FMT_ANSI:
		return funcID_time_FMT_ANSI
	case funcName_time_FMT_RFC3339:
		return funcID_time_FMT_RFC3339
	case funcName_time_FMT_RFC3339Nano:
		return funcID_time_FMT_RFC3339Nano
	case funcName_time_FMT_Ruby:
		return funcID_time_FMT_Ruby
	case funcName_time_FMT_Stamp:
		return funcID_time_FMT_Stamp
	case funcName_time_FMT_Unix:
		return funcID_time_FMT_Unix
	case funcName_time_Floor:
		return funcID_time_Floor
	case funcName_time_FloorDuration:
		return funcID_time_FloorDuration
	case funcName_time_Format:
		return funcID_time_Format
	case funcName_time_Hour:
		return funcID_time_Hour
	case funcName_time_Microsecond:
		return funcID_time_Microsecond
	case funcName_time_Millisecond:
		return funcID_time_Millisecond
	case funcName_time_Minute:
		return funcID_time_Minute
	case funcName_time_Nanosecond:
		return funcID_time_Nanosecond
	case funcName_time_Now:
		return funcID_time_Now
	case funcName_time_Parse:
		return funcID_time_Parse
	case funcName_time_ParseDuration:
		return funcID_time_ParseDuration
	case funcName_time_Round:
		return funcID_time_Round
	case funcName_time_RoundDuration:
		return funcID_time_RoundDuration
	case funcName_time_Second:
		return funcID_time_Second
	case funcName_time_Since:
		return funcID_time_Since
	case funcName_time_Unix:
		return funcID_time_Unix
	case funcName_time_Until:
		return funcID_time_Until
	case funcName_time_Week:
		return funcID_time_Week
	case funcName_time_ZoneName:
		return funcID_time_ZoneName
	case funcName_time_ZoneOffset:
		return funcID_time_ZoneOffset
	case funcName_title:
		return funcID_title
	case funcName_toBytes:
		return funcID_toBytes
	case funcName_toDate:
		return funcID_toDate
	case funcName_toDecimal:
		return funcID_toDecimal
	case funcName_toJson:
		return funcID_toJson
	case funcName_toLower:
		return funcID_toLower
	case funcName_toPrettyJson:
		return funcID_toPrettyJson
	case funcName_toRawJson:
		return funcID_toRawJson
	case funcName_toString:
		return funcID_toString
	case funcName_toStrings:
		return funcID_toStrings
	case funcName_toUpper:
		return funcID_toUpper
	case funcName_toYaml:
		return funcID_toYaml
	case funcName_totp:
		return funcID_totp
	case funcName_trim:
		return funcID_trim
	case funcName_trimAll:
		return funcID_trimAll
	case funcName_trimPrefix:
		return funcID_trimPrefix
	case funcName_trimSpace:
		return funcID_trimSpace
	case funcName_trimSuffix:
		return funcID_trimSuffix
	case funcName_trimall:
		return funcID_trimall
	case funcName_trunc:
		return funcID_trunc
	case funcName_tuple:
		return funcID_tuple
	case funcName_typeIs:
		return funcID_typeIs
	case funcName_typeIsLike:
		return funcID_typeIsLike
	case funcName_typeOf:
		return funcID_typeOf
	case funcName_uniq:
		return funcID_uniq
	case funcName_unixEpoch:
		return funcID_unixEpoch
	case funcName_unset:
		return funcID_unset
	case funcName_until:
		return funcID_until
	case funcName_untilStep:
		return funcID_untilStep
	case funcName_untitle:
		return funcID_untitle
	case funcName_upper:
		return funcID_upper
	case funcName_urlJoin:
		return funcID_urlJoin
	case funcName_urlParse:
		return funcID_urlParse
	case funcName_uuid_IsValid:
		return funcID_uuid_IsValid
	case funcName_uuid_Nil:
		return funcID_uuid_Nil
	case funcName_uuid_Parse:
		return funcID_uuid_Parse
	case funcName_uuid_V1:
		return funcID_uuid_V1
	case funcName_uuid_V4:
		return funcID_uuid_V4
	// case funcName_values:
	// 	return funcID_values
	case funcName_without:
		return funcID_without
	case funcName_wrap:
		return funcID_wrap
	case funcName_wrapWith:
		return funcID_wrapWith
	case funcName_yq:
		return funcID_yq

	// end of static funcs

	// start of contextual funcs
	case funcName_dukkha_CacheDir:
		return funcID_dukkha_CacheDir
	case funcName_dukkha_CrossPlatform:
		return funcID_dukkha_CrossPlatform
	case funcName_dukkha_Self:
		return funcID_dukkha_Self
	case funcName_dukkha_Set:
		return funcID_dukkha_Set
	case funcName_dukkha_SetValue:
		return funcID_dukkha_SetValue
	case funcName_dukkha_WorkDir:
		return funcID_dukkha_WorkDir
	// case funcName_env:
	// 	return funcID_env
	case funcName_filepath_Abs:
		return funcID_filepath_Abs
	case funcName_filepath_Base:
		return funcID_filepath_Base
	case funcName_filepath_Clean:
		return funcID_filepath_Clean
	case funcName_filepath_Dir:
		return funcID_filepath_Dir
	case funcName_filepath_Ext:
		return funcID_filepath_Ext
	case funcName_filepath_FromSlash:
		return funcID_filepath_FromSlash
	case funcName_filepath_Glob:
		return funcID_filepath_Glob
	case funcName_filepath_IsAbs:
		return funcID_filepath_IsAbs
	case funcName_filepath_Join:
		return funcID_filepath_Join
	case funcName_filepath_Match:
		return funcID_filepath_Match
	case funcName_filepath_Rel:
		return funcID_filepath_Rel
	case funcName_filepath_Split:
		return funcID_filepath_Split
	case funcName_filepath_ToSlash:
		return funcID_filepath_ToSlash
	case funcName_filepath_VolumeName:
		return funcID_filepath_VolumeName
	case funcName_fromYaml:
		return funcID_fromYaml
	case funcName_getDefaultImageTag:
		return funcID_getDefaultImageTag
	case funcName_getDefaultManifestTag:
		return funcID_getDefaultManifestTag
	case funcName_git:
		return funcID_git
	case funcName_host:
		return funcID_host
	case funcName_matrix:
		return funcID_matrix
	case funcName_os_AppendFile:
		return funcID_os_AppendFile
	case funcName_os_Lookup:
		return funcID_os_Lookup
	case funcName_os_LookupFile:
		return funcID_os_LookupFile
	case funcName_os_MkdirAll:
		return funcID_os_MkdirAll
	case funcName_os_ReadFile:
		return funcID_os_ReadFile
	case funcName_os_Stderr:
		return funcID_os_Stderr
	case funcName_os_Stdin:
		return funcID_os_Stdin
	case funcName_os_Stdout:
		return funcID_os_Stdout
	case funcName_os_UserCacheDir:
		return funcID_os_UserCacheDir
	case funcName_os_UserConfigDir:
		return funcID_os_UserConfigDir
	case funcName_os_UserHomeDir:
		return funcID_os_UserHomeDir
	case funcName_os_WriteFile:
		return funcID_os_WriteFile
	case funcName_setDefaultImageTag:
		return funcID_setDefaultImageTag
	case funcName_setDefaultManifestTag:
		return funcID_setDefaultManifestTag
	case funcName_state_Failed:
		return funcID_state_Failed
	case funcName_state_Succeeded:
		return funcID_state_Succeeded
	case funcName_values:
		return funcID_values

	// end of contextual funcs

	// start of placeholder funcs
	case funcName_include:
		return funcID_include
	case funcName_var:
		return funcID_var

	// end of placeholder funcs

	default:
		return _unknown_template_func
	}
}

func (id funcID) String() string {
	switch id {

	// start of static funcs
	case funcID_abbrev:
		return funcName_abbrev
	case funcID_abbrevboth:
		return funcName_abbrevboth
	case funcID_add:
		return funcName_add
	case funcID_add1:
		return funcName_add1
	case funcID_add1f:
		return funcName_add1f
	case funcID_addPrefix:
		return funcName_addPrefix
	case funcID_addSuffix:
		return funcName_addSuffix
	case funcID_addf:
		return funcName_addf
	case funcID_adler32sum:
		return funcName_adler32sum
	case funcID_ago:
		return funcName_ago
	case funcID_all:
		return funcName_all
	case funcID_any:
		return funcName_any
	case funcID_append:
		return funcName_append
	case funcID_archconv_AlpineArch:
		return funcName_archconv_AlpineArch
	case funcID_archconv_AlpineTripleName:
		return funcName_archconv_AlpineTripleName
	case funcID_archconv_DebianArch:
		return funcName_archconv_DebianArch
	case funcID_archconv_DebianTripleName:
		return funcName_archconv_DebianTripleName
	case funcID_archconv_DockerArch:
		return funcName_archconv_DockerArch
	case funcID_archconv_DockerArchVariant:
		return funcName_archconv_DockerArchVariant
	case funcID_archconv_DockerHubArch:
		return funcName_archconv_DockerHubArch
	case funcID_archconv_DockerOS:
		return funcName_archconv_DockerOS
	case funcID_archconv_DockerPlatformArch:
		return funcName_archconv_DockerPlatformArch
	case funcID_archconv_GNUArch:
		return funcName_archconv_GNUArch
	case funcID_archconv_GNUTripleName:
		return funcName_archconv_GNUTripleName
	case funcID_archconv_GolangArch:
		return funcName_archconv_GolangArch
	case funcID_archconv_GolangOS:
		return funcName_archconv_GolangOS
	case funcID_archconv_HF:
		return funcName_archconv_HF
	case funcID_archconv_HardFloadArch:
		return funcName_archconv_HardFloadArch
	case funcID_archconv_LLVMArch:
		return funcName_archconv_LLVMArch
	case funcID_archconv_LLVMTripleName:
		return funcName_archconv_LLVMTripleName
	case funcID_archconv_OciArch:
		return funcName_archconv_OciArch
	case funcID_archconv_OciArchVariant:
		return funcName_archconv_OciArchVariant
	case funcID_archconv_OciOS:
		return funcName_archconv_OciOS
	case funcID_archconv_QemuArch:
		return funcName_archconv_QemuArch
	case funcID_archconv_SF:
		return funcName_archconv_SF
	case funcID_archconv_SimpleArch:
		return funcName_archconv_SimpleArch
	case funcID_archconv_SoftFloadArch:
		return funcName_archconv_SoftFloadArch
	case funcID_atoi:
		return funcName_atoi
	case funcID_b32dec:
		return funcName_b32dec
	case funcID_b32enc:
		return funcName_b32enc
	case funcID_b64dec:
		return funcName_b64dec
	case funcID_b64enc:
		return funcName_b64enc
	case funcID_base:
		return funcName_base
	case funcID_bcrypt:
		return funcName_bcrypt
	case funcID_biggest:
		return funcName_biggest
	case funcID_bool:
		return funcName_bool
	case funcID_buildCustomCert:
		return funcName_buildCustomCert
	case funcID_camelcase:
		return funcName_camelcase
	case funcID_cat:
		return funcName_cat
	case funcID_ceil:
		return funcName_ceil
	case funcID_chunk:
		return funcName_chunk
	case funcID_clean:
		return funcName_clean
	case funcID_coalesce:
		return funcName_coalesce
	case funcID_coll_Append:
		return funcName_coll_Append
	case funcID_coll_Dict:
		return funcName_coll_Dict
	case funcID_coll_Flatten:
		return funcName_coll_Flatten
	case funcID_coll_Has:
		return funcName_coll_Has
	case funcID_coll_Keys:
		return funcName_coll_Keys
	case funcID_coll_Merge:
		return funcName_coll_Merge
	case funcID_coll_Omit:
		return funcName_coll_Omit
	case funcID_coll_Pick:
		return funcName_coll_Pick
	case funcID_coll_Prepend:
		return funcName_coll_Prepend
	case funcID_coll_Reverse:
		return funcName_coll_Reverse
	case funcID_coll_Slice:
		return funcName_coll_Slice
	case funcID_coll_Sort:
		return funcName_coll_Sort
	case funcID_coll_Uniq:
		return funcName_coll_Uniq
	case funcID_coll_Values:
		return funcName_coll_Values
	case funcID_compact:
		return funcName_compact
	case funcID_concat:
		return funcName_concat
	case funcID_contains:
		return funcName_contains
	case funcID_conv_Atoi:
		return funcName_conv_Atoi
	case funcID_conv_Bool:
		return funcName_conv_Bool
	case funcID_conv_Default:
		return funcName_conv_Default
	case funcID_conv_Dict:
		return funcName_conv_Dict
	case funcID_conv_Has:
		return funcName_conv_Has
	case funcID_conv_Join:
		return funcName_conv_Join
	case funcID_conv_ParseFloat:
		return funcName_conv_ParseFloat
	case funcID_conv_ParseInt:
		return funcName_conv_ParseInt
	case funcID_conv_ParseUint:
		return funcName_conv_ParseUint
	case funcID_conv_Slice:
		return funcName_conv_Slice
	case funcID_conv_ToBool:
		return funcName_conv_ToBool
	case funcID_conv_ToBools:
		return funcName_conv_ToBools
	case funcID_conv_ToFloat64:
		return funcName_conv_ToFloat64
	case funcID_conv_ToFloat64s:
		return funcName_conv_ToFloat64s
	case funcID_conv_ToInt:
		return funcName_conv_ToInt
	case funcID_conv_ToInt64:
		return funcName_conv_ToInt64
	case funcID_conv_ToInt64s:
		return funcName_conv_ToInt64s
	case funcID_conv_ToInts:
		return funcName_conv_ToInts
	case funcID_conv_ToString:
		return funcName_conv_ToString
	case funcID_conv_ToStrings:
		return funcName_conv_ToStrings
	case funcID_conv_URL:
		return funcName_conv_URL
	case funcID_crypto_Bcrypt:
		return funcName_crypto_Bcrypt
	case funcID_crypto_PBKDF2:
		return funcName_crypto_PBKDF2
	case funcID_crypto_RSADecrypt:
		return funcName_crypto_RSADecrypt
	case funcID_crypto_RSADecryptBytes:
		return funcName_crypto_RSADecryptBytes
	case funcID_crypto_RSADerivePublicKey:
		return funcName_crypto_RSADerivePublicKey
	case funcID_crypto_RSAEncrypt:
		return funcName_crypto_RSAEncrypt
	case funcID_crypto_RSAGenerateKey:
		return funcName_crypto_RSAGenerateKey
	case funcID_crypto_SHA1:
		return funcName_crypto_SHA1
	case funcID_crypto_SHA224:
		return funcName_crypto_SHA224
	case funcID_crypto_SHA256:
		return funcName_crypto_SHA256
	case funcID_crypto_SHA384:
		return funcName_crypto_SHA384
	case funcID_crypto_SHA512:
		return funcName_crypto_SHA512
	case funcID_crypto_SHA512_224:
		return funcName_crypto_SHA512_224
	case funcID_crypto_SHA512_256:
		return funcName_crypto_SHA512_256
	case funcID_crypto_WPAPSK:
		return funcName_crypto_WPAPSK
	case funcID_date:
		return funcName_date
	case funcID_dateInZone:
		return funcName_dateInZone
	case funcID_dateModify:
		return funcName_dateModify
	case funcID_date_in_zone:
		return funcName_date_in_zone
	case funcID_date_modify:
		return funcName_date_modify
	case funcID_decryptAES:
		return funcName_decryptAES
	case funcID_deepCopy:
		return funcName_deepCopy
	case funcID_deepEqual:
		return funcName_deepEqual
	case funcID_default:
		return funcName_default
	case funcID_derivePassword:
		return funcName_derivePassword
	case funcID_dict:
		return funcName_dict
	case funcID_dig:
		return funcName_dig
	case funcID_dir:
		return funcName_dir
	case funcID_div:
		return funcName_div
	case funcID_divf:
		return funcName_divf
	case funcID_duration:
		return funcName_duration
	case funcID_durationRound:
		return funcName_durationRound
	case funcID_empty:
		return funcName_empty
	case funcID_encryptAES:
		return funcName_encryptAES
	case funcID_env:
		return funcName_env
	case funcID_expandenv:
		return funcName_expandenv
	case funcID_ext:
		return funcName_ext
	case funcID_fail:
		return funcName_fail
	case funcID_file_Exists:
		return funcName_file_Exists
	case funcID_file_IsDir:
		return funcName_file_IsDir
	case funcID_file_Read:
		return funcName_file_Read
	case funcID_file_ReadDir:
		return funcName_file_ReadDir
	case funcID_file_Stat:
		return funcName_file_Stat
	case funcID_file_Walk:
		return funcName_file_Walk
	case funcID_file_Write:
		return funcName_file_Write
	case funcID_first:
		return funcName_first
	case funcID_flatten:
		return funcName_flatten
	case funcID_float64:
		return funcName_float64
	case funcID_floor:
		return funcName_floor
	case funcID_fromJson:
		return funcName_fromJson
	case funcID_genCA:
		return funcName_genCA
	case funcID_genCAWithKey:
		return funcName_genCAWithKey
	case funcID_genPrivateKey:
		return funcName_genPrivateKey
	case funcID_genSelfSignedCert:
		return funcName_genSelfSignedCert
	case funcID_genSelfSignedCertWithKey:
		return funcName_genSelfSignedCertWithKey
	case funcID_genSignedCert:
		return funcName_genSignedCert
	case funcID_genSignedCertWithKey:
		return funcName_genSignedCertWithKey
	case funcID_get:
		return funcName_get
	case funcID_getHostByName:
		return funcName_getHostByName
	case funcID_has:
		return funcName_has
	case funcID_hasKey:
		return funcName_hasKey
	case funcID_hasPrefix:
		return funcName_hasPrefix
	case funcID_hasSuffix:
		return funcName_hasSuffix
	case funcID_htmlDate:
		return funcName_htmlDate
	case funcID_htmlDateInZone:
		return funcName_htmlDateInZone
	case funcID_htpasswd:
		return funcName_htpasswd
	case funcID_indent:
		return funcName_indent
	case funcID_initial:
		return funcName_initial
	case funcID_initials:
		return funcName_initials
	case funcID_int:
		return funcName_int
	case funcID_int64:
		return funcName_int64
	case funcID_isAbs:
		return funcName_isAbs
	case funcID_join:
		return funcName_join
	case funcID_jq:
		return funcName_jq
	case funcID_jqObj:
		return funcName_jqObj
	case funcID_kebabcase:
		return funcName_kebabcase
	case funcID_keys:
		return funcName_keys
	case funcID_kindIs:
		return funcName_kindIs
	case funcID_kindOf:
		return funcName_kindOf
	case funcID_last:
		return funcName_last
	case funcID_list:
		return funcName_list
	case funcID_lower:
		return funcName_lower
	case funcID_math_Abs:
		return funcName_math_Abs
	case funcID_math_Add:
		return funcName_math_Add
	case funcID_math_Ceil:
		return funcName_math_Ceil
	case funcID_math_Div:
		return funcName_math_Div
	case funcID_math_Floor:
		return funcName_math_Floor
	case funcID_math_IsFloat:
		return funcName_math_IsFloat
	case funcID_math_IsInt:
		return funcName_math_IsInt
	case funcID_math_IsNum:
		return funcName_math_IsNum
	case funcID_math_Max:
		return funcName_math_Max
	case funcID_math_Min:
		return funcName_math_Min
	case funcID_math_Mul:
		return funcName_math_Mul
	case funcID_math_Pow:
		return funcName_math_Pow
	case funcID_math_Rem:
		return funcName_math_Rem
	case funcID_math_Round:
		return funcName_math_Round
	case funcID_math_Seq:
		return funcName_math_Seq
	case funcID_math_Sub:
		return funcName_math_Sub
	case funcID_max:
		return funcName_max
	case funcID_maxf:
		return funcName_maxf
	case funcID_md5sum:
		return funcName_md5sum
	case funcID_merge:
		return funcName_merge
	case funcID_mergeOverwrite:
		return funcName_mergeOverwrite
	case funcID_min:
		return funcName_min
	case funcID_minf:
		return funcName_minf
	case funcID_mod:
		return funcName_mod
	case funcID_mul:
		return funcName_mul
	case funcID_mulf:
		return funcName_mulf
	case funcID_mustAppend:
		return funcName_mustAppend
	case funcID_mustChunk:
		return funcName_mustChunk
	case funcID_mustCompact:
		return funcName_mustCompact
	case funcID_mustDateModify:
		return funcName_mustDateModify
	case funcID_mustDeepCopy:
		return funcName_mustDeepCopy
	case funcID_mustFirst:
		return funcName_mustFirst
	case funcID_mustFromJson:
		return funcName_mustFromJson
	case funcID_mustHas:
		return funcName_mustHas
	case funcID_mustInitial:
		return funcName_mustInitial
	case funcID_mustLast:
		return funcName_mustLast
	case funcID_mustMerge:
		return funcName_mustMerge
	case funcID_mustMergeOverwrite:
		return funcName_mustMergeOverwrite
	case funcID_mustPrepend:
		return funcName_mustPrepend
	case funcID_mustPush:
		return funcName_mustPush
	case funcID_mustRegexFind:
		return funcName_mustRegexFind
	case funcID_mustRegexFindAll:
		return funcName_mustRegexFindAll
	case funcID_mustRegexMatch:
		return funcName_mustRegexMatch
	case funcID_mustRegexReplaceAll:
		return funcName_mustRegexReplaceAll
	case funcID_mustRegexReplaceAllLiteral:
		return funcName_mustRegexReplaceAllLiteral
	case funcID_mustRegexSplit:
		return funcName_mustRegexSplit
	case funcID_mustRest:
		return funcName_mustRest
	case funcID_mustReverse:
		return funcName_mustReverse
	case funcID_mustSlice:
		return funcName_mustSlice
	case funcID_mustToDate:
		return funcName_mustToDate
	case funcID_mustToJson:
		return funcName_mustToJson
	case funcID_mustToPrettyJson:
		return funcName_mustToPrettyJson
	case funcID_mustToRawJson:
		return funcName_mustToRawJson
	case funcID_mustUniq:
		return funcName_mustUniq
	case funcID_mustWithout:
		return funcName_mustWithout
	case funcID_must_date_modify:
		return funcName_must_date_modify
	case funcID_net_LookupCNAME:
		return funcName_net_LookupCNAME
	case funcID_net_LookupIP:
		return funcName_net_LookupIP
	case funcID_net_LookupIPs:
		return funcName_net_LookupIPs
	case funcID_net_LookupSRV:
		return funcName_net_LookupSRV
	case funcID_net_LookupSRVs:
		return funcName_net_LookupSRVs
	case funcID_net_LookupTXT:
		return funcName_net_LookupTXT
	case funcID_nindent:
		return funcName_nindent
	case funcID_nospace:
		return funcName_nospace
	case funcID_now_Add:
		return funcName_now_Add
	case funcID_now_AddDate:
		return funcName_now_AddDate
	case funcID_now_After:
		return funcName_now_After
	case funcID_now_AppendFormat:
		return funcName_now_AppendFormat
	case funcID_now_Before:
		return funcName_now_Before
	case funcID_now_Clock:
		return funcName_now_Clock
	case funcID_now_Date:
		return funcName_now_Date
	case funcID_now_Day:
		return funcName_now_Day
	case funcID_now_Equal:
		return funcName_now_Equal
	case funcID_now_Format:
		return funcName_now_Format
	case funcID_now_GoString:
		return funcName_now_GoString
	case funcID_now_GobEncode:
		return funcName_now_GobEncode
	case funcID_now_Hour:
		return funcName_now_Hour
	case funcID_now_ISOWeek:
		return funcName_now_ISOWeek
	case funcID_now_In:
		return funcName_now_In
	case funcID_now_IsDST:
		return funcName_now_IsDST
	case funcID_now_IsZero:
		return funcName_now_IsZero
	case funcID_now_Local:
		return funcName_now_Local
	case funcID_now_Location:
		return funcName_now_Location
	case funcID_now_MarshalBinary:
		return funcName_now_MarshalBinary
	case funcID_now_MarshalJSON:
		return funcName_now_MarshalJSON
	case funcID_now_MarshalText:
		return funcName_now_MarshalText
	case funcID_now_Minute:
		return funcName_now_Minute
	case funcID_now_Month:
		return funcName_now_Month
	case funcID_now_Nanosecond:
		return funcName_now_Nanosecond
	case funcID_now_Round:
		return funcName_now_Round
	case funcID_now_Second:
		return funcName_now_Second
	case funcID_now_String:
		return funcName_now_String
	case funcID_now_Sub:
		return funcName_now_Sub
	case funcID_now_Truncate:
		return funcName_now_Truncate
	case funcID_now_UTC:
		return funcName_now_UTC
	case funcID_now_Unix:
		return funcName_now_Unix
	case funcID_now_UnixMicro:
		return funcName_now_UnixMicro
	case funcID_now_UnixMilli:
		return funcName_now_UnixMilli
	case funcID_now_UnixNano:
		return funcName_now_UnixNano
	case funcID_now_Weekday:
		return funcName_now_Weekday
	case funcID_now_Year:
		return funcName_now_Year
	case funcID_now_YearDay:
		return funcName_now_YearDay
	case funcID_now_Zone:
		return funcName_now_Zone
	case funcID_omit:
		return funcName_omit
	case funcID_osBase:
		return funcName_osBase
	case funcID_osClean:
		return funcName_osClean
	case funcID_osDir:
		return funcName_osDir
	case funcID_osExt:
		return funcName_osExt
	case funcID_osIsAbs:
		return funcName_osIsAbs
	case funcID_path_Base:
		return funcName_path_Base
	case funcID_path_Clean:
		return funcName_path_Clean
	case funcID_path_Dir:
		return funcName_path_Dir
	case funcID_path_Ext:
		return funcName_path_Ext
	case funcID_path_IsAbs:
		return funcName_path_IsAbs
	case funcID_path_Join:
		return funcName_path_Join
	case funcID_path_Match:
		return funcName_path_Match
	case funcID_path_Split:
		return funcName_path_Split
	case funcID_pick:
		return funcName_pick
	case funcID_pluck:
		return funcName_pluck
	case funcID_plural:
		return funcName_plural
	case funcID_pow:
		return funcName_pow
	case funcID_prepend:
		return funcName_prepend
	case funcID_push:
		return funcName_push
	case funcID_quote:
		return funcName_quote
	case funcID_randAlpha:
		return funcName_randAlpha
	case funcID_randAlphaNum:
		return funcName_randAlphaNum
	case funcID_randAscii:
		return funcName_randAscii
	case funcID_randBytes:
		return funcName_randBytes
	case funcID_randInt:
		return funcName_randInt
	case funcID_randNumeric:
		return funcName_randNumeric
	case funcID_random_ASCII:
		return funcName_random_ASCII
	case funcID_random_Alpha:
		return funcName_random_Alpha
	case funcID_random_AlphaNum:
		return funcName_random_AlphaNum
	case funcID_random_Float:
		return funcName_random_Float
	case funcID_random_Item:
		return funcName_random_Item
	case funcID_random_Number:
		return funcName_random_Number
	case funcID_random_String:
		return funcName_random_String
	case funcID_regexFind:
		return funcName_regexFind
	case funcID_regexFindAll:
		return funcName_regexFindAll
	case funcID_regexMatch:
		return funcName_regexMatch
	case funcID_regexQuoteMeta:
		return funcName_regexQuoteMeta
	case funcID_regexReplaceAll:
		return funcName_regexReplaceAll
	case funcID_regexReplaceAllLiteral:
		return funcName_regexReplaceAllLiteral
	case funcID_regexSplit:
		return funcName_regexSplit
	case funcID_regexp_Find:
		return funcName_regexp_Find
	case funcID_regexp_FindAll:
		return funcName_regexp_FindAll
	case funcID_regexp_Match:
		return funcName_regexp_Match
	case funcID_regexp_QuoteMeta:
		return funcName_regexp_QuoteMeta
	case funcID_regexp_Replace:
		return funcName_regexp_Replace
	case funcID_regexp_ReplaceLiteral:
		return funcName_regexp_ReplaceLiteral
	case funcID_regexp_Split:
		return funcName_regexp_Split
	case funcID_rem:
		return funcName_rem
	case funcID_removePrefix:
		return funcName_removePrefix
	case funcID_removeSuffix:
		return funcName_removeSuffix
	case funcID_repeat:
		return funcName_repeat
	case funcID_replace:
		return funcName_replace
	case funcID_replaceAll:
		return funcName_replaceAll
	case funcID_rest:
		return funcName_rest
	case funcID_reverse:
		return funcName_reverse
	case funcID_round:
		return funcName_round
	case funcID_semver:
		return funcName_semver
	case funcID_semverCompare:
		return funcName_semverCompare
	case funcID_seq:
		return funcName_seq
	case funcID_set:
		return funcName_set
	case funcID_sha1sum:
		return funcName_sha1sum
	case funcID_sha256sum:
		return funcName_sha256sum
	case funcID_shellQuote:
		return funcName_shellQuote
	case funcID_shuffle:
		return funcName_shuffle
	case funcID_slice:
		return funcName_slice
	case funcID_snakecase:
		return funcName_snakecase
	case funcID_sockaddr_Attr:
		return funcName_sockaddr_Attr
	case funcID_sockaddr_Exclude:
		return funcName_sockaddr_Exclude
	case funcID_sockaddr_GetAllInterfaces:
		return funcName_sockaddr_GetAllInterfaces
	case funcID_sockaddr_GetDefaultInterfaces:
		return funcName_sockaddr_GetDefaultInterfaces
	case funcID_sockaddr_GetInterfaceIP:
		return funcName_sockaddr_GetInterfaceIP
	case funcID_sockaddr_GetInterfaceIPs:
		return funcName_sockaddr_GetInterfaceIPs
	case funcID_sockaddr_GetPrivateIP:
		return funcName_sockaddr_GetPrivateIP
	case funcID_sockaddr_GetPrivateIPs:
		return funcName_sockaddr_GetPrivateIPs
	case funcID_sockaddr_GetPrivateInterfaces:
		return funcName_sockaddr_GetPrivateInterfaces
	case funcID_sockaddr_GetPublicIP:
		return funcName_sockaddr_GetPublicIP
	case funcID_sockaddr_GetPublicIPs:
		return funcName_sockaddr_GetPublicIPs
	case funcID_sockaddr_GetPublicInterfaces:
		return funcName_sockaddr_GetPublicInterfaces
	case funcID_sockaddr_Include:
		return funcName_sockaddr_Include
	case funcID_sockaddr_Join:
		return funcName_sockaddr_Join
	case funcID_sockaddr_Limit:
		return funcName_sockaddr_Limit
	case funcID_sockaddr_Math:
		return funcName_sockaddr_Math
	case funcID_sockaddr_Offset:
		return funcName_sockaddr_Offset
	case funcID_sockaddr_Sort:
		return funcName_sockaddr_Sort
	case funcID_sockaddr_Unique:
		return funcName_sockaddr_Unique
	case funcID_sort:
		return funcName_sort
	case funcID_sortAlpha:
		return funcName_sortAlpha
	case funcID_split:
		return funcName_split
	case funcID_splitList:
		return funcName_splitList
	case funcID_splitN:
		return funcName_splitN
	case funcID_splitn:
		return funcName_splitn
	case funcID_squote:
		return funcName_squote
	case funcID_strconv_Unquote:
		return funcName_strconv_Unquote
	case funcID_strings_Abbrev:
		return funcName_strings_Abbrev
	case funcID_strings_AddPrefix:
		return funcName_strings_AddPrefix
	case funcID_strings_AddSuffix:
		return funcName_strings_AddSuffix
	case funcID_strings_CamelCase:
		return funcName_strings_CamelCase
	case funcID_strings_Contains:
		return funcName_strings_Contains
	case funcID_strings_HasPrefix:
		return funcName_strings_HasPrefix
	case funcID_strings_HasSuffix:
		return funcName_strings_HasSuffix
	case funcID_strings_Indent:
		return funcName_strings_Indent
	case funcID_strings_JQ:
		return funcName_strings_JQ
	case funcID_strings_JQObj:
		return funcName_strings_JQObj
	case funcID_strings_KebabCase:
		return funcName_strings_KebabCase
	case funcID_strings_Quote:
		return funcName_strings_Quote
	case funcID_strings_RemovePrefix:
		return funcName_strings_RemovePrefix
	case funcID_strings_RemoveSuffix:
		return funcName_strings_RemoveSuffix
	case funcID_strings_Repeat:
		return funcName_strings_Repeat
	case funcID_strings_ReplaceAll:
		return funcName_strings_ReplaceAll
	case funcID_strings_RuneCount:
		return funcName_strings_RuneCount
	case funcID_strings_ShellQuote:
		return funcName_strings_ShellQuote
	case funcID_strings_Slug:
		return funcName_strings_Slug
	case funcID_strings_SnakeCase:
		return funcName_strings_SnakeCase
	case funcID_strings_Split:
		return funcName_strings_Split
	case funcID_strings_SplitN:
		return funcName_strings_SplitN
	case funcID_strings_Squote:
		return funcName_strings_Squote
	case funcID_strings_Title:
		return funcName_strings_Title
	case funcID_strings_ToLower:
		return funcName_strings_ToLower
	case funcID_strings_ToUpper:
		return funcName_strings_ToUpper
	case funcID_strings_Trim:
		return funcName_strings_Trim
	case funcID_strings_TrimPrefix:
		return funcName_strings_TrimPrefix
	case funcID_strings_TrimSpace:
		return funcName_strings_TrimSpace
	case funcID_strings_TrimSuffix:
		return funcName_strings_TrimSuffix
	case funcID_strings_Trunc:
		return funcName_strings_Trunc
	case funcID_strings_WordWrap:
		return funcName_strings_WordWrap
	case funcID_strings_YQ:
		return funcName_strings_YQ
	case funcID_sub:
		return funcName_sub
	case funcID_subf:
		return funcName_subf
	case funcID_substr:
		return funcName_substr
	case funcID_swapcase:
		return funcName_swapcase
	case funcID_ternary:
		return funcName_ternary
	case funcID_time_Add:
		return funcName_time_Add
	case funcID_time_Ceil:
		return funcName_time_Ceil
	case funcID_time_CeilDuration:
		return funcName_time_CeilDuration
	case funcID_time_Day:
		return funcName_time_Day
	case funcID_time_FMT_ANSI:
		return funcName_time_FMT_ANSI
	case funcID_time_FMT_RFC3339:
		return funcName_time_FMT_RFC3339
	case funcID_time_FMT_RFC3339Nano:
		return funcName_time_FMT_RFC3339Nano
	case funcID_time_FMT_Ruby:
		return funcName_time_FMT_Ruby
	case funcID_time_FMT_Stamp:
		return funcName_time_FMT_Stamp
	case funcID_time_FMT_Unix:
		return funcName_time_FMT_Unix
	case funcID_time_Floor:
		return funcName_time_Floor
	case funcID_time_FloorDuration:
		return funcName_time_FloorDuration
	case funcID_time_Format:
		return funcName_time_Format
	case funcID_time_Hour:
		return funcName_time_Hour
	case funcID_time_Microsecond:
		return funcName_time_Microsecond
	case funcID_time_Millisecond:
		return funcName_time_Millisecond
	case funcID_time_Minute:
		return funcName_time_Minute
	case funcID_time_Nanosecond:
		return funcName_time_Nanosecond
	case funcID_time_Now:
		return funcName_time_Now
	case funcID_time_Parse:
		return funcName_time_Parse
	case funcID_time_ParseDuration:
		return funcName_time_ParseDuration
	case funcID_time_Round:
		return funcName_time_Round
	case funcID_time_RoundDuration:
		return funcName_time_RoundDuration
	case funcID_time_Second:
		return funcName_time_Second
	case funcID_time_Since:
		return funcName_time_Since
	case funcID_time_Unix:
		return funcName_time_Unix
	case funcID_time_Until:
		return funcName_time_Until
	case funcID_time_Week:
		return funcName_time_Week
	case funcID_time_ZoneName:
		return funcName_time_ZoneName
	case funcID_time_ZoneOffset:
		return funcName_time_ZoneOffset
	case funcID_title:
		return funcName_title
	case funcID_toBytes:
		return funcName_toBytes
	case funcID_toDate:
		return funcName_toDate
	case funcID_toDecimal:
		return funcName_toDecimal
	case funcID_toJson:
		return funcName_toJson
	case funcID_toLower:
		return funcName_toLower
	case funcID_toPrettyJson:
		return funcName_toPrettyJson
	case funcID_toRawJson:
		return funcName_toRawJson
	case funcID_toString:
		return funcName_toString
	case funcID_toStrings:
		return funcName_toStrings
	case funcID_toUpper:
		return funcName_toUpper
	case funcID_toYaml:
		return funcName_toYaml
	case funcID_totp:
		return funcName_totp
	case funcID_trim:
		return funcName_trim
	case funcID_trimAll:
		return funcName_trimAll
	case funcID_trimPrefix:
		return funcName_trimPrefix
	case funcID_trimSpace:
		return funcName_trimSpace
	case funcID_trimSuffix:
		return funcName_trimSuffix
	case funcID_trimall:
		return funcName_trimall
	case funcID_trunc:
		return funcName_trunc
	case funcID_tuple:
		return funcName_tuple
	case funcID_typeIs:
		return funcName_typeIs
	case funcID_typeIsLike:
		return funcName_typeIsLike
	case funcID_typeOf:
		return funcName_typeOf
	case funcID_uniq:
		return funcName_uniq
	case funcID_unixEpoch:
		return funcName_unixEpoch
	case funcID_unset:
		return funcName_unset
	case funcID_until:
		return funcName_until
	case funcID_untilStep:
		return funcName_untilStep
	case funcID_untitle:
		return funcName_untitle
	case funcID_upper:
		return funcName_upper
	case funcID_urlJoin:
		return funcName_urlJoin
	case funcID_urlParse:
		return funcName_urlParse
	case funcID_uuid_IsValid:
		return funcName_uuid_IsValid
	case funcID_uuid_Nil:
		return funcName_uuid_Nil
	case funcID_uuid_Parse:
		return funcName_uuid_Parse
	case funcID_uuid_V1:
		return funcName_uuid_V1
	case funcID_uuid_V4:
		return funcName_uuid_V4
	case funcID_values:
		return funcName_values
	case funcID_without:
		return funcName_without
	case funcID_wrap:
		return funcName_wrap
	case funcID_wrapWith:
		return funcName_wrapWith
	case funcID_yq:
		return funcName_yq

	// end of static funcs

	// start of contextual funcs
	case funcID_dukkha_CacheDir:
		return funcName_dukkha_CacheDir
	case funcID_dukkha_CrossPlatform:
		return funcName_dukkha_CrossPlatform
	case funcID_dukkha_Self:
		return funcName_dukkha_Self
	case funcID_dukkha_Set:
		return funcName_dukkha_Set
	case funcID_dukkha_SetValue:
		return funcName_dukkha_SetValue
	case funcID_dukkha_WorkDir:
		return funcName_dukkha_WorkDir
	// case funcID_env:
	// 	return funcName_env
	case funcID_filepath_Abs:
		return funcName_filepath_Abs
	case funcID_filepath_Base:
		return funcName_filepath_Base
	case funcID_filepath_Clean:
		return funcName_filepath_Clean
	case funcID_filepath_Dir:
		return funcName_filepath_Dir
	case funcID_filepath_Ext:
		return funcName_filepath_Ext
	case funcID_filepath_FromSlash:
		return funcName_filepath_FromSlash
	case funcID_filepath_Glob:
		return funcName_filepath_Glob
	case funcID_filepath_IsAbs:
		return funcName_filepath_IsAbs
	case funcID_filepath_Join:
		return funcName_filepath_Join
	case funcID_filepath_Match:
		return funcName_filepath_Match
	case funcID_filepath_Rel:
		return funcName_filepath_Rel
	case funcID_filepath_Split:
		return funcName_filepath_Split
	case funcID_filepath_ToSlash:
		return funcName_filepath_ToSlash
	case funcID_filepath_VolumeName:
		return funcName_filepath_VolumeName
	case funcID_fromYaml:
		return funcName_fromYaml
	case funcID_getDefaultImageTag:
		return funcName_getDefaultImageTag
	case funcID_getDefaultManifestTag:
		return funcName_getDefaultManifestTag
	case funcID_git:
		return funcName_git
	case funcID_host:
		return funcName_host
	case funcID_matrix:
		return funcName_matrix
	case funcID_os_AppendFile:
		return funcName_os_AppendFile
	case funcID_os_Lookup:
		return funcName_os_Lookup
	case funcID_os_LookupFile:
		return funcName_os_LookupFile
	case funcID_os_MkdirAll:
		return funcName_os_MkdirAll
	case funcID_os_ReadFile:
		return funcName_os_ReadFile
	case funcID_os_Stderr:
		return funcName_os_Stderr
	case funcID_os_Stdin:
		return funcName_os_Stdin
	case funcID_os_Stdout:
		return funcName_os_Stdout
	case funcID_os_UserCacheDir:
		return funcName_os_UserCacheDir
	case funcID_os_UserConfigDir:
		return funcName_os_UserConfigDir
	case funcID_os_UserHomeDir:
		return funcName_os_UserHomeDir
	case funcID_os_WriteFile:
		return funcName_os_WriteFile
	case funcID_setDefaultImageTag:
		return funcName_setDefaultImageTag
	case funcID_setDefaultManifestTag:
		return funcName_setDefaultManifestTag
	case funcID_state_Failed:
		return funcName_state_Failed
	case funcID_state_Succeeded:
		return funcName_state_Succeeded
	// case funcID_values:
	// 	return funcName_values

	// end of contextual funcs

	// start of placeholder funcs
	case funcID_include:
		return funcName_include
	case funcID_var:
		return funcName_var

	// end of placeholder funcs
	default:
		return ""
	}
}

const (
	_unknown_template_func funcID = iota

	// start of static funcs
	funcID_abbrev
	funcID_abbrevboth
	funcID_add
	funcID_add1
	funcID_add1f
	funcID_addPrefix
	funcID_addSuffix
	funcID_addf
	funcID_adler32sum
	funcID_ago
	funcID_all
	funcID_any
	funcID_append
	funcID_archconv_AlpineArch
	funcID_archconv_AlpineTripleName
	funcID_archconv_DebianArch
	funcID_archconv_DebianTripleName
	funcID_archconv_DockerArch
	funcID_archconv_DockerArchVariant
	funcID_archconv_DockerHubArch
	funcID_archconv_DockerOS
	funcID_archconv_DockerPlatformArch
	funcID_archconv_GNUArch
	funcID_archconv_GNUTripleName
	funcID_archconv_GolangArch
	funcID_archconv_GolangOS
	funcID_archconv_HF
	funcID_archconv_HardFloadArch
	funcID_archconv_LLVMArch
	funcID_archconv_LLVMTripleName
	funcID_archconv_OciArch
	funcID_archconv_OciArchVariant
	funcID_archconv_OciOS
	funcID_archconv_QemuArch
	funcID_archconv_SF
	funcID_archconv_SimpleArch
	funcID_archconv_SoftFloadArch
	funcID_atoi
	funcID_b32dec
	funcID_b32enc
	funcID_b64dec
	funcID_b64enc
	funcID_base
	funcID_bcrypt
	funcID_biggest
	funcID_bool
	funcID_buildCustomCert
	funcID_camelcase
	funcID_cat
	funcID_ceil
	funcID_chunk
	funcID_clean
	funcID_coalesce
	funcID_coll_Append
	funcID_coll_Dict
	funcID_coll_Flatten
	funcID_coll_Has
	funcID_coll_Keys
	funcID_coll_Merge
	funcID_coll_Omit
	funcID_coll_Pick
	funcID_coll_Prepend
	funcID_coll_Reverse
	funcID_coll_Slice
	funcID_coll_Sort
	funcID_coll_Uniq
	funcID_coll_Values
	funcID_compact
	funcID_concat
	funcID_contains
	funcID_conv_Atoi
	funcID_conv_Bool
	funcID_conv_Default
	funcID_conv_Dict
	funcID_conv_Has
	funcID_conv_Join
	funcID_conv_ParseFloat
	funcID_conv_ParseInt
	funcID_conv_ParseUint
	funcID_conv_Slice
	funcID_conv_ToBool
	funcID_conv_ToBools
	funcID_conv_ToFloat64
	funcID_conv_ToFloat64s
	funcID_conv_ToInt
	funcID_conv_ToInt64
	funcID_conv_ToInt64s
	funcID_conv_ToInts
	funcID_conv_ToString
	funcID_conv_ToStrings
	funcID_conv_URL
	funcID_crypto_Bcrypt
	funcID_crypto_PBKDF2
	funcID_crypto_RSADecrypt
	funcID_crypto_RSADecryptBytes
	funcID_crypto_RSADerivePublicKey
	funcID_crypto_RSAEncrypt
	funcID_crypto_RSAGenerateKey
	funcID_crypto_SHA1
	funcID_crypto_SHA224
	funcID_crypto_SHA256
	funcID_crypto_SHA384
	funcID_crypto_SHA512
	funcID_crypto_SHA512_224
	funcID_crypto_SHA512_256
	funcID_crypto_WPAPSK
	funcID_date
	funcID_dateInZone
	funcID_dateModify
	funcID_date_in_zone
	funcID_date_modify
	funcID_decryptAES
	funcID_deepCopy
	funcID_deepEqual
	funcID_default
	funcID_derivePassword
	funcID_dict
	funcID_dig
	funcID_dir
	funcID_div
	funcID_divf
	funcID_duration
	funcID_durationRound
	funcID_empty
	funcID_encryptAES
	// funcID_env
	funcID_expandenv
	funcID_ext
	funcID_fail
	funcID_file_Exists
	funcID_file_IsDir
	funcID_file_Read
	funcID_file_ReadDir
	funcID_file_Stat
	funcID_file_Walk
	funcID_file_Write
	funcID_first
	funcID_flatten
	funcID_float64
	funcID_floor
	funcID_fromJson
	funcID_genCA
	funcID_genCAWithKey
	funcID_genPrivateKey
	funcID_genSelfSignedCert
	funcID_genSelfSignedCertWithKey
	funcID_genSignedCert
	funcID_genSignedCertWithKey
	funcID_get
	funcID_getHostByName
	funcID_has
	funcID_hasKey
	funcID_hasPrefix
	funcID_hasSuffix
	funcID_htmlDate
	funcID_htmlDateInZone
	funcID_htpasswd
	funcID_indent
	funcID_initial
	funcID_initials
	funcID_int
	funcID_int64
	funcID_isAbs
	funcID_join
	funcID_jq
	funcID_jqObj
	funcID_kebabcase
	funcID_keys
	funcID_kindIs
	funcID_kindOf
	funcID_last
	funcID_list
	funcID_lower
	funcID_math_Abs
	funcID_math_Add
	funcID_math_Ceil
	funcID_math_Div
	funcID_math_Floor
	funcID_math_IsFloat
	funcID_math_IsInt
	funcID_math_IsNum
	funcID_math_Max
	funcID_math_Min
	funcID_math_Mul
	funcID_math_Pow
	funcID_math_Rem
	funcID_math_Round
	funcID_math_Seq
	funcID_math_Sub
	funcID_max
	funcID_maxf
	funcID_md5sum
	funcID_merge
	funcID_mergeOverwrite
	funcID_min
	funcID_minf
	funcID_mod
	funcID_mul
	funcID_mulf
	funcID_mustAppend
	funcID_mustChunk
	funcID_mustCompact
	funcID_mustDateModify
	funcID_mustDeepCopy
	funcID_mustFirst
	funcID_mustFromJson
	funcID_mustHas
	funcID_mustInitial
	funcID_mustLast
	funcID_mustMerge
	funcID_mustMergeOverwrite
	funcID_mustPrepend
	funcID_mustPush
	funcID_mustRegexFind
	funcID_mustRegexFindAll
	funcID_mustRegexMatch
	funcID_mustRegexReplaceAll
	funcID_mustRegexReplaceAllLiteral
	funcID_mustRegexSplit
	funcID_mustRest
	funcID_mustReverse
	funcID_mustSlice
	funcID_mustToDate
	funcID_mustToJson
	funcID_mustToPrettyJson
	funcID_mustToRawJson
	funcID_mustUniq
	funcID_mustWithout
	funcID_must_date_modify
	funcID_net_LookupCNAME
	funcID_net_LookupIP
	funcID_net_LookupIPs
	funcID_net_LookupSRV
	funcID_net_LookupSRVs
	funcID_net_LookupTXT
	funcID_nindent
	funcID_nospace
	funcID_now_Add
	funcID_now_AddDate
	funcID_now_After
	funcID_now_AppendFormat
	funcID_now_Before
	funcID_now_Clock
	funcID_now_Date
	funcID_now_Day
	funcID_now_Equal
	funcID_now_Format
	funcID_now_GoString
	funcID_now_GobEncode
	funcID_now_Hour
	funcID_now_ISOWeek
	funcID_now_In
	funcID_now_IsDST
	funcID_now_IsZero
	funcID_now_Local
	funcID_now_Location
	funcID_now_MarshalBinary
	funcID_now_MarshalJSON
	funcID_now_MarshalText
	funcID_now_Minute
	funcID_now_Month
	funcID_now_Nanosecond
	funcID_now_Round
	funcID_now_Second
	funcID_now_String
	funcID_now_Sub
	funcID_now_Truncate
	funcID_now_UTC
	funcID_now_Unix
	funcID_now_UnixMicro
	funcID_now_UnixMilli
	funcID_now_UnixNano
	funcID_now_Weekday
	funcID_now_Year
	funcID_now_YearDay
	funcID_now_Zone
	funcID_omit
	funcID_osBase
	funcID_osClean
	funcID_osDir
	funcID_osExt
	funcID_osIsAbs
	funcID_path_Base
	funcID_path_Clean
	funcID_path_Dir
	funcID_path_Ext
	funcID_path_IsAbs
	funcID_path_Join
	funcID_path_Match
	funcID_path_Split
	funcID_pick
	funcID_pluck
	funcID_plural
	funcID_pow
	funcID_prepend
	funcID_push
	funcID_quote
	funcID_randAlpha
	funcID_randAlphaNum
	funcID_randAscii
	funcID_randBytes
	funcID_randInt
	funcID_randNumeric
	funcID_random_ASCII
	funcID_random_Alpha
	funcID_random_AlphaNum
	funcID_random_Float
	funcID_random_Item
	funcID_random_Number
	funcID_random_String
	funcID_regexFind
	funcID_regexFindAll
	funcID_regexMatch
	funcID_regexQuoteMeta
	funcID_regexReplaceAll
	funcID_regexReplaceAllLiteral
	funcID_regexSplit
	funcID_regexp_Find
	funcID_regexp_FindAll
	funcID_regexp_Match
	funcID_regexp_QuoteMeta
	funcID_regexp_Replace
	funcID_regexp_ReplaceLiteral
	funcID_regexp_Split
	funcID_rem
	funcID_removePrefix
	funcID_removeSuffix
	funcID_repeat
	funcID_replace
	funcID_replaceAll
	funcID_rest
	funcID_reverse
	funcID_round
	funcID_semver
	funcID_semverCompare
	funcID_seq
	funcID_set
	funcID_sha1sum
	funcID_sha256sum
	funcID_shellQuote
	funcID_shuffle
	funcID_slice
	funcID_snakecase
	funcID_sockaddr_Attr
	funcID_sockaddr_Exclude
	funcID_sockaddr_GetAllInterfaces
	funcID_sockaddr_GetDefaultInterfaces
	funcID_sockaddr_GetInterfaceIP
	funcID_sockaddr_GetInterfaceIPs
	funcID_sockaddr_GetPrivateIP
	funcID_sockaddr_GetPrivateIPs
	funcID_sockaddr_GetPrivateInterfaces
	funcID_sockaddr_GetPublicIP
	funcID_sockaddr_GetPublicIPs
	funcID_sockaddr_GetPublicInterfaces
	funcID_sockaddr_Include
	funcID_sockaddr_Join
	funcID_sockaddr_Limit
	funcID_sockaddr_Math
	funcID_sockaddr_Offset
	funcID_sockaddr_Sort
	funcID_sockaddr_Unique
	funcID_sort
	funcID_sortAlpha
	funcID_split
	funcID_splitList
	funcID_splitN
	funcID_splitn
	funcID_squote
	funcID_strconv_Unquote
	funcID_strings_Abbrev
	funcID_strings_AddPrefix
	funcID_strings_AddSuffix
	funcID_strings_CamelCase
	funcID_strings_Contains
	funcID_strings_HasPrefix
	funcID_strings_HasSuffix
	funcID_strings_Indent
	funcID_strings_JQ
	funcID_strings_JQObj
	funcID_strings_KebabCase
	funcID_strings_Quote
	funcID_strings_RemovePrefix
	funcID_strings_RemoveSuffix
	funcID_strings_Repeat
	funcID_strings_ReplaceAll
	funcID_strings_RuneCount
	funcID_strings_ShellQuote
	funcID_strings_Slug
	funcID_strings_SnakeCase
	funcID_strings_Split
	funcID_strings_SplitN
	funcID_strings_Squote
	funcID_strings_Title
	funcID_strings_ToLower
	funcID_strings_ToUpper
	funcID_strings_Trim
	funcID_strings_TrimPrefix
	funcID_strings_TrimSpace
	funcID_strings_TrimSuffix
	funcID_strings_Trunc
	funcID_strings_WordWrap
	funcID_strings_YQ
	funcID_sub
	funcID_subf
	funcID_substr
	funcID_swapcase
	funcID_ternary
	funcID_time_Add
	funcID_time_Ceil
	funcID_time_CeilDuration
	funcID_time_Day
	funcID_time_FMT_ANSI
	funcID_time_FMT_RFC3339
	funcID_time_FMT_RFC3339Nano
	funcID_time_FMT_Ruby
	funcID_time_FMT_Stamp
	funcID_time_FMT_Unix
	funcID_time_Floor
	funcID_time_FloorDuration
	funcID_time_Format
	funcID_time_Hour
	funcID_time_Microsecond
	funcID_time_Millisecond
	funcID_time_Minute
	funcID_time_Nanosecond
	funcID_time_Now
	funcID_time_Parse
	funcID_time_ParseDuration
	funcID_time_Round
	funcID_time_RoundDuration
	funcID_time_Second
	funcID_time_Since
	funcID_time_Unix
	funcID_time_Until
	funcID_time_Week
	funcID_time_ZoneName
	funcID_time_ZoneOffset
	funcID_title
	funcID_toBytes
	funcID_toDate
	funcID_toDecimal
	funcID_toJson
	funcID_toLower
	funcID_toPrettyJson
	funcID_toRawJson
	funcID_toString
	funcID_toStrings
	funcID_toUpper
	funcID_toYaml
	funcID_totp
	funcID_trim
	funcID_trimAll
	funcID_trimPrefix
	funcID_trimSpace
	funcID_trimSuffix
	funcID_trimall
	funcID_trunc
	funcID_tuple
	funcID_typeIs
	funcID_typeIsLike
	funcID_typeOf
	funcID_uniq
	funcID_unixEpoch
	funcID_unset
	funcID_until
	funcID_untilStep
	funcID_untitle
	funcID_upper
	funcID_urlJoin
	funcID_urlParse
	funcID_uuid_IsValid
	funcID_uuid_Nil
	funcID_uuid_Parse
	funcID_uuid_V1
	funcID_uuid_V4
	funcID_values
	funcID_without
	funcID_wrap
	funcID_wrapWith
	funcID_yq

	// end of static funcs

	// start of contextual funcs
	funcID_dukkha_CacheDir
	funcID_dukkha_CrossPlatform
	funcID_dukkha_Self
	funcID_dukkha_Set
	funcID_dukkha_SetValue
	funcID_dukkha_WorkDir
	funcID_env
	funcID_filepath_Abs
	funcID_filepath_Base
	funcID_filepath_Clean
	funcID_filepath_Dir
	funcID_filepath_Ext
	funcID_filepath_FromSlash
	funcID_filepath_Glob
	funcID_filepath_IsAbs
	funcID_filepath_Join
	funcID_filepath_Match
	funcID_filepath_Rel
	funcID_filepath_Split
	funcID_filepath_ToSlash
	funcID_filepath_VolumeName
	funcID_fromYaml
	funcID_getDefaultImageTag
	funcID_getDefaultManifestTag
	funcID_git
	funcID_host
	funcID_matrix
	funcID_os_AppendFile
	funcID_os_Lookup
	funcID_os_LookupFile
	funcID_os_MkdirAll
	funcID_os_ReadFile
	funcID_os_Stderr
	funcID_os_Stdin
	funcID_os_Stdout
	funcID_os_UserCacheDir
	funcID_os_UserConfigDir
	funcID_os_UserHomeDir
	funcID_os_WriteFile
	funcID_setDefaultImageTag
	funcID_setDefaultManifestTag
	funcID_state_Failed
	funcID_state_Succeeded
	// funcID_values

	// end of contextual funcs

	// start of placeholder funcs
	funcID_include
	funcID_var

	// end of placeholder funcs

	funcID_COUNT
)

const (
	funcID_LAST_STATIC_FUNC      = funcID_yq
	funcID_LAST_CONTEXTUAL_FUNC  = funcID_values
	funcID_LAST_Placeholder_FUNC = funcID_var
)

const (
	// start of static funcs
	funcName_abbrev                        = "abbrev"
	funcName_abbrevboth                    = "abbrevboth"
	funcName_add                           = "add"
	funcName_add1                          = "add1"
	funcName_add1f                         = "add1f"
	funcName_addPrefix                     = "addPrefix"
	funcName_addSuffix                     = "addSuffix"
	funcName_addf                          = "addf"
	funcName_adler32sum                    = "adler32sum"
	funcName_ago                           = "ago"
	funcName_all                           = "all"
	funcName_any                           = "any"
	funcName_append                        = "append"
	funcName_archconv_AlpineArch           = "archconv.AlpineArch"
	funcName_archconv_AlpineTripleName     = "archconv.AlpineTripleName"
	funcName_archconv_DebianArch           = "archconv.DebianArch"
	funcName_archconv_DebianTripleName     = "archconv.DebianTripleName"
	funcName_archconv_DockerArch           = "archconv.DockerArch"
	funcName_archconv_DockerArchVariant    = "archconv.DockerArchVariant"
	funcName_archconv_DockerHubArch        = "archconv.DockerHubArch"
	funcName_archconv_DockerOS             = "archconv.DockerOS"
	funcName_archconv_DockerPlatformArch   = "archconv.DockerPlatformArch"
	funcName_archconv_GNUArch              = "archconv.GNUArch"
	funcName_archconv_GNUTripleName        = "archconv.GNUTripleName"
	funcName_archconv_GolangArch           = "archconv.GolangArch"
	funcName_archconv_GolangOS             = "archconv.GolangOS"
	funcName_archconv_HF                   = "archconv.HF"
	funcName_archconv_HardFloadArch        = "archconv.HardFloadArch"
	funcName_archconv_LLVMArch             = "archconv.LLVMArch"
	funcName_archconv_LLVMTripleName       = "archconv.LLVMTripleName"
	funcName_archconv_OciArch              = "archconv.OciArch"
	funcName_archconv_OciArchVariant       = "archconv.OciArchVariant"
	funcName_archconv_OciOS                = "archconv.OciOS"
	funcName_archconv_QemuArch             = "archconv.QemuArch"
	funcName_archconv_SF                   = "archconv.SF"
	funcName_archconv_SimpleArch           = "archconv.SimpleArch"
	funcName_archconv_SoftFloadArch        = "archconv.SoftFloadArch"
	funcName_atoi                          = "atoi"
	funcName_b32dec                        = "b32dec"
	funcName_b32enc                        = "b32enc"
	funcName_b64dec                        = "b64dec"
	funcName_b64enc                        = "b64enc"
	funcName_base                          = "base"
	funcName_bcrypt                        = "bcrypt"
	funcName_biggest                       = "biggest"
	funcName_bool                          = "bool"
	funcName_buildCustomCert               = "buildCustomCert"
	funcName_camelcase                     = "camelcase"
	funcName_cat                           = "cat"
	funcName_ceil                          = "ceil"
	funcName_chunk                         = "chunk"
	funcName_clean                         = "clean"
	funcName_coalesce                      = "coalesce"
	funcName_coll_Append                   = "coll.Append"
	funcName_coll_Dict                     = "coll.Dict"
	funcName_coll_Flatten                  = "coll.Flatten"
	funcName_coll_Has                      = "coll.Has"
	funcName_coll_Keys                     = "coll.Keys"
	funcName_coll_Merge                    = "coll.Merge"
	funcName_coll_Omit                     = "coll.Omit"
	funcName_coll_Pick                     = "coll.Pick"
	funcName_coll_Prepend                  = "coll.Prepend"
	funcName_coll_Reverse                  = "coll.Reverse"
	funcName_coll_Slice                    = "coll.Slice"
	funcName_coll_Sort                     = "coll.Sort"
	funcName_coll_Uniq                     = "coll.Uniq"
	funcName_coll_Values                   = "coll.Values"
	funcName_compact                       = "compact"
	funcName_concat                        = "concat"
	funcName_contains                      = "contains"
	funcName_conv_Atoi                     = "conv.Atoi"
	funcName_conv_Bool                     = "conv.Bool"
	funcName_conv_Default                  = "conv.Default"
	funcName_conv_Dict                     = "conv.Dict"
	funcName_conv_Has                      = "conv.Has"
	funcName_conv_Join                     = "conv.Join"
	funcName_conv_ParseFloat               = "conv.ParseFloat"
	funcName_conv_ParseInt                 = "conv.ParseInt"
	funcName_conv_ParseUint                = "conv.ParseUint"
	funcName_conv_Slice                    = "conv.Slice"
	funcName_conv_ToBool                   = "conv.ToBool"
	funcName_conv_ToBools                  = "conv.ToBools"
	funcName_conv_ToFloat64                = "conv.ToFloat64"
	funcName_conv_ToFloat64s               = "conv.ToFloat64s"
	funcName_conv_ToInt                    = "conv.ToInt"
	funcName_conv_ToInt64                  = "conv.ToInt64"
	funcName_conv_ToInt64s                 = "conv.ToInt64s"
	funcName_conv_ToInts                   = "conv.ToInts"
	funcName_conv_ToString                 = "conv.ToString"
	funcName_conv_ToStrings                = "conv.ToStrings"
	funcName_conv_URL                      = "conv.URL"
	funcName_crypto_Bcrypt                 = "crypto.Bcrypt"
	funcName_crypto_PBKDF2                 = "crypto.PBKDF2"
	funcName_crypto_RSADecrypt             = "crypto.RSADecrypt"
	funcName_crypto_RSADecryptBytes        = "crypto.RSADecryptBytes"
	funcName_crypto_RSADerivePublicKey     = "crypto.RSADerivePublicKey"
	funcName_crypto_RSAEncrypt             = "crypto.RSAEncrypt"
	funcName_crypto_RSAGenerateKey         = "crypto.RSAGenerateKey"
	funcName_crypto_SHA1                   = "crypto.SHA1"
	funcName_crypto_SHA224                 = "crypto.SHA224"
	funcName_crypto_SHA256                 = "crypto.SHA256"
	funcName_crypto_SHA384                 = "crypto.SHA384"
	funcName_crypto_SHA512                 = "crypto.SHA512"
	funcName_crypto_SHA512_224             = "crypto.SHA512_224"
	funcName_crypto_SHA512_256             = "crypto.SHA512_256"
	funcName_crypto_WPAPSK                 = "crypto.WPAPSK"
	funcName_date                          = "date"
	funcName_dateInZone                    = "dateInZone"
	funcName_dateModify                    = "dateModify"
	funcName_date_in_zone                  = "date_in_zone"
	funcName_date_modify                   = "date_modify"
	funcName_decryptAES                    = "decryptAES"
	funcName_deepCopy                      = "deepCopy"
	funcName_deepEqual                     = "deepEqual"
	funcName_default                       = "default"
	funcName_derivePassword                = "derivePassword"
	funcName_dict                          = "dict"
	funcName_dig                           = "dig"
	funcName_dir                           = "dir"
	funcName_div                           = "div"
	funcName_divf                          = "divf"
	funcName_duration                      = "duration"
	funcName_durationRound                 = "durationRound"
	funcName_empty                         = "empty"
	funcName_encryptAES                    = "encryptAES"
	funcName_env                           = "env"
	funcName_expandenv                     = "expandenv"
	funcName_ext                           = "ext"
	funcName_fail                          = "fail"
	funcName_file_Exists                   = "file.Exists"
	funcName_file_IsDir                    = "file.IsDir"
	funcName_file_Read                     = "file.Read"
	funcName_file_ReadDir                  = "file.ReadDir"
	funcName_file_Stat                     = "file.Stat"
	funcName_file_Walk                     = "file.Walk"
	funcName_file_Write                    = "file.Write"
	funcName_first                         = "first"
	funcName_flatten                       = "flatten"
	funcName_float64                       = "float64"
	funcName_floor                         = "floor"
	funcName_fromJson                      = "fromJson"
	funcName_genCA                         = "genCA"
	funcName_genCAWithKey                  = "genCAWithKey"
	funcName_genPrivateKey                 = "genPrivateKey"
	funcName_genSelfSignedCert             = "genSelfSignedCert"
	funcName_genSelfSignedCertWithKey      = "genSelfSignedCertWithKey"
	funcName_genSignedCert                 = "genSignedCert"
	funcName_genSignedCertWithKey          = "genSignedCertWithKey"
	funcName_get                           = "get"
	funcName_getHostByName                 = "getHostByName"
	funcName_has                           = "has"
	funcName_hasKey                        = "hasKey"
	funcName_hasPrefix                     = "hasPrefix"
	funcName_hasSuffix                     = "hasSuffix"
	funcName_htmlDate                      = "htmlDate"
	funcName_htmlDateInZone                = "htmlDateInZone"
	funcName_htpasswd                      = "htpasswd"
	funcName_indent                        = "indent"
	funcName_initial                       = "initial"
	funcName_initials                      = "initials"
	funcName_int                           = "int"
	funcName_int64                         = "int64"
	funcName_isAbs                         = "isAbs"
	funcName_join                          = "join"
	funcName_jq                            = "jq"
	funcName_jqObj                         = "jqObj"
	funcName_kebabcase                     = "kebabcase"
	funcName_keys                          = "keys"
	funcName_kindIs                        = "kindIs"
	funcName_kindOf                        = "kindOf"
	funcName_last                          = "last"
	funcName_list                          = "list"
	funcName_lower                         = "lower"
	funcName_math_Abs                      = "math.Abs"
	funcName_math_Add                      = "math.Add"
	funcName_math_Ceil                     = "math.Ceil"
	funcName_math_Div                      = "math.Div"
	funcName_math_Floor                    = "math.Floor"
	funcName_math_IsFloat                  = "math.IsFloat"
	funcName_math_IsInt                    = "math.IsInt"
	funcName_math_IsNum                    = "math.IsNum"
	funcName_math_Max                      = "math.Max"
	funcName_math_Min                      = "math.Min"
	funcName_math_Mul                      = "math.Mul"
	funcName_math_Pow                      = "math.Pow"
	funcName_math_Rem                      = "math.Rem"
	funcName_math_Round                    = "math.Round"
	funcName_math_Seq                      = "math.Seq"
	funcName_math_Sub                      = "math.Sub"
	funcName_max                           = "max"
	funcName_maxf                          = "maxf"
	funcName_md5sum                        = "md5sum"
	funcName_merge                         = "merge"
	funcName_mergeOverwrite                = "mergeOverwrite"
	funcName_min                           = "min"
	funcName_minf                          = "minf"
	funcName_mod                           = "mod"
	funcName_mul                           = "mul"
	funcName_mulf                          = "mulf"
	funcName_mustAppend                    = "mustAppend"
	funcName_mustChunk                     = "mustChunk"
	funcName_mustCompact                   = "mustCompact"
	funcName_mustDateModify                = "mustDateModify"
	funcName_mustDeepCopy                  = "mustDeepCopy"
	funcName_mustFirst                     = "mustFirst"
	funcName_mustFromJson                  = "mustFromJson"
	funcName_mustHas                       = "mustHas"
	funcName_mustInitial                   = "mustInitial"
	funcName_mustLast                      = "mustLast"
	funcName_mustMerge                     = "mustMerge"
	funcName_mustMergeOverwrite            = "mustMergeOverwrite"
	funcName_mustPrepend                   = "mustPrepend"
	funcName_mustPush                      = "mustPush"
	funcName_mustRegexFind                 = "mustRegexFind"
	funcName_mustRegexFindAll              = "mustRegexFindAll"
	funcName_mustRegexMatch                = "mustRegexMatch"
	funcName_mustRegexReplaceAll           = "mustRegexReplaceAll"
	funcName_mustRegexReplaceAllLiteral    = "mustRegexReplaceAllLiteral"
	funcName_mustRegexSplit                = "mustRegexSplit"
	funcName_mustRest                      = "mustRest"
	funcName_mustReverse                   = "mustReverse"
	funcName_mustSlice                     = "mustSlice"
	funcName_mustToDate                    = "mustToDate"
	funcName_mustToJson                    = "mustToJson"
	funcName_mustToPrettyJson              = "mustToPrettyJson"
	funcName_mustToRawJson                 = "mustToRawJson"
	funcName_mustUniq                      = "mustUniq"
	funcName_mustWithout                   = "mustWithout"
	funcName_must_date_modify              = "must_date_modify"
	funcName_net_LookupCNAME               = "net.LookupCNAME"
	funcName_net_LookupIP                  = "net.LookupIP"
	funcName_net_LookupIPs                 = "net.LookupIPs"
	funcName_net_LookupSRV                 = "net.LookupSRV"
	funcName_net_LookupSRVs                = "net.LookupSRVs"
	funcName_net_LookupTXT                 = "net.LookupTXT"
	funcName_nindent                       = "nindent"
	funcName_nospace                       = "nospace"
	funcName_now_Add                       = "now.Add"
	funcName_now_AddDate                   = "now.AddDate"
	funcName_now_After                     = "now.After"
	funcName_now_AppendFormat              = "now.AppendFormat"
	funcName_now_Before                    = "now.Before"
	funcName_now_Clock                     = "now.Clock"
	funcName_now_Date                      = "now.Date"
	funcName_now_Day                       = "now.Day"
	funcName_now_Equal                     = "now.Equal"
	funcName_now_Format                    = "now.Format"
	funcName_now_GoString                  = "now.GoString"
	funcName_now_GobEncode                 = "now.GobEncode"
	funcName_now_Hour                      = "now.Hour"
	funcName_now_ISOWeek                   = "now.ISOWeek"
	funcName_now_In                        = "now.In"
	funcName_now_IsDST                     = "now.IsDST"
	funcName_now_IsZero                    = "now.IsZero"
	funcName_now_Local                     = "now.Local"
	funcName_now_Location                  = "now.Location"
	funcName_now_MarshalBinary             = "now.MarshalBinary"
	funcName_now_MarshalJSON               = "now.MarshalJSON"
	funcName_now_MarshalText               = "now.MarshalText"
	funcName_now_Minute                    = "now.Minute"
	funcName_now_Month                     = "now.Month"
	funcName_now_Nanosecond                = "now.Nanosecond"
	funcName_now_Round                     = "now.Round"
	funcName_now_Second                    = "now.Second"
	funcName_now_String                    = "now.String"
	funcName_now_Sub                       = "now.Sub"
	funcName_now_Truncate                  = "now.Truncate"
	funcName_now_UTC                       = "now.UTC"
	funcName_now_Unix                      = "now.Unix"
	funcName_now_UnixMicro                 = "now.UnixMicro"
	funcName_now_UnixMilli                 = "now.UnixMilli"
	funcName_now_UnixNano                  = "now.UnixNano"
	funcName_now_Weekday                   = "now.Weekday"
	funcName_now_Year                      = "now.Year"
	funcName_now_YearDay                   = "now.YearDay"
	funcName_now_Zone                      = "now.Zone"
	funcName_omit                          = "omit"
	funcName_osBase                        = "osBase"
	funcName_osClean                       = "osClean"
	funcName_osDir                         = "osDir"
	funcName_osExt                         = "osExt"
	funcName_osIsAbs                       = "osIsAbs"
	funcName_path_Base                     = "path.Base"
	funcName_path_Clean                    = "path.Clean"
	funcName_path_Dir                      = "path.Dir"
	funcName_path_Ext                      = "path.Ext"
	funcName_path_IsAbs                    = "path.IsAbs"
	funcName_path_Join                     = "path.Join"
	funcName_path_Match                    = "path.Match"
	funcName_path_Split                    = "path.Split"
	funcName_pick                          = "pick"
	funcName_pluck                         = "pluck"
	funcName_plural                        = "plural"
	funcName_pow                           = "pow"
	funcName_prepend                       = "prepend"
	funcName_push                          = "push"
	funcName_quote                         = "quote"
	funcName_randAlpha                     = "randAlpha"
	funcName_randAlphaNum                  = "randAlphaNum"
	funcName_randAscii                     = "randAscii"
	funcName_randBytes                     = "randBytes"
	funcName_randInt                       = "randInt"
	funcName_randNumeric                   = "randNumeric"
	funcName_random_ASCII                  = "random.ASCII"
	funcName_random_Alpha                  = "random.Alpha"
	funcName_random_AlphaNum               = "random.AlphaNum"
	funcName_random_Float                  = "random.Float"
	funcName_random_Item                   = "random.Item"
	funcName_random_Number                 = "random.Number"
	funcName_random_String                 = "random.String"
	funcName_regexFind                     = "regexFind"
	funcName_regexFindAll                  = "regexFindAll"
	funcName_regexMatch                    = "regexMatch"
	funcName_regexQuoteMeta                = "regexQuoteMeta"
	funcName_regexReplaceAll               = "regexReplaceAll"
	funcName_regexReplaceAllLiteral        = "regexReplaceAllLiteral"
	funcName_regexSplit                    = "regexSplit"
	funcName_regexp_Find                   = "regexp.Find"
	funcName_regexp_FindAll                = "regexp.FindAll"
	funcName_regexp_Match                  = "regexp.Match"
	funcName_regexp_QuoteMeta              = "regexp.QuoteMeta"
	funcName_regexp_Replace                = "regexp.Replace"
	funcName_regexp_ReplaceLiteral         = "regexp.ReplaceLiteral"
	funcName_regexp_Split                  = "regexp.Split"
	funcName_rem                           = "rem"
	funcName_removePrefix                  = "removePrefix"
	funcName_removeSuffix                  = "removeSuffix"
	funcName_repeat                        = "repeat"
	funcName_replace                       = "replace"
	funcName_replaceAll                    = "replaceAll"
	funcName_rest                          = "rest"
	funcName_reverse                       = "reverse"
	funcName_round                         = "round"
	funcName_semver                        = "semver"
	funcName_semverCompare                 = "semverCompare"
	funcName_seq                           = "seq"
	funcName_set                           = "set"
	funcName_sha1sum                       = "sha1sum"
	funcName_sha256sum                     = "sha256sum"
	funcName_shellQuote                    = "shellQuote"
	funcName_shuffle                       = "shuffle"
	funcName_slice                         = "slice"
	funcName_snakecase                     = "snakecase"
	funcName_sockaddr_Attr                 = "sockaddr.Attr"
	funcName_sockaddr_Exclude              = "sockaddr.Exclude"
	funcName_sockaddr_GetAllInterfaces     = "sockaddr.GetAllInterfaces"
	funcName_sockaddr_GetDefaultInterfaces = "sockaddr.GetDefaultInterfaces"
	funcName_sockaddr_GetInterfaceIP       = "sockaddr.GetInterfaceIP"
	funcName_sockaddr_GetInterfaceIPs      = "sockaddr.GetInterfaceIPs"
	funcName_sockaddr_GetPrivateIP         = "sockaddr.GetPrivateIP"
	funcName_sockaddr_GetPrivateIPs        = "sockaddr.GetPrivateIPs"
	funcName_sockaddr_GetPrivateInterfaces = "sockaddr.GetPrivateInterfaces"
	funcName_sockaddr_GetPublicIP          = "sockaddr.GetPublicIP"
	funcName_sockaddr_GetPublicIPs         = "sockaddr.GetPublicIPs"
	funcName_sockaddr_GetPublicInterfaces  = "sockaddr.GetPublicInterfaces"
	funcName_sockaddr_Include              = "sockaddr.Include"
	funcName_sockaddr_Join                 = "sockaddr.Join"
	funcName_sockaddr_Limit                = "sockaddr.Limit"
	funcName_sockaddr_Math                 = "sockaddr.Math"
	funcName_sockaddr_Offset               = "sockaddr.Offset"
	funcName_sockaddr_Sort                 = "sockaddr.Sort"
	funcName_sockaddr_Unique               = "sockaddr.Unique"
	funcName_sort                          = "sort"
	funcName_sortAlpha                     = "sortAlpha"
	funcName_split                         = "split"
	funcName_splitList                     = "splitList"
	funcName_splitN                        = "splitN"
	funcName_splitn                        = "splitn"
	funcName_squote                        = "squote"
	funcName_strconv_Unquote               = "strconv.Unquote"
	funcName_strings_Abbrev                = "strings.Abbrev"
	funcName_strings_AddPrefix             = "strings.AddPrefix"
	funcName_strings_AddSuffix             = "strings.AddSuffix"
	funcName_strings_CamelCase             = "strings.CamelCase"
	funcName_strings_Contains              = "strings.Contains"
	funcName_strings_HasPrefix             = "strings.HasPrefix"
	funcName_strings_HasSuffix             = "strings.HasSuffix"
	funcName_strings_Indent                = "strings.Indent"
	funcName_strings_JQ                    = "strings.JQ"
	funcName_strings_JQObj                 = "strings.JQObj"
	funcName_strings_KebabCase             = "strings.KebabCase"
	funcName_strings_Quote                 = "strings.Quote"
	funcName_strings_RemovePrefix          = "strings.RemovePrefix"
	funcName_strings_RemoveSuffix          = "strings.RemoveSuffix"
	funcName_strings_Repeat                = "strings.Repeat"
	funcName_strings_ReplaceAll            = "strings.ReplaceAll"
	funcName_strings_RuneCount             = "strings.RuneCount"
	funcName_strings_ShellQuote            = "strings.ShellQuote"
	funcName_strings_Slug                  = "strings.Slug"
	funcName_strings_SnakeCase             = "strings.SnakeCase"
	funcName_strings_Split                 = "strings.Split"
	funcName_strings_SplitN                = "strings.SplitN"
	funcName_strings_Squote                = "strings.Squote"
	funcName_strings_Title                 = "strings.Title"
	funcName_strings_ToLower               = "strings.ToLower"
	funcName_strings_ToUpper               = "strings.ToUpper"
	funcName_strings_Trim                  = "strings.Trim"
	funcName_strings_TrimPrefix            = "strings.TrimPrefix"
	funcName_strings_TrimSpace             = "strings.TrimSpace"
	funcName_strings_TrimSuffix            = "strings.TrimSuffix"
	funcName_strings_Trunc                 = "strings.Trunc"
	funcName_strings_WordWrap              = "strings.WordWrap"
	funcName_strings_YQ                    = "strings.YQ"
	funcName_sub                           = "sub"
	funcName_subf                          = "subf"
	funcName_substr                        = "substr"
	funcName_swapcase                      = "swapcase"
	funcName_ternary                       = "ternary"
	funcName_time_Add                      = "time.Add"
	funcName_time_Ceil                     = "time.Ceil"
	funcName_time_CeilDuration             = "time.CeilDuration"
	funcName_time_Day                      = "time.Day"
	funcName_time_FMT_ANSI                 = "time.FMT_ANSI"
	funcName_time_FMT_RFC3339              = "time.FMT_RFC3339"
	funcName_time_FMT_RFC3339Nano          = "time.FMT_RFC3339Nano"
	funcName_time_FMT_Ruby                 = "time.FMT_Ruby"
	funcName_time_FMT_Stamp                = "time.FMT_Stamp"
	funcName_time_FMT_Unix                 = "time.FMT_Unix"
	funcName_time_Floor                    = "time.Floor"
	funcName_time_FloorDuration            = "time.FloorDuration"
	funcName_time_Format                   = "time.Format"
	funcName_time_Hour                     = "time.Hour"
	funcName_time_Microsecond              = "time.Microsecond"
	funcName_time_Millisecond              = "time.Millisecond"
	funcName_time_Minute                   = "time.Minute"
	funcName_time_Nanosecond               = "time.Nanosecond"
	funcName_time_Now                      = "time.Now"
	funcName_time_Parse                    = "time.Parse"
	funcName_time_ParseDuration            = "time.ParseDuration"
	funcName_time_Round                    = "time.Round"
	funcName_time_RoundDuration            = "time.RoundDuration"
	funcName_time_Second                   = "time.Second"
	funcName_time_Since                    = "time.Since"
	funcName_time_Unix                     = "time.Unix"
	funcName_time_Until                    = "time.Until"
	funcName_time_Week                     = "time.Week"
	funcName_time_ZoneName                 = "time.ZoneName"
	funcName_time_ZoneOffset               = "time.ZoneOffset"
	funcName_title                         = "title"
	funcName_toBytes                       = "toBytes"
	funcName_toDate                        = "toDate"
	funcName_toDecimal                     = "toDecimal"
	funcName_toJson                        = "toJson"
	funcName_toLower                       = "toLower"
	funcName_toPrettyJson                  = "toPrettyJson"
	funcName_toRawJson                     = "toRawJson"
	funcName_toString                      = "toString"
	funcName_toStrings                     = "toStrings"
	funcName_toUpper                       = "toUpper"
	funcName_toYaml                        = "toYaml"
	funcName_totp                          = "totp"
	funcName_trim                          = "trim"
	funcName_trimAll                       = "trimAll"
	funcName_trimPrefix                    = "trimPrefix"
	funcName_trimSpace                     = "trimSpace"
	funcName_trimSuffix                    = "trimSuffix"
	funcName_trimall                       = "trimall"
	funcName_trunc                         = "trunc"
	funcName_tuple                         = "tuple"
	funcName_typeIs                        = "typeIs"
	funcName_typeIsLike                    = "typeIsLike"
	funcName_typeOf                        = "typeOf"
	funcName_uniq                          = "uniq"
	funcName_unixEpoch                     = "unixEpoch"
	funcName_unset                         = "unset"
	funcName_until                         = "until"
	funcName_untilStep                     = "untilStep"
	funcName_untitle                       = "untitle"
	funcName_upper                         = "upper"
	funcName_urlJoin                       = "urlJoin"
	funcName_urlParse                      = "urlParse"
	funcName_uuid_IsValid                  = "uuid.IsValid"
	funcName_uuid_Nil                      = "uuid.Nil"
	funcName_uuid_Parse                    = "uuid.Parse"
	funcName_uuid_V1                       = "uuid.V1"
	funcName_uuid_V4                       = "uuid.V4"
	// funcName_values                        = "values"
	funcName_without  = "without"
	funcName_wrap     = "wrap"
	funcName_wrapWith = "wrapWith"
	funcName_yq       = "yq"

	// end of static funcs

	// start of contextual funcs
	funcName_dukkha_CacheDir      = "dukkha.CacheDir"
	funcName_dukkha_CrossPlatform = "dukkha.CrossPlatform"
	funcName_dukkha_Self          = "dukkha.Self"
	funcName_dukkha_Set           = "dukkha.Set"
	funcName_dukkha_SetValue      = "dukkha.SetValue"
	funcName_dukkha_WorkDir       = "dukkha.WorkDir"
	// funcName_env                   = "env"
	funcName_filepath_Abs          = "filepath.Abs"
	funcName_filepath_Base         = "filepath.Base"
	funcName_filepath_Clean        = "filepath.Clean"
	funcName_filepath_Dir          = "filepath.Dir"
	funcName_filepath_Ext          = "filepath.Ext"
	funcName_filepath_FromSlash    = "filepath.FromSlash"
	funcName_filepath_Glob         = "filepath.Glob"
	funcName_filepath_IsAbs        = "filepath.IsAbs"
	funcName_filepath_Join         = "filepath.Join"
	funcName_filepath_Match        = "filepath.Match"
	funcName_filepath_Rel          = "filepath.Rel"
	funcName_filepath_Split        = "filepath.Split"
	funcName_filepath_ToSlash      = "filepath.ToSlash"
	funcName_filepath_VolumeName   = "filepath.VolumeName"
	funcName_fromYaml              = "fromYaml"
	funcName_getDefaultImageTag    = "getDefaultImageTag"
	funcName_getDefaultManifestTag = "getDefaultManifestTag"
	funcName_git                   = "git"
	funcName_host                  = "host"
	funcName_matrix                = "matrix"
	funcName_os_AppendFile         = "os.AppendFile"
	funcName_os_Lookup             = "os.Lookup"
	funcName_os_LookupFile         = "os.LookupFile"
	funcName_os_MkdirAll           = "os.MkdirAll"
	funcName_os_ReadFile           = "os.ReadFile"
	funcName_os_Stderr             = "os.Stderr"
	funcName_os_Stdin              = "os.Stdin"
	funcName_os_Stdout             = "os.Stdout"
	funcName_os_UserCacheDir       = "os.UserCacheDir"
	funcName_os_UserConfigDir      = "os.UserConfigDir"
	funcName_os_UserHomeDir        = "os.UserHomeDir"
	funcName_os_WriteFile          = "os.WriteFile"
	funcName_setDefaultImageTag    = "setDefaultImageTag"
	funcName_setDefaultManifestTag = "setDefaultManifestTag"
	funcName_state_Failed          = "state.Failed"
	funcName_state_Succeeded       = "state.Succeeded"
	funcName_values                = "values"

	// end of contextual funcs

	// start of placeholder funcs
	funcName_include = "include"
	funcName_var     = "var"

	// end of placeholder funcs
)
