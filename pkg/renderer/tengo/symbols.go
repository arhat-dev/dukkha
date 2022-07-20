// code generated by funcs_test.go, DO NOT EDIT.

package tengo

import (
	"github.com/d5/tengo/v2"

	tu "arhat.dev/dukkha/pkg/templateutils"
)

var static_symbols = [tu.FuncID_LAST_Static_FUNC + 1]tengo.Symbol{
	tu.FuncID_add:                         {Name: "add", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_add)},
	tu.FuncID_add1:                        {Name: "add1", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_add1)},
	tu.FuncID_addPrefix:                   {Name: "addPrefix", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_addPrefix)},
	tu.FuncID_addSuffix:                   {Name: "addSuffix", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_addSuffix)},
	tu.FuncID_all:                         {Name: "all", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_all)},
	tu.FuncID_and:                         {Name: "and", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_and)},
	tu.FuncID_any:                         {Name: "any", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_any)},
	tu.FuncID_append:                      {Name: "append", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_append)},
	tu.FuncID_archconv:                    {Name: "archconv", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_archconv)},
	tu.FuncID_archconv_AlpineArch:         {Name: "archconv.AlpineArch", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_archconv_AlpineArch)},
	tu.FuncID_archconv_AlpineTripleName:   {Name: "archconv.AlpineTripleName", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_archconv_AlpineTripleName)},
	tu.FuncID_archconv_DebianArch:         {Name: "archconv.DebianArch", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_archconv_DebianArch)},
	tu.FuncID_archconv_DebianTripleName:   {Name: "archconv.DebianTripleName", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_archconv_DebianTripleName)},
	tu.FuncID_archconv_DockerArch:         {Name: "archconv.DockerArch", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_archconv_DockerArch)},
	tu.FuncID_archconv_DockerArchVariant:  {Name: "archconv.DockerArchVariant", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_archconv_DockerArchVariant)},
	tu.FuncID_archconv_DockerHubArch:      {Name: "archconv.DockerHubArch", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_archconv_DockerHubArch)},
	tu.FuncID_archconv_DockerOS:           {Name: "archconv.DockerOS", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_archconv_DockerOS)},
	tu.FuncID_archconv_DockerPlatformArch: {Name: "archconv.DockerPlatformArch", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_archconv_DockerPlatformArch)},
	tu.FuncID_archconv_GNUArch:            {Name: "archconv.GNUArch", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_archconv_GNUArch)},
	tu.FuncID_archconv_GNUTripleName:      {Name: "archconv.GNUTripleName", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_archconv_GNUTripleName)},
	tu.FuncID_archconv_GolangArch:         {Name: "archconv.GolangArch", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_archconv_GolangArch)},
	tu.FuncID_archconv_GolangOS:           {Name: "archconv.GolangOS", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_archconv_GolangOS)},
	tu.FuncID_archconv_HF:                 {Name: "archconv.HF", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_archconv_HF)},
	tu.FuncID_archconv_HardFloadArch:      {Name: "archconv.HardFloadArch", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_archconv_HardFloadArch)},
	tu.FuncID_archconv_LLVMArch:           {Name: "archconv.LLVMArch", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_archconv_LLVMArch)},
	tu.FuncID_archconv_LLVMTripleName:     {Name: "archconv.LLVMTripleName", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_archconv_LLVMTripleName)},
	tu.FuncID_archconv_OciArch:            {Name: "archconv.OciArch", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_archconv_OciArch)},
	tu.FuncID_archconv_OciArchVariant:     {Name: "archconv.OciArchVariant", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_archconv_OciArchVariant)},
	tu.FuncID_archconv_OciOS:              {Name: "archconv.OciOS", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_archconv_OciOS)},
	tu.FuncID_archconv_QemuArch:           {Name: "archconv.QemuArch", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_archconv_QemuArch)},
	tu.FuncID_archconv_SF:                 {Name: "archconv.SF", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_archconv_SF)},
	tu.FuncID_archconv_SimpleArch:         {Name: "archconv.SimpleArch", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_archconv_SimpleArch)},
	tu.FuncID_archconv_SoftFloadArch:      {Name: "archconv.SoftFloadArch", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_archconv_SoftFloadArch)},
	tu.FuncID_base64:                      {Name: "base64", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_base64)},
	tu.FuncID_call:                        {Name: "call", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_call)},
	tu.FuncID_close:                       {Name: "close", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_close)},
	tu.FuncID_coll:                        {Name: "coll", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_coll)},
	tu.FuncID_coll_Append:                 {Name: "coll.Append", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_coll_Append)},
	tu.FuncID_coll_Bools:                  {Name: "coll.Bools", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_coll_Bools)},
	tu.FuncID_coll_Dup:                    {Name: "coll.Dup", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_coll_Dup)},
	tu.FuncID_coll_Flatten:                {Name: "coll.Flatten", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_coll_Flatten)},
	tu.FuncID_coll_Floats:                 {Name: "coll.Floats", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_coll_Floats)},
	tu.FuncID_coll_HasAll:                 {Name: "coll.HasAll", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_coll_HasAll)},
	tu.FuncID_coll_HasAny:                 {Name: "coll.HasAny", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_coll_HasAny)},
	tu.FuncID_coll_Index:                  {Name: "coll.Index", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_coll_Index)},
	tu.FuncID_coll_Ints:                   {Name: "coll.Ints", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_coll_Ints)},
	tu.FuncID_coll_Keys:                   {Name: "coll.Keys", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_coll_Keys)},
	tu.FuncID_coll_List:                   {Name: "coll.List", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_coll_List)},
	tu.FuncID_coll_MapAnyAny:              {Name: "coll.MapAnyAny", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_coll_MapAnyAny)},
	tu.FuncID_coll_MapStringAny:           {Name: "coll.MapStringAny", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_coll_MapStringAny)},
	tu.FuncID_coll_Merge:                  {Name: "coll.Merge", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_coll_Merge)},
	tu.FuncID_coll_Omit:                   {Name: "coll.Omit", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_coll_Omit)},
	tu.FuncID_coll_Pick:                   {Name: "coll.Pick", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_coll_Pick)},
	tu.FuncID_coll_Prepend:                {Name: "coll.Prepend", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_coll_Prepend)},
	tu.FuncID_coll_Push:                   {Name: "coll.Push", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_coll_Push)},
	tu.FuncID_coll_Reverse:                {Name: "coll.Reverse", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_coll_Reverse)},
	tu.FuncID_coll_Slice:                  {Name: "coll.Slice", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_coll_Slice)},
	tu.FuncID_coll_Sort:                   {Name: "coll.Sort", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_coll_Sort)},
	tu.FuncID_coll_Strings:                {Name: "coll.Strings", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_coll_Strings)},
	tu.FuncID_coll_Uints:                  {Name: "coll.Uints", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_coll_Uints)},
	tu.FuncID_coll_Unique:                 {Name: "coll.Unique", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_coll_Unique)},
	tu.FuncID_coll_Values:                 {Name: "coll.Values", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_coll_Values)},
	tu.FuncID_contains:                    {Name: "contains", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_contains)},
	tu.FuncID_cred:                        {Name: "cred", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_cred)},
	tu.FuncID_cred_Htpasswd:               {Name: "cred.Htpasswd", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_cred_Htpasswd)},
	tu.FuncID_cred_Totp:                   {Name: "cred.Totp", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_cred_Totp)},
	tu.FuncID_default:                     {Name: "default", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_default)},
	tu.FuncID_dict:                        {Name: "dict", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_dict)},
	tu.FuncID_div:                         {Name: "div", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_div)},
	tu.FuncID_dns:                         {Name: "dns", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_dns)},
	tu.FuncID_dns_CNAME:                   {Name: "dns.CNAME", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_dns_CNAME)},
	tu.FuncID_dns_HOST:                    {Name: "dns.HOST", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_dns_HOST)},
	tu.FuncID_dns_IP:                      {Name: "dns.IP", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_dns_IP)},
	tu.FuncID_dns_SRV:                     {Name: "dns.SRV", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_dns_SRV)},
	tu.FuncID_dns_TXT:                     {Name: "dns.TXT", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_dns_TXT)},
	tu.FuncID_double:                      {Name: "double", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_double)},
	tu.FuncID_dup:                         {Name: "dup", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_dup)},
	tu.FuncID_enc:                         {Name: "enc", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_enc)},
	tu.FuncID_enc_Base32:                  {Name: "enc.Base32", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_enc_Base32)},
	tu.FuncID_enc_Base64:                  {Name: "enc.Base64", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_enc_Base64)},
	tu.FuncID_enc_Hex:                     {Name: "enc.Hex", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_enc_Hex)},
	tu.FuncID_enc_JSON:                    {Name: "enc.JSON", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_enc_JSON)},
	tu.FuncID_enc_YAML:                    {Name: "enc.YAML", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_enc_YAML)},
	tu.FuncID_eq:                          {Name: "eq", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_eq)},
	tu.FuncID_ge:                          {Name: "ge", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_ge)},
	tu.FuncID_gt:                          {Name: "gt", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_gt)},
	tu.FuncID_half:                        {Name: "half", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_half)},
	tu.FuncID_has:                         {Name: "has", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_has)},
	tu.FuncID_hasAny:                      {Name: "hasAny", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_hasAny)},
	tu.FuncID_hasPrefix:                   {Name: "hasPrefix", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_hasPrefix)},
	tu.FuncID_hasSuffix:                   {Name: "hasSuffix", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_hasSuffix)},
	tu.FuncID_hash:                        {Name: "hash", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_hash)},
	tu.FuncID_hash_ADLER32:                {Name: "hash.ADLER32", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_hash_ADLER32)},
	tu.FuncID_hash_Bcrypt:                 {Name: "hash.Bcrypt", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_hash_Bcrypt)},
	tu.FuncID_hash_CRC32:                  {Name: "hash.CRC32", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_hash_CRC32)},
	tu.FuncID_hash_CRC64:                  {Name: "hash.CRC64", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_hash_CRC64)},
	tu.FuncID_hash_MD4:                    {Name: "hash.MD4", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_hash_MD4)},
	tu.FuncID_hash_MD5:                    {Name: "hash.MD5", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_hash_MD5)},
	tu.FuncID_hash_RIPEMD160:              {Name: "hash.RIPEMD160", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_hash_RIPEMD160)},
	tu.FuncID_hash_SHA1:                   {Name: "hash.SHA1", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_hash_SHA1)},
	tu.FuncID_hash_SHA224:                 {Name: "hash.SHA224", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_hash_SHA224)},
	tu.FuncID_hash_SHA256:                 {Name: "hash.SHA256", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_hash_SHA256)},
	tu.FuncID_hash_SHA384:                 {Name: "hash.SHA384", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_hash_SHA384)},
	tu.FuncID_hash_SHA512:                 {Name: "hash.SHA512", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_hash_SHA512)},
	tu.FuncID_hash_SHA512_224:             {Name: "hash.SHA512_224", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_hash_SHA512_224)},
	tu.FuncID_hash_SHA512_256:             {Name: "hash.SHA512_256", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_hash_SHA512_256)},
	tu.FuncID_hex:                         {Name: "hex", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_hex)},
	tu.FuncID_html:                        {Name: "html", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_html)},
	tu.FuncID_indent:                      {Name: "indent", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_indent)},
	tu.FuncID_index:                       {Name: "index", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_index)},
	tu.FuncID_js:                          {Name: "js", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_js)},
	tu.FuncID_le:                          {Name: "le", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_le)},
	tu.FuncID_len:                         {Name: "len", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_len)},
	tu.FuncID_list:                        {Name: "list", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_list)},
	tu.FuncID_lower:                       {Name: "lower", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_lower)},
	tu.FuncID_lt:                          {Name: "lt", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_lt)},
	tu.FuncID_math:                        {Name: "math", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_math)},
	tu.FuncID_math_Abs:                    {Name: "math.Abs", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_math_Abs)},
	tu.FuncID_math_Add:                    {Name: "math.Add", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_math_Add)},
	tu.FuncID_math_Add1:                   {Name: "math.Add1", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_math_Add1)},
	tu.FuncID_math_Ceil:                   {Name: "math.Ceil", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_math_Ceil)},
	tu.FuncID_math_Div:                    {Name: "math.Div", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_math_Div)},
	tu.FuncID_math_Double:                 {Name: "math.Double", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_math_Double)},
	tu.FuncID_math_Floor:                  {Name: "math.Floor", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_math_Floor)},
	tu.FuncID_math_Half:                   {Name: "math.Half", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_math_Half)},
	tu.FuncID_math_Log10:                  {Name: "math.Log10", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_math_Log10)},
	tu.FuncID_math_Log2:                   {Name: "math.Log2", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_math_Log2)},
	tu.FuncID_math_LogE:                   {Name: "math.LogE", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_math_LogE)},
	tu.FuncID_math_Max:                    {Name: "math.Max", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_math_Max)},
	tu.FuncID_math_Min:                    {Name: "math.Min", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_math_Min)},
	tu.FuncID_math_Mod:                    {Name: "math.Mod", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_math_Mod)},
	tu.FuncID_math_Mul:                    {Name: "math.Mul", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_math_Mul)},
	tu.FuncID_math_Pow:                    {Name: "math.Pow", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_math_Pow)},
	tu.FuncID_math_Round:                  {Name: "math.Round", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_math_Round)},
	tu.FuncID_math_Seq:                    {Name: "math.Seq", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_math_Seq)},
	tu.FuncID_math_Sub:                    {Name: "math.Sub", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_math_Sub)},
	tu.FuncID_math_Sub1:                   {Name: "math.Sub1", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_math_Sub1)},
	tu.FuncID_max:                         {Name: "max", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_max)},
	tu.FuncID_md5:                         {Name: "md5", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_md5)},
	tu.FuncID_min:                         {Name: "min", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_min)},
	tu.FuncID_mod:                         {Name: "mod", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_mod)},
	tu.FuncID_mul:                         {Name: "mul", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_mul)},
	tu.FuncID_ne:                          {Name: "ne", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_ne)},
	tu.FuncID_nindent:                     {Name: "nindent", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_nindent)},
	tu.FuncID_not:                         {Name: "not", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_not)},
	tu.FuncID_now:                         {Name: "now", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_now)},
	tu.FuncID_omit:                        {Name: "omit", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_omit)},
	tu.FuncID_or:                          {Name: "or", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_or)},
	tu.FuncID_path:                        {Name: "path", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_path)},
	tu.FuncID_path_Base:                   {Name: "path.Base", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_path_Base)},
	tu.FuncID_path_Clean:                  {Name: "path.Clean", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_path_Clean)},
	tu.FuncID_path_Dir:                    {Name: "path.Dir", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_path_Dir)},
	tu.FuncID_path_Ext:                    {Name: "path.Ext", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_path_Ext)},
	tu.FuncID_path_IsAbs:                  {Name: "path.IsAbs", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_path_IsAbs)},
	tu.FuncID_path_Join:                   {Name: "path.Join", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_path_Join)},
	tu.FuncID_path_Match:                  {Name: "path.Match", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_path_Match)},
	tu.FuncID_path_Split:                  {Name: "path.Split", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_path_Split)},
	tu.FuncID_pick:                        {Name: "pick", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_pick)},
	tu.FuncID_prepend:                     {Name: "prepend", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_prepend)},
	tu.FuncID_print:                       {Name: "print", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_print)},
	tu.FuncID_printf:                      {Name: "printf", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_printf)},
	tu.FuncID_println:                     {Name: "println", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_println)},
	tu.FuncID_quote:                       {Name: "quote", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_quote)},
	tu.FuncID_re:                          {Name: "re", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_re)},
	tu.FuncID_re_Find:                     {Name: "re.Find", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_re_Find)},
	tu.FuncID_re_FindAll:                  {Name: "re.FindAll", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_re_FindAll)},
	tu.FuncID_re_Match:                    {Name: "re.Match", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_re_Match)},
	tu.FuncID_re_QuoteMeta:                {Name: "re.QuoteMeta", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_re_QuoteMeta)},
	tu.FuncID_re_Replace:                  {Name: "re.Replace", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_re_Replace)},
	tu.FuncID_re_ReplaceLiteral:           {Name: "re.ReplaceLiteral", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_re_ReplaceLiteral)},
	tu.FuncID_re_Split:                    {Name: "re.Split", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_re_Split)},
	tu.FuncID_removePrefix:                {Name: "removePrefix", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_removePrefix)},
	tu.FuncID_removeSuffix:                {Name: "removeSuffix", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_removeSuffix)},
	tu.FuncID_replaceAll:                  {Name: "replaceAll", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_replaceAll)},
	tu.FuncID_seq:                         {Name: "seq", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_seq)},
	tu.FuncID_sha1:                        {Name: "sha1", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_sha1)},
	tu.FuncID_sha256:                      {Name: "sha256", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_sha256)},
	tu.FuncID_sha512:                      {Name: "sha512", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_sha512)},
	tu.FuncID_slice:                       {Name: "slice", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_slice)},
	tu.FuncID_sockaddr:                    {Name: "sockaddr", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_sockaddr)},
	tu.FuncID_sockaddr_AllInterfaces:      {Name: "sockaddr.AllInterfaces", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_sockaddr_AllInterfaces)},
	tu.FuncID_sockaddr_Attr:               {Name: "sockaddr.Attr", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_sockaddr_Attr)},
	tu.FuncID_sockaddr_DefaultInterfaces:  {Name: "sockaddr.DefaultInterfaces", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_sockaddr_DefaultInterfaces)},
	tu.FuncID_sockaddr_Exclude:            {Name: "sockaddr.Exclude", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_sockaddr_Exclude)},
	tu.FuncID_sockaddr_Include:            {Name: "sockaddr.Include", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_sockaddr_Include)},
	tu.FuncID_sockaddr_InterfaceIP:        {Name: "sockaddr.InterfaceIP", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_sockaddr_InterfaceIP)},
	tu.FuncID_sockaddr_Join:               {Name: "sockaddr.Join", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_sockaddr_Join)},
	tu.FuncID_sockaddr_Limit:              {Name: "sockaddr.Limit", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_sockaddr_Limit)},
	tu.FuncID_sockaddr_Math:               {Name: "sockaddr.Math", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_sockaddr_Math)},
	tu.FuncID_sockaddr_Offset:             {Name: "sockaddr.Offset", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_sockaddr_Offset)},
	tu.FuncID_sockaddr_PrivateIP:          {Name: "sockaddr.PrivateIP", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_sockaddr_PrivateIP)},
	tu.FuncID_sockaddr_PrivateInterfaces:  {Name: "sockaddr.PrivateInterfaces", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_sockaddr_PrivateInterfaces)},
	tu.FuncID_sockaddr_PublicIP:           {Name: "sockaddr.PublicIP", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_sockaddr_PublicIP)},
	tu.FuncID_sockaddr_PublicInterfaces:   {Name: "sockaddr.PublicInterfaces", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_sockaddr_PublicInterfaces)},
	tu.FuncID_sockaddr_Sort:               {Name: "sockaddr.Sort", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_sockaddr_Sort)},
	tu.FuncID_sockaddr_Unique:             {Name: "sockaddr.Unique", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_sockaddr_Unique)},
	tu.FuncID_sort:                        {Name: "sort", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_sort)},
	tu.FuncID_split:                       {Name: "split", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_split)},
	tu.FuncID_splitN:                      {Name: "splitN", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_splitN)},
	tu.FuncID_squote:                      {Name: "squote", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_squote)},
	tu.FuncID_stringList:                  {Name: "stringList", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_stringList)},
	tu.FuncID_strings:                     {Name: "strings", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings)},
	tu.FuncID_strings_Abbrev:              {Name: "strings.Abbrev", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_Abbrev)},
	tu.FuncID_strings_AddPrefix:           {Name: "strings.AddPrefix", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_AddPrefix)},
	tu.FuncID_strings_AddSuffix:           {Name: "strings.AddSuffix", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_AddSuffix)},
	tu.FuncID_strings_CamelCase:           {Name: "strings.CamelCase", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_CamelCase)},
	tu.FuncID_strings_Contains:            {Name: "strings.Contains", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_Contains)},
	tu.FuncID_strings_ContainsAny:         {Name: "strings.ContainsAny", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_ContainsAny)},
	tu.FuncID_strings_DoubleQuote:         {Name: "strings.DoubleQuote", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_DoubleQuote)},
	tu.FuncID_strings_HasPrefix:           {Name: "strings.HasPrefix", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_HasPrefix)},
	tu.FuncID_strings_HasSuffix:           {Name: "strings.HasSuffix", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_HasSuffix)},
	tu.FuncID_strings_Indent:              {Name: "strings.Indent", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_Indent)},
	tu.FuncID_strings_Initials:            {Name: "strings.Initials", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_Initials)},
	tu.FuncID_strings_Join:                {Name: "strings.Join", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_Join)},
	tu.FuncID_strings_KebabCase:           {Name: "strings.KebabCase", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_KebabCase)},
	tu.FuncID_strings_Lower:               {Name: "strings.Lower", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_Lower)},
	tu.FuncID_strings_NIndent:             {Name: "strings.NIndent", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_NIndent)},
	tu.FuncID_strings_NoSpace:             {Name: "strings.NoSpace", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_NoSpace)},
	tu.FuncID_strings_RemovePrefix:        {Name: "strings.RemovePrefix", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_RemovePrefix)},
	tu.FuncID_strings_RemoveSuffix:        {Name: "strings.RemoveSuffix", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_RemoveSuffix)},
	tu.FuncID_strings_Repeat:              {Name: "strings.Repeat", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_Repeat)},
	tu.FuncID_strings_ReplaceAll:          {Name: "strings.ReplaceAll", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_ReplaceAll)},
	tu.FuncID_strings_RuneCount:           {Name: "strings.RuneCount", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_RuneCount)},
	tu.FuncID_strings_ShellQuote:          {Name: "strings.ShellQuote", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_ShellQuote)},
	tu.FuncID_strings_Shuffle:             {Name: "strings.Shuffle", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_Shuffle)},
	tu.FuncID_strings_SingleQuote:         {Name: "strings.SingleQuote", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_SingleQuote)},
	tu.FuncID_strings_Slug:                {Name: "strings.Slug", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_Slug)},
	tu.FuncID_strings_SnakeCase:           {Name: "strings.SnakeCase", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_SnakeCase)},
	tu.FuncID_strings_Split:               {Name: "strings.Split", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_Split)},
	tu.FuncID_strings_SplitN:              {Name: "strings.SplitN", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_SplitN)},
	tu.FuncID_strings_Substr:              {Name: "strings.Substr", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_Substr)},
	tu.FuncID_strings_SwapCase:            {Name: "strings.SwapCase", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_SwapCase)},
	tu.FuncID_strings_Title:               {Name: "strings.Title", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_Title)},
	tu.FuncID_strings_Trim:                {Name: "strings.Trim", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_Trim)},
	tu.FuncID_strings_TrimLeft:            {Name: "strings.TrimLeft", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_TrimLeft)},
	tu.FuncID_strings_TrimPrefix:          {Name: "strings.TrimPrefix", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_TrimPrefix)},
	tu.FuncID_strings_TrimRight:           {Name: "strings.TrimRight", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_TrimRight)},
	tu.FuncID_strings_TrimSpace:           {Name: "strings.TrimSpace", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_TrimSpace)},
	tu.FuncID_strings_TrimSuffix:          {Name: "strings.TrimSuffix", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_TrimSuffix)},
	tu.FuncID_strings_Unquote:             {Name: "strings.Unquote", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_Unquote)},
	tu.FuncID_strings_Untitle:             {Name: "strings.Untitle", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_Untitle)},
	tu.FuncID_strings_Upper:               {Name: "strings.Upper", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_Upper)},
	tu.FuncID_strings_WordWrap:            {Name: "strings.WordWrap", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_strings_WordWrap)},
	tu.FuncID_sub:                         {Name: "sub", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_sub)},
	tu.FuncID_sub1:                        {Name: "sub1", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_sub1)},
	tu.FuncID_time:                        {Name: "time", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time)},
	tu.FuncID_time_Add:                    {Name: "time.Add", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time_Add)},
	tu.FuncID_time_Ceil:                   {Name: "time.Ceil", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time_Ceil)},
	tu.FuncID_time_CeilDuration:           {Name: "time.CeilDuration", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time_CeilDuration)},
	tu.FuncID_time_Day:                    {Name: "time.Day", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time_Day)},
	tu.FuncID_time_FMT_ANSI:               {Name: "time.FMT_ANSI", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time_FMT_ANSI)},
	tu.FuncID_time_FMT_Clock:              {Name: "time.FMT_Clock", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time_FMT_Clock)},
	tu.FuncID_time_FMT_Date:               {Name: "time.FMT_Date", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time_FMT_Date)},
	tu.FuncID_time_FMT_DateTime:           {Name: "time.FMT_DateTime", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time_FMT_DateTime)},
	tu.FuncID_time_FMT_RFC3339:            {Name: "time.FMT_RFC3339", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time_FMT_RFC3339)},
	tu.FuncID_time_FMT_RFC3339Nano:        {Name: "time.FMT_RFC3339Nano", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time_FMT_RFC3339Nano)},
	tu.FuncID_time_FMT_Ruby:               {Name: "time.FMT_Ruby", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time_FMT_Ruby)},
	tu.FuncID_time_FMT_Stamp:              {Name: "time.FMT_Stamp", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time_FMT_Stamp)},
	tu.FuncID_time_FMT_Unix:               {Name: "time.FMT_Unix", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time_FMT_Unix)},
	tu.FuncID_time_Floor:                  {Name: "time.Floor", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time_Floor)},
	tu.FuncID_time_FloorDuration:          {Name: "time.FloorDuration", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time_FloorDuration)},
	tu.FuncID_time_Format:                 {Name: "time.Format", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time_Format)},
	tu.FuncID_time_Hour:                   {Name: "time.Hour", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time_Hour)},
	tu.FuncID_time_Microsecond:            {Name: "time.Microsecond", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time_Microsecond)},
	tu.FuncID_time_Millisecond:            {Name: "time.Millisecond", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time_Millisecond)},
	tu.FuncID_time_Minute:                 {Name: "time.Minute", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time_Minute)},
	tu.FuncID_time_Nanosecond:             {Name: "time.Nanosecond", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time_Nanosecond)},
	tu.FuncID_time_Now:                    {Name: "time.Now", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time_Now)},
	tu.FuncID_time_Parse:                  {Name: "time.Parse", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time_Parse)},
	tu.FuncID_time_ParseDuration:          {Name: "time.ParseDuration", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time_ParseDuration)},
	tu.FuncID_time_Round:                  {Name: "time.Round", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time_Round)},
	tu.FuncID_time_RoundDuration:          {Name: "time.RoundDuration", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time_RoundDuration)},
	tu.FuncID_time_Second:                 {Name: "time.Second", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time_Second)},
	tu.FuncID_time_Since:                  {Name: "time.Since", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time_Since)},
	tu.FuncID_time_Until:                  {Name: "time.Until", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time_Until)},
	tu.FuncID_time_Week:                   {Name: "time.Week", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time_Week)},
	tu.FuncID_time_ZoneName:               {Name: "time.ZoneName", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time_ZoneName)},
	tu.FuncID_time_ZoneOffset:             {Name: "time.ZoneOffset", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_time_ZoneOffset)},
	tu.FuncID_title:                       {Name: "title", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_title)},
	tu.FuncID_toJson:                      {Name: "toJson", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_toJson)},
	tu.FuncID_toString:                    {Name: "toString", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_toString)},
	tu.FuncID_toYaml:                      {Name: "toYaml", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_toYaml)},
	tu.FuncID_totp:                        {Name: "totp", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_totp)},
	tu.FuncID_trim:                        {Name: "trim", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_trim)},
	tu.FuncID_trimPrefix:                  {Name: "trimPrefix", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_trimPrefix)},
	tu.FuncID_trimSpace:                   {Name: "trimSpace", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_trimSpace)},
	tu.FuncID_trimSuffix:                  {Name: "trimSuffix", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_trimSuffix)},
	tu.FuncID_type:                        {Name: "type", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_type)},
	tu.FuncID_type_AllTrue:                {Name: "type.AllTrue", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_type_AllTrue)},
	tu.FuncID_type_AnyTrue:                {Name: "type.AnyTrue", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_type_AnyTrue)},
	tu.FuncID_type_Close:                  {Name: "type.Close", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_type_Close)},
	tu.FuncID_type_Default:                {Name: "type.Default", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_type_Default)},
	tu.FuncID_type_FirstNoneZero:          {Name: "type.FirstNoneZero", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_type_FirstNoneZero)},
	tu.FuncID_type_IsBool:                 {Name: "type.IsBool", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_type_IsBool)},
	tu.FuncID_type_IsFloat:                {Name: "type.IsFloat", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_type_IsFloat)},
	tu.FuncID_type_IsInt:                  {Name: "type.IsInt", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_type_IsInt)},
	tu.FuncID_type_IsNum:                  {Name: "type.IsNum", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_type_IsNum)},
	tu.FuncID_type_IsZero:                 {Name: "type.IsZero", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_type_IsZero)},
	tu.FuncID_type_ToBool:                 {Name: "type.ToBool", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_type_ToBool)},
	tu.FuncID_type_ToFloat:                {Name: "type.ToFloat", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_type_ToFloat)},
	tu.FuncID_type_ToInt:                  {Name: "type.ToInt", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_type_ToInt)},
	tu.FuncID_type_ToString:               {Name: "type.ToString", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_type_ToString)},
	tu.FuncID_type_ToStrings:              {Name: "type.ToStrings", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_type_ToStrings)},
	tu.FuncID_type_ToUint:                 {Name: "type.ToUint", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_type_ToUint)},
	tu.FuncID_uniq:                        {Name: "uniq", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_uniq)},
	tu.FuncID_upper:                       {Name: "upper", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_upper)},
	tu.FuncID_urlquery:                    {Name: "urlquery", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_urlquery)},
	tu.FuncID_uuid:                        {Name: "uuid", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_uuid)},
	tu.FuncID_uuid_IsValid:                {Name: "uuid.IsValid", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_uuid_IsValid)},
	tu.FuncID_uuid_New:                    {Name: "uuid.New", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_uuid_New)},
	tu.FuncID_uuid_V1:                     {Name: "uuid.V1", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_uuid_V1)},
	tu.FuncID_uuid_V4:                     {Name: "uuid.V4", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_uuid_V4)},
	tu.FuncID_uuid_Zero:                   {Name: "uuid.Zero", Scope: tengo.ScopeGlobal, Index: int(tu.FuncID_uuid_Zero)},
}

var static_funcs = [...]tengo.UserFunction{
	// TODO: map template funcs to tengo funcs
	// {Name: "", Value: tengo.CallableFunc, Index: 0},
}
