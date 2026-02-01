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
		schema.Category = extractCategoryFromPath(schema.Path)
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
				schema.Request.Properties = g.extractProperties(page, section)
				fmt.Printf("Found %d request properties\n", len(schema.Request.Properties))
			}
		} else if strings.Contains(headerText, "response") {
			if schema.Response == nil {
				schema.Response = &SchemaObject{Name: "Response"}
				schema.Response.Properties = g.extractProperties(page, section)
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
