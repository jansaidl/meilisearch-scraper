package src

import "encoding/xml"

type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	URLs    []URL    `xml:"url"`
}

type URL struct {
	Loc string `xml:"loc"`
}

type Config struct {
	Selectors struct {
		Lvl0 SelectorConfig `json:"lvl0"`
		Lvl1 string         `json:"lvl1"`
		Lvl2 string         `json:"lvl2"`
		Lvl3 string         `json:"lvl3"`
		Lvl4 string         `json:"lvl4"`
		Lvl5 string         `json:"lvl5"`
		Lvl6 string         `json:"lvl6"`
		Text string         `json:"text"`
	} `json:"selectors"`
}

type SelectorConfig struct {
	Selector     string `json:"selector"`
	Global       bool   `json:"global"`
	DefaultValue string `json:"default_value"`
}

type Document struct {
	Anchor             string  `json:"anchor"`
	Content            *string `json:"content"`
	URL                string  `json:"url"`
	ObjectID           string  `json:"objectID"`
	HierarchyLvl0      *string `json:"hierarchy_lvl0"`
	HierarchyLvl1      *string `json:"hierarchy_lvl1"`
	HierarchyLvl2      *string `json:"hierarchy_lvl2"`
	HierarchyLvl3      *string `json:"hierarchy_lvl3"`
	HierarchyLvl4      *string `json:"hierarchy_lvl4"`
	HierarchyLvl5      *string `json:"hierarchy_lvl5"`
	HierarchyLvl6      *string `json:"hierarchy_lvl6"`
	HierarchyRadioLvl0 *string `json:"hierarchy_radio_lvl0"`
	HierarchyRadioLvl1 *string `json:"hierarchy_radio_lvl1"`
	HierarchyRadioLvl2 *string `json:"hierarchy_radio_lvl2"`
	HierarchyRadioLvl3 *string `json:"hierarchy_radio_lvl3"`
	HierarchyRadioLvl4 *string `json:"hierarchy_radio_lvl4"`
	HierarchyRadioLvl5 *string `json:"hierarchy_radio_lvl5"`
}
