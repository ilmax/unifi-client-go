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

type unionRewriteTarget struct {
	SchemaName         string
	SyntheticBaseName  string
	OverrideMappedRefs map[string]mappedRefOverride
}

type mappedRefOverride struct {
	SyntheticSchemaName string
	CurrentParentRef    string
}

// RewriteDiscriminatorUnions rewrites known discriminator-based inheritance
// patterns into oneOf unions that oapi-codegen can model correctly.
func RewriteDiscriminatorUnions(doc []byte) ([]byte, error) {
	var spec map[string]any
	if err := json.Unmarshal(doc, &spec); err != nil {
		return nil, fmt.Errorf("decode spec: %w", err)
	}

	schemas, err := schemaObjects(spec)
	if err != nil {
		return nil, err
	}

	changed := false
	for _, target := range rewriteTargets() {
		targetChanged, err := rewriteDiscriminatorAllOfUnion(schemas, target)
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

func rewriteTargets() []unionRewriteTarget {
	return []unionRewriteTarget{
		{SchemaName: createOrUpdateDNSPolicySchema, SyntheticBaseName: createOrUpdateDNSPolicyBase},
		{SchemaName: dnsPolicySchema, SyntheticBaseName: dnsPolicyBase},
		{SchemaName: "Firewall policy action", SyntheticBaseName: "Firewall policy action base"},
		{SchemaName: "Firewall policy source traffic filter", SyntheticBaseName: "Firewall policy source traffic filter base"},
		{SchemaName: "Firewall policy destination traffic filter", SyntheticBaseName: "Firewall policy destination traffic filter base"},
		{SchemaName: "Firewall policy port filter", SyntheticBaseName: "Firewall policy port filter base"},
		{SchemaName: "Firewall policy IP address filter", SyntheticBaseName: "Firewall policy IP address filter base"},
		{SchemaName: "Firewall policy IP protocol scope", SyntheticBaseName: "Firewall policy IP protocol scope base"},
		{SchemaName: "Firewall schedule", SyntheticBaseName: "Firewall schedule base"},
		{
			SchemaName:        "Firewall policy IPv4 protocol",
			SyntheticBaseName: "Firewall policy IPv4 protocol base",
			OverrideMappedRefs: map[string]mappedRefOverride{
				schemaRef("Firewall policy IPv4 protocol number"): {
					SyntheticSchemaName: "Firewall policy IPv4 protocol number (IPv4 base)",
					CurrentParentRef:    schemaRef("Firewall policy IPv6 protocol"),
				},
			},
		},
		{SchemaName: "Firewall policy IPv6 protocol", SyntheticBaseName: "Firewall policy IPv6 protocol base"},
		{SchemaName: "Firewall policy IPv4 and IPv6 protocol", SyntheticBaseName: "Firewall policy IPv4 and IPv6 protocol base"},
		{SchemaName: "Firewall policy IPv4 named protocol", SyntheticBaseName: "Firewall policy IPv4 named protocol base"},
		{SchemaName: "Firewall policy IPv6 named protocol", SyntheticBaseName: "Firewall policy IPv6 named protocol base"},
		{SchemaName: "Firewall policy IPv4 and IPv6 named protocol", SyntheticBaseName: "Firewall policy IPv4 and IPv6 named protocol base"},
		{SchemaName: "Firewall policy IPv4 protocol preset", SyntheticBaseName: "Firewall policy IPv4 protocol preset base"},
		{SchemaName: "Firewall policy IPv6 protocol preset", SyntheticBaseName: "Firewall policy IPv6 protocol preset base"},
		{SchemaName: "Firewall policy IPv4 and IPv6 protocol preset", SyntheticBaseName: "Firewall policy IPv4 and IPv6 protocol preset base"},
	}
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

func rewriteDiscriminatorAllOfUnion(schemas map[string]any, target unionRewriteTarget) (bool, error) {
	schema, exists, err := maybeSchemaObject(schemas, target.SchemaName)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}

	if _, alreadyUnion := schema["oneOf"]; alreadyUnion {
		return false, nil
	}

	discriminator, ok := schema["discriminator"].(map[string]any)
	if !ok {
		return false, fmt.Errorf("schema %q is missing a discriminator", target.SchemaName)
	}

	mapping, ok := discriminator["mapping"].(map[string]any)
	if !ok || len(mapping) == 0 {
		return false, fmt.Errorf("schema %q has an invalid discriminator mapping", target.SchemaName)
	}

	if _, exists := schemas[target.SyntheticBaseName]; exists {
		return false, fmt.Errorf("synthetic schema %q already exists", target.SyntheticBaseName)
	}

	syntheticBase := deepCopyMap(schema)
	delete(syntheticBase, "discriminator")
	schemas[target.SyntheticBaseName] = syntheticBase

	originalRef := schemaRef(target.SchemaName)
	syntheticRef := schemaRef(target.SyntheticBaseName)

	for _, ref := range sortedUniqueMappingRefs(mapping) {
		childSchema, exists, err := maybeSchemaObject(schemas, schemaNameFromRef(ref))
		if err != nil {
			return false, err
		}
		if !exists {
			return false, fmt.Errorf("mapped schema %q not found", schemaNameFromRef(ref))
		}

		replaced, err := rewriteAllOfRef(childSchema, originalRef, syntheticRef)
		if err != nil {
			return false, fmt.Errorf("rewrite %q parent ref: %w", schemaNameFromRef(ref), err)
		}
		if replaced {
			continue
		}

		override, ok := target.OverrideMappedRefs[ref]
		if !ok {
			return false, fmt.Errorf("schema %q does not inherit from %q", schemaNameFromRef(ref), target.SchemaName)
		}

		overrideRef, err := cloneSchemaWithParentRefOverride(schemas, ref, override, syntheticRef)
		if err != nil {
			return false, err
		}
		replaceMappingRef(mapping, ref, overrideRef)
	}

	delete(schema, "properties")
	delete(schema, "required")
	schema["oneOf"] = refsAsSchemaRefs(sortedUniqueMappingRefs(mapping))

	return true, nil
}

func cloneSchemaWithParentRefOverride(schemas map[string]any, sourceRef string, override mappedRefOverride, newParentRef string) (string, error) {
	if _, exists := schemas[override.SyntheticSchemaName]; exists {
		return "", fmt.Errorf("synthetic schema %q already exists", override.SyntheticSchemaName)
	}

	sourceSchema, exists, err := maybeSchemaObject(schemas, schemaNameFromRef(sourceRef))
	if err != nil {
		return "", err
	}
	if !exists {
		return "", fmt.Errorf("mapped schema %q not found", schemaNameFromRef(sourceRef))
	}

	clonedSchema := deepCopyMap(sourceSchema)
	replaced, err := rewriteAllOfRef(clonedSchema, override.CurrentParentRef, newParentRef)
	if err != nil {
		return "", fmt.Errorf("rewrite %q parent ref: %w", override.SyntheticSchemaName, err)
	}
	if !replaced {
		return "", fmt.Errorf("schema %q does not inherit from override parent %q", schemaNameFromRef(sourceRef), schemaNameFromRef(override.CurrentParentRef))
	}

	schemas[override.SyntheticSchemaName] = clonedSchema
	return schemaRef(override.SyntheticSchemaName), nil
}

func replaceMappingRef(mapping map[string]any, oldRef, newRef string) {
	for key, rawRef := range mapping {
		ref, ok := rawRef.(string)
		if ok && ref == oldRef {
			mapping[key] = newRef
		}
	}
}

func maybeSchemaObject(schemas map[string]any, name string) (map[string]any, bool, error) {
	rawSchema, ok := schemas[name]
	if !ok {
		return nil, false, nil
	}

	schema, ok := rawSchema.(map[string]any)
	if !ok {
		return nil, false, fmt.Errorf("schema %q is invalid", name)
	}

	return schema, true, nil
}

func schemaObject(schemas map[string]any, name string) (map[string]any, error) {
	schema, exists, err := maybeSchemaObject(schemas, name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("schema %q not found", name)
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
