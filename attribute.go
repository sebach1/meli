package meli

type Attribute struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`

	AttributeGroupId   string `json:"attribute_group_id,omitempty"`
	AttributeGroupName string `json:"attribute_group_name,omitempty"`

	Tags []Tag `json:"tags,omitempty"`

	Values         []*Value `json:"values,omitempty"`
	ValueName      string   `json:"value_name,omitempty"`
	ValueMaxLength int      `json:"value_max_length,omitempty"`
	ValueType      string   `json:"value_type,omitempty"`
	ValueId        string   `json:"value_id,omitempty"`
}

// Value is the entity which represents the possible option values for the parent attr
type Value struct {
	Id     string      `json:"id,omitempty"`
	Name   string      `json:"name,omitempty"`
	Struct interface{} `json:"struct,omitempty"`
}

type Tag map[string]bool

func (attr *Attribute) tagValue(tagName string) bool {
	for _, tag := range attr.Tags {
		if val, ok := tag[tagName]; ok {
			return val
		}
	}
	return false
}
