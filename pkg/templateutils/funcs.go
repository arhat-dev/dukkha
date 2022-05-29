package templateutils

import (
	"reflect"

	"arhat.dev/dukkha/pkg/dukkha"
)

// nolint:gocyclo
func FuncNameToFuncID(name string) FuncID {
	switch name {

	// start of static funcs
	case FuncName_add:
		return FuncID_add
	case FuncName_add1:
		return FuncID_add1
	case FuncName_addPrefix:
		return FuncID_addPrefix
	case FuncName_addSuffix:
		return FuncID_addSuffix
	case FuncName_all:
		return FuncID_all
	case FuncName_and:
		return FuncID_and
	case FuncName_any:
		return FuncID_any
	case FuncName_append:
		return FuncID_append
	case FuncName_archconv:
		return FuncID_archconv
	case FuncName_archconv_AlpineArch:
		return FuncID_archconv_AlpineArch
	case FuncName_archconv_AlpineTripleName:
		return FuncID_archconv_AlpineTripleName
	case FuncName_archconv_DebianArch:
		return FuncID_archconv_DebianArch
	case FuncName_archconv_DebianTripleName:
		return FuncID_archconv_DebianTripleName
	case FuncName_archconv_DockerArch:
		return FuncID_archconv_DockerArch
	case FuncName_archconv_DockerArchVariant:
		return FuncID_archconv_DockerArchVariant
	case FuncName_archconv_DockerHubArch:
		return FuncID_archconv_DockerHubArch
	case FuncName_archconv_DockerOS:
		return FuncID_archconv_DockerOS
	case FuncName_archconv_DockerPlatformArch:
		return FuncID_archconv_DockerPlatformArch
	case FuncName_archconv_GNUArch:
		return FuncID_archconv_GNUArch
	case FuncName_archconv_GNUTripleName:
		return FuncID_archconv_GNUTripleName
	case FuncName_archconv_GolangArch:
		return FuncID_archconv_GolangArch
	case FuncName_archconv_GolangOS:
		return FuncID_archconv_GolangOS
	case FuncName_archconv_HF:
		return FuncID_archconv_HF
	case FuncName_archconv_HardFloadArch:
		return FuncID_archconv_HardFloadArch
	case FuncName_archconv_LLVMArch:
		return FuncID_archconv_LLVMArch
	case FuncName_archconv_LLVMTripleName:
		return FuncID_archconv_LLVMTripleName
	case FuncName_archconv_OciArch:
		return FuncID_archconv_OciArch
	case FuncName_archconv_OciArchVariant:
		return FuncID_archconv_OciArchVariant
	case FuncName_archconv_OciOS:
		return FuncID_archconv_OciOS
	case FuncName_archconv_QemuArch:
		return FuncID_archconv_QemuArch
	case FuncName_archconv_SF:
		return FuncID_archconv_SF
	case FuncName_archconv_SimpleArch:
		return FuncID_archconv_SimpleArch
	case FuncName_archconv_SoftFloadArch:
		return FuncID_archconv_SoftFloadArch
	case FuncName_base64:
		return FuncID_base64
	case FuncName_call:
		return FuncID_call
	case FuncName_close:
		return FuncID_close
	case FuncName_coll:
		return FuncID_coll
	case FuncName_coll_Append:
		return FuncID_coll_Append
	case FuncName_coll_Bools:
		return FuncID_coll_Bools
	case FuncName_coll_Dup:
		return FuncID_coll_Dup
	case FuncName_coll_Flatten:
		return FuncID_coll_Flatten
	case FuncName_coll_Floats:
		return FuncID_coll_Floats
	case FuncName_coll_HasAll:
		return FuncID_coll_HasAll
	case FuncName_coll_HasAny:
		return FuncID_coll_HasAny
	case FuncName_coll_Index:
		return FuncID_coll_Index
	case FuncName_coll_Ints:
		return FuncID_coll_Ints
	case FuncName_coll_Keys:
		return FuncID_coll_Keys
	case FuncName_coll_List:
		return FuncID_coll_List
	case FuncName_coll_MapAnyAny:
		return FuncID_coll_MapAnyAny
	case FuncName_coll_MapStringAny:
		return FuncID_coll_MapStringAny
	case FuncName_coll_Merge:
		return FuncID_coll_Merge
	case FuncName_coll_Omit:
		return FuncID_coll_Omit
	case FuncName_coll_Pick:
		return FuncID_coll_Pick
	case FuncName_coll_Prepend:
		return FuncID_coll_Prepend
	case FuncName_coll_Push:
		return FuncID_coll_Push
	case FuncName_coll_Reverse:
		return FuncID_coll_Reverse
	case FuncName_coll_Slice:
		return FuncID_coll_Slice
	case FuncName_coll_Sort:
		return FuncID_coll_Sort
	case FuncName_coll_Strings:
		return FuncID_coll_Strings
	case FuncName_coll_Uints:
		return FuncID_coll_Uints
	case FuncName_coll_Unique:
		return FuncID_coll_Unique
	case FuncName_coll_Values:
		return FuncID_coll_Values
	case FuncName_contains:
		return FuncID_contains
	case FuncName_cred:
		return FuncID_cred
	case FuncName_cred_Htpasswd:
		return FuncID_cred_Htpasswd
	case FuncName_cred_Totp:
		return FuncID_cred_Totp
	case FuncName_default:
		return FuncID_default
	case FuncName_dict:
		return FuncID_dict
	case FuncName_div:
		return FuncID_div
	case FuncName_dns:
		return FuncID_dns
	case FuncName_dns_CNAME:
		return FuncID_dns_CNAME
	case FuncName_dns_HOST:
		return FuncID_dns_HOST
	case FuncName_dns_IP:
		return FuncID_dns_IP
	case FuncName_dns_SRV:
		return FuncID_dns_SRV
	case FuncName_dns_TXT:
		return FuncID_dns_TXT
	case FuncName_double:
		return FuncID_double
	case FuncName_dup:
		return FuncID_dup
	case FuncName_enc:
		return FuncID_enc
	case FuncName_enc_Base32:
		return FuncID_enc_Base32
	case FuncName_enc_Base64:
		return FuncID_enc_Base64
	case FuncName_enc_Hex:
		return FuncID_enc_Hex
	case FuncName_enc_JSON:
		return FuncID_enc_JSON
	case FuncName_enc_YAML:
		return FuncID_enc_YAML
	case FuncName_eq:
		return FuncID_eq
	case FuncName_ge:
		return FuncID_ge
	case FuncName_gt:
		return FuncID_gt
	case FuncName_half:
		return FuncID_half
	case FuncName_has:
		return FuncID_has
	case FuncName_hasAny:
		return FuncID_hasAny
	case FuncName_hasPrefix:
		return FuncID_hasPrefix
	case FuncName_hasSuffix:
		return FuncID_hasSuffix
	case FuncName_hash:
		return FuncID_hash
	case FuncName_hash_ADLER32:
		return FuncID_hash_ADLER32
	case FuncName_hash_Bcrypt:
		return FuncID_hash_Bcrypt
	case FuncName_hash_CRC32:
		return FuncID_hash_CRC32
	case FuncName_hash_CRC64:
		return FuncID_hash_CRC64
	case FuncName_hash_MD4:
		return FuncID_hash_MD4
	case FuncName_hash_MD5:
		return FuncID_hash_MD5
	case FuncName_hash_RIPEMD160:
		return FuncID_hash_RIPEMD160
	case FuncName_hash_SHA1:
		return FuncID_hash_SHA1
	case FuncName_hash_SHA224:
		return FuncID_hash_SHA224
	case FuncName_hash_SHA256:
		return FuncID_hash_SHA256
	case FuncName_hash_SHA384:
		return FuncID_hash_SHA384
	case FuncName_hash_SHA512:
		return FuncID_hash_SHA512
	case FuncName_hash_SHA512_224:
		return FuncID_hash_SHA512_224
	case FuncName_hash_SHA512_256:
		return FuncID_hash_SHA512_256
	case FuncName_hex:
		return FuncID_hex
	case FuncName_html:
		return FuncID_html
	case FuncName_indent:
		return FuncID_indent
	case FuncName_index:
		return FuncID_index
	case FuncName_js:
		return FuncID_js
	case FuncName_le:
		return FuncID_le
	case FuncName_len:
		return FuncID_len
	case FuncName_list:
		return FuncID_list
	case FuncName_lower:
		return FuncID_lower
	case FuncName_lt:
		return FuncID_lt
	case FuncName_math:
		return FuncID_math
	case FuncName_math_Abs:
		return FuncID_math_Abs
	case FuncName_math_Add:
		return FuncID_math_Add
	case FuncName_math_Add1:
		return FuncID_math_Add1
	case FuncName_math_Ceil:
		return FuncID_math_Ceil
	case FuncName_math_Div:
		return FuncID_math_Div
	case FuncName_math_Double:
		return FuncID_math_Double
	case FuncName_math_Floor:
		return FuncID_math_Floor
	case FuncName_math_Half:
		return FuncID_math_Half
	case FuncName_math_Log10:
		return FuncID_math_Log10
	case FuncName_math_Log2:
		return FuncID_math_Log2
	case FuncName_math_LogE:
		return FuncID_math_LogE
	case FuncName_math_Max:
		return FuncID_math_Max
	case FuncName_math_Min:
		return FuncID_math_Min
	case FuncName_math_Mod:
		return FuncID_math_Mod
	case FuncName_math_Mul:
		return FuncID_math_Mul
	case FuncName_math_Pow:
		return FuncID_math_Pow
	case FuncName_math_Round:
		return FuncID_math_Round
	case FuncName_math_Seq:
		return FuncID_math_Seq
	case FuncName_math_Sub:
		return FuncID_math_Sub
	case FuncName_math_Sub1:
		return FuncID_math_Sub1
	case FuncName_max:
		return FuncID_max
	case FuncName_md5:
		return FuncID_md5
	case FuncName_min:
		return FuncID_min
	case FuncName_mod:
		return FuncID_mod
	case FuncName_mul:
		return FuncID_mul
	case FuncName_ne:
		return FuncID_ne
	case FuncName_nindent:
		return FuncID_nindent
	case FuncName_not:
		return FuncID_not
	case FuncName_now:
		return FuncID_now
	case FuncName_omit:
		return FuncID_omit
	case FuncName_or:
		return FuncID_or
	case FuncName_path:
		return FuncID_path
	case FuncName_path_Base:
		return FuncID_path_Base
	case FuncName_path_Clean:
		return FuncID_path_Clean
	case FuncName_path_Dir:
		return FuncID_path_Dir
	case FuncName_path_Ext:
		return FuncID_path_Ext
	case FuncName_path_IsAbs:
		return FuncID_path_IsAbs
	case FuncName_path_Join:
		return FuncID_path_Join
	case FuncName_path_Match:
		return FuncID_path_Match
	case FuncName_path_Split:
		return FuncID_path_Split
	case FuncName_pick:
		return FuncID_pick
	case FuncName_prepend:
		return FuncID_prepend
	case FuncName_print:
		return FuncID_print
	case FuncName_printf:
		return FuncID_printf
	case FuncName_println:
		return FuncID_println
	case FuncName_quote:
		return FuncID_quote
	case FuncName_re:
		return FuncID_re
	case FuncName_re_Find:
		return FuncID_re_Find
	case FuncName_re_FindAll:
		return FuncID_re_FindAll
	case FuncName_re_Match:
		return FuncID_re_Match
	case FuncName_re_QuoteMeta:
		return FuncID_re_QuoteMeta
	case FuncName_re_Replace:
		return FuncID_re_Replace
	case FuncName_re_ReplaceLiteral:
		return FuncID_re_ReplaceLiteral
	case FuncName_re_Split:
		return FuncID_re_Split
	case FuncName_removePrefix:
		return FuncID_removePrefix
	case FuncName_removeSuffix:
		return FuncID_removeSuffix
	case FuncName_replaceAll:
		return FuncID_replaceAll
	case FuncName_seq:
		return FuncID_seq
	case FuncName_sha1:
		return FuncID_sha1
	case FuncName_sha256:
		return FuncID_sha256
	case FuncName_sha512:
		return FuncID_sha512
	case FuncName_slice:
		return FuncID_slice
	case FuncName_sockaddr:
		return FuncID_sockaddr
	case FuncName_sockaddr_AllInterfaces:
		return FuncID_sockaddr_AllInterfaces
	case FuncName_sockaddr_Attr:
		return FuncID_sockaddr_Attr
	case FuncName_sockaddr_DefaultInterfaces:
		return FuncID_sockaddr_DefaultInterfaces
	case FuncName_sockaddr_Exclude:
		return FuncID_sockaddr_Exclude
	case FuncName_sockaddr_Include:
		return FuncID_sockaddr_Include
	case FuncName_sockaddr_InterfaceIP:
		return FuncID_sockaddr_InterfaceIP
	case FuncName_sockaddr_Join:
		return FuncID_sockaddr_Join
	case FuncName_sockaddr_Limit:
		return FuncID_sockaddr_Limit
	case FuncName_sockaddr_Math:
		return FuncID_sockaddr_Math
	case FuncName_sockaddr_Offset:
		return FuncID_sockaddr_Offset
	case FuncName_sockaddr_PrivateIP:
		return FuncID_sockaddr_PrivateIP
	case FuncName_sockaddr_PrivateInterfaces:
		return FuncID_sockaddr_PrivateInterfaces
	case FuncName_sockaddr_PublicIP:
		return FuncID_sockaddr_PublicIP
	case FuncName_sockaddr_PublicInterfaces:
		return FuncID_sockaddr_PublicInterfaces
	case FuncName_sockaddr_Sort:
		return FuncID_sockaddr_Sort
	case FuncName_sockaddr_Unique:
		return FuncID_sockaddr_Unique
	case FuncName_sort:
		return FuncID_sort
	case FuncName_split:
		return FuncID_split
	case FuncName_splitN:
		return FuncID_splitN
	case FuncName_squote:
		return FuncID_squote
	case FuncName_stringList:
		return FuncID_stringList
	case FuncName_strings:
		return FuncID_strings
	case FuncName_strings_Abbrev:
		return FuncID_strings_Abbrev
	case FuncName_strings_AddPrefix:
		return FuncID_strings_AddPrefix
	case FuncName_strings_AddSuffix:
		return FuncID_strings_AddSuffix
	case FuncName_strings_CamelCase:
		return FuncID_strings_CamelCase
	case FuncName_strings_Contains:
		return FuncID_strings_Contains
	case FuncName_strings_ContainsAny:
		return FuncID_strings_ContainsAny
	case FuncName_strings_DoubleQuote:
		return FuncID_strings_DoubleQuote
	case FuncName_strings_HasPrefix:
		return FuncID_strings_HasPrefix
	case FuncName_strings_HasSuffix:
		return FuncID_strings_HasSuffix
	case FuncName_strings_Indent:
		return FuncID_strings_Indent
	case FuncName_strings_Initials:
		return FuncID_strings_Initials
	case FuncName_strings_Join:
		return FuncID_strings_Join
	case FuncName_strings_KebabCase:
		return FuncID_strings_KebabCase
	case FuncName_strings_Lower:
		return FuncID_strings_Lower
	case FuncName_strings_NIndent:
		return FuncID_strings_NIndent
	case FuncName_strings_NoSpace:
		return FuncID_strings_NoSpace
	case FuncName_strings_RemovePrefix:
		return FuncID_strings_RemovePrefix
	case FuncName_strings_RemoveSuffix:
		return FuncID_strings_RemoveSuffix
	case FuncName_strings_Repeat:
		return FuncID_strings_Repeat
	case FuncName_strings_ReplaceAll:
		return FuncID_strings_ReplaceAll
	case FuncName_strings_RuneCount:
		return FuncID_strings_RuneCount
	case FuncName_strings_ShellQuote:
		return FuncID_strings_ShellQuote
	case FuncName_strings_Shuffle:
		return FuncID_strings_Shuffle
	case FuncName_strings_SingleQuote:
		return FuncID_strings_SingleQuote
	case FuncName_strings_Slug:
		return FuncID_strings_Slug
	case FuncName_strings_SnakeCase:
		return FuncID_strings_SnakeCase
	case FuncName_strings_Split:
		return FuncID_strings_Split
	case FuncName_strings_SplitN:
		return FuncID_strings_SplitN
	case FuncName_strings_Substr:
		return FuncID_strings_Substr
	case FuncName_strings_SwapCase:
		return FuncID_strings_SwapCase
	case FuncName_strings_Title:
		return FuncID_strings_Title
	case FuncName_strings_Trim:
		return FuncID_strings_Trim
	case FuncName_strings_TrimLeft:
		return FuncID_strings_TrimLeft
	case FuncName_strings_TrimPrefix:
		return FuncID_strings_TrimPrefix
	case FuncName_strings_TrimRight:
		return FuncID_strings_TrimRight
	case FuncName_strings_TrimSpace:
		return FuncID_strings_TrimSpace
	case FuncName_strings_TrimSuffix:
		return FuncID_strings_TrimSuffix
	case FuncName_strings_Unquote:
		return FuncID_strings_Unquote
	case FuncName_strings_Untitle:
		return FuncID_strings_Untitle
	case FuncName_strings_Upper:
		return FuncID_strings_Upper
	case FuncName_strings_WordWrap:
		return FuncID_strings_WordWrap
	case FuncName_sub:
		return FuncID_sub
	case FuncName_sub1:
		return FuncID_sub1
	case FuncName_time:
		return FuncID_time
	case FuncName_time_Add:
		return FuncID_time_Add
	case FuncName_time_Ceil:
		return FuncID_time_Ceil
	case FuncName_time_CeilDuration:
		return FuncID_time_CeilDuration
	case FuncName_time_Day:
		return FuncID_time_Day
	case FuncName_time_FMT_ANSI:
		return FuncID_time_FMT_ANSI
	case FuncName_time_FMT_Clock:
		return FuncID_time_FMT_Clock
	case FuncName_time_FMT_Date:
		return FuncID_time_FMT_Date
	case FuncName_time_FMT_DateTime:
		return FuncID_time_FMT_DateTime
	case FuncName_time_FMT_RFC3339:
		return FuncID_time_FMT_RFC3339
	case FuncName_time_FMT_RFC3339Nano:
		return FuncID_time_FMT_RFC3339Nano
	case FuncName_time_FMT_Ruby:
		return FuncID_time_FMT_Ruby
	case FuncName_time_FMT_Stamp:
		return FuncID_time_FMT_Stamp
	case FuncName_time_FMT_Unix:
		return FuncID_time_FMT_Unix
	case FuncName_time_Floor:
		return FuncID_time_Floor
	case FuncName_time_FloorDuration:
		return FuncID_time_FloorDuration
	case FuncName_time_Format:
		return FuncID_time_Format
	case FuncName_time_Hour:
		return FuncID_time_Hour
	case FuncName_time_Microsecond:
		return FuncID_time_Microsecond
	case FuncName_time_Millisecond:
		return FuncID_time_Millisecond
	case FuncName_time_Minute:
		return FuncID_time_Minute
	case FuncName_time_Nanosecond:
		return FuncID_time_Nanosecond
	case FuncName_time_Now:
		return FuncID_time_Now
	case FuncName_time_Parse:
		return FuncID_time_Parse
	case FuncName_time_ParseDuration:
		return FuncID_time_ParseDuration
	case FuncName_time_Round:
		return FuncID_time_Round
	case FuncName_time_RoundDuration:
		return FuncID_time_RoundDuration
	case FuncName_time_Second:
		return FuncID_time_Second
	case FuncName_time_Since:
		return FuncID_time_Since
	case FuncName_time_Until:
		return FuncID_time_Until
	case FuncName_time_Week:
		return FuncID_time_Week
	case FuncName_time_ZoneName:
		return FuncID_time_ZoneName
	case FuncName_time_ZoneOffset:
		return FuncID_time_ZoneOffset
	case FuncName_title:
		return FuncID_title
	case FuncName_toJson:
		return FuncID_toJson
	case FuncName_toString:
		return FuncID_toString
	case FuncName_toYaml:
		return FuncID_toYaml
	case FuncName_totp:
		return FuncID_totp
	case FuncName_trim:
		return FuncID_trim
	case FuncName_trimPrefix:
		return FuncID_trimPrefix
	case FuncName_trimSpace:
		return FuncID_trimSpace
	case FuncName_trimSuffix:
		return FuncID_trimSuffix
	case FuncName_type:
		return FuncID_type
	case FuncName_type_AllTrue:
		return FuncID_type_AllTrue
	case FuncName_type_AnyTrue:
		return FuncID_type_AnyTrue
	case FuncName_type_Close:
		return FuncID_type_Close
	case FuncName_type_Default:
		return FuncID_type_Default
	case FuncName_type_FirstNoneZero:
		return FuncID_type_FirstNoneZero
	case FuncName_type_IsBool:
		return FuncID_type_IsBool
	case FuncName_type_IsFloat:
		return FuncID_type_IsFloat
	case FuncName_type_IsInt:
		return FuncID_type_IsInt
	case FuncName_type_IsNum:
		return FuncID_type_IsNum
	case FuncName_type_IsZero:
		return FuncID_type_IsZero
	case FuncName_type_ToBool:
		return FuncID_type_ToBool
	case FuncName_type_ToFloat:
		return FuncID_type_ToFloat
	case FuncName_type_ToInt:
		return FuncID_type_ToInt
	case FuncName_type_ToString:
		return FuncID_type_ToString
	case FuncName_type_ToStrings:
		return FuncID_type_ToStrings
	case FuncName_type_ToUint:
		return FuncID_type_ToUint
	case FuncName_uniq:
		return FuncID_uniq
	case FuncName_upper:
		return FuncID_upper
	case FuncName_urlquery:
		return FuncID_urlquery
	case FuncName_uuid:
		return FuncID_uuid
	case FuncName_uuid_IsValid:
		return FuncID_uuid_IsValid
	case FuncName_uuid_New:
		return FuncID_uuid_New
	case FuncName_uuid_V1:
		return FuncID_uuid_V1
	case FuncName_uuid_V4:
		return FuncID_uuid_V4
	case FuncName_uuid_Zero:
		return FuncID_uuid_Zero

	// end of static funcs

	// start of contextual funcs
	case FuncName_VALUE:
		return FuncID_VALUE
	case FuncName_dukkha:
		return FuncID_dukkha
	case FuncName_dukkha_CacheDir:
		return FuncID_dukkha_CacheDir
	case FuncName_dukkha_CrossPlatform:
		return FuncID_dukkha_CrossPlatform
	case FuncName_dukkha_FromJson:
		return FuncID_dukkha_FromJson
	case FuncName_dukkha_FromYaml:
		return FuncID_dukkha_FromYaml
	case FuncName_dukkha_JQ:
		return FuncID_dukkha_JQ
	case FuncName_dukkha_JQObj:
		return FuncID_dukkha_JQObj
	case FuncName_dukkha_Self:
		return FuncID_dukkha_Self
	case FuncName_dukkha_Set:
		return FuncID_dukkha_Set
	case FuncName_dukkha_SetValue:
		return FuncID_dukkha_SetValue
	case FuncName_dukkha_WorkDir:
		return FuncID_dukkha_WorkDir
	case FuncName_dukkha_YQ:
		return FuncID_dukkha_YQ
	case FuncName_dukkha_YQObj:
		return FuncID_dukkha_YQObj
	case FuncName_env:
		return FuncID_env
	case FuncName_eval:
		return FuncID_eval
	case FuncName_eval_Env:
		return FuncID_eval_Env
	case FuncName_eval_Shell:
		return FuncID_eval_Shell
	case FuncName_eval_Template:
		return FuncID_eval_Template
	case FuncName_find:
		return FuncID_find
	case FuncName_fromJson:
		return FuncID_fromJson
	case FuncName_fromYaml:
		return FuncID_fromYaml
	case FuncName_fs:
		return FuncID_fs
	case FuncName_fs_Abs:
		return FuncID_fs_Abs
	case FuncName_fs_AppendFile:
		return FuncID_fs_AppendFile
	case FuncName_fs_Base:
		return FuncID_fs_Base
	case FuncName_fs_Clean:
		return FuncID_fs_Clean
	case FuncName_fs_Dir:
		return FuncID_fs_Dir
	case FuncName_fs_Exists:
		return FuncID_fs_Exists
	case FuncName_fs_Ext:
		return FuncID_fs_Ext
	case FuncName_fs_Find:
		return FuncID_fs_Find
	case FuncName_fs_FromSlash:
		return FuncID_fs_FromSlash
	case FuncName_fs_Glob:
		return FuncID_fs_Glob
	case FuncName_fs_IsAbs:
		return FuncID_fs_IsAbs
	case FuncName_fs_IsCharDevice:
		return FuncID_fs_IsCharDevice
	case FuncName_fs_IsDevice:
		return FuncID_fs_IsDevice
	case FuncName_fs_IsDir:
		return FuncID_fs_IsDir
	case FuncName_fs_IsFIFO:
		return FuncID_fs_IsFIFO
	case FuncName_fs_IsOther:
		return FuncID_fs_IsOther
	case FuncName_fs_IsSocket:
		return FuncID_fs_IsSocket
	case FuncName_fs_IsSymlink:
		return FuncID_fs_IsSymlink
	case FuncName_fs_Join:
		return FuncID_fs_Join
	case FuncName_fs_Lookup:
		return FuncID_fs_Lookup
	case FuncName_fs_LookupFile:
		return FuncID_fs_LookupFile
	case FuncName_fs_Match:
		return FuncID_fs_Match
	case FuncName_fs_Mkdir:
		return FuncID_fs_Mkdir
	case FuncName_fs_OpenFile:
		return FuncID_fs_OpenFile
	case FuncName_fs_ReadDir:
		return FuncID_fs_ReadDir
	case FuncName_fs_ReadFile:
		return FuncID_fs_ReadFile
	case FuncName_fs_Rel:
		return FuncID_fs_Rel
	case FuncName_fs_Split:
		return FuncID_fs_Split
	case FuncName_fs_ToSlash:
		return FuncID_fs_ToSlash
	case FuncName_fs_Touch:
		return FuncID_fs_Touch
	case FuncName_fs_UserCacheDir:
		return FuncID_fs_UserCacheDir
	case FuncName_fs_UserConfigDir:
		return FuncID_fs_UserConfigDir
	case FuncName_fs_UserHomeDir:
		return FuncID_fs_UserHomeDir
	case FuncName_fs_VolumeName:
		return FuncID_fs_VolumeName
	case FuncName_fs_WriteFile:
		return FuncID_fs_WriteFile
	case FuncName_git:
		return FuncID_git
	case FuncName_host:
		return FuncID_host
	case FuncName_jq:
		return FuncID_jq
	case FuncName_jqObj:
		return FuncID_jqObj
	case FuncName_matrix:
		return FuncID_matrix
	case FuncName_mkdir:
		return FuncID_mkdir
	case FuncName_os:
		return FuncID_os
	case FuncName_os_Stderr:
		return FuncID_os_Stderr
	case FuncName_os_Stdin:
		return FuncID_os_Stdin
	case FuncName_os_Stdout:
		return FuncID_os_Stdout
	case FuncName_state:
		return FuncID_state
	case FuncName_state_Failed:
		return FuncID_state_Failed
	case FuncName_state_Succeeded:
		return FuncID_state_Succeeded
	case FuncName_tag:
		return FuncID_tag
	case FuncName_tag_ImageName:
		return FuncID_tag_ImageName
	case FuncName_tag_ImageTag:
		return FuncID_tag_ImageTag
	case FuncName_tag_ManifestName:
		return FuncID_tag_ManifestName
	case FuncName_tag_ManifestTag:
		return FuncID_tag_ManifestTag
	case FuncName_touch:
		return FuncID_touch
	case FuncName_values:
		return FuncID_values
	case FuncName_write:
		return FuncID_write
	case FuncName_yq:
		return FuncID_yq
	case FuncName_yqObj:
		return FuncID_yqObj

	// end of contextual funcs

	// start of placeholder funcs
	case FuncName_include:
		return FuncID_include
	case FuncName_var:
		return FuncID_var

	// end of placeholder funcs

	default:
		return _unknown_template_func
	}
}

