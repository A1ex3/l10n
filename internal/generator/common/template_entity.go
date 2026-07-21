package common

type LocalizationData struct {
	PackageName     string
	Namespace       string
	Languages       []Language
	ClassName       string
	BaseClassName   string
	DefaultLangCode string
	CurrentLangCode string
	Methods         []Method
	MethodsByName   map[string]Method
}

type Language struct {
	Name         string
	Code         string
	Translations map[string]string
}

type Method struct {
	Name        string
	Parameters  []Parameter
	Description string
	Example     string
}

type Parameter struct {
	Name string
	Type string // can be "string", "float", or "int"
}
