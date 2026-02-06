package typegen

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

// scrape extracts API schema from a documentation page.
func (g *Generator) scrape(url string) (*APISchema, error) {
	fmt.Println("Opening page:", url)
	page := g.browser.MustPage(url)
	defer page.MustClose()

	fmt.Println("Waiting for page to load...")
	page.MustWaitLoad()
	time.Sleep(2 * time.Second)

	fmt.Println("Expanding all nested objects...")
	g.expandAllSections(page)

	fmt.Println("Extracting content...")

	schema := &APISchema{}

	// Extract endpoint name
	if el, err := page.Element("article h1, main h1, [class*='EndpointTitle'], [class*='PageTitle']"); err == nil {
		schema.Endpoint = el.MustText()
	}

	// Extract method
	methodEls, _ := page.Elements("[class*='MethodBadge']")
	if len(methodEls) > 0 {
		schema.Method = methodEls[0].MustText()
		fmt.Println("Method:", schema.Method)
	}

	// Extract path
	pathEls, _ := page.Elements("[class*='EndpointPath__Path']")
	if len(pathEls) > 0 {
		schema.Path = pathEls[0].MustText()
		fmt.Println("Path:", schema.Path)

		// Extract category from path
		schema.Category = sanitizeGoPackageName(extractCategoryFromPath(schema.Path))
		if schema.Category != "" {
			fmt.Println("Category:", schema.Category)
		}
	}

	// Derive endpoint name if needed
	if schema.Endpoint == "" || schema.Endpoint == "Developer" {
		schema.Endpoint = deriveEndpointName(url)
	}
	fmt.Println("Endpoint:", schema.Endpoint)

	// Extract sections
	sections, _ := page.Elements("section")
	fmt.Printf("Found %d sections\n", len(sections))

	for i, section := range sections {
		headerEl, err := section.Element("h2, h3, [class*='SectionTitle']")
		var headerText string

		if err != nil {
			text, _ := section.Text()
			textLower := strings.ToLower(text)
			if len(text) > 100 {
				fmt.Printf("Section %d: no header, text preview: %s...\n", i, text[:100])
			}

			if strings.Contains(textLower, "path parameters") {
				props := g.extractPropertiesFromText(page, section, "path parameters")
				if len(props) > 0 {
					schema.PathParams = props
					fmt.Printf("Found %d path parameters (from text)\n", len(schema.PathParams))
				}
			}
			if strings.Contains(textLower, "request body") {
				props := g.extractPropertiesFromText(page, section, "request body")
				if len(props) > 0 && schema.Request == nil {
					schema.Request = &SchemaObject{Name: "Request"}
					schema.Request.Properties = props
					if disc, variants := g.extractVariantsForSection(page, section, props); len(variants) > 1 {
						schema.Request.VariantDiscriminator = disc
						schema.Request.Variants = variants
						fmt.Printf("Found %d request variants for %s (from text)\n", len(variants), disc)
					}
					g.stripPathParamsFromRequest(schema)
					fmt.Printf("Found %d request properties (from text)\n", len(schema.Request.Properties))
				}
			}
			continue
		}

		headerText = strings.ToLower(strings.TrimSpace(headerEl.MustText()))
		fmt.Printf("Section %d header: %s\n", i, headerText)

		if strings.Contains(headerText, "path parameters") || strings.Contains(headerText, "path params") {
			schema.PathParams = g.extractProperties(page, section)
			fmt.Printf("Found %d path parameters\n", len(schema.PathParams))
		} else if strings.Contains(headerText, "request body") || strings.Contains(headerText, "request") {
			if schema.Request == nil {
				schema.Request = &SchemaObject{Name: "Request"}
				props := g.extractProperties(page, section)
				schema.Request.Properties = props
				if disc, variants := g.extractVariantsForSection(page, section, props); len(variants) > 1 {
					schema.Request.VariantDiscriminator = disc
					schema.Request.Variants = variants
					fmt.Printf("Found %d request variants for %s\n", len(variants), disc)
				}
				g.stripPathParamsFromRequest(schema)
				fmt.Printf("Found %d request properties\n", len(schema.Request.Properties))
			}
		} else if strings.Contains(headerText, "response") {
			if schema.Response == nil {
				schema.Response = &SchemaObject{Name: "Response"}
				props := g.extractProperties(page, section)
				schema.Response.Properties = props
				if disc, variants := g.extractVariantsForSection(page, section, props); len(variants) > 1 {
					schema.Response.VariantDiscriminator = disc
					schema.Response.Variants = variants
					fmt.Printf("Found %d response variants for %s\n", len(variants), disc)
				}
				fmt.Printf("Found %d response properties\n", len(schema.Response.Properties))
			}
		}
	}

	return schema, nil
}

