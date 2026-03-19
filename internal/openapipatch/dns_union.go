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

	changed := false
	targets := []struct {
		name     string
		baseName string
	}{
		{name: createOrUpdateDNSPolicySchema, baseName: createOrUpdateDNSPolicyBase},
		{name: dnsPolicySchema, baseName: dnsPolicyBase},
	}

	for _, target := range targets {
		targetChanged, err := rewriteDiscriminatorAllOfUnion(schemas, target.name, target.baseName)
		if err != nil {
			return nil, err
		}
		changed = changed || targetChanged
	}

	if !changed {
		return doc, nil
	}

	out, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode patched spec: %w", err)
	}

	return out, nil
}

func schemaObjects(spec map[string]any) (map[string]any, error) {
	components, ok := spec["components"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("spec.components is missing or invalid")
	}

	schemas, ok := components["schemas"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("spec.components.schemas is missing or invalid")
	}

	for name, rawSchema := range schemas {
		_, ok := rawSchema.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("schema %q is invalid", name)
		}
	}

	return schemas, nil
}

func rewriteDiscriminatorAllOfUnion(schemas map[string]any, schemaName, syntheticBaseName string) (bool, error) {
	schema, err := schemaObject(schemas, schemaName)
	if err != nil {
		return false, err
	}

	if _, alreadyUnion := schema["oneOf"]; alreadyUnion {
		return false, nil
	}

	discriminator, ok := schema["discriminator"].(map[string]any)
	if !ok {
		return false, fmt.Errorf("schema %q is missing a discriminator", schemaName)
	}

	mapping, ok := discriminator["mapping"].(map[string]any)
	if !ok || len(mapping) == 0 {
		return false, fmt.Errorf("schema %q has an invalid discriminator mapping", schemaName)
	}

	if _, exists := schemas[syntheticBaseName]; exists {
		return false, fmt.Errorf("synthetic schema %q already exists", syntheticBaseName)
	}

	syntheticBase := deepCopyMap(schema)
	delete(syntheticBase, "discriminator")
	schemas[syntheticBaseName] = syntheticBase

	originalRef := schemaRef(schemaName)
	syntheticRef := schemaRef(syntheticBaseName)
	oneOfRefs := sortedUniqueMappingRefs(mapping)

	for _, ref := range oneOfRefs {
		childName := schemaNameFromRef(ref)
		childSchema, err := schemaObject(schemas, childName)
		if err != nil {
			return false, err
		}

		replaced, err := rewriteAllOfRef(childSchema, originalRef, syntheticRef)
		if err != nil {
			return false, fmt.Errorf("rewrite %q parent ref: %w", childName, err)
		}
		if !replaced {
			return false, fmt.Errorf("schema %q does not inherit from %q", childName, schemaName)
		}
	}

	delete(schema, "properties")
	delete(schema, "required")
	schema["oneOf"] = refsAsSchemaRefs(oneOfRefs)

	return true, nil
}

func schemaObject(schemas map[string]any, name string) (map[string]any, error) {
	rawSchema, ok := schemas[name]
	if !ok {
		return nil, fmt.Errorf("schema %q not found", name)
	}

	schema, ok := rawSchema.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("schema %q is invalid", name)
	}

	return schema, nil
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
