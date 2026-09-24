package openapipatch

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/pb33f/libopenapi/orderedmap"
)

// RewriteDiscriminatorUnions applies only the compatibility rewrite still
// needed by oapi-codegen for the current UniFi specification.
func RewriteDiscriminatorUnions(doc []byte) ([]byte, error) {
	openapiDoc, err := libopenapi.NewDocument(doc)
	if err != nil {
		return nil, fmt.Errorf("parse spec: %w", err)
	}

	model, err := openapiDoc.BuildV3Model()
	if err != nil {
		return nil, fmt.Errorf("build v3 model: %w", err)
	}
	if model == nil {
		return nil, fmt.Errorf("build v3 model: empty document model")
	}
	if model.Model.Components == nil || model.Model.Components.Schemas == nil {
		return nil, fmt.Errorf("spec.components.schemas is missing or invalid")
	}

	schemas := model.Model.Components.Schemas

	changed, err := simplifyEntityMetadataSchemas(schemas)
	if err != nil {
		return nil, err
	}

	if !changed {
		return doc, nil
	}

	out, _, _, renderErr := openapiDoc.RenderAndReload()
	if renderErr != nil {
		return nil, fmt.Errorf("render patched spec: %w", renderErr)
	}

	normalized, _, err := normalizeSchemaNames(out)
	if err != nil {
		return nil, fmt.Errorf("normalize schema names: %w", err)
	}

	return normalized, nil
}

func simplifyEntityMetadataSchemas(schemas *orderedmap.Map[string, *base.SchemaProxy]) (bool, error) {
	userDefinedSchemaName := "User defined entity metadata"
	entitySchemaName := "Entity metadata"

	userDefinedProxy, ok := schemas.Get(userDefinedSchemaName)
	if !ok || userDefinedProxy == nil {
		return false, nil
	}

	userDefinedSchema := userDefinedProxy.Schema()
	if userDefinedSchema == nil {
		return false, fmt.Errorf("schema %q is invalid: %v", userDefinedSchemaName, userDefinedProxy.GetBuildError())
	}

	if len(userDefinedSchema.AllOf) == 1 && userDefinedSchema.AllOf[0] != nil && userDefinedSchema.AllOf[0].IsReference() && userDefinedSchema.AllOf[0].GetReference() == schemaRef(entitySchemaName) {
		return false, nil
	}

	userDefinedSchema.AllOf = []*base.SchemaProxy{base.CreateSchemaProxyRef(schemaRef(entitySchemaName))}
	userDefinedSchema.OneOf = nil
	userDefinedSchema.AnyOf = nil

	return true, nil
}

func normalizeSchemaNames(doc []byte) ([]byte, bool, error) {
	var spec map[string]any
	if err := json.Unmarshal(doc, &spec); err != nil {
		return nil, false, fmt.Errorf("decode spec: %w", err)
	}

	components, ok := spec["components"].(map[string]any)
	if !ok {
		return doc, false, nil
	}
	schemas, ok := components["schemas"].(map[string]any)
	if !ok || len(schemas) == 0 {
		return doc, false, nil
	}

	renameMap := make(map[string]string)
	usedNames := make(map[string]struct{}, len(schemas))
	for name := range schemas {
		usedNames[name] = struct{}{}
	}

	for name := range schemas {
		if isValidSchemaName(name) {
			continue
		}
		newName := sanitizeSchemaName(name)
		if newName == "" {
			newName = "Schema"
		}
		if newName == name {
			continue
		}
		candidate := newName
		for suffix := 1; ; suffix++ {
			if _, exists := usedNames[candidate]; !exists {
				break
			}
			candidate = fmt.Sprintf("%s_%d", newName, suffix)
		}
		usedNames[candidate] = struct{}{}
		renameMap[name] = candidate
	}

	if len(renameMap) == 0 {
		return doc, false, nil
	}

	updatedSchemas := make(map[string]any, len(schemas))
	for name, schema := range schemas {
		if newName, ok := renameMap[name]; ok {
			updatedSchemas[newName] = schema
			continue
		}
		updatedSchemas[name] = schema
	}
	components["schemas"] = updatedSchemas

	updateRefs(spec, renameMap)

	updated, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return nil, false, fmt.Errorf("encode updated spec: %w", err)
	}

	return updated, true, nil
}

func updateRefs(node any, renameMap map[string]string) {
	switch typed := node.(type) {
	case map[string]any:
		for key, value := range typed {
			if key == "$ref" {
				if ref, ok := value.(string); ok {
					if newRef, updated := updateSchemaRef(ref, renameMap); updated {
						typed[key] = newRef
					}
				}
				continue
			}
			if ref, ok := value.(string); ok {
				if newRef, updated := updateSchemaRef(ref, renameMap); updated {
					typed[key] = newRef
					continue
				}
			}
			updateRefs(value, renameMap)
		}
	case []any:
		for idx := range typed {
			updateRefs(typed[idx], renameMap)
		}
	}
}

func updateSchemaRef(ref string, renameMap map[string]string) (string, bool) {
	const prefix = "#/components/schemas/"
	if !strings.HasPrefix(ref, prefix) {
		return ref, false
	}
	name := strings.TrimPrefix(ref, prefix)
	newName, ok := renameMap[name]
	if !ok {
		return ref, false
	}
	return prefix + newName, true
}

func isValidSchemaName(name string) bool {
	if name == "" {
		return false
	}
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= 'A' && r <= 'Z':
		case r >= '0' && r <= '9':
		case r == '.' || r == '-' || r == '_':
		default:
			return false
		}
	}
	return true
}

func sanitizeSchemaName(name string) string {
	var out strings.Builder
	lastUnderscore := false
	for _, r := range name {
		valid := (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '.' || r == '-' || r == '_'
		if valid {
			out.WriteRune(r)
			lastUnderscore = false
			continue
		}
		if !lastUnderscore {
			out.WriteByte('_')
			lastUnderscore = true
		}
	}

	return strings.Trim(out.String(), "_.-")
}

func schemaRef(name string) string {
	return "#/components/schemas/" + name
}