// expandAllSections clicks all expand buttons to reveal nested object properties.
func (g *Generator) expandAllSections(page *rod.Page) {
	maxIterations := 5

	for i := 0; i < maxIterations; i++ {
		expandButtons, err := page.Elements("[class*='ExpandButton']")
		if err != nil {
			break
		}

		clicked := 0
		for _, btn := range expandButtons {
			btnText := btn.MustText()
			if btnText == "Expand" {
				visible, _ := btn.Visible()
				if visible {
					btn.Click("left", 1)
					clicked++
					time.Sleep(200 * time.Millisecond)
				}
			}
		}

		if clicked == 0 {
			break
		}

		time.Sleep(300 * time.Millisecond)
		fmt.Printf("Expanded %d sections (iteration %d)\n", clicked, i+1)
	}
}

// extractProperties extracts properties from a section.
func (g *Generator) extractProperties(page *rod.Page, section *rod.Element) []Property {
	g.resetHashDepth()

	rows, err := section.Elements("[class*='PropertyRow']")
	if err != nil || len(rows) == 0 {
		return nil
	}

	var flatProps []propertyRowInfo
	for _, row := range rows {
		info := g.extractRowInfo(row)
		if info.Name != "" {
			flatProps = append(flatProps, info)
		}
	}

	return g.buildPropertyHierarchy(flatProps)
}

// extractPropertiesFromText extracts properties from a section using text markers.
func (g *Generator) extractPropertiesFromText(page *rod.Page, section *rod.Element, marker string) []Property {
	g.resetHashDepth()

	rows, err := section.Elements("[class*='PropertyRow']")
	if err != nil || len(rows) == 0 {
		return nil
	}

	sectionText, _ := section.Text()
	sectionTextLower := strings.ToLower(sectionText)

	markerPos := strings.Index(sectionTextLower, marker)
	if markerPos == -1 {
		return nil
	}

	nextMarkers := []string{"path parameters", "request body", "responses", "query parameters"}
	nextPos := len(sectionText)
	for _, nm := range nextMarkers {
		if nm == marker {
			continue
		}
		pos := strings.Index(sectionTextLower[markerPos+len(marker):], nm)
		if pos != -1 && pos+markerPos+len(marker) < nextPos {
			nextPos = pos + markerPos + len(marker)
		}
	}

	relevantText := sectionTextLower[markerPos:nextPos]

	var flatProps []propertyRowInfo
	for _, row := range rows {
		info := g.extractRowInfo(row)
		if info.Name == "" {
			continue
		}
		if !strings.Contains(relevantText, strings.ToLower(info.Name)) {
			continue
		}
		flatProps = append(flatProps, info)
	}

	return g.buildPropertyHierarchy(flatProps)
}

// extractRowInfo extracts property information from a single row.
func (g *Generator) extractRowInfo(row *rod.Element) propertyRowInfo {
	info := propertyRowInfo{}

	if nameEl, err := row.Element("[class*='PropertyName']"); err == nil {
		info.Name = strings.TrimSpace(nameEl.MustText())
	}

	if info.Name == "" {
		return info
	}

	info.Type = g.extractTypeFromRow(row)

	if _, err := row.Element("[class*='RequiredBadge']"); err == nil {
		info.Required = true
	}

	if descEl, err := row.Element("[class*='PropertyDescription']"); err == nil {
		info.Description = strings.TrimSpace(descEl.MustText())
	}

	info.Enum = g.extractEnumValues(row)

	if strings.Contains(strings.ToLower(info.Type), "array") {
		info.IsArray = true
	}

	if strings.Contains(strings.ToLower(info.Type), "object") {
		info.IsObject = true
	}

	info.Depth = g.detectDepth(row)

	return info
}

// extractTypeFromRow extracts the type from a property row.
func (g *Generator) extractTypeFromRow(row *rod.Element) string {
	if typeEl, err := row.Element("[class*='PropertyType']"); err == nil {
		typeText := strings.TrimSpace(typeEl.MustText())
		if typeText != "" {
			return typeText
		}
	}

	if selectEl, err := row.Element("select"); err == nil {
		options, _ := selectEl.Elements("option")
		if len(options) > 0 {
			return "string"
		}
	}

	if inputEl, err := row.Element("input[type='number']"); err == nil {
		_ = inputEl
		return "integer"
	}

	if _, err := row.Element("input[type='checkbox']"); err == nil {
		return "boolean"
	}

	if _, err := row.Element("input[type='text']"); err == nil {
		return "string"
	}

	return ""
}

