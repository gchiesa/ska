package templateprovider

func ByType(templateType TemplateType, name string) TemplateService {
	var ts TemplateService
	switch templateType {
	case SprigTemplateService:
		ts = NewSprigTemplate(name)
	case JinjaTemplateService:
		ts = NewJinjaTemplate(name)
	}
	return ts
}