// nolint:gocyclo
func (id FuncID) String() string {
	switch id {

	// start of static funcs
	case FuncID_add:
		return FuncName_add
	case FuncID_add1:
		return FuncName_add1
	case FuncID_addPrefix:
		return FuncName_addPrefix
	case FuncID_addSuffix:
		return FuncName_addSuffix
	case FuncID_all:
		return FuncName_all
	case FuncID_and:
		return FuncName_and
	case FuncID_any:
		return FuncName_any
	case FuncID_append:
		return FuncName_append
	case FuncID_archconv:
		return FuncName_archconv
	case FuncID_archconv_AlpineArch:
		return FuncName_archconv_AlpineArch
	case FuncID_archconv_AlpineTripleName:
		return FuncName_archconv_AlpineTripleName
	case FuncID_archconv_DebianArch:
		return FuncName_archconv_DebianArch
	case FuncID_archconv_DebianTripleName:
		return FuncName_archconv_DebianTripleName
	case FuncID_archconv_DockerArch:
		return FuncName_archconv_DockerArch
	case FuncID_archconv_DockerArchVariant:
		return FuncName_archconv_DockerArchVariant
	case FuncID_archconv_DockerHubArch:
		return FuncName_archconv_DockerHubArch
	case FuncID_archconv_DockerOS:
		return FuncName_archconv_DockerOS
	case FuncID_archconv_DockerPlatformArch:
		return FuncName_archconv_DockerPlatformArch
	case FuncID_archconv_GNUArch:
		return FuncName_archconv_GNUArch
	case FuncID_archconv_GNUTripleName:
		return FuncName_archconv_GNUTripleName
	case FuncID_archconv_GolangArch:
		return FuncName_archconv_GolangArch
	case FuncID_archconv_GolangOS:
		return FuncName_archconv_GolangOS
	case FuncID_archconv_HF:
		return FuncName_archconv_HF
	case FuncID_archconv_HardFloadArch:
		return FuncName_archconv_HardFloadArch
	case FuncID_archconv_LLVMArch:
		return FuncName_archconv_LLVMArch
	case FuncID_archconv_LLVMTripleName:
		return FuncName_archconv_LLVMTripleName
	case FuncID_archconv_OciArch:
		return FuncName_archconv_OciArch
	case FuncID_archconv_OciArchVariant:
		return FuncName_archconv_OciArchVariant
	case FuncID_archconv_OciOS:
		return FuncName_archconv_OciOS
	case FuncID_archconv_QemuArch:
		return FuncName_archconv_QemuArch
	case FuncID_archconv_SF:
		return FuncName_archconv_SF
	case FuncID_archconv_SimpleArch:
		return FuncName_archconv_SimpleArch
	case FuncID_archconv_SoftFloadArch:
		return FuncName_archconv_SoftFloadArch
	case FuncID_base64:
		return FuncName_base64
	case FuncID_call:
		return FuncName_call
	case FuncID_close:
		return FuncName_close
	case FuncID_coll:
		return FuncName_coll
	case FuncID_coll_Append:
		return FuncName_coll_Append
	case FuncID_coll_Bools:
		return FuncName_coll_Bools
	case FuncID_coll_Dup:
		return FuncName_coll_Dup
	case FuncID_coll_Flatten:
		return FuncName_coll_Flatten
	case FuncID_coll_Floats:
		return FuncName_coll_Floats
	case FuncID_coll_HasAll:
		return FuncName_coll_HasAll
	case FuncID_coll_HasAny:
		return FuncName_coll_HasAny
	case FuncID_coll_Index:
		return FuncName_coll_Index
	case FuncID_coll_Ints:
		return FuncName_coll_Ints
	case FuncID_coll_Keys:
		return FuncName_coll_Keys
	case FuncID_coll_List:
		return FuncName_coll_List
	case FuncID_coll_MapAnyAny:
		return FuncName_coll_MapAnyAny
	case FuncID_coll_MapStringAny:
		return FuncName_coll_MapStringAny
	case FuncID_coll_Merge:
		return FuncName_coll_Merge
	case FuncID_coll_Omit:
		return FuncName_coll_Omit
	case FuncID_coll_Pick:
		return FuncName_coll_Pick
	case FuncID_coll_Prepend:
		return FuncName_coll_Prepend
	case FuncID_coll_Push:
		return FuncName_coll_Push
	case FuncID_coll_Reverse:
		return FuncName_coll_Reverse
	case FuncID_coll_Slice:
		return FuncName_coll_Slice
	case FuncID_coll_Sort:
		return FuncName_coll_Sort
	case FuncID_coll_Strings:
		return FuncName_coll_Strings
	case FuncID_coll_Uints:
		return FuncName_coll_Uints
	case FuncID_coll_Unique:
		return FuncName_coll_Unique
	case FuncID_coll_Values:
		return FuncName_coll_Values
	case FuncID_contains:
		return FuncName_contains
	case FuncID_cred:
		return FuncName_cred
	case FuncID_cred_Htpasswd:
		return FuncName_cred_Htpasswd
	case FuncID_cred_Totp:
		return FuncName_cred_Totp
	case FuncID_default:
		return FuncName_default
	case FuncID_dict:
		return FuncName_dict
	case FuncID_div:
		return FuncName_div
	case FuncID_dns:
		return FuncName_dns
	case FuncID_dns_CNAME:
		return FuncName_dns_CNAME
	case FuncID_dns_HOST:
		return FuncName_dns_HOST
	case FuncID_dns_IP:
		return FuncName_dns_IP
	case FuncID_dns_SRV:
		return FuncName_dns_SRV
	case FuncID_dns_TXT:
		return FuncName_dns_TXT
	case FuncID_double:
		return FuncName_double
	case FuncID_dup:
		return FuncName_dup
	case FuncID_enc:
		return FuncName_enc
	case FuncID_enc_Base32:
		return FuncName_enc_Base32
	case FuncID_enc_Base64:
		return FuncName_enc_Base64
	case FuncID_enc_Hex:
		return FuncName_enc_Hex
	case FuncID_enc_JSON:
		return FuncName_enc_JSON
	case FuncID_enc_YAML:
		return FuncName_enc_YAML
	case FuncID_eq:
		return FuncName_eq
	case FuncID_ge:
		return FuncName_ge
	case FuncID_gt:
		return FuncName_gt
	case FuncID_half:
		return FuncName_half
	case FuncID_has:
		return FuncName_has
	case FuncID_hasAny:
		return FuncName_hasAny
	case FuncID_hasPrefix:
		return FuncName_hasPrefix
	case FuncID_hasSuffix:
		return FuncName_hasSuffix
	case FuncID_hash:
		return FuncName_hash
	case FuncID_hash_ADLER32:
		return FuncName_hash_ADLER32
	case FuncID_hash_Bcrypt:
		return FuncName_hash_Bcrypt
	case FuncID_hash_CRC32:
		return FuncName_hash_CRC32
	case FuncID_hash_CRC64:
		return FuncName_hash_CRC64
	case FuncID_hash_MD4:
		return FuncName_hash_MD4
	case FuncID_hash_MD5:
		return FuncName_hash_MD5
	case FuncID_hash_RIPEMD160:
		return FuncName_hash_RIPEMD160
	case FuncID_hash_SHA1:
		return FuncName_hash_SHA1
	case FuncID_hash_SHA224:
		return FuncName_hash_SHA224
	case FuncID_hash_SHA256:
		return FuncName_hash_SHA256
	case FuncID_hash_SHA384:
		return FuncName_hash_SHA384
	case FuncID_hash_SHA512:
		return FuncName_hash_SHA512
	case FuncID_hash_SHA512_224:
		return FuncName_hash_SHA512_224
	case FuncID_hash_SHA512_256:
		return FuncName_hash_SHA512_256
	case FuncID_hex:
		return FuncName_hex
	case FuncID_html:
		return FuncName_html
	case FuncID_indent:
		return FuncName_indent
	case FuncID_index:
		return FuncName_index
	case FuncID_js:
		return FuncName_js
	case FuncID_le:
		return FuncName_le
	case FuncID_len:
		return FuncName_len
	case FuncID_list:
		return FuncName_list
	case FuncID_lower:
		return FuncName_lower
	case FuncID_lt:
		return FuncName_lt
	case FuncID_math:
		return FuncName_math
	case FuncID_math_Abs:
		return FuncName_math_Abs
	case FuncID_math_Add:
		return FuncName_math_Add
	case FuncID_math_Add1:
		return FuncName_math_Add1
	case FuncID_math_Ceil:
		return FuncName_math_Ceil
	case FuncID_math_Div:
		return FuncName_math_Div
	case FuncID_math_Double:
		return FuncName_math_Double
	case FuncID_math_Floor:
		return FuncName_math_Floor
	case FuncID_math_Half:
		return FuncName_math_Half
	case FuncID_math_Log10:
		return FuncName_math_Log10
	case FuncID_math_Log2:
		return FuncName_math_Log2
	case FuncID_math_LogE:
		return FuncName_math_LogE
	case FuncID_math_Max:
		return FuncName_math_Max
	case FuncID_math_Min:
		return FuncName_math_Min
	case FuncID_math_Mod:
		return FuncName_math_Mod
	case FuncID_math_Mul:
		return FuncName_math_Mul
	case FuncID_math_Pow:
		return FuncName_math_Pow
	case FuncID_math_Round:
		return FuncName_math_Round
	case FuncID_math_Seq:
		return FuncName_math_Seq
	case FuncID_math_Sub:
		return FuncName_math_Sub
	case FuncID_math_Sub1:
		return FuncName_math_Sub1
	case FuncID_max:
		return FuncName_max
	case FuncID_md5:
		return FuncName_md5
	case FuncID_min:
		return FuncName_min
	case FuncID_mod:
		return FuncName_mod
	case FuncID_mul:
		return FuncName_mul
	case FuncID_ne:
		return FuncName_ne
	case FuncID_nindent:
		return FuncName_nindent
	case FuncID_not:
		return FuncName_not
	case FuncID_now:
		return FuncName_now
	case FuncID_omit:
		return FuncName_omit
	case FuncID_or:
		return FuncName_or
	case FuncID_path:
		return FuncName_path
	case FuncID_path_Base:
		return FuncName_path_Base
	case FuncID_path_Clean:
		return FuncName_path_Clean
	case FuncID_path_Dir:
		return FuncName_path_Dir
	case FuncID_path_Ext:
		return FuncName_path_Ext
	case FuncID_path_IsAbs:
		return FuncName_path_IsAbs
	case FuncID_path_Join:
		return FuncName_path_Join
	case FuncID_path_Match:
		return FuncName_path_Match
	case FuncID_path_Split:
		return FuncName_path_Split
	case FuncID_pick:
		return FuncName_pick
	case FuncID_prepend:
		return FuncName_prepend
	case FuncID_print:
		return FuncName_print
	case FuncID_printf:
		return FuncName_printf
	case FuncID_println:
		return FuncName_println
	case FuncID_quote:
		return FuncName_quote
	case FuncID_re:
		return FuncName_re
	case FuncID_re_Find:
		return FuncName_re_Find
	case FuncID_re_FindAll:
		return FuncName_re_FindAll
	case FuncID_re_Match:
		return FuncName_re_Match
	case FuncID_re_QuoteMeta:
		return FuncName_re_QuoteMeta
	case FuncID_re_Replace:
		return FuncName_re_Replace
	case FuncID_re_ReplaceLiteral:
		return FuncName_re_ReplaceLiteral
	case FuncID_re_Split:
		return FuncName_re_Split
	case FuncID_removePrefix:
		return FuncName_removePrefix
	case FuncID_removeSuffix:
		return FuncName_removeSuffix
	case FuncID_replaceAll:
		return FuncName_replaceAll
	case FuncID_seq:
		return FuncName_seq
	case FuncID_sha1:
		return FuncName_sha1
	case FuncID_sha256:
		return FuncName_sha256
	case FuncID_sha512:
		return FuncName_sha512
	case FuncID_slice:
		return FuncName_slice
	case FuncID_sockaddr:
		return FuncName_sockaddr
	case FuncID_sockaddr_AllInterfaces:
		return FuncName_sockaddr_AllInterfaces
	case FuncID_sockaddr_Attr:
		return FuncName_sockaddr_Attr
	case FuncID_sockaddr_DefaultInterfaces:
		return FuncName_sockaddr_DefaultInterfaces
	case FuncID_sockaddr_Exclude:
		return FuncName_sockaddr_Exclude
	case FuncID_sockaddr_Include:
		return FuncName_sockaddr_Include
	case FuncID_sockaddr_InterfaceIP:
		return FuncName_sockaddr_InterfaceIP
	case FuncID_sockaddr_Join:
		return FuncName_sockaddr_Join
	case FuncID_sockaddr_Limit:
		return FuncName_sockaddr_Limit
	case FuncID_sockaddr_Math:
		return FuncName_sockaddr_Math
	case FuncID_sockaddr_Offset:
		return FuncName_sockaddr_Offset
	case FuncID_sockaddr_PrivateIP:
		return FuncName_sockaddr_PrivateIP
	case FuncID_sockaddr_PrivateInterfaces:
		return FuncName_sockaddr_PrivateInterfaces
	case FuncID_sockaddr_PublicIP:
		return FuncName_sockaddr_PublicIP
	case FuncID_sockaddr_PublicInterfaces:
		return FuncName_sockaddr_PublicInterfaces
	case FuncID_sockaddr_Sort:
		return FuncName_sockaddr_Sort
	case FuncID_sockaddr_Unique:
		return FuncName_sockaddr_Unique
	case FuncID_sort:
		return FuncName_sort
	case FuncID_split:
		return FuncName_split
	case FuncID_splitN:
		return FuncName_splitN
	case FuncID_squote:
		return FuncName_squote
	case FuncID_stringList:
		return FuncName_stringList
	case FuncID_strings:
		return FuncName_strings
	case FuncID_strings_Abbrev:
		return FuncName_strings_Abbrev
	case FuncID_strings_AddPrefix:
		return FuncName_strings_AddPrefix
	case FuncID_strings_AddSuffix:
		return FuncName_strings_AddSuffix
	case FuncID_strings_CamelCase:
		return FuncName_strings_CamelCase
	case FuncID_strings_Contains:
		return FuncName_strings_Contains
	case FuncID_strings_ContainsAny:
		return FuncName_strings_ContainsAny
	case FuncID_strings_DoubleQuote:
		return FuncName_strings_DoubleQuote
	case FuncID_strings_HasPrefix:
		return FuncName_strings_HasPrefix
	case FuncID_strings_HasSuffix:
		return FuncName_strings_HasSuffix
	case FuncID_strings_Indent:
		return FuncName_strings_Indent
	case FuncID_strings_Initials:
		return FuncName_strings_Initials
	case FuncID_strings_Join:
		return FuncName_strings_Join
	case FuncID_strings_KebabCase:
		return FuncName_strings_KebabCase
	case FuncID_strings_Lower:
		return FuncName_strings_Lower
	case FuncID_strings_NIndent:
		return FuncName_strings_NIndent
	case FuncID_strings_NoSpace:
		return FuncName_strings_NoSpace
	case FuncID_strings_RemovePrefix:
		return FuncName_strings_RemovePrefix
	case FuncID_strings_RemoveSuffix:
		return FuncName_strings_RemoveSuffix
	case FuncID_strings_Repeat:
		return FuncName_strings_Repeat
	case FuncID_strings_ReplaceAll:
		return FuncName_strings_ReplaceAll
	case FuncID_strings_RuneCount:
		return FuncName_strings_RuneCount
	case FuncID_strings_ShellQuote:
		return FuncName_strings_ShellQuote
	case FuncID_strings_Shuffle:
		return FuncName_strings_Shuffle
	case FuncID_strings_SingleQuote:
		return FuncName_strings_SingleQuote
	case FuncID_strings_Slug:
		return FuncName_strings_Slug
	case FuncID_strings_SnakeCase:
		return FuncName_strings_SnakeCase
	case FuncID_strings_Split:
		return FuncName_strings_Split
	case FuncID_strings_SplitN:
		return FuncName_strings_SplitN
	case FuncID_strings_Substr:
		return FuncName_strings_Substr
	case FuncID_strings_SwapCase:
		return FuncName_strings_SwapCase
	case FuncID_strings_Title:
		return FuncName_strings_Title
	case FuncID_strings_Trim:
		return FuncName_strings_Trim
	case FuncID_strings_TrimLeft:
		return FuncName_strings_TrimLeft
	case FuncID_strings_TrimPrefix:
		return FuncName_strings_TrimPrefix
	case FuncID_strings_TrimRight:
		return FuncName_strings_TrimRight
	case FuncID_strings_TrimSpace:
		return FuncName_strings_TrimSpace
	case FuncID_strings_TrimSuffix:
		return FuncName_strings_TrimSuffix
	case FuncID_strings_Unquote:
		return FuncName_strings_Unquote
	case FuncID_strings_Untitle:
		return FuncName_strings_Untitle
	case FuncID_strings_Upper:
		return FuncName_strings_Upper
	case FuncID_strings_WordWrap:
		return FuncName_strings_WordWrap
	case FuncID_sub:
		return FuncName_sub
	case FuncID_sub1:
		return FuncName_sub1
	case FuncID_time:
		return FuncName_time
	case FuncID_time_Add:
		return FuncName_time_Add
	case FuncID_time_Ceil:
		return FuncName_time_Ceil
	case FuncID_time_CeilDuration:
		return FuncName_time_CeilDuration
	case FuncID_time_Day:
		return FuncName_time_Day
	case FuncID_time_FMT_ANSI:
		return FuncName_time_FMT_ANSI
	case FuncID_time_FMT_Clock:
		return FuncName_time_FMT_Clock
	case FuncID_time_FMT_Date:
		return FuncName_time_FMT_Date
	case FuncID_time_FMT_DateTime:
		return FuncName_time_FMT_DateTime
	case FuncID_time_FMT_RFC3339:
		return FuncName_time_FMT_RFC3339
	case FuncID_time_FMT_RFC3339Nano:
		return FuncName_time_FMT_RFC3339Nano
	case FuncID_time_FMT_Ruby:
		return FuncName_time_FMT_Ruby
	case FuncID_time_FMT_Stamp:
		return FuncName_time_FMT_Stamp
	case FuncID_time_FMT_Unix:
		return FuncName_time_FMT_Unix
	case FuncID_time_Floor:
		return FuncName_time_Floor
	case FuncID_time_FloorDuration:
		return FuncName_time_FloorDuration
	case FuncID_time_Format:
		return FuncName_time_Format
	case FuncID_time_Hour:
		return FuncName_time_Hour
	case FuncID_time_Microsecond:
		return FuncName_time_Microsecond
	case FuncID_time_Millisecond:
		return FuncName_time_Millisecond
	case FuncID_time_Minute:
		return FuncName_time_Minute
	case FuncID_time_Nanosecond:
		return FuncName_time_Nanosecond
	case FuncID_time_Now:
		return FuncName_time_Now
	case FuncID_time_Parse:
		return FuncName_time_Parse
	case FuncID_time_ParseDuration:
		return FuncName_time_ParseDuration
	case FuncID_time_Round:
		return FuncName_time_Round
	case FuncID_time_RoundDuration:
		return FuncName_time_RoundDuration
	case FuncID_time_Second:
		return FuncName_time_Second
	case FuncID_time_Since:
		return FuncName_time_Since
	case FuncID_time_Until:
		return FuncName_time_Until
	case FuncID_time_Week:
		return FuncName_time_Week
	case FuncID_time_ZoneName:
		return FuncName_time_ZoneName
	case FuncID_time_ZoneOffset:
		return FuncName_time_ZoneOffset
	case FuncID_title:
		return FuncName_title
	case FuncID_toJson:
		return FuncName_toJson
	case FuncID_toString:
		return FuncName_toString
	case FuncID_toYaml:
		return FuncName_toYaml
	case FuncID_totp:
		return FuncName_totp
	case FuncID_trim:
		return FuncName_trim
	case FuncID_trimPrefix:
		return FuncName_trimPrefix
	case FuncID_trimSpace:
		return FuncName_trimSpace
	case FuncID_trimSuffix:
		return FuncName_trimSuffix
	case FuncID_type:
		return FuncName_type
	case FuncID_type_AllTrue:
		return FuncName_type_AllTrue
	case FuncID_type_AnyTrue:
		return FuncName_type_AnyTrue
	case FuncID_type_Close:
		return FuncName_type_Close
	case FuncID_type_Default:
		return FuncName_type_Default
	case FuncID_type_FirstNoneZero:
		return FuncName_type_FirstNoneZero
	case FuncID_type_IsBool:
		return FuncName_type_IsBool
	case FuncID_type_IsFloat:
		return FuncName_type_IsFloat
	case FuncID_type_IsInt:
		return FuncName_type_IsInt
	case FuncID_type_IsNum:
		return FuncName_type_IsNum
	case FuncID_type_IsZero:
		return FuncName_type_IsZero
	case FuncID_type_ToBool:
		return FuncName_type_ToBool
	case FuncID_type_ToFloat:
		return FuncName_type_ToFloat
	case FuncID_type_ToInt:
		return FuncName_type_ToInt
	case FuncID_type_ToString:
		return FuncName_type_ToString
	case FuncID_type_ToStrings:
		return FuncName_type_ToStrings
	case FuncID_type_ToUint:
		return FuncName_type_ToUint
	case FuncID_uniq:
		return FuncName_uniq
	case FuncID_upper:
		return FuncName_upper
	case FuncID_urlquery:
		return FuncName_urlquery
	case FuncID_uuid:
		return FuncName_uuid
	case FuncID_uuid_IsValid:
		return FuncName_uuid_IsValid
	case FuncID_uuid_New:
		return FuncName_uuid_New
	case FuncID_uuid_V1:
		return FuncName_uuid_V1
	case FuncID_uuid_V4:
		return FuncName_uuid_V4
	case FuncID_uuid_Zero:
		return FuncName_uuid_Zero

	// end of static funcs

	// start of contextual funcs
	case FuncID_VALUE:
		return FuncName_VALUE
	case FuncID_dukkha:
		return FuncName_dukkha
	case FuncID_dukkha_CacheDir:
		return FuncName_dukkha_CacheDir
	case FuncID_dukkha_CrossPlatform:
		return FuncName_dukkha_CrossPlatform
	case FuncID_dukkha_FromJson:
		return FuncName_dukkha_FromJson
	case FuncID_dukkha_FromYaml:
		return FuncName_dukkha_FromYaml
	case FuncID_dukkha_JQ:
		return FuncName_dukkha_JQ
	case FuncID_dukkha_JQObj:
		return FuncName_dukkha_JQObj
	case FuncID_dukkha_Self:
		return FuncName_dukkha_Self
	case FuncID_dukkha_Set:
		return FuncName_dukkha_Set
	case FuncID_dukkha_SetValue:
		return FuncName_dukkha_SetValue
	case FuncID_dukkha_WorkDir:
		return FuncName_dukkha_WorkDir
	case FuncID_dukkha_YQ:
		return FuncName_dukkha_YQ
	case FuncID_dukkha_YQObj:
		return FuncName_dukkha_YQObj
	case FuncID_env:
		return FuncName_env
	case FuncID_eval:
		return FuncName_eval
	case FuncID_eval_Env:
		return FuncName_eval_Env
	case FuncID_eval_Shell:
		return FuncName_eval_Shell
	case FuncID_eval_Template:
		return FuncName_eval_Template
	case FuncID_find:
		return FuncName_find
	case FuncID_fromJson:
		return FuncName_fromJson
	case FuncID_fromYaml:
		return FuncName_fromYaml
	case FuncID_fs:
		return FuncName_fs
	case FuncID_fs_Abs:
		return FuncName_fs_Abs
	case FuncID_fs_AppendFile:
		return FuncName_fs_AppendFile
	case FuncID_fs_Base:
		return FuncName_fs_Base
	case FuncID_fs_Clean:
		return FuncName_fs_Clean
	case FuncID_fs_Dir:
		return FuncName_fs_Dir
	case FuncID_fs_Exists:
		return FuncName_fs_Exists
	case FuncID_fs_Ext:
		return FuncName_fs_Ext
	case FuncID_fs_Find:
		return FuncName_fs_Find
	case FuncID_fs_FromSlash:
		return FuncName_fs_FromSlash
	case FuncID_fs_Glob:
		return FuncName_fs_Glob
	case FuncID_fs_IsAbs:
		return FuncName_fs_IsAbs
	case FuncID_fs_IsCharDevice:
		return FuncName_fs_IsCharDevice
	case FuncID_fs_IsDevice:
		return FuncName_fs_IsDevice
	case FuncID_fs_IsDir:
		return FuncName_fs_IsDir
	case FuncID_fs_IsFIFO:
		return FuncName_fs_IsFIFO
	case FuncID_fs_IsOther:
		return FuncName_fs_IsOther
	case FuncID_fs_IsSocket:
		return FuncName_fs_IsSocket
	case FuncID_fs_IsSymlink:
		return FuncName_fs_IsSymlink
	case FuncID_fs_Join:
		return FuncName_fs_Join
	case FuncID_fs_Lookup:
		return FuncName_fs_Lookup
	case FuncID_fs_LookupFile:
		return FuncName_fs_LookupFile
	case FuncID_fs_Match:
		return FuncName_fs_Match
	case FuncID_fs_Mkdir:
		return FuncName_fs_Mkdir
	case FuncID_fs_OpenFile:
		return FuncName_fs_OpenFile
	case FuncID_fs_ReadDir:
		return FuncName_fs_ReadDir
	case FuncID_fs_ReadFile:
		return FuncName_fs_ReadFile
	case FuncID_fs_Rel:
		return FuncName_fs_Rel
	case FuncID_fs_Split:
		return FuncName_fs_Split
	case FuncID_fs_ToSlash:
		return FuncName_fs_ToSlash
	case FuncID_fs_Touch:
		return FuncName_fs_Touch
	case FuncID_fs_UserCacheDir:
		return FuncName_fs_UserCacheDir
	case FuncID_fs_UserConfigDir:
		return FuncName_fs_UserConfigDir
	case FuncID_fs_UserHomeDir:
		return FuncName_fs_UserHomeDir
	case FuncID_fs_VolumeName:
		return FuncName_fs_VolumeName
	case FuncID_fs_WriteFile:
		return FuncName_fs_WriteFile
	case FuncID_git:
		return FuncName_git
	case FuncID_host:
		return FuncName_host
	case FuncID_jq:
		return FuncName_jq
	case FuncID_jqObj:
		return FuncName_jqObj
	case FuncID_matrix:
		return FuncName_matrix
	case FuncID_mkdir:
		return FuncName_mkdir
	case FuncID_os:
		return FuncName_os
	case FuncID_os_Stderr:
		return FuncName_os_Stderr
	case FuncID_os_Stdin:
		return FuncName_os_Stdin
	case FuncID_os_Stdout:
		return FuncName_os_Stdout
	case FuncID_state:
		return FuncName_state
	case FuncID_state_Failed:
		return FuncName_state_Failed
	case FuncID_state_Succeeded:
		return FuncName_state_Succeeded
	case FuncID_tag:
		return FuncName_tag
	case FuncID_tag_ImageName:
		return FuncName_tag_ImageName
	case FuncID_tag_ImageTag:
		return FuncName_tag_ImageTag
	case FuncID_tag_ManifestName:
		return FuncName_tag_ManifestName
	case FuncID_tag_ManifestTag:
		return FuncName_tag_ManifestTag
	case FuncID_touch:
		return FuncName_touch
	case FuncID_values:
		return FuncName_values
	case FuncID_write:
		return FuncName_write
	case FuncID_yq:
		return FuncName_yq
	case FuncID_yqObj:
		return FuncName_yqObj

	// end of contextual funcs

	// start of placeholder funcs
	case FuncID_include:
		return FuncName_include
	case FuncID_var:
		return FuncName_var

	// end of placeholder funcs
	default:
		return ""
	}
}

