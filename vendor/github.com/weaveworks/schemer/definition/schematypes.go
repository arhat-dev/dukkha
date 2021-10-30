package definition

//nolint:golint,goconst
func setTypeOrRef(def *Definition, typeName string) {
	switch typeName {
	case "string":
		def.Type = "string"
	case "bool":
		def.Type = "boolean"
	case "int", "int8", "int16", "int64", "int32",
		"uint", "uint16", "uint32", "uint64", "uintptr":
		def.Type = "integer"
	case "float64", "float32":
		def.Type = "number"
	case "byte":
		def.Type = "string"
		def.ContentEncoding = "base64"
	default:
		def.Ref = DefPrefix + typeName
	}
}

func setDefaultForNonPointerType(def *Definition, typeName string) {
	// It only really makes sense to set default for bools
	// For strings or numbers, the empty value typically has a
	// different semantic meaning
	switch typeName {
	case "bool":
		def.Default = "false"
	}
}
