package mangodex

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ResponseType : Interface for API responses.
type ResponseType interface {
	GetResult() string
}

// Response : Plain response struct containing only Result field.
type Response struct {
	Result string `json:"result"`
}

func (r *Response) GetResult() string {
	return r.Result
}

// Relationship : Struct containing relationships, with optional attributes for the relation.
type Relationship struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"`
	Attributes interface{} `json:"attributes"`
}

func (a *Relationship) UnmarshalJSON(data []byte) error {
	// Check for the type of the relationship, then unmarshal accordingly.
	typ := struct {
		ID         string          `json:"id"`
		Type       string          `json:"type"`
		Attributes json.RawMessage `json:"attributes"`
	}{}
	if err := json.Unmarshal(data, &typ); err != nil {
		return err
	}

	var err error
	switch typ.Type {
	case MangaRel:
		a.Attributes = &MangaAttributes{}
	case AuthorRel:
		a.Attributes = &AuthorAttributes{}
	case ScanlationGroupRel:
		a.Attributes = &ScanlationGroupAttributes{}
	default:
		a.Attributes = &json.RawMessage{}
	}

	a.ID = typ.ID
	a.Type = typ.Type
	if typ.Attributes != nil {
		if err = json.Unmarshal(typ.Attributes, a.Attributes); err != nil {
			return fmt.Errorf("error unmarshalling relationship of type %s: %s, %s", 
				typ.Type, err.Error(), string(data))
		}
	}
	return err
}

// LocalisedStrings : A struct wrapping around a map containing each localised string.
type LocalisedStrings struct {
	Values map[string]string
}

func (l *LocalisedStrings) UnmarshalJSON(data []byte) error {
	l.Values = map[string]string{}

	// First try if can unmarshal directly.
	if err := json.Unmarshal(data, &l.Values); err == nil {
		return nil
	}

	// If fail, try to unmarshal to array of maps.
	var locals []map[string]string
	if err := json.Unmarshal(data, &locals); err != nil {
		return fmt.Errorf("error unmarshalling localisedstring: %s", err.Error())
	}

	// If pass, then add each item in the array to flatten to one map.
	for _, entry := range locals {
		for key, value := range entry {
			l.Values[key] = value
		}
	}
	return nil
}

// GetLocalString : Get the localised string for a particular language code.
// If the required string is not found, it will return the first entry, or an empty string otherwise.
func (l *LocalisedStrings) GetLocalString(langCode string) string {
	// If we cannot find the required code, then return first value.
	if s, ok := l.Values[langCode]; !ok {
		v := ""
		for _, value := range l.Values {
			v = value
			break
		}
		return v
	} else {
		return s
	}
}

// Tag : Struct containing information on a tag.
type Tag struct {
	ID            string         `json:"id"`
	Type          string         `json:"type"`
	Attributes    TagAttributes  `json:"attributes"`
	Relationships []Relationship `json:"relationships"`
}

// GetName : Get name of the tag.
func (t *Tag) GetName(langCode string) string {
	return t.Attributes.Name.GetLocalString(langCode)
}

// TagAttributes : Attributes for a Tag.
type TagAttributes struct {
	Name        LocalisedStrings `json:"name"`
	Description LocalisedStrings `json:"description"`
	Group       string           `json:"group"`
	Version     int              `json:"version"`
}

// ErrorResponse : Typical response for errored requests.
type ErrorResponse struct {
	Result string  `json:"result"`
	Errors []Error `json:"errors"`
}

func (er *ErrorResponse) GetResult() string {
	return er.Result
}

// GetErrors : Get the errors for this particular request.
func (er *ErrorResponse) GetErrors() string {
	var errors strings.Builder
	for _, err := range er.Errors {
		errors.WriteString(fmt.Sprintf("%s: %s\n", err.Title, err.Detail))
	}
	return errors.String()
}

// Error : Struct containing details of an error.
type Error struct {
	ID     string `json:"id"`
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}