const (
	_unknown_template_func FuncID = iota

	// start of static funcs
	FuncID_add                         // func(...Number) (Number, error)
	FuncID_add1                        // func(Number) (Number, error)
	FuncID_addPrefix                   // func(...String) (string, error)
	FuncID_addSuffix                   // func(...String) (string, error)
	FuncID_all                         // func(...any) bool
	FuncID_and                         // func(any, ...any) any
	FuncID_any                         // func(...any) bool
	FuncID_append                      // func(...any) (any, error)
	FuncID_archconv                    // func() archconvNS
	FuncID_archconv_AlpineArch         // func(String) string
	FuncID_archconv_AlpineTripleName   // func(String) string
	FuncID_archconv_DebianArch         // func(String) string
	FuncID_archconv_DebianTripleName   // func(String, ...String) string
	FuncID_archconv_DockerArch         // func(String) string
	FuncID_archconv_DockerArchVariant  // func(String) string
	FuncID_archconv_DockerHubArch      // func(String, ...String) string
	FuncID_archconv_DockerOS           // func(String) string
	FuncID_archconv_DockerPlatformArch // func(String) string
	FuncID_archconv_GNUArch            // func(String) string
	FuncID_archconv_GNUTripleName      // func(String, ...String) string
	FuncID_archconv_GolangArch         // func(String) string
	FuncID_archconv_GolangOS           // func(String) string
	FuncID_archconv_HF                 // func(String) string
	FuncID_archconv_HardFloadArch      // func(String) string
	FuncID_archconv_LLVMArch           // func(String) string
	FuncID_archconv_LLVMTripleName     // func(String, ...String) string
	FuncID_archconv_OciArch            // func(String) string
	FuncID_archconv_OciArchVariant     // func(String) string
	FuncID_archconv_OciOS              // func(String) string
	FuncID_archconv_QemuArch           // func(String) string
	FuncID_archconv_SF                 // func(String) string
	FuncID_archconv_SimpleArch         // func(String) string
	FuncID_archconv_SoftFloadArch      // func(String) string
	FuncID_base64                      // func(...any) (string, error)
	FuncID_call                        // func(any, ...any) (any, error)
	FuncID_close                       // func(any) (None, error)
	FuncID_coll                        // func() collNS
	FuncID_coll_Append                 // func(...any) (any, error)
	FuncID_coll_Bools                  // func(...Bool) ([]bool, error)
	FuncID_coll_Dup                    // func(any) (any, error)
	FuncID_coll_Flatten                // func(...any) ([]any, error)
	FuncID_coll_Floats                 // func(...Number) ([]float64, error)
	FuncID_coll_HasAll                 // func(...any) (bool, error)
	FuncID_coll_HasAny                 // func(...any) (bool, error)
	FuncID_coll_Index                  // func(any, any) (any, error)
	FuncID_coll_Ints                   // func(...Number) ([]int64, error)
	FuncID_coll_Keys                   // func(Map) ([]string, error)
	FuncID_coll_List                   // func(...any) []any
	FuncID_coll_MapAnyAny              // func(...any) map[any]any
	FuncID_coll_MapStringAny           // func(...any) (map[string]any, error)
	FuncID_coll_Merge                  // func(...Map) (Map, error)
	FuncID_coll_Omit                   // func(...any) (any, error)
	FuncID_coll_Pick                   // func(...any) (any, error)
	FuncID_coll_Prepend                // func(...any) (any, error)
	FuncID_coll_Push                   // func(...any) (any, error)
	FuncID_coll_Reverse                // func(any) (any, error)
	FuncID_coll_Slice                  // func(...any) (Slice, error)
	FuncID_coll_Sort                   // func(Slice) (Slice, error)
	FuncID_coll_Strings                // func(...String) ([]string, error)
	FuncID_coll_Uints                  // func(...Number) ([]uint64, error)
	FuncID_coll_Unique                 // func(Slice) (any, error)
	FuncID_coll_Values                 // func(Map) ([]any, error)
	FuncID_contains                    // func(String, String) (bool, error)
	FuncID_cred                        // func() credentialNS
	FuncID_cred_Htpasswd               // func(String, String) (string, error)
	FuncID_cred_Totp                   // func(...any) (string, error)
	FuncID_default                     // func(any, any) any
	FuncID_dict                        // func(...any) (map[string]any, error)
	FuncID_div                         // func(...Number) (Number, error)
	FuncID_dns                         // func() dnsNS
	FuncID_dns_CNAME                   // func(...String) (string, error)
	FuncID_dns_HOST                    // func(...String) ([]string, error)
	FuncID_dns_IP                      // func(...String) ([]string, error)
	FuncID_dns_SRV                     // func(...String) ([]*net.SRV, error)
	FuncID_dns_TXT                     // func(...String) ([]string, error)
	FuncID_double                      // func(Number) (Number, error)
	FuncID_dup                         // func(any) (any, error)
	FuncID_enc                         // func() encNS
	FuncID_enc_Base32                  // func(...any) (string, error)
	FuncID_enc_Base64                  // func(...any) (string, error)
	FuncID_enc_Hex                     // func(...any) (string, error)
	FuncID_enc_JSON                    // func(...any) (string, error)
	FuncID_enc_YAML                    // func(...any) (string, error)
	FuncID_eq                          // func(any, ...any) (bool, error)
	FuncID_ge                          // func(any, any) (bool, error)
	FuncID_gt                          // func(any, any) (bool, error)
	FuncID_half                        // func(Number) (Number, error)
	FuncID_has                         // func(...any) (bool, error)
	FuncID_hasAny                      // func(...any) (bool, error)
	FuncID_hasPrefix                   // func(String, String) (bool, error)
	FuncID_hasSuffix                   // func(String, String) (bool, error)
	FuncID_hash                        // func() hashNS
	FuncID_hash_ADLER32                // func(...any) (string, error)
	FuncID_hash_Bcrypt                 // func(...any) (string, error)
	FuncID_hash_CRC32                  // func(...any) (string, error)
	FuncID_hash_CRC64                  // func(...any) (string, error)
	FuncID_hash_MD4                    // func(...any) (string, error)
	FuncID_hash_MD5                    // func(...any) (string, error)
	FuncID_hash_RIPEMD160              // func(...any) (string, error)
	FuncID_hash_SHA1                   // func(...any) (string, error)
	FuncID_hash_SHA224                 // func(...any) (string, error)
	FuncID_hash_SHA256                 // func(...any) (string, error)
	FuncID_hash_SHA384                 // func(...any) (string, error)
	FuncID_hash_SHA512                 // func(...any) (string, error)
	FuncID_hash_SHA512_224             // func(...any) (string, error)
	FuncID_hash_SHA512_256             // func(...any) (string, error)
	FuncID_hex                         // func(...any) (string, error)
	FuncID_html                        // func(...any) string
	FuncID_indent                      // func(...any) (string, error)
	FuncID_index                       // func(any, any) (any, error)
	FuncID_js                          // func(...any) string
	FuncID_le                          // func(any, any) (bool, error)
	FuncID_len                         // func(any) (int, error)
	FuncID_list                        // func(...any) []any
	FuncID_lower                       // func(String) (string, error)
	FuncID_lt                          // func(any, any) (bool, error)
	FuncID_math                        // func() mathNS
	FuncID_math_Abs                    // func(Number) (Number, error)
	FuncID_math_Add                    // func(...Number) (Number, error)
	FuncID_math_Add1                   // func(Number) (Number, error)
	FuncID_math_Ceil                   // func(Number) (float64, error)
	FuncID_math_Div                    // func(...Number) (Number, error)
	FuncID_math_Double                 // func(Number) (Number, error)
	FuncID_math_Floor                  // func(Number) (float64, error)
	FuncID_math_Half                   // func(Number) (Number, error)
	FuncID_math_Log10                  // func(Number) (float64, error)
	FuncID_math_Log2                   // func(Number) (float64, error)
	FuncID_math_LogE                   // func(Number) (float64, error)
	FuncID_math_Max                    // func(...Number) (any, error)
	FuncID_math_Min                    // func(...Number) (any, error)
	FuncID_math_Mod                    // func(...Number) (Number, error)
	FuncID_math_Mul                    // func(...Number) (Number, error)
	FuncID_math_Pow                    // func(...Number) (Number, error)
	FuncID_math_Round                  // func(Number) (float64, error)
	FuncID_math_Seq                    // func(...Number) ([]int64, error)
	FuncID_math_Sub                    // func(...Number) (Number, error)
	FuncID_math_Sub1                   // func(Number) (Number, error)
	FuncID_max                         // func(...Number) (any, error)
	FuncID_md5                         // func(...any) (string, error)
	FuncID_min                         // func(...Number) (any, error)
	FuncID_mod                         // func(...Number) (Number, error)
	FuncID_mul                         // func(...Number) (Number, error)
	FuncID_ne                          // func(any, any) (bool, error)
	FuncID_nindent                     // func(...any) (string, error)
	FuncID_not                         // func(any) bool
	FuncID_now                         // func(...String) (time.Time, error)
	FuncID_omit                        // func(...any) (any, error)
	FuncID_or                          // func(any, ...any) any
	FuncID_path                        // func() pathNS
	FuncID_path_Base                   // func(String) string
	FuncID_path_Clean                  // func(String) string
	FuncID_path_Dir                    // func(String) string
	FuncID_path_Ext                    // func(String) string
	FuncID_path_IsAbs                  // func(String) bool
	FuncID_path_Join                   // func(...String) (string, error)
	FuncID_path_Match                  // func(String, String) (bool, error)
	FuncID_path_Split                  // func(String) []string
	FuncID_pick                        // func(...any) (any, error)
	FuncID_prepend                     // func(...any) (any, error)
	FuncID_print                       // func(...any) string
	FuncID_printf                      // func(String, ...any) string
	FuncID_println                     // func(...any) string
	FuncID_quote                       // func(String) (string, error)
	FuncID_re                          // func() regexpNS
	FuncID_re_Find                     // func(...String) (string, error)
	FuncID_re_FindAll                  // func(...String) ([]string, error)
	FuncID_re_Match                    // func(...String) (bool, error)
	FuncID_re_QuoteMeta                // func(String) string
	FuncID_re_Replace                  // func(...String) (string, error)
	FuncID_re_ReplaceLiteral           // func(...String) (string, error)
	FuncID_re_Split                    // func(...String) ([]string, error)
	FuncID_removePrefix                // func(...String) (string, error)
	FuncID_removeSuffix                // func(...String) (string, error)
	FuncID_replaceAll                  // func(String, String, String) (string, error)
	FuncID_seq                         // func(...Number) ([]int64, error)
	FuncID_sha1                        // func(...any) (string, error)
	FuncID_sha256                      // func(...any) (string, error)
	FuncID_sha512                      // func(...any) (string, error)
	FuncID_slice                       // func(...any) (Slice, error)
	FuncID_sockaddr                    // func() sockaddrNS
	FuncID_sockaddr_AllInterfaces      // func() (sockaddr.IfAddrs, error)
	FuncID_sockaddr_Attr               // func(String, any) (string, error)
	FuncID_sockaddr_DefaultInterfaces  // func() (sockaddr.IfAddrs, error)
	FuncID_sockaddr_Exclude            // func(String, String, sockaddr.IfAddrs) (sockaddr.IfAddrs, error)
	FuncID_sockaddr_Include            // func(String, String, sockaddr.IfAddrs) (sockaddr.IfAddrs, error)
	FuncID_sockaddr_InterfaceIP        // func(String) (string, error)
	FuncID_sockaddr_Join               // func(String, String, sockaddr.IfAddrs) (string, error)
	FuncID_sockaddr_Limit              // func(Number, sockaddr.IfAddrs) (sockaddr.IfAddrs, error)
	FuncID_sockaddr_Math               // func(String, String, sockaddr.IfAddrs) (sockaddr.IfAddrs, error)
	FuncID_sockaddr_Offset             // func(Number, sockaddr.IfAddrs) (sockaddr.IfAddrs, error)
	FuncID_sockaddr_PrivateIP          // func() (string, error)
	FuncID_sockaddr_PrivateInterfaces  // func() (sockaddr.IfAddrs, error)
	FuncID_sockaddr_PublicIP           // func() (string, error)
	FuncID_sockaddr_PublicInterfaces   // func() (sockaddr.IfAddrs, error)
	FuncID_sockaddr_Sort               // func(String, sockaddr.IfAddrs) (sockaddr.IfAddrs, error)
	FuncID_sockaddr_Unique             // func(String, sockaddr.IfAddrs) (sockaddr.IfAddrs, error)
	FuncID_sort                        // func(Slice) (Slice, error)
	FuncID_split                       // func(String, String) ([]string, error)
	FuncID_splitN                      // func(String, Number, String) ([]string, error)
	FuncID_squote                      // func(String) (string, error)
	FuncID_stringList                  // func(...String) ([]string, error)
	FuncID_strings                     // func() stringsNS
	FuncID_strings_Abbrev              // func(...any) (string, error)
	FuncID_strings_AddPrefix           // func(...String) (string, error)
	FuncID_strings_AddSuffix           // func(...String) (string, error)
	FuncID_strings_CamelCase           // func(String) (string, error)
	FuncID_strings_Contains            // func(String, String) (bool, error)
	FuncID_strings_ContainsAny         // func(String, String) (bool, error)
	FuncID_strings_DoubleQuote         // func(String) (string, error)
	FuncID_strings_HasPrefix           // func(String, String) (bool, error)
	FuncID_strings_HasSuffix           // func(String, String) (bool, error)
	FuncID_strings_Indent              // func(...any) (string, error)
	FuncID_strings_Initials            // func(String) (string, error)
	FuncID_strings_Join                // func(any, any) (string, error)
	FuncID_strings_KebabCase           // func(String) (string, error)
	FuncID_strings_Lower               // func(String) (string, error)
	FuncID_strings_NIndent             // func(...any) (string, error)
	FuncID_strings_NoSpace             // func(String) (string, error)
	FuncID_strings_RemovePrefix        // func(...String) (string, error)
	FuncID_strings_RemoveSuffix        // func(...String) (string, error)
	FuncID_strings_Repeat              // func(Number, String) (string, error)
	FuncID_strings_ReplaceAll          // func(String, String, String) (string, error)
	FuncID_strings_RuneCount           // func(...String) (int, error)
	FuncID_strings_ShellQuote          // func(...String) (string, error)
	FuncID_strings_Shuffle             // func(String) (string, error)
	FuncID_strings_SingleQuote         // func(String) (string, error)
	FuncID_strings_Slug                // func(String) (string, error)
	FuncID_strings_SnakeCase           // func(String) (string, error)
	FuncID_strings_Split               // func(String, String) ([]string, error)
	FuncID_strings_SplitN              // func(String, Number, String) ([]string, error)
	FuncID_strings_Substr              // func(...any) (string, error)
	FuncID_strings_SwapCase            // func(String) (string, error)
	FuncID_strings_Title               // func(String) (string, error)
	FuncID_strings_Trim                // func(String, String) (string, error)
	FuncID_strings_TrimLeft            // func(String, String) (string, error)
	FuncID_strings_TrimPrefix          // func(String, String) (string, error)
	FuncID_strings_TrimRight           // func(String, String) (string, error)
	FuncID_strings_TrimSpace           // func(String) (string, error)
	FuncID_strings_TrimSuffix          // func(String, String) (string, error)
	FuncID_strings_Unquote             // func(String) (string, error)
	FuncID_strings_Untitle             // func(String) (string, error)
	FuncID_strings_Upper               // func(String) (string, error)
	FuncID_strings_WordWrap            // func(...any) (string, error)
	FuncID_sub                         // func(...Number) (Number, error)
	FuncID_sub1                        // func(Number) (Number, error)
	FuncID_time                        // func() timeNS
	FuncID_time_Add                    // func(Duration, Time) (time.Time, error)
	FuncID_time_Ceil                   // func(Duration, Time) (time.Time, error)
	FuncID_time_CeilDuration           // func(Duration, Duration) (time.Duration, error)
	FuncID_time_Day                    // func(...Number) (time.Duration, error)
	FuncID_time_FMT_ANSI               // func() string
	FuncID_time_FMT_Clock              // func() string
	FuncID_time_FMT_Date               // func() string
	FuncID_time_FMT_DateTime           // func() string
	FuncID_time_FMT_RFC3339            // func() string
	FuncID_time_FMT_RFC3339Nano        // func() string
	FuncID_time_FMT_Ruby               // func() string
	FuncID_time_FMT_Stamp              // func() string
	FuncID_time_FMT_Unix               // func() string
	FuncID_time_Floor                  // func(Duration, Time) (time.Time, error)
	FuncID_time_FloorDuration          // func(Duration, Duration) (time.Duration, error)
	FuncID_time_Format                 // func(...any) (string, error)
	FuncID_time_Hour                   // func(...Number) (time.Duration, error)
	FuncID_time_Microsecond            // func(...Number) (time.Duration, error)
	FuncID_time_Millisecond            // func(...Number) (time.Duration, error)
	FuncID_time_Minute                 // func(...Number) (time.Duration, error)
	FuncID_time_Nanosecond             // func(...Number) (time.Duration, error)
	FuncID_time_Now                    // func(...String) (time.Time, error)
	FuncID_time_Parse                  // func(...any) (time.Time, error)
	FuncID_time_ParseDuration          // func(String) (time.Duration, error)
	FuncID_time_Round                  // func(Duration, Time) (time.Time, error)
	FuncID_time_RoundDuration          // func(Duration, Duration) (time.Duration, error)
	FuncID_time_Second                 // func(...Number) (time.Duration, error)
	FuncID_time_Since                  // func(...Time) (time.Duration, error)
	FuncID_time_Until                  // func(...Time) (time.Duration, error)
	FuncID_time_Week                   // func(...Number) (time.Duration, error)
	FuncID_time_ZoneName               // func(...any) (string, error)
	FuncID_time_ZoneOffset             // func(...any) (int, error)
	FuncID_title                       // func(String) (string, error)
	FuncID_toJson                      // func(...any) (string, error)
	FuncID_toString                    // func(any) (string, error)
	FuncID_toYaml                      // func(...any) (string, error)
	FuncID_totp                        // func(...any) (string, error)
	FuncID_trim                        // func(String, String) (string, error)
	FuncID_trimPrefix                  // func(String, String) (string, error)
	FuncID_trimSpace                   // func(String) (string, error)
	FuncID_trimSuffix                  // func(String, String) (string, error)
	FuncID_type                        // func() typeNS
	FuncID_type_AllTrue                // func(...any) bool
	FuncID_type_AnyTrue                // func(...any) bool
	FuncID_type_Close                  // func(any) (None, error)
	FuncID_type_Default                // func(any, any) any
	FuncID_type_FirstNoneZero          // func(...any) any
	FuncID_type_IsBool                 // func(any) bool
	FuncID_type_IsFloat                // func(any) bool
	FuncID_type_IsInt                  // func(any) bool
	FuncID_type_IsNum                  // func(any) bool
	FuncID_type_IsZero                 // func(any) bool
	FuncID_type_ToBool                 // func(any) (bool, error)
	FuncID_type_ToFloat                // func(any) (float64, error)
	FuncID_type_ToInt                  // func(any) (int64, error)
	FuncID_type_ToString               // func(any) (string, error)
	FuncID_type_ToStrings              // func(any) ([]string, error)
	FuncID_type_ToUint                 // func(any) (uint64, error)
	FuncID_uniq                        // func(Slice) (any, error)
	FuncID_upper                       // func(String) (string, error)
	FuncID_urlquery                    // func(...any) string
	FuncID_uuid                        // func() uuidNS
	FuncID_uuid_IsValid                // func(String) bool
	FuncID_uuid_New                    // func(...Number) (string, error)
	FuncID_uuid_V1                     // func() (string, error)
	FuncID_uuid_V4                     // func() (string, error)
	FuncID_uuid_Zero                   // func() string

	// end of static funcs

	// start of contextual funcs
	FuncID_VALUE                // func() any
	FuncID_dukkha               // func() dukkhaNS
	FuncID_dukkha_CacheDir      // func() string
	FuncID_dukkha_CrossPlatform // func(...String) bool
	FuncID_dukkha_FromJson      // func(Bytes) (any, error)
	FuncID_dukkha_FromYaml      // func(Bytes) (any, error)
	FuncID_dukkha_JQ            // func(...any) (string, error)
	FuncID_dukkha_JQObj         // func(...any) (any, error)
	FuncID_dukkha_Self          // func(...String) (struct { Stdout string; Stderr string }, error)
	FuncID_dukkha_Set           // func(String, any) (any, error)
	FuncID_dukkha_SetValue      // func(String, any) (any, error)
	FuncID_dukkha_WorkDir       // func() string
	FuncID_dukkha_YQ            // func(...any) (string, error)
	FuncID_dukkha_YQObj         // func(...any) (any, error)
	FuncID_env                  // func() map[string]string
	FuncID_eval                 // func() evalNS
	FuncID_eval_Env             // func(...String) (string, error)
	FuncID_eval_Shell           // func(String, ...Bytes) (struct { Stdout string; Stderr string }, error)
	FuncID_eval_Template        // func(String) (string, error)
	FuncID_find                 // func(String, ...String) ([]string, error)
	FuncID_fromJson             // func(Bytes) (any, error)
	FuncID_fromYaml             // func(Bytes) (any, error)
	FuncID_fs                   // func() fsNS
	FuncID_fs_Abs               // func(String) (string, error)
	FuncID_fs_AppendFile        // func(String, ...any) (None, error)
	FuncID_fs_Base              // func(String) string
	FuncID_fs_Clean             // func(String) string
	FuncID_fs_Dir               // func(String) string
	FuncID_fs_Exists            // func(String) (bool, error)
	FuncID_fs_Ext               // func(String) string
	FuncID_fs_Find              // func(String, ...String) ([]string, error)
	FuncID_fs_FromSlash         // func(String) string
	FuncID_fs_Glob              // func(String) ([]string, error)
	FuncID_fs_IsAbs             // func(String) bool
	FuncID_fs_IsCharDevice      // func(String) bool
	FuncID_fs_IsDevice          // func(String) bool
	FuncID_fs_IsDir             // func(String) bool
	FuncID_fs_IsFIFO            // func(String) bool
	FuncID_fs_IsOther           // func(String) bool
	FuncID_fs_IsSocket          // func(String) bool
	FuncID_fs_IsSymlink         // func(String) bool
	FuncID_fs_Join              // func(...String) (string, error)
	FuncID_fs_Lookup            // func(...String) (string, error)
	FuncID_fs_LookupFile        // func(...String) (string, error)
	FuncID_fs_Match             // func(String, String) (bool, error)
	FuncID_fs_Mkdir             // func(...String) (None, error)
	FuncID_fs_OpenFile          // func(...String) (*os.File, error)
	FuncID_fs_ReadDir           // func(String) ([]string, error)
	FuncID_fs_ReadFile          // func(String) (string, error)
	FuncID_fs_Rel               // func(String, String) (string, error)
	FuncID_fs_Split             // func(String) (string, string)
	FuncID_fs_ToSlash           // func(String) string
	FuncID_fs_Touch             // func(String) (None, error)
	FuncID_fs_UserCacheDir      // func() string
	FuncID_fs_UserConfigDir     // func() string
	FuncID_fs_UserHomeDir       // func() string
	FuncID_fs_VolumeName        // func(String) string
	FuncID_fs_WriteFile         // func(String, ...any) (None, error)
	FuncID_git                  // func() map[string]string
	FuncID_host                 // func() map[string]string
	FuncID_jq                   // func(...any) (string, error)
	FuncID_jqObj                // func(...any) (any, error)
	FuncID_matrix               // func() map[string]string
	FuncID_mkdir                // func(...String) (None, error)
	FuncID_os                   // func() osNS
	FuncID_os_Stderr            // func() io.Writer
	FuncID_os_Stdin             // func() io.Reader
	FuncID_os_Stdout            // func() io.Writer
	FuncID_state                // func() stateNS
	FuncID_state_Failed         // func() bool
	FuncID_state_Succeeded      // func() bool
	FuncID_tag                  // func() tagNS
	FuncID_tag_ImageName        // func(...String) (string, error)
	FuncID_tag_ImageTag         // func(...String) (string, error)
	FuncID_tag_ManifestName     // func(...String) (string, error)
	FuncID_tag_ManifestTag      // func(...String) (string, error)
	FuncID_touch                // func(String) (None, error)
	FuncID_values               // func() map[string]any
	FuncID_write                // func(String, ...any) (None, error)
	FuncID_yq                   // func(...any) (string, error)
	FuncID_yqObj                // func(...any) (any, error)

	// end of contextual funcs

	// start of placeholder funcs
	FuncID_include // func(String, any) (string, error)
	FuncID_var     // func() map[string]any

	// end of placeholder funcs

	FuncID_COUNT
)

