package templateprovider

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