// extractEnumValues extracts enum values from various input elements.
func (g *Generator) extractEnumValues(row *rod.Element) []string {
	var enums []string

	radios, _ := row.Elements("label[class*='Radio']")
	for _, radio := range radios {
		text := strings.TrimSpace(radio.MustText())
		if text != "" && !strings.Contains(strings.ToLower(text), "discriminator") {
			enums = append(enums, text)
		}
	}

	if selectEl, err := row.Element("select"); err == nil {
		options, _ := selectEl.Elements("option")
		for _, opt := range options {
			val, _ := opt.Attribute("value")
			if val != nil && *val != "" {
				enums = append(enums, *val)
			} else {
				text := strings.TrimSpace(opt.MustText())
				if text != "" && text != "Select..." && text != "Choose..." {
					enums = append(enums, text)
				}
			}
		}
	}

	if datalist, err := row.Element("datalist"); err == nil {
		options, _ := datalist.Elements("option")
		for _, opt := range options {
			val, _ := opt.Attribute("value")
			if val != nil && *val != "" {
				enums = append(enums, *val)
			}
		}
	}

	return enums
}

func (g *Generator) extractVariantsForSection(page *rod.Page, section *rod.Element, props []Property) (string, map[string][]Property) {
	discriminator := g.findDiscriminatorProperty(props)
	if discriminator == nil || len(discriminator.Enum) < 2 {
		return "", nil
	}

	variants := make(map[string][]Property)
	for _, option := range discriminator.Enum {
		if g.clickVariantOption(section, option) {
			time.Sleep(200 * time.Millisecond)
		}
		g.expandAllSections(page)

		variantProps := g.extractProperties(page, section)
		if len(variantProps) > 0 {
			variants[option] = variantProps
		}
	}

	if len(variants) < 2 {
		return "", nil
	}

	if !variantsAreDistinct(variants) {
		return "", nil
	}

	return discriminator.Name, variants
}

func (g *Generator) stripPathParamsFromRequest(schema *APISchema) {
	if schema == nil || schema.Request == nil || len(schema.PathParams) == 0 {
		return
	}

	pathParamNames := make(map[string]bool)
	for _, param := range schema.PathParams {
		pathParamNames[strings.ToLower(param.Name)] = true
	}

	schema.Request.Properties = filterPropertiesByName(schema.Request.Properties, pathParamNames)
	if len(schema.Request.Variants) > 0 {
		for key, props := range schema.Request.Variants {
			schema.Request.Variants[key] = filterPropertiesByName(props, pathParamNames)
		}
	}
}

func filterPropertiesByName(props []Property, excluded map[string]bool) []Property {
	if len(props) == 0 {
		return props
	}

	filtered := make([]Property, 0, len(props))
	for _, prop := range props {
		if excluded[strings.ToLower(prop.Name)] {
			continue
		}
		if len(prop.Children) > 0 {
			prop.Children = filterPropertiesByName(prop.Children, excluded)
		}
		filtered = append(filtered, prop)
	}

	return filtered
}

func (g *Generator) findDiscriminatorProperty(props []Property) *Property {
	for i := range props {
		prop := &props[i]
		if len(prop.Enum) < 2 {
			continue
		}
		if prop.IsArray || prop.IsObject {
			continue
		}
		return prop
	}
	return nil
}

func (g *Generator) clickVariantOption(section *rod.Element, value string) bool {
	normalized := normalizeToggleText(value)
	selectors := []string{
		"button",
		"[role='tab']",
		"[role='radio']",
		"label",
		"[class*='Toggle']",
		"[class*='Segment']",
		"[class*='Tab']",
	}

	var candidates []*rod.Element
	for _, selector := range selectors {
		els, _ := section.Elements(selector)
		if len(els) > 0 {
			candidates = append(candidates, els...)
		}
	}

	for _, el := range candidates {
		visible, _ := el.Visible()
		if !visible {
			continue
		}

		text := strings.TrimSpace(el.MustText())
		if text != "" && normalizeToggleText(text) == normalized {
			el.Click("left", 1)
			return true
		}

		if val, _ := el.Attribute("value"); val != nil {
			if normalizeToggleText(*val) == normalized {
				el.Click("left", 1)
				return true
			}
		}
		if val, _ := el.Attribute("data-value"); val != nil {
			if normalizeToggleText(*val) == normalized {
				el.Click("left", 1)
				return true
			}
		}
	}

	return false
}

func normalizeToggleText(text string) string {
	text = strings.ToLower(strings.TrimSpace(text))
	replacer := strings.NewReplacer(" ", "", "_", "", "-", "")
	return replacer.Replace(text)
}