const (
	FuncID_LAST_Static_FUNC      = FuncID_uuid_Zero
	FuncID_LAST_Contextual_FUNC  = FuncID_yqObj
	FuncID_LAST_Placeholder_FUNC = FuncID_var
)

const (
	// start of static funcs
	FuncName_add                         = "add"
	FuncName_add1                        = "add1"
	FuncName_addPrefix                   = "addPrefix"
	FuncName_addSuffix                   = "addSuffix"
	FuncName_all                         = "all"
	FuncName_and                         = "and"
	FuncName_any                         = "any"
	FuncName_append                      = "append"
	FuncName_archconv                    = "archconv"
	FuncName_archconv_AlpineArch         = "archconv.AlpineArch"
	FuncName_archconv_AlpineTripleName   = "archconv.AlpineTripleName"
	FuncName_archconv_DebianArch         = "archconv.DebianArch"
	FuncName_archconv_DebianTripleName   = "archconv.DebianTripleName"
	FuncName_archconv_DockerArch         = "archconv.DockerArch"
	FuncName_archconv_DockerArchVariant  = "archconv.DockerArchVariant"
	FuncName_archconv_DockerHubArch      = "archconv.DockerHubArch"
	FuncName_archconv_DockerOS           = "archconv.DockerOS"
	FuncName_archconv_DockerPlatformArch = "archconv.DockerPlatformArch"
	FuncName_archconv_GNUArch            = "archconv.GNUArch"
	FuncName_archconv_GNUTripleName      = "archconv.GNUTripleName"
	FuncName_archconv_GolangArch         = "archconv.GolangArch"
	FuncName_archconv_GolangOS           = "archconv.GolangOS"
	FuncName_archconv_HF                 = "archconv.HF"
	FuncName_archconv_HardFloadArch      = "archconv.HardFloadArch"
	FuncName_archconv_LLVMArch           = "archconv.LLVMArch"
	FuncName_archconv_LLVMTripleName     = "archconv.LLVMTripleName"
	FuncName_archconv_OciArch            = "archconv.OciArch"
	FuncName_archconv_OciArchVariant     = "archconv.OciArchVariant"
	FuncName_archconv_OciOS              = "archconv.OciOS"
	FuncName_archconv_QemuArch           = "archconv.QemuArch"
	FuncName_archconv_SF                 = "archconv.SF"
	FuncName_archconv_SimpleArch         = "archconv.SimpleArch"
	FuncName_archconv_SoftFloadArch      = "archconv.SoftFloadArch"
	FuncName_base64                      = "base64"
	FuncName_call                        = "call"
	FuncName_close                       = "close"
	FuncName_coll                        = "coll"
	FuncName_coll_Append                 = "coll.Append"
	FuncName_coll_Bools                  = "coll.Bools"
	FuncName_coll_Dup                    = "coll.Dup"
	FuncName_coll_Flatten                = "coll.Flatten"
	FuncName_coll_Floats                 = "coll.Floats"
	FuncName_coll_HasAll                 = "coll.HasAll"
	FuncName_coll_HasAny                 = "coll.HasAny"
	FuncName_coll_Index                  = "coll.Index"
	FuncName_coll_Ints                   = "coll.Ints"
	FuncName_coll_Keys                   = "coll.Keys"
	FuncName_coll_List                   = "coll.List"
	FuncName_coll_MapAnyAny              = "coll.MapAnyAny"
	FuncName_coll_MapStringAny           = "coll.MapStringAny"
	FuncName_coll_Merge                  = "coll.Merge"
	FuncName_coll_Omit                   = "coll.Omit"
	FuncName_coll_Pick                   = "coll.Pick"
	FuncName_coll_Prepend                = "coll.Prepend"
	FuncName_coll_Push                   = "coll.Push"
	FuncName_coll_Reverse                = "coll.Reverse"
	FuncName_coll_Slice                  = "coll.Slice"
	FuncName_coll_Sort                   = "coll.Sort"
	FuncName_coll_Strings                = "coll.Strings"
	FuncName_coll_Uints                  = "coll.Uints"
	FuncName_coll_Unique                 = "coll.Unique"
	FuncName_coll_Values                 = "coll.Values"
	FuncName_contains                    = "contains"
	FuncName_cred                        = "cred"
	FuncName_cred_Htpasswd               = "cred.Htpasswd"
	FuncName_cred_Totp                   = "cred.Totp"
	FuncName_default                     = "default"
	FuncName_dict                        = "dict"
	FuncName_div                         = "div"
	FuncName_dns                         = "dns"
	FuncName_dns_CNAME                   = "dns.CNAME"
	FuncName_dns_HOST                    = "dns.HOST"
	FuncName_dns_IP                      = "dns.IP"
	FuncName_dns_SRV                     = "dns.SRV"
	FuncName_dns_TXT                     = "dns.TXT"
	FuncName_double                      = "double"
	FuncName_dup                         = "dup"
	FuncName_enc                         = "enc"
	FuncName_enc_Base32                  = "enc.Base32"
	FuncName_enc_Base64                  = "enc.Base64"
	FuncName_enc_Hex                     = "enc.Hex"
	FuncName_enc_JSON                    = "enc.JSON"
	FuncName_enc_YAML                    = "enc.YAML"
	FuncName_eq                          = "eq"
	FuncName_ge                          = "ge"
	FuncName_gt                          = "gt"
	FuncName_half                        = "half"
	FuncName_has                         = "has"
	FuncName_hasAny                      = "hasAny"
	FuncName_hasPrefix                   = "hasPrefix"
	FuncName_hasSuffix                   = "hasSuffix"
	FuncName_hash                        = "hash"
	FuncName_hash_ADLER32                = "hash.ADLER32"
	FuncName_hash_Bcrypt                 = "hash.Bcrypt"
	FuncName_hash_CRC32                  = "hash.CRC32"
	FuncName_hash_CRC64                  = "hash.CRC64"
	FuncName_hash_MD4                    = "hash.MD4"
	FuncName_hash_MD5                    = "hash.MD5"
	FuncName_hash_RIPEMD160              = "hash.RIPEMD160"
	FuncName_hash_SHA1                   = "hash.SHA1"
	FuncName_hash_SHA224                 = "hash.SHA224"
	FuncName_hash_SHA256                 = "hash.SHA256"
	FuncName_hash_SHA384                 = "hash.SHA384"
	FuncName_hash_SHA512                 = "hash.SHA512"
	FuncName_hash_SHA512_224             = "hash.SHA512_224"
	FuncName_hash_SHA512_256             = "hash.SHA512_256"
	FuncName_hex                         = "hex"
	FuncName_html                        = "html"
	FuncName_indent                      = "indent"
	FuncName_index                       = "index"
	FuncName_js                          = "js"
	FuncName_le                          = "le"
	FuncName_len                         = "len"
	FuncName_list                        = "list"
	FuncName_lower                       = "lower"
	FuncName_lt                          = "lt"
	FuncName_math                        = "math"
	FuncName_math_Abs                    = "math.Abs"
	FuncName_math_Add                    = "math.Add"
	FuncName_math_Add1                   = "math.Add1"
	FuncName_math_Ceil                   = "math.Ceil"
	FuncName_math_Div                    = "math.Div"
	FuncName_math_Double                 = "math.Double"
	FuncName_math_Floor                  = "math.Floor"
	FuncName_math_Half                   = "math.Half"
	FuncName_math_Log10                  = "math.Log10"
	FuncName_math_Log2                   = "math.Log2"
	FuncName_math_LogE                   = "math.LogE"
	FuncName_math_Max                    = "math.Max"
	FuncName_math_Min                    = "math.Min"
	FuncName_math_Mod                    = "math.Mod"
	FuncName_math_Mul                    = "math.Mul"
	FuncName_math_Pow                    = "math.Pow"
	FuncName_math_Round                  = "math.Round"
	FuncName_math_Seq                    = "math.Seq"
	FuncName_math_Sub                    = "math.Sub"
	FuncName_math_Sub1                   = "math.Sub1"
	FuncName_max                         = "max"
	FuncName_md5                         = "md5"
	FuncName_min                         = "min"
	FuncName_mod                         = "mod"
	FuncName_mul                         = "mul"
	FuncName_ne                          = "ne"
	FuncName_nindent                     = "nindent"
	FuncName_not                         = "not"
	FuncName_now                         = "now"
	FuncName_omit                        = "omit"
	FuncName_or                          = "or"
	FuncName_path                        = "path"
	FuncName_path_Base                   = "path.Base"
	FuncName_path_Clean                  = "path.Clean"
	FuncName_path_Dir                    = "path.Dir"
	FuncName_path_Ext                    = "path.Ext"
	FuncName_path_IsAbs                  = "path.IsAbs"
	FuncName_path_Join                   = "path.Join"
	FuncName_path_Match                  = "path.Match"
	FuncName_path_Split                  = "path.Split"
	FuncName_pick                        = "pick"
	FuncName_prepend                     = "prepend"
	FuncName_print                       = "print"
	FuncName_printf                      = "printf"
	FuncName_println                     = "println"
	FuncName_quote                       = "quote"
	FuncName_re                          = "re"
	FuncName_re_Find                     = "re.Find"
	FuncName_re_FindAll                  = "re.FindAll"
	FuncName_re_Match                    = "re.Match"
	FuncName_re_QuoteMeta                = "re.QuoteMeta"
	FuncName_re_Replace                  = "re.Replace"
	FuncName_re_ReplaceLiteral           = "re.ReplaceLiteral"
	FuncName_re_Split                    = "re.Split"
	FuncName_removePrefix                = "removePrefix"
	FuncName_removeSuffix                = "removeSuffix"
	FuncName_replaceAll                  = "replaceAll"
	FuncName_seq                         = "seq"
	FuncName_sha1                        = "sha1"
	FuncName_sha256                      = "sha256"
	FuncName_sha512                      = "sha512"
	FuncName_slice                       = "slice"
	FuncName_sockaddr                    = "sockaddr"
	FuncName_sockaddr_AllInterfaces      = "sockaddr.AllInterfaces"
	FuncName_sockaddr_Attr               = "sockaddr.Attr"
	FuncName_sockaddr_DefaultInterfaces  = "sockaddr.DefaultInterfaces"
	FuncName_sockaddr_Exclude            = "sockaddr.Exclude"
	FuncName_sockaddr_Include            = "sockaddr.Include"
	FuncName_sockaddr_InterfaceIP        = "sockaddr.InterfaceIP"
	FuncName_sockaddr_Join               = "sockaddr.Join"
	FuncName_sockaddr_Limit              = "sockaddr.Limit"
	FuncName_sockaddr_Math               = "sockaddr.Math"
	FuncName_sockaddr_Offset             = "sockaddr.Offset"
	FuncName_sockaddr_PrivateIP          = "sockaddr.PrivateIP"
	FuncName_sockaddr_PrivateInterfaces  = "sockaddr.PrivateInterfaces"
	FuncName_sockaddr_PublicIP           = "sockaddr.PublicIP"
	FuncName_sockaddr_PublicInterfaces   = "sockaddr.PublicInterfaces"
	FuncName_sockaddr_Sort               = "sockaddr.Sort"
	FuncName_sockaddr_Unique             = "sockaddr.Unique"
	FuncName_sort                        = "sort"
	FuncName_split                       = "split"
	FuncName_splitN                      = "splitN"
	FuncName_squote                      = "squote"
	FuncName_stringList                  = "stringList"
	FuncName_strings                     = "strings"
	FuncName_strings_Abbrev              = "strings.Abbrev"
	FuncName_strings_AddPrefix           = "strings.AddPrefix"
	FuncName_strings_AddSuffix           = "strings.AddSuffix"
	FuncName_strings_CamelCase           = "strings.CamelCase"
	FuncName_strings_Contains            = "strings.Contains"
	FuncName_strings_ContainsAny         = "strings.ContainsAny"
	FuncName_strings_DoubleQuote         = "strings.DoubleQuote"
	FuncName_strings_HasPrefix           = "strings.HasPrefix"
	FuncName_strings_HasSuffix           = "strings.HasSuffix"
	FuncName_strings_Indent              = "strings.Indent"
	FuncName_strings_Initials            = "strings.Initials"
	FuncName_strings_Join                = "strings.Join"
	FuncName_strings_KebabCase           = "strings.KebabCase"
	FuncName_strings_Lower               = "strings.Lower"
	FuncName_strings_NIndent             = "strings.NIndent"
	FuncName_strings_NoSpace             = "strings.NoSpace"
	FuncName_strings_RemovePrefix        = "strings.RemovePrefix"
	FuncName_strings_RemoveSuffix        = "strings.RemoveSuffix"
	FuncName_strings_Repeat              = "strings.Repeat"
	FuncName_strings_ReplaceAll          = "strings.ReplaceAll"
	FuncName_strings_RuneCount           = "strings.RuneCount"
	FuncName_strings_ShellQuote          = "strings.ShellQuote"
	FuncName_strings_Shuffle             = "strings.Shuffle"
	FuncName_strings_SingleQuote         = "strings.SingleQuote"
	FuncName_strings_Slug                = "strings.Slug"
	FuncName_strings_SnakeCase           = "strings.SnakeCase"
	FuncName_strings_Split               = "strings.Split"
	FuncName_strings_SplitN              = "strings.SplitN"
	FuncName_strings_Substr              = "strings.Substr"
	FuncName_strings_SwapCase            = "strings.SwapCase"
	FuncName_strings_Title               = "strings.Title"
	FuncName_strings_Trim                = "strings.Trim"
	FuncName_strings_TrimLeft            = "strings.TrimLeft"
	FuncName_strings_TrimPrefix          = "strings.TrimPrefix"
	FuncName_strings_TrimRight           = "strings.TrimRight"
	FuncName_strings_TrimSpace           = "strings.TrimSpace"
	FuncName_strings_TrimSuffix          = "strings.TrimSuffix"
	FuncName_strings_Unquote             = "strings.Unquote"
	FuncName_strings_Untitle             = "strings.Untitle"
	FuncName_strings_Upper               = "strings.Upper"
	FuncName_strings_WordWrap            = "strings.WordWrap"
	FuncName_sub                         = "sub"
	FuncName_sub1                        = "sub1"
	FuncName_time                        = "time"
	FuncName_time_Add                    = "time.Add"
	FuncName_time_Ceil                   = "time.Ceil"
	FuncName_time_CeilDuration           = "time.CeilDuration"
	FuncName_time_Day                    = "time.Day"
	FuncName_time_FMT_ANSI               = "time.FMT_ANSI"
	FuncName_time_FMT_Clock              = "time.FMT_Clock"
	FuncName_time_FMT_Date               = "time.FMT_Date"
	FuncName_time_FMT_DateTime           = "time.FMT_DateTime"
	FuncName_time_FMT_RFC3339            = "time.FMT_RFC3339"
	FuncName_time_FMT_RFC3339Nano        = "time.FMT_RFC3339Nano"
	FuncName_time_FMT_Ruby               = "time.FMT_Ruby"
	FuncName_time_FMT_Stamp              = "time.FMT_Stamp"
	FuncName_time_FMT_Unix               = "time.FMT_Unix"
	FuncName_time_Floor                  = "time.Floor"
	FuncName_time_FloorDuration          = "time.FloorDuration"
	FuncName_time_Format                 = "time.Format"
	FuncName_time_Hour                   = "time.Hour"
	FuncName_time_Microsecond            = "time.Microsecond"
	FuncName_time_Millisecond            = "time.Millisecond"
	FuncName_time_Minute                 = "time.Minute"
	FuncName_time_Nanosecond             = "time.Nanosecond"
	FuncName_time_Now                    = "time.Now"
	FuncName_time_Parse                  = "time.Parse"
	FuncName_time_ParseDuration          = "time.ParseDuration"
	FuncName_time_Round                  = "time.Round"
	FuncName_time_RoundDuration          = "time.RoundDuration"
	FuncName_time_Second                 = "time.Second"
	FuncName_time_Since                  = "time.Since"
	FuncName_time_Until                  = "time.Until"
	FuncName_time_Week                   = "time.Week"
	FuncName_time_ZoneName               = "time.ZoneName"
	FuncName_time_ZoneOffset             = "time.ZoneOffset"
	FuncName_title                       = "title"
	FuncName_toJson                      = "toJson"
	FuncName_toString                    = "toString"
	FuncName_toYaml                      = "toYaml"
	FuncName_totp                        = "totp"
	FuncName_trim                        = "trim"
	FuncName_trimPrefix                  = "trimPrefix"
	FuncName_trimSpace                   = "trimSpace"
	FuncName_trimSuffix                  = "trimSuffix"
	FuncName_type                        = "type"
	FuncName_type_AllTrue                = "type.AllTrue"
	FuncName_type_AnyTrue                = "type.AnyTrue"
	FuncName_type_Close                  = "type.Close"
	FuncName_type_Default                = "type.Default"
	FuncName_type_FirstNoneZero          = "type.FirstNoneZero"
	FuncName_type_IsBool                 = "type.IsBool"
	FuncName_type_IsFloat                = "type.IsFloat"
	FuncName_type_IsInt                  = "type.IsInt"
	FuncName_type_IsNum                  = "type.IsNum"
	FuncName_type_IsZero                 = "type.IsZero"
	FuncName_type_ToBool                 = "type.ToBool"
	FuncName_type_ToFloat                = "type.ToFloat"
	FuncName_type_ToInt                  = "type.ToInt"
	FuncName_type_ToString               = "type.ToString"
	FuncName_type_ToStrings              = "type.ToStrings"
	FuncName_type_ToUint                 = "type.ToUint"
	FuncName_uniq                        = "uniq"
	FuncName_upper                       = "upper"
	FuncName_urlquery                    = "urlquery"
	FuncName_uuid                        = "uuid"
	FuncName_uuid_IsValid                = "uuid.IsValid"
	FuncName_uuid_New                    = "uuid.New"
	FuncName_uuid_V1                     = "uuid.V1"
	FuncName_uuid_V4                     = "uuid.V4"
	FuncName_uuid_Zero                   = "uuid.Zero"

	// end of static funcs

	// start of contextual funcs
	FuncName_VALUE                = "VALUE"
	FuncName_dukkha               = "dukkha"
	FuncName_dukkha_CacheDir      = "dukkha.CacheDir"
	FuncName_dukkha_CrossPlatform = "dukkha.CrossPlatform"
	FuncName_dukkha_FromJson      = "dukkha.FromJson"
	FuncName_dukkha_FromYaml      = "dukkha.FromYaml"
	FuncName_dukkha_JQ            = "dukkha.JQ"
	FuncName_dukkha_JQObj         = "dukkha.JQObj"
	FuncName_dukkha_Self          = "dukkha.Self"
	FuncName_dukkha_Set           = "dukkha.Set"
	FuncName_dukkha_SetValue      = "dukkha.SetValue"
	FuncName_dukkha_WorkDir       = "dukkha.WorkDir"
	FuncName_dukkha_YQ            = "dukkha.YQ"
	FuncName_dukkha_YQObj         = "dukkha.YQObj"
	FuncName_env                  = "env"
	FuncName_eval                 = "eval"
	FuncName_eval_Env             = "eval.Env"
	FuncName_eval_Shell           = "eval.Shell"
	FuncName_eval_Template        = "eval.Template"
	FuncName_find                 = "find"
	FuncName_fromJson             = "fromJson"
	FuncName_fromYaml             = "fromYaml"
	FuncName_fs                   = "fs"
	FuncName_fs_Abs               = "fs.Abs"
	FuncName_fs_AppendFile        = "fs.AppendFile"
	FuncName_fs_Base              = "fs.Base"
	FuncName_fs_Clean             = "fs.Clean"
	FuncName_fs_Dir               = "fs.Dir"
	FuncName_fs_Exists            = "fs.Exists"
	FuncName_fs_Ext               = "fs.Ext"
	FuncName_fs_Find              = "fs.Find"
	FuncName_fs_FromSlash         = "fs.FromSlash"
	FuncName_fs_Glob              = "fs.Glob"
	FuncName_fs_IsAbs             = "fs.IsAbs"
	FuncName_fs_IsCharDevice      = "fs.IsCharDevice"
	FuncName_fs_IsDevice          = "fs.IsDevice"
	FuncName_fs_IsDir             = "fs.IsDir"
	FuncName_fs_IsFIFO            = "fs.IsFIFO"
	FuncName_fs_IsOther           = "fs.IsOther"
	FuncName_fs_IsSocket          = "fs.IsSocket"
	FuncName_fs_IsSymlink         = "fs.IsSymlink"
	FuncName_fs_Join              = "fs.Join"
	FuncName_fs_Lookup            = "fs.Lookup"
	FuncName_fs_LookupFile        = "fs.LookupFile"
	FuncName_fs_Match             = "fs.Match"
	FuncName_fs_Mkdir             = "fs.Mkdir"
	FuncName_fs_OpenFile          = "fs.OpenFile"
	FuncName_fs_ReadDir           = "fs.ReadDir"
	FuncName_fs_ReadFile          = "fs.ReadFile"
	FuncName_fs_Rel               = "fs.Rel"
	FuncName_fs_Split             = "fs.Split"
	FuncName_fs_ToSlash           = "fs.ToSlash"
	FuncName_fs_Touch             = "fs.Touch"
	FuncName_fs_UserCacheDir      = "fs.UserCacheDir"
	FuncName_fs_UserConfigDir     = "fs.UserConfigDir"
	FuncName_fs_UserHomeDir       = "fs.UserHomeDir"
	FuncName_fs_VolumeName        = "fs.VolumeName"
	FuncName_fs_WriteFile         = "fs.WriteFile"
	FuncName_git                  = "git"
	FuncName_host                 = "host"
	FuncName_jq                   = "jq"
	FuncName_jqObj                = "jqObj"
	FuncName_matrix               = "matrix"
	FuncName_mkdir                = "mkdir"
	FuncName_os                   = "os"
	FuncName_os_Stderr            = "os.Stderr"
	FuncName_os_Stdin             = "os.Stdin"
	FuncName_os_Stdout            = "os.Stdout"
	FuncName_state                = "state"
	FuncName_state_Failed         = "state.Failed"
	FuncName_state_Succeeded      = "state.Succeeded"
	FuncName_tag                  = "tag"
	FuncName_tag_ImageName        = "tag.ImageName"
	FuncName_tag_ImageTag         = "tag.ImageTag"
	FuncName_tag_ManifestName     = "tag.ManifestName"
	FuncName_tag_ManifestTag      = "tag.ManifestTag"
	FuncName_touch                = "touch"
	FuncName_values               = "values"
	FuncName_write                = "write"
	FuncName_yq                   = "yq"
	FuncName_yqObj                = "yqObj"

	// end of contextual funcs

	// start of placeholder funcs
	FuncName_include = "include"
	FuncName_var     = "var"

	// end of placeholder funcs
)

