package types

import (
	"bytes"
	"embed"
	"text/template"
)

var templater *template.Template

//go:embed templates
var templateFS embed.FS

func init() {
	var err error
	templater, err = template.ParseFS(templateFS, "templates/*")
	if err != nil {
		panic(err)
	}
}

type Data struct {
	ActiveSkills ActiveSkills `json:"active_skills"`
}

type ActiveSkills map[string]ActiveSkill

type ActiveSkill struct {
	DisplayName        string            `json:"display_name"`
	Description        string            `json:"description"`
	Types              []string          `json:"types"`
	WeaponRestrictions []string          `json:"weapon_restrictions"`
	IsManuallyCasted   bool              `json:"is_manually_casted"`
	StatConversions    map[string]string `json:"stat_conversions"`
	Icon               string            `json:"icon"`
	StatTranslations   string            `json:"stat_translations"`
	Index              int               `json:"index"`
}

func (a ActiveSkill) MarshalText() ([]byte, error) {
	var b bytes.Buffer
	if err := templater.Execute(&b, a); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