func variantsAreDistinct(variants map[string][]Property) bool {
	var baseline []Property
	for _, props := range variants {
		baseline = props
		break
	}

	for _, props := range variants {
		if !propertiesEqual(baseline, props) {
			return true
		}
	}

	return false
}

func propertiesEqual(a, b []Property) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !propertyEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

func propertyEqual(a, b Property) bool {
	if strings.ToLower(a.Name) != strings.ToLower(b.Name) {
		return false
	}
	if a.Type != b.Type || a.Required != b.Required || a.IsArray != b.IsArray || a.IsObject != b.IsObject {
		return false
	}
	if len(a.Enum) != len(b.Enum) {
		return false
	}
	for i := range a.Enum {
		if a.Enum[i] != b.Enum[i] {
			return false
		}
	}
	if len(a.Children) != len(b.Children) {
		return false
	}
	for i := range a.Children {
		if !propertyEqual(a.Children[i], b.Children[i]) {
			return false
		}
	}
	return true
}

// detectDepth determines the nesting depth of a property row.
func (g *Generator) detectDepth(row *rod.Element) int {
	classAttr, err := row.Attribute("class")
	if err != nil || classAttr == nil {
		return 0
	}

	className := *classAttr
	parts := strings.Fields(className)
	if len(parts) < 2 {
		return 0
	}

	hash := parts[len(parts)-1]
	return g.getDepthFromHash(hash)
}

// getDepthFromHash returns the depth for a given CSS hash.
func (g *Generator) getDepthFromHash(hash string) int {
	if depth, ok := g.hashDepth[hash]; ok {
		return depth
	}

	g.hashDepth[hash] = g.nextDepth
	g.nextDepth++
	return g.hashDepth[hash]
}

// resetHashDepth resets the hash-to-depth mapping for a new section.
func (g *Generator) resetHashDepth() {
	g.hashDepth = make(map[string]int)
	g.nextDepth = 0
}

// buildPropertyHierarchy builds a hierarchical property structure from a flat list.
func (g *Generator) buildPropertyHierarchy(flatProps []propertyRowInfo) []Property {
	if len(flatProps) == 0 {
		return nil
	}

	var result []Property
	i := 0

	for i < len(flatProps) {
		prop, consumed := g.buildPropertyWithChildren(flatProps, i, 0)
		if prop.Name != "" {
			result = append(result, prop)
		}
		i += consumed
		if consumed == 0 {
			i++
		}
	}

	return result
}

// buildPropertyWithChildren recursively builds a property with its children.
func (g *Generator) buildPropertyWithChildren(flatProps []propertyRowInfo, startIdx int, expectedDepth int) (Property, int) {
	if startIdx >= len(flatProps) {
		return Property{}, 0
	}

	info := flatProps[startIdx]

	if info.Depth < expectedDepth {
		return Property{}, 0
	}

	prop := Property{
		Name:        info.Name,
		Type:        info.Type,
		Description: info.Description,
		Required:    info.Required,
		Enum:        info.Enum,
		IsArray:     info.IsArray,
	}

	consumed := 1

	if info.IsObject {
		childDepth := info.Depth + 1

		for startIdx+consumed < len(flatProps) {
			nextInfo := flatProps[startIdx+consumed]

			if nextInfo.Depth <= info.Depth {
				break
			}

			if nextInfo.Depth == childDepth {
				childProp, childConsumed := g.buildPropertyWithChildren(flatProps, startIdx+consumed, childDepth)
				if childProp.Name != "" {
					prop.Children = append(prop.Children, childProp)
				}
				consumed += childConsumed
				if childConsumed == 0 {
					consumed++
				}
			} else if nextInfo.Depth > childDepth {
				consumed++
			} else {
				break
			}
		}
	}

	return prop, consumed
}

// isEndpointURL checks if a URL looks like an API endpoint documentation page.
func isEndpointURL(url string) bool {
	lowerURL := strings.ToLower(url)

	versionPattern := regexp.MustCompile(`/v\d+(\.\d+)*(/|$)`)
	if !versionPattern.MatchString(lowerURL) {
		return false
	}

	if strings.HasSuffix(lowerURL, "/overview") ||
		strings.HasSuffix(lowerURL, "/index") ||
		strings.HasSuffix(lowerURL, "/") {
		return false
	}

	docPages := []string{
		"gettingstarted", "filtering", "error-handling", "errorhandling",
		"authentication", "authorization", "introduction", "overview",
		"changelog", "migration", "generic-information", "quick_start", "quickstart",
	}

	if strings.Contains(lowerURL, "connector") {
		return false
	}

	for _, doc := range docPages {
		if strings.HasSuffix(lowerURL, "/"+doc) {
			return false
		}
	}

	return true
}