var staticFuncs = [FuncID_LAST_Static_FUNC]any{
	FuncID_add - 1:                         ns_math.Add,
	FuncID_add1 - 1:                        ns_math.Add1,
	FuncID_addPrefix - 1:                   ns_strings.AddPrefix,
	FuncID_addSuffix - 1:                   ns_strings.AddSuffix,
	FuncID_all - 1:                         ns_type.AllTrue,
	FuncID_and - 1:                         ns_golang.And,
	FuncID_any - 1:                         ns_type.AnyTrue,
	FuncID_append - 1:                      ns_coll.Append,
	FuncID_archconv - 1:                    get_ns_archconv,
	FuncID_archconv_AlpineArch - 1:         ns_archconv.AlpineArch,
	FuncID_archconv_AlpineTripleName - 1:   ns_archconv.AlpineTripleName,
	FuncID_archconv_DebianArch - 1:         ns_archconv.DebianArch,
	FuncID_archconv_DebianTripleName - 1:   ns_archconv.DebianTripleName,
	FuncID_archconv_DockerArch - 1:         ns_archconv.DockerArch,
	FuncID_archconv_DockerArchVariant - 1:  ns_archconv.DockerArchVariant,
	FuncID_archconv_DockerHubArch - 1:      ns_archconv.DockerHubArch,
	FuncID_archconv_DockerOS - 1:           ns_archconv.DockerOS,
	FuncID_archconv_DockerPlatformArch - 1: ns_archconv.DockerPlatformArch,
	FuncID_archconv_GNUArch - 1:            ns_archconv.GNUArch,
	FuncID_archconv_GNUTripleName - 1:      ns_archconv.GNUTripleName,
	FuncID_archconv_GolangArch - 1:         ns_archconv.GolangArch,
	FuncID_archconv_GolangOS - 1:           ns_archconv.GolangOS,
	FuncID_archconv_HF - 1:                 ns_archconv.HF,
	FuncID_archconv_HardFloadArch - 1:      ns_archconv.HardFloadArch,
	FuncID_archconv_LLVMArch - 1:           ns_archconv.LLVMArch,
	FuncID_archconv_LLVMTripleName - 1:     ns_archconv.LLVMTripleName,
	FuncID_archconv_OciArch - 1:            ns_archconv.OciArch,
	FuncID_archconv_OciArchVariant - 1:     ns_archconv.OciArchVariant,
	FuncID_archconv_OciOS - 1:              ns_archconv.OciOS,
	FuncID_archconv_QemuArch - 1:           ns_archconv.QemuArch,
	FuncID_archconv_SF - 1:                 ns_archconv.SF,
	FuncID_archconv_SimpleArch - 1:         ns_archconv.SimpleArch,
	FuncID_archconv_SoftFloadArch - 1:      ns_archconv.SoftFloadArch,
	FuncID_base64 - 1:                      ns_enc.Base64,
	FuncID_call - 1:                        ns_golang.Call,
	FuncID_close - 1:                       ns_type.Close,
	FuncID_coll - 1:                        get_ns_coll,
	FuncID_coll_Append - 1:                 ns_coll.Append,
	FuncID_coll_Bools - 1:                  ns_coll.Bools,
	FuncID_coll_Dup - 1:                    ns_coll.Dup,
	FuncID_coll_Flatten - 1:                ns_coll.Flatten,
	FuncID_coll_Floats - 1:                 ns_coll.Floats,
	FuncID_coll_HasAll - 1:                 ns_coll.HasAll,
	FuncID_coll_HasAny - 1:                 ns_coll.HasAny,
	FuncID_coll_Index - 1:                  ns_coll.Index,
	FuncID_coll_Ints - 1:                   ns_coll.Ints,
	FuncID_coll_Keys - 1:                   ns_coll.Keys,
	FuncID_coll_List - 1:                   ns_coll.List,
	FuncID_coll_MapAnyAny - 1:              ns_coll.MapAnyAny,
	FuncID_coll_MapStringAny - 1:           ns_coll.MapStringAny,
	FuncID_coll_Merge - 1:                  ns_coll.Merge,
	FuncID_coll_Omit - 1:                   ns_coll.Omit,
	FuncID_coll_Pick - 1:                   ns_coll.Pick,
	FuncID_coll_Prepend - 1:                ns_coll.Prepend,
	FuncID_coll_Push - 1:                   ns_coll.Push,
	FuncID_coll_Reverse - 1:                ns_coll.Reverse,
	FuncID_coll_Slice - 1:                  ns_coll.Slice,
	FuncID_coll_Sort - 1:                   ns_coll.Sort,
	FuncID_coll_Strings - 1:                ns_coll.Strings,
	FuncID_coll_Uints - 1:                  ns_coll.Uints,
	FuncID_coll_Unique - 1:                 ns_coll.Unique,
	FuncID_coll_Values - 1:                 ns_coll.Values,
	FuncID_contains - 1:                    ns_strings.Contains,
	FuncID_cred - 1:                        get_ns_cred,
	FuncID_cred_Htpasswd - 1:               ns_cred.Htpasswd,
	FuncID_cred_Totp - 1:                   ns_cred.Totp,
	FuncID_default - 1:                     ns_type.Default,
	FuncID_dict - 1:                        ns_coll.MapStringAny,
	FuncID_div - 1:                         ns_math.Div,
	FuncID_dns - 1:                         get_ns_dns,
	FuncID_dns_CNAME - 1:                   ns_dns.CNAME,
	FuncID_dns_HOST - 1:                    ns_dns.HOST,
	FuncID_dns_IP - 1:                      ns_dns.IP,
	FuncID_dns_SRV - 1:                     ns_dns.SRV,
	FuncID_dns_TXT - 1:                     ns_dns.TXT,
	FuncID_double - 1:                      ns_math.Double,
	FuncID_dup - 1:                         ns_coll.Dup,
	FuncID_enc - 1:                         get_ns_enc,
	FuncID_enc_Base32 - 1:                  ns_enc.Base32,
	FuncID_enc_Base64 - 1:                  ns_enc.Base64,
	FuncID_enc_Hex - 1:                     ns_enc.Hex,
	FuncID_enc_JSON - 1:                    ns_enc.JSON,
	FuncID_enc_YAML - 1:                    ns_enc.YAML,
	FuncID_eq - 1:                          ns_golang.Eq,
	FuncID_ge - 1:                          ns_golang.Ge,
	FuncID_gt - 1:                          ns_golang.Gt,
	FuncID_half - 1:                        ns_math.Half,
	FuncID_has - 1:                         ns_coll.HasAll,
	FuncID_hasAny - 1:                      ns_coll.HasAny,
	FuncID_hasPrefix - 1:                   ns_strings.HasPrefix,
	FuncID_hasSuffix - 1:                   ns_strings.HasSuffix,
	FuncID_hash - 1:                        get_ns_hash,
	FuncID_hash_ADLER32 - 1:                ns_hash.ADLER32,
	FuncID_hash_Bcrypt - 1:                 ns_hash.Bcrypt,
	FuncID_hash_CRC32 - 1:                  ns_hash.CRC32,
	FuncID_hash_CRC64 - 1:                  ns_hash.CRC64,
	FuncID_hash_MD4 - 1:                    ns_hash.MD4,
	FuncID_hash_MD5 - 1:                    ns_hash.MD5,
	FuncID_hash_RIPEMD160 - 1:              ns_hash.RIPEMD160,
	FuncID_hash_SHA1 - 1:                   ns_hash.SHA1,
	FuncID_hash_SHA224 - 1:                 ns_hash.SHA224,
	FuncID_hash_SHA256 - 1:                 ns_hash.SHA256,
	FuncID_hash_SHA384 - 1:                 ns_hash.SHA384,
	FuncID_hash_SHA512 - 1:                 ns_hash.SHA512,
	FuncID_hash_SHA512_224 - 1:             ns_hash.SHA512_224,
	FuncID_hash_SHA512_256 - 1:             ns_hash.SHA512_256,
	FuncID_hex - 1:                         ns_enc.Hex,
	FuncID_html - 1:                        ns_golang.HTMLEscaper,
	FuncID_indent - 1:                      ns_strings.Indent,
	FuncID_index - 1:                       ns_coll.Index,
	FuncID_js - 1:                          ns_golang.JSEscaper,
	FuncID_le - 1:                          ns_golang.Le,
	FuncID_len - 1:                         ns_golang.Length,
	FuncID_list - 1:                        ns_coll.List,
	FuncID_lower - 1:                       ns_strings.Lower,
	FuncID_lt - 1:                          ns_golang.Lt,
	FuncID_math - 1:                        get_ns_math,
	FuncID_math_Abs - 1:                    ns_math.Abs,
	FuncID_math_Add - 1:                    ns_math.Add,
	FuncID_math_Add1 - 1:                   ns_math.Add1,
	FuncID_math_Ceil - 1:                   ns_math.Ceil,
	FuncID_math_Div - 1:                    ns_math.Div,
	FuncID_math_Double - 1:                 ns_math.Double,
	FuncID_math_Floor - 1:                  ns_math.Floor,
	FuncID_math_Half - 1:                   ns_math.Half,
	FuncID_math_Log10 - 1:                  ns_math.Log10,
	FuncID_math_Log2 - 1:                   ns_math.Log2,
	FuncID_math_LogE - 1:                   ns_math.LogE,
	FuncID_math_Max - 1:                    ns_math.Max,
	FuncID_math_Min - 1:                    ns_math.Min,
	FuncID_math_Mod - 1:                    ns_math.Mod,
	FuncID_math_Mul - 1:                    ns_math.Mul,
	FuncID_math_Pow - 1:                    ns_math.Pow,
	FuncID_math_Round - 1:                  ns_math.Round,
	FuncID_math_Seq - 1:                    ns_math.Seq,
	FuncID_math_Sub - 1:                    ns_math.Sub,
	FuncID_math_Sub1 - 1:                   ns_math.Sub1,
	FuncID_max - 1:                         ns_math.Max,
	FuncID_md5 - 1:                         ns_hash.MD5,
	FuncID_min - 1:                         ns_math.Min,
	FuncID_mod - 1:                         ns_math.Mod,
	FuncID_mul - 1:                         ns_math.Mul,
	FuncID_ne - 1:                          ns_golang.Ne,
	FuncID_nindent - 1:                     ns_strings.NIndent,
	FuncID_not - 1:                         ns_golang.Not,
	FuncID_now - 1:                         ns_time.Now,
	FuncID_omit - 1:                        ns_coll.Omit,
	FuncID_or - 1:                          ns_golang.Or,
	FuncID_path - 1:                        get_ns_path,
	FuncID_path_Base - 1:                   ns_path.Base,
	FuncID_path_Clean - 1:                  ns_path.Clean,
	FuncID_path_Dir - 1:                    ns_path.Dir,
	FuncID_path_Ext - 1:                    ns_path.Ext,
	FuncID_path_IsAbs - 1:                  ns_path.IsAbs,
	FuncID_path_Join - 1:                   ns_path.Join,
	FuncID_path_Match - 1:                  ns_path.Match,
	FuncID_path_Split - 1:                  ns_path.Split,
	FuncID_pick - 1:                        ns_coll.Pick,
	FuncID_prepend - 1:                     ns_coll.Prepend,
	FuncID_print - 1:                       ns_golang.Sprint,
	FuncID_printf - 1:                      ns_golang.Sprintf,
	FuncID_println - 1:                     ns_golang.Sprintln,
	FuncID_quote - 1:                       ns_strings.DoubleQuote,
	FuncID_re - 1:                          get_ns_re,
	FuncID_re_Find - 1:                     ns_re.Find,
	FuncID_re_FindAll - 1:                  ns_re.FindAll,
	FuncID_re_Match - 1:                    ns_re.Match,
	FuncID_re_QuoteMeta - 1:                ns_re.QuoteMeta,
	FuncID_re_Replace - 1:                  ns_re.Replace,
	FuncID_re_ReplaceLiteral - 1:           ns_re.ReplaceLiteral,
	FuncID_re_Split - 1:                    ns_re.Split,
	FuncID_removePrefix - 1:                ns_strings.RemovePrefix,
	FuncID_removeSuffix - 1:                ns_strings.RemoveSuffix,
	FuncID_replaceAll - 1:                  ns_strings.ReplaceAll,
	FuncID_seq - 1:                         ns_math.Seq,
	FuncID_sha1 - 1:                        ns_hash.SHA1,
	FuncID_sha256 - 1:                      ns_hash.SHA256,
	FuncID_sha512 - 1:                      ns_hash.SHA512,
	FuncID_slice - 1:                       ns_coll.Slice,
	FuncID_sockaddr - 1:                    get_ns_sockaddr,
	FuncID_sockaddr_AllInterfaces - 1:      ns_sockaddr.AllInterfaces,
	FuncID_sockaddr_Attr - 1:               ns_sockaddr.Attr,
	FuncID_sockaddr_DefaultInterfaces - 1:  ns_sockaddr.DefaultInterfaces,
	FuncID_sockaddr_Exclude - 1:            ns_sockaddr.Exclude,
	FuncID_sockaddr_Include - 1:            ns_sockaddr.Include,
	FuncID_sockaddr_InterfaceIP - 1:        ns_sockaddr.InterfaceIP,
	FuncID_sockaddr_Join - 1:               ns_sockaddr.Join,
	FuncID_sockaddr_Limit - 1:              ns_sockaddr.Limit,
	FuncID_sockaddr_Math - 1:               ns_sockaddr.Math,
	FuncID_sockaddr_Offset - 1:             ns_sockaddr.Offset,
	FuncID_sockaddr_PrivateIP - 1:          ns_sockaddr.PrivateIP,
	FuncID_sockaddr_PrivateInterfaces - 1:  ns_sockaddr.PrivateInterfaces,
	FuncID_sockaddr_PublicIP - 1:           ns_sockaddr.PublicIP,
	FuncID_sockaddr_PublicInterfaces - 1:   ns_sockaddr.PublicInterfaces,
	FuncID_sockaddr_Sort - 1:               ns_sockaddr.Sort,
	FuncID_sockaddr_Unique - 1:             ns_sockaddr.Unique,
	FuncID_sort - 1:                        ns_coll.Sort,
	FuncID_split - 1:                       ns_strings.Split,
	FuncID_splitN - 1:                      ns_strings.SplitN,
	FuncID_squote - 1:                      ns_strings.SingleQuote,
	FuncID_stringList - 1:                  ns_coll.Strings,
	FuncID_strings - 1:                     get_ns_strings,
	FuncID_strings_Abbrev - 1:              ns_strings.Abbrev,
	FuncID_strings_AddPrefix - 1:           ns_strings.AddPrefix,
	FuncID_strings_AddSuffix - 1:           ns_strings.AddSuffix,
	FuncID_strings_CamelCase - 1:           ns_strings.CamelCase,
	FuncID_strings_Contains - 1:            ns_strings.Contains,
	FuncID_strings_ContainsAny - 1:         ns_strings.ContainsAny,
	FuncID_strings_DoubleQuote - 1:         ns_strings.DoubleQuote,
	FuncID_strings_HasPrefix - 1:           ns_strings.HasPrefix,
	FuncID_strings_HasSuffix - 1:           ns_strings.HasSuffix,
	FuncID_strings_Indent - 1:              ns_strings.Indent,
	FuncID_strings_Initials - 1:            ns_strings.Initials,
	FuncID_strings_Join - 1:                ns_strings.Join,
	FuncID_strings_KebabCase - 1:           ns_strings.KebabCase,
	FuncID_strings_Lower - 1:               ns_strings.Lower,
	FuncID_strings_NIndent - 1:             ns_strings.NIndent,
	FuncID_strings_NoSpace - 1:             ns_strings.NoSpace,
	FuncID_strings_RemovePrefix - 1:        ns_strings.RemovePrefix,
	FuncID_strings_RemoveSuffix - 1:        ns_strings.RemoveSuffix,
	FuncID_strings_Repeat - 1:              ns_strings.Repeat,
	FuncID_strings_ReplaceAll - 1:          ns_strings.ReplaceAll,
	FuncID_strings_RuneCount - 1:           ns_strings.RuneCount,
	FuncID_strings_ShellQuote - 1:          ns_strings.ShellQuote,
	FuncID_strings_Shuffle - 1:             ns_strings.Shuffle,
	FuncID_strings_SingleQuote - 1:         ns_strings.SingleQuote,
	FuncID_strings_Slug - 1:                ns_strings.Slug,
	FuncID_strings_SnakeCase - 1:           ns_strings.SnakeCase,
	FuncID_strings_Split - 1:               ns_strings.Split,
	FuncID_strings_SplitN - 1:              ns_strings.SplitN,
	FuncID_strings_Substr - 1:              ns_strings.Substr,
	FuncID_strings_SwapCase - 1:            ns_strings.SwapCase,
	FuncID_strings_Title - 1:               ns_strings.Title,
	FuncID_strings_Trim - 1:                ns_strings.Trim,
	FuncID_strings_TrimLeft - 1:            ns_strings.TrimLeft,
	FuncID_strings_TrimPrefix - 1:          ns_strings.TrimPrefix,
	FuncID_strings_TrimRight - 1:           ns_strings.TrimRight,
	FuncID_strings_TrimSpace - 1:           ns_strings.TrimSpace,
	FuncID_strings_TrimSuffix - 1:          ns_strings.TrimSuffix,
	FuncID_strings_Unquote - 1:             ns_strings.Unquote,
	FuncID_strings_Untitle - 1:             ns_strings.Untitle,
	FuncID_strings_Upper - 1:               ns_strings.Upper,
	FuncID_strings_WordWrap - 1:            ns_strings.WordWrap,
	FuncID_sub - 1:                         ns_math.Sub,
	FuncID_sub1 - 1:                        ns_math.Sub1,
	FuncID_time - 1:                        get_ns_time,
	FuncID_time_Add - 1:                    ns_time.Add,
	FuncID_time_Ceil - 1:                   ns_time.Ceil,
	FuncID_time_CeilDuration - 1:           ns_time.CeilDuration,
	FuncID_time_Day - 1:                    ns_time.Day,
	FuncID_time_FMT_ANSI - 1:               ns_time.FMT_ANSI,
	FuncID_time_FMT_Clock - 1:              ns_time.FMT_Clock,
	FuncID_time_FMT_Date - 1:               ns_time.FMT_Date,
	FuncID_time_FMT_DateTime - 1:           ns_time.FMT_DateTime,
	FuncID_time_FMT_RFC3339 - 1:            ns_time.FMT_RFC3339,
	FuncID_time_FMT_RFC3339Nano - 1:        ns_time.FMT_RFC3339Nano,
	FuncID_time_FMT_Ruby - 1:               ns_time.FMT_Ruby,
	FuncID_time_FMT_Stamp - 1:              ns_time.FMT_Stamp,
	FuncID_time_FMT_Unix - 1:               ns_time.FMT_Unix,
	FuncID_time_Floor - 1:                  ns_time.Floor,
	FuncID_time_FloorDuration - 1:          ns_time.FloorDuration,
	FuncID_time_Format - 1:                 ns_time.Format,
	FuncID_time_Hour - 1:                   ns_time.Hour,
	FuncID_time_Microsecond - 1:            ns_time.Microsecond,
	FuncID_time_Millisecond - 1:            ns_time.Millisecond,
	FuncID_time_Minute - 1:                 ns_time.Minute,
	FuncID_time_Nanosecond - 1:             ns_time.Nanosecond,
	FuncID_time_Now - 1:                    ns_time.Now,
	FuncID_time_Parse - 1:                  ns_time.Parse,
	FuncID_time_ParseDuration - 1:          ns_time.ParseDuration,
	FuncID_time_Round - 1:                  ns_time.Round,
	FuncID_time_RoundDuration - 1:          ns_time.RoundDuration,
	FuncID_time_Second - 1:                 ns_time.Second,
	FuncID_time_Since - 1:                  ns_time.Since,
	FuncID_time_Until - 1:                  ns_time.Until,
	FuncID_time_Week - 1:                   ns_time.Week,
	FuncID_time_ZoneName - 1:               ns_time.ZoneName,
	FuncID_time_ZoneOffset - 1:             ns_time.ZoneOffset,
	FuncID_title - 1:                       ns_strings.Title,
	FuncID_toJson - 1:                      ns_enc.JSON,
	FuncID_toString - 1:                    ns_type.ToString,
	FuncID_toYaml - 1:                      ns_enc.YAML,
	FuncID_totp - 1:                        ns_cred.Totp,
	FuncID_trim - 1:                        ns_strings.Trim,
	FuncID_trimPrefix - 1:                  ns_strings.TrimPrefix,
	FuncID_trimSpace - 1:                   ns_strings.TrimSpace,
	FuncID_trimSuffix - 1:                  ns_strings.TrimSuffix,
	FuncID_type - 1:                        get_ns_type,
	FuncID_type_AllTrue - 1:                ns_type.AllTrue,
	FuncID_type_AnyTrue - 1:                ns_type.AnyTrue,
	FuncID_type_Close - 1:                  ns_type.Close,
	FuncID_type_Default - 1:                ns_type.Default,
	FuncID_type_FirstNoneZero - 1:          ns_type.FirstNoneZero,
	FuncID_type_IsBool - 1:                 ns_type.IsBool,
	FuncID_type_IsFloat - 1:                ns_type.IsFloat,
	FuncID_type_IsInt - 1:                  ns_type.IsInt,
	FuncID_type_IsNum - 1:                  ns_type.IsNum,
	FuncID_type_IsZero - 1:                 ns_type.IsZero,
	FuncID_type_ToBool - 1:                 ns_type.ToBool,
	FuncID_type_ToFloat - 1:                ns_type.ToFloat,
	FuncID_type_ToInt - 1:                  ns_type.ToInt,
	FuncID_type_ToString - 1:               ns_type.ToString,
	FuncID_type_ToStrings - 1:              ns_type.ToStrings,
	FuncID_type_ToUint - 1:                 ns_type.ToUint,
	FuncID_uniq - 1:                        ns_coll.Unique,
	FuncID_upper - 1:                       ns_strings.Upper,
	FuncID_urlquery - 1:                    ns_golang.URLQueryEscaper,
	FuncID_uuid - 1:                        get_ns_uuid,
	FuncID_uuid_IsValid - 1:                ns_uuid.IsValid,
	FuncID_uuid_New - 1:                    ns_uuid.New,
	FuncID_uuid_V1 - 1:                     ns_uuid.V1,
	FuncID_uuid_V4 - 1:                     ns_uuid.V4,
	FuncID_uuid_Zero - 1:                   ns_uuid.Zero,
}

