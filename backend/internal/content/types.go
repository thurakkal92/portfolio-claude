package content

import "encoding/json"

type Payload struct {
	Locale       string         `json:"locale"`
	SiteSettings SiteSettings   `json:"siteSettings"`
	Hero         Hero           `json:"hero"`
	About        About          `json:"about"`
	Skills       []SkillGroup   `json:"skills"`
	Projects     []Project      `json:"projects"`
	Experiences  []Experience   `json:"experiences"`
	Contact      Contact        `json:"contact"`
	SocialLinks  []SocialLink   `json:"socialLinks"`
}

type SiteSettings struct {
	SiteTitle       string  `json:"siteTitle"`
	SiteDescription string  `json:"siteDescription"`
	OGImagePath     *string `json:"ogImagePath,omitempty"`
	CVPath          *string `json:"cvPath,omitempty"`
}

type Hero struct {
	Eyebrow            string `json:"eyebrow"`
	Name               string `json:"name"`
	Subtitle           string `json:"subtitle"`
	CTAPrimaryLabel    string `json:"ctaPrimaryLabel"`
	CTASecondaryLabel  string `json:"ctaSecondaryLabel"`
}

type QuickFact struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Kind  string `json:"kind"` // "text" | "status"
}

type About struct {
	Heading    string      `json:"heading"`
	BodyMD     string      `json:"bodyMd"`
	QuickFacts []QuickFact `json:"quickFacts"`
}

type SkillGroup struct {
	Slug  string   `json:"slug"`
	Icon  string   `json:"icon"`
	Title string   `json:"title"`
	Items []string `json:"items"`
}

type ProjectHighlight struct {
	Icon string `json:"icon"`
	Text string `json:"text"`
}

type Project struct {
	Slug        string             `json:"slug"`
	Company     string             `json:"company"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	ImagePath   string             `json:"imagePath"`
	LiveURL     *string            `json:"liveUrl,omitempty"`
	SourceURL   *string            `json:"sourceUrl,omitempty"`
	Tags        []string           `json:"tags"`
	Highlights  []ProjectHighlight `json:"highlights"`
}

type Experience struct {
	Company   string  `json:"company"`
	Role      string  `json:"role"`
	StartDate string  `json:"startDate"` // ISO YYYY-MM-DD
	EndDate   *string `json:"endDate,omitempty"`
	Location  *string `json:"location,omitempty"`
	Summary   string  `json:"summary"`
}

type Contact struct {
	Heading           string `json:"heading"`
	Subheading        string `json:"subheading"`
	FormNameLabel     string `json:"formNameLabel"`
	FormEmailLabel    string `json:"formEmailLabel"`
	FormMessageLabel  string `json:"formMessageLabel"`
	FormSubmitLabel   string `json:"formSubmitLabel"`
	SuccessMessage    string `json:"successMessage"`
	ErrorMessage      string `json:"errorMessage"`
}

type SocialLink struct {
	Kind         string `json:"kind"`
	Href         string `json:"href"`
	DisplayLabel string `json:"displayLabel"`
}

// rawJSON helper for scanning JSONB into typed slices.
func unmarshalJSONB(b []byte, dst any) error {
	if len(b) == 0 {
		return nil
	}
	return json.Unmarshal(b, dst)
}
