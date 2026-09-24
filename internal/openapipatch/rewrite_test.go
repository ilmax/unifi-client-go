package openapipatch

import (
	"encoding/json"
	"testing"
)

func TestRewriteDiscriminatorUnionsSimplifiesEntityMetadata(t *testing.T) {
	doc := []byte(`{
		"openapi": "3.1.0",
		"components": {
			"schemas": {
				"Entity metadata": {
					"type": "object",
					"discriminator": {"propertyName": "origin"}
				},
				"User defined entity metadata": {
					"allOf": [
						{"$ref": "#/components/schemas/Entity metadata"},
						{"$ref": "#/components/schemas/User or system defined entity metadata"}
					]
				},
				"User or system defined entity metadata": {
					"type": "object"
				}
			}
		}
	}`)

	patched, err := RewriteDiscriminatorUnions(doc)
	if err != nil {
		t.Fatalf("RewriteDiscriminatorUnions() error = %v", err)
	}

	var spec struct {
		Components struct {
			Schemas map[string]struct {
				AllOf []struct {
					Ref string `json:"$ref"`
				} `json:"allOf"`
			} `json:"schemas"`
		} `json:"components"`
	}
	if err := json.Unmarshal(patched, &spec); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	metadata := spec.Components.Schemas["User_defined_entity_metadata"]
	if len(metadata.AllOf) != 1 {
		t.Fatalf("User defined entity metadata allOf length = %d, want 1", len(metadata.AllOf))
	}
	if metadata.AllOf[0].Ref != "#/components/schemas/Entity_metadata" {
		t.Fatalf("User defined entity metadata ref = %q", metadata.AllOf[0].Ref)
	}
}
