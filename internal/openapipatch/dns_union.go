package openapipatch

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

const (
	createOrUpdateDNSPolicySchema = "Create or update DNS policy"
	createOrUpdateDNSPolicyBase   = "Create or update DNS policy base"
	dnsPolicySchema               = "DNS policy"
	dnsPolicyBase                 = "DNS policy base"
)

// RewriteDNSUnions rewrites the DNS policy schemas so discriminator-based
// inheritance becomes a proper oneOf union that oapi-codegen can model.
func RewriteDNSUnions(doc []byte) ([]byte, error) {
	var spec map[string]any
	if err := json.Unmarshal(doc, &spec); err != nil {
		return nil, fmt.Errorf("decode spec: %w", err)
	}

	schemas, err := schemaObjects(spec)
	if err != nil {
		return nil, err
	}

	targets := []struct {
		name     string
		baseName string
	}{
		{name: createOrUpdateDNSPolicySchema, baseName: createOrUpdateDNSPolicyBase},
		{name: dnsPolicySchema, baseName: dnsPolicyBase},
	}

	for _, target := range targets {
		if err := rewriteDiscriminatorAllOfUnion(schemas, target.name, target.baseName); err != nil {
			return nil, err
		}
	}

	out, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode patched spec: %w", err)
	}

	return out, nil
}

func schemaObjects(spec map[string]any) (map[string]map[string]any, error) {
	components, ok := spec["components"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("spec.components is missing or invalid")
	}

	rawSchemas, ok := components["schemas"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("spec.components.schemas is missing or invalid")
	}

	schemas := make(map[string]map[string]any, len(rawSchemas))
	for name, rawSchema := range rawSchemas {
		schema, ok := rawSchema.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("schema %q is invalid", name)
		}
		schemas[name] = schema
	}

	return schemas, nil
}

func rewriteDiscriminatorAllOfUnion(schemas map[string]map[string]any, schemaName, syntheticBaseName string) error {
	schema, ok := schemas[schemaName]
	if !ok {
		return fmt.Errorf("schema %q not found", schemaName)
	}

	if _, alreadyUnion := schema["oneOf"]; alreadyUnion {
		return nil
	}

	discriminator, ok := schema["discriminator"].(map[string]any)
	if !ok {
		return fmt.Errorf("schema %q is missing a discriminator", schemaName)
	}

	mapping, ok := discriminator["mapping"].(map[string]any)
	if !ok || len(mapping) == 0 {
		return fmt.Errorf("schema %q has an invalid discriminator mapping", schemaName)
	}

	if _, exists := schemas[syntheticBaseName]; exists {
		return fmt.Errorf("synthetic schema %q already exists", syntheticBaseName)
	}

	syntheticBase := deepCopyMap(schema)
	delete(syntheticBase, "discriminator")
	schemas[syntheticBaseName] = syntheticBase

	originalRef := schemaRef(schemaName)
	syntheticRef := schemaRef(syntheticBaseName)
	oneOfRefs := sortedUniqueMappingRefs(mapping)

	for _, ref := range oneOfRefs {
		childName := schemaNameFromRef(ref)
		childSchema, ok := schemas[childName]
		if !ok {
			return fmt.Errorf("mapped schema %q not found", childName)
		}

		replaced, err := rewriteAllOfRef(childSchema, originalRef, syntheticRef)
		if err != nil {
			return fmt.Errorf("rewrite %q parent ref: %w", childName, err)
		}
		if !replaced {
			return fmt.Errorf("schema %q does not inherit from %q", childName, schemaName)
		}
	}

	delete(schema, "properties")
	delete(schema, "required")
	schema["oneOf"] = refsAsSchemaRefs(oneOfRefs)

	return nil
}

func rewriteAllOfRef(schema map[string]any, oldRef, newRef string) (bool, error) {
	rawAllOf, ok := schema["allOf"].([]any)
	if !ok {
		return false, fmt.Errorf("allOf is missing or invalid")
	}

	replaced := false
	for _, rawItem := range rawAllOf {
		item, ok := rawItem.(map[string]any)
		if !ok {
			return false, fmt.Errorf("allOf entry is invalid")
		}
		if itemRef, ok := item["$ref"].(string); ok && itemRef == oldRef {
			item["$ref"] = newRef
			replaced = true
		}
	}

	return replaced, nil
}

func refsAsSchemaRefs(refs []string) []any {
	out := make([]any, 0, len(refs))
	for _, ref := range refs {
		out = append(out, map[string]any{"$ref": ref})
	}
	return out
}

func sortedUniqueMappingRefs(mapping map[string]any) []string {
	refSet := make(map[string]struct{}, len(mapping))
	for _, rawRef := range mapping {
		ref, ok := rawRef.(string)
		if ok {
			refSet[ref] = struct{}{}
		}
	}

	refs := make([]string, 0, len(refSet))
	for ref := range refSet {
		refs = append(refs, ref)
	}
	sort.Strings(refs)
	return refs
}

func schemaRef(name string) string {
	return "#/components/schemas/" + name
}

func schemaNameFromRef(ref string) string {
	return strings.TrimPrefix(ref, "#/components/schemas/")
}

func deepCopyMap(in map[string]any) map[string]any {
	out := make(map[string]any, len(in))
	for key, value := range in {
		out[key] = deepCopyValue(value)
	}
	return out
}

func deepCopySlice(in []any) []any {
	out := make([]any, len(in))
	for idx, value := range in {
		out[idx] = deepCopyValue(value)
	}
	return out
}

func deepCopyValue(value any) any {
	switch v := value.(type) {
	case map[string]any:
		return deepCopyMap(v)
	case []any:
		return deepCopySlice(v)
	default:
		return v
	}
}
