package openapipatch

import (
	"encoding/json"
	"testing"
)

func TestRewriteDNSUnions(t *testing.T) {
	t.Parallel()

	patched, err := RewriteDNSUnions([]byte(testSpec))
	if err != nil {
		t.Fatalf("RewriteDNSUnions() error = %v", err)
	}

	var spec map[string]any
	if err := json.Unmarshal(patched, &spec); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	schemas, err := schemaObjects(spec)
	if err != nil {
		t.Fatalf("schemaObjects() error = %v", err)
	}

	assertUnionSchema(t, schemas[createOrUpdateDNSPolicySchema], []string{
		schemaRef("IntegrationDnsARecordCreateUpdateDto"),
		schemaRef("IntegrationDnsForwardDomainPolicyCreateUpdateDto"),
	})
	assertUnionSchema(t, schemas[dnsPolicySchema], []string{
		schemaRef("IntegrationDnsARecordDto"),
		schemaRef("IntegrationDnsForwardDomainPolicyDto"),
	})

	if _, ok := schemas[createOrUpdateDNSPolicyBase]["discriminator"]; ok {
		t.Fatalf("synthetic create/update base should not keep discriminator")
	}
	if _, ok := schemas[dnsPolicyBase]["discriminator"]; ok {
		t.Fatalf("synthetic response base should not keep discriminator")
	}

	assertAllOfParentRef(t, schemas["IntegrationDnsARecordCreateUpdateDto"], schemaRef(createOrUpdateDNSPolicyBase))
	assertAllOfParentRef(t, schemas["IntegrationDnsForwardDomainPolicyCreateUpdateDto"], schemaRef(createOrUpdateDNSPolicyBase))
	assertAllOfParentRef(t, schemas["IntegrationDnsARecordDto"], schemaRef(dnsPolicyBase))
	assertAllOfParentRef(t, schemas["IntegrationDnsForwardDomainPolicyDto"], schemaRef(dnsPolicyBase))
}

func assertUnionSchema(t *testing.T, schema map[string]any, wantRefs []string) {
	t.Helper()

	if schema == nil {
		t.Fatalf("schema is nil")
	}
	if _, ok := schema["properties"]; ok {
		t.Fatalf("union schema should not keep properties")
	}
	if _, ok := schema["required"]; ok {
		t.Fatalf("union schema should not keep required")
	}

	rawOneOf, ok := schema["oneOf"].([]any)
	if !ok {
		t.Fatalf("union schema is missing oneOf")
	}

	gotRefs := make([]string, 0, len(rawOneOf))
	for _, rawItem := range rawOneOf {
		item, ok := rawItem.(map[string]any)
		if !ok {
			t.Fatalf("oneOf item has invalid type %T", rawItem)
		}
		ref, ok := item["$ref"].(string)
		if !ok {
			t.Fatalf("oneOf item is missing $ref")
		}
		gotRefs = append(gotRefs, ref)
	}

	if len(gotRefs) != len(wantRefs) {
		t.Fatalf("oneOf length = %d, want %d", len(gotRefs), len(wantRefs))
	}
	for idx, wantRef := range wantRefs {
		if gotRefs[idx] != wantRef {
			t.Fatalf("oneOf[%d] = %q, want %q", idx, gotRefs[idx], wantRef)
		}
	}
}

func assertAllOfParentRef(t *testing.T, schema map[string]any, wantRef string) {
	t.Helper()

	rawAllOf, ok := schema["allOf"].([]any)
	if !ok {
		t.Fatalf("schema missing allOf")
	}
	if len(rawAllOf) == 0 {
		t.Fatalf("schema allOf is empty")
	}

	firstItem, ok := rawAllOf[0].(map[string]any)
	if !ok {
		t.Fatalf("allOf[0] has invalid type %T", rawAllOf[0])
	}
	gotRef, ok := firstItem["$ref"].(string)
	if !ok {
		t.Fatalf("allOf[0] is missing $ref")
	}
	if gotRef != wantRef {
		t.Fatalf("allOf[0].$ref = %q, want %q", gotRef, wantRef)
	}
}

const testSpec = `{
  "openapi": "3.1.0",
  "components": {
    "schemas": {
      "Create or update DNS policy": {
        "type": "object",
        "discriminator": {
          "propertyName": "type",
          "mapping": {
            "A_RECORD": "#/components/schemas/IntegrationDnsARecordCreateUpdateDto",
            "FORWARD_DOMAIN": "#/components/schemas/IntegrationDnsForwardDomainPolicyCreateUpdateDto"
          }
        },
        "properties": {
          "type": {
            "type": "string"
          },
          "enabled": {
            "type": "boolean"
          }
        },
        "required": [
          "enabled",
          "type"
        ]
      },
      "IntegrationDnsARecordCreateUpdateDto": {
        "allOf": [
          {
            "$ref": "#/components/schemas/Create or update DNS policy"
          },
          {
            "type": "object",
            "properties": {
              "name": {
                "type": "string"
              }
            }
          }
        ]
      },
      "IntegrationDnsForwardDomainPolicyCreateUpdateDto": {
        "allOf": [
          {
            "$ref": "#/components/schemas/Create or update DNS policy"
          },
          {
            "type": "object",
            "properties": {
              "domain": {
                "type": "string"
              }
            }
          }
        ]
      },
      "DNS policy": {
        "type": "object",
        "discriminator": {
          "propertyName": "type",
          "mapping": {
            "A_RECORD": "#/components/schemas/IntegrationDnsARecordDto",
            "FORWARD_DOMAIN": "#/components/schemas/IntegrationDnsForwardDomainPolicyDto"
          }
        },
        "properties": {
          "type": {
            "type": "string"
          },
          "enabled": {
            "type": "boolean"
          }
        },
        "required": [
          "enabled",
          "type"
        ]
      },
      "IntegrationDnsARecordDto": {
        "allOf": [
          {
            "$ref": "#/components/schemas/DNS policy"
          },
          {
            "type": "object",
            "properties": {
              "name": {
                "type": "string"
              }
            }
          }
        ]
      },
      "IntegrationDnsForwardDomainPolicyDto": {
        "allOf": [
          {
            "$ref": "#/components/schemas/DNS policy"
          },
          {
            "type": "object",
            "properties": {
              "domain": {
                "type": "string"
              }
            }
          }
        ]
      }
    }
  }
}`
