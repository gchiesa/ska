package templateprovider

// ByType returns a TemplateService implementation by TemplateType.
func ByType(templateType TemplateType, name string) TemplateService {
	var ts TemplateService
	switch templateType {
	case SprigTemplateType:
		ts = NewSprigTemplate(name)
	case JinjaTemplateType:
		ts = NewJinjaTemplate(name)
	}
	return ts
}

// GetTypeFromString parses a template type string ("sprig" or "jinja") and
// returns the corresponding TemplateType. It returns -1 for unknown values.
func GetTypeFromString(templateType string) TemplateType {
	switch templateType {
	case "sprig":
		return SprigTemplateType
	case "jinja":
		return JinjaTemplateType
	default:
		return -1
	}
}