func createContextualFuncs(rc dukkha.RenderingContext) *ContextualFuncs {
	var (
		ns_dukkha = createDukkhaNS(rc)
		ns_fs     = createFSNS(rc)
		ns_os     = createOSNS(rc)
		ns_eval   = createEvalNS(rc)
		ns_tag    = createTagNS(rc)
		ns_state  = createStateNS(rc)
		ns_misc   = createMiscNS(rc)
	)

	get_ns_dukkha := func() dukkhaNS { return ns_dukkha }
	get_ns_fs := func() fsNS { return ns_fs }
	get_ns_os := func() osNS { return ns_os }
	get_ns_eval := func() evalNS { return ns_eval }
	get_ns_tag := func() tagNS { return ns_tag }
	get_ns_state := func() stateNS { return ns_state }

	return &ContextualFuncs{
		FuncID_VALUE - FuncID_LAST_Static_FUNC - 1:                reflect.ValueOf(ns_misc.VALUE),
		FuncID_dukkha - FuncID_LAST_Static_FUNC - 1:               reflect.ValueOf(get_ns_dukkha),
		FuncID_dukkha_CacheDir - FuncID_LAST_Static_FUNC - 1:      reflect.ValueOf(ns_dukkha.CacheDir),
		FuncID_dukkha_CrossPlatform - FuncID_LAST_Static_FUNC - 1: reflect.ValueOf(ns_dukkha.CrossPlatform),
		FuncID_dukkha_FromJson - FuncID_LAST_Static_FUNC - 1:      reflect.ValueOf(ns_dukkha.FromJson),
		FuncID_dukkha_FromYaml - FuncID_LAST_Static_FUNC - 1:      reflect.ValueOf(ns_dukkha.FromYaml),
		FuncID_dukkha_JQ - FuncID_LAST_Static_FUNC - 1:            reflect.ValueOf(ns_dukkha.JQ),
		FuncID_dukkha_JQObj - FuncID_LAST_Static_FUNC - 1:         reflect.ValueOf(ns_dukkha.JQObj),
		FuncID_dukkha_Self - FuncID_LAST_Static_FUNC - 1:          reflect.ValueOf(ns_dukkha.Self),
		FuncID_dukkha_Set - FuncID_LAST_Static_FUNC - 1:           reflect.ValueOf(ns_dukkha.Set),
		FuncID_dukkha_SetValue - FuncID_LAST_Static_FUNC - 1:      reflect.ValueOf(ns_dukkha.SetValue),
		FuncID_dukkha_WorkDir - FuncID_LAST_Static_FUNC - 1:       reflect.ValueOf(ns_dukkha.WorkDir),
		FuncID_dukkha_YQ - FuncID_LAST_Static_FUNC - 1:            reflect.ValueOf(ns_dukkha.YQ),
		FuncID_dukkha_YQObj - FuncID_LAST_Static_FUNC - 1:         reflect.ValueOf(ns_dukkha.YQObj),
		FuncID_env - FuncID_LAST_Static_FUNC - 1:                  reflect.ValueOf(ns_misc.Env),
		FuncID_eval - FuncID_LAST_Static_FUNC - 1:                 reflect.ValueOf(get_ns_eval),
		FuncID_eval_Env - FuncID_LAST_Static_FUNC - 1:             reflect.ValueOf(ns_eval.Env),
		FuncID_eval_Shell - FuncID_LAST_Static_FUNC - 1:           reflect.ValueOf(ns_eval.Shell),
		FuncID_eval_Template - FuncID_LAST_Static_FUNC - 1:        reflect.ValueOf(ns_eval.Template),
		FuncID_find - FuncID_LAST_Static_FUNC - 1:                 reflect.ValueOf(ns_fs.Find),
		FuncID_fromJson - FuncID_LAST_Static_FUNC - 1:             reflect.ValueOf(ns_dukkha.FromJson),
		FuncID_fromYaml - FuncID_LAST_Static_FUNC - 1:             reflect.ValueOf(ns_dukkha.FromYaml),
		FuncID_fs - FuncID_LAST_Static_FUNC - 1:                   reflect.ValueOf(get_ns_fs),
		FuncID_fs_Abs - FuncID_LAST_Static_FUNC - 1:               reflect.ValueOf(ns_fs.Abs),
		FuncID_fs_AppendFile - FuncID_LAST_Static_FUNC - 1:        reflect.ValueOf(ns_fs.AppendFile),
		FuncID_fs_Base - FuncID_LAST_Static_FUNC - 1:              reflect.ValueOf(ns_fs.Base),
		FuncID_fs_Clean - FuncID_LAST_Static_FUNC - 1:             reflect.ValueOf(ns_fs.Clean),
		FuncID_fs_Dir - FuncID_LAST_Static_FUNC - 1:               reflect.ValueOf(ns_fs.Dir),
		FuncID_fs_Exists - FuncID_LAST_Static_FUNC - 1:            reflect.ValueOf(ns_fs.Exists),
		FuncID_fs_Ext - FuncID_LAST_Static_FUNC - 1:               reflect.ValueOf(ns_fs.Ext),
		FuncID_fs_Find - FuncID_LAST_Static_FUNC - 1:              reflect.ValueOf(ns_fs.Find),
		FuncID_fs_FromSlash - FuncID_LAST_Static_FUNC - 1:         reflect.ValueOf(ns_fs.FromSlash),
		FuncID_fs_Glob - FuncID_LAST_Static_FUNC - 1:              reflect.ValueOf(ns_fs.Glob),
		FuncID_fs_IsAbs - FuncID_LAST_Static_FUNC - 1:             reflect.ValueOf(ns_fs.IsAbs),
		FuncID_fs_IsCharDevice - FuncID_LAST_Static_FUNC - 1:      reflect.ValueOf(ns_fs.IsCharDevice),
		FuncID_fs_IsDevice - FuncID_LAST_Static_FUNC - 1:          reflect.ValueOf(ns_fs.IsDevice),
		FuncID_fs_IsDir - FuncID_LAST_Static_FUNC - 1:             reflect.ValueOf(ns_fs.IsDir),
		FuncID_fs_IsFIFO - FuncID_LAST_Static_FUNC - 1:            reflect.ValueOf(ns_fs.IsFIFO),
		FuncID_fs_IsOther - FuncID_LAST_Static_FUNC - 1:           reflect.ValueOf(ns_fs.IsOther),
		FuncID_fs_IsSocket - FuncID_LAST_Static_FUNC - 1:          reflect.ValueOf(ns_fs.IsSocket),
		FuncID_fs_IsSymlink - FuncID_LAST_Static_FUNC - 1:         reflect.ValueOf(ns_fs.IsSymlink),
		FuncID_fs_Join - FuncID_LAST_Static_FUNC - 1:              reflect.ValueOf(ns_fs.Join),
		FuncID_fs_Lookup - FuncID_LAST_Static_FUNC - 1:            reflect.ValueOf(ns_fs.Lookup),
		FuncID_fs_LookupFile - FuncID_LAST_Static_FUNC - 1:        reflect.ValueOf(ns_fs.LookupFile),
		FuncID_fs_Match - FuncID_LAST_Static_FUNC - 1:             reflect.ValueOf(ns_fs.Match),
		FuncID_fs_Mkdir - FuncID_LAST_Static_FUNC - 1:             reflect.ValueOf(ns_fs.Mkdir),
		FuncID_fs_OpenFile - FuncID_LAST_Static_FUNC - 1:          reflect.ValueOf(ns_fs.OpenFile),
		FuncID_fs_ReadDir - FuncID_LAST_Static_FUNC - 1:           reflect.ValueOf(ns_fs.ReadDir),
		FuncID_fs_ReadFile - FuncID_LAST_Static_FUNC - 1:          reflect.ValueOf(ns_fs.ReadFile),
		FuncID_fs_Rel - FuncID_LAST_Static_FUNC - 1:               reflect.ValueOf(ns_fs.Rel),
		FuncID_fs_Split - FuncID_LAST_Static_FUNC - 1:             reflect.ValueOf(ns_fs.Split),
		FuncID_fs_ToSlash - FuncID_LAST_Static_FUNC - 1:           reflect.ValueOf(ns_fs.ToSlash),
		FuncID_fs_Touch - FuncID_LAST_Static_FUNC - 1:             reflect.ValueOf(ns_fs.Touch),
		FuncID_fs_UserCacheDir - FuncID_LAST_Static_FUNC - 1:      reflect.ValueOf(ns_fs.UserCacheDir),
		FuncID_fs_UserConfigDir - FuncID_LAST_Static_FUNC - 1:     reflect.ValueOf(ns_fs.UserConfigDir),
		FuncID_fs_UserHomeDir - FuncID_LAST_Static_FUNC - 1:       reflect.ValueOf(ns_fs.UserHomeDir),
		FuncID_fs_VolumeName - FuncID_LAST_Static_FUNC - 1:        reflect.ValueOf(ns_fs.VolumeName),
		FuncID_fs_WriteFile - FuncID_LAST_Static_FUNC - 1:         reflect.ValueOf(ns_fs.WriteFile),
		FuncID_git - FuncID_LAST_Static_FUNC - 1:                  reflect.ValueOf(ns_misc.Git),
		FuncID_host - FuncID_LAST_Static_FUNC - 1:                 reflect.ValueOf(ns_misc.Host),
		FuncID_jq - FuncID_LAST_Static_FUNC - 1:                   reflect.ValueOf(ns_dukkha.JQ),
		FuncID_jqObj - FuncID_LAST_Static_FUNC - 1:                reflect.ValueOf(ns_dukkha.JQObj),
		FuncID_matrix - FuncID_LAST_Static_FUNC - 1:               reflect.ValueOf(ns_misc.Matrix),
		FuncID_mkdir - FuncID_LAST_Static_FUNC - 1:                reflect.ValueOf(ns_fs.Mkdir),
		FuncID_os - FuncID_LAST_Static_FUNC - 1:                   reflect.ValueOf(get_ns_os),
		FuncID_os_Stderr - FuncID_LAST_Static_FUNC - 1:            reflect.ValueOf(ns_os.Stderr),
		FuncID_os_Stdin - FuncID_LAST_Static_FUNC - 1:             reflect.ValueOf(ns_os.Stdin),
		FuncID_os_Stdout - FuncID_LAST_Static_FUNC - 1:            reflect.ValueOf(ns_os.Stdout),
		FuncID_state - FuncID_LAST_Static_FUNC - 1:                reflect.ValueOf(get_ns_state),
		FuncID_state_Failed - FuncID_LAST_Static_FUNC - 1:         reflect.ValueOf(ns_state.Failed),
		FuncID_state_Succeeded - FuncID_LAST_Static_FUNC - 1:      reflect.ValueOf(ns_state.Succeeded),
		FuncID_tag - FuncID_LAST_Static_FUNC - 1:                  reflect.ValueOf(get_ns_tag),
		FuncID_tag_ImageName - FuncID_LAST_Static_FUNC - 1:        reflect.ValueOf(ns_tag.ImageName),
		FuncID_tag_ImageTag - FuncID_LAST_Static_FUNC - 1:         reflect.ValueOf(ns_tag.ImageTag),
		FuncID_tag_ManifestName - FuncID_LAST_Static_FUNC - 1:     reflect.ValueOf(ns_tag.ManifestName),
		FuncID_tag_ManifestTag - FuncID_LAST_Static_FUNC - 1:      reflect.ValueOf(ns_tag.ManifestTag),
		FuncID_touch - FuncID_LAST_Static_FUNC - 1:                reflect.ValueOf(ns_fs.Touch),
		FuncID_values - FuncID_LAST_Static_FUNC - 1:               reflect.ValueOf(ns_misc.Values),
		FuncID_write - FuncID_LAST_Static_FUNC - 1:                reflect.ValueOf(ns_fs.WriteFile),
		FuncID_yq - FuncID_LAST_Static_FUNC - 1:                   reflect.ValueOf(ns_dukkha.YQ),
		FuncID_yqObj - FuncID_LAST_Static_FUNC - 1:                reflect.ValueOf(ns_dukkha.YQObj),
	}
}
