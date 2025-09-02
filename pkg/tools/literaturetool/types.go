package literaturetool

import "time"

// Article represents literature information from various providers.
type Article struct {
	ID           string        `json:"id"`
	Source       string        `json:"source"`
	PMID         string        `json:"pmid"`
	PMCID        string        `json:"pmcid,omitempty"`
	DOI          string        `json:"doi,omitempty"`
	Title        string        `json:"title"`
	AuthorString string        `json:"author_string"`
	Authors      []Author      `json:"authors"`
	Abstract     string        `json:"abstract"`
	Journal      Journal       `json:"journal"`
	PubYear      string        `json:"pub_year"`
	PageInfo     string        `json:"page_info,omitempty"`
	Keywords     []string      `json:"keywords,omitempty"`
	IsOpenAccess bool          `json:"is_open_access"`
	HasPDF       bool          `json:"has_pdf"`
	License      string        `json:"license,omitempty"`
	CitedByCount int           `json:"cited_by_count"`
	Language     string        `json:"language,omitempty"`
	PubTypes     []string      `json:"pub_types,omitempty"`
	MeshHeadings []MeshHeading `json:"mesh_headings,omitempty"`
	Chemicals    []Chemical    `json:"chemicals,omitempty"`
	Grants       []Grant       `json:"grants,omitempty"`
	PublishDate  *time.Time    `json:"publish_date,omitempty"`
	CreationDate *time.Time    `json:"creation_date,omitempty"`
	RevisionDate *time.Time    `json:"revision_date,omitempty"`
}

// Author represents author information.
type Author struct {
	FullName     string        `json:"full_name"`
	FirstName    string        `json:"first_name"`
	LastName     string        `json:"last_name"`
	Initials     string        `json:"initials"`
	ORCID        string        `json:"orcid,omitempty"`
	Affiliations []Affiliation `json:"affiliations,omitempty"`
}

// Affiliation represents author affiliation information.
type Affiliation struct {
	Affiliation string `json:"affiliation"`
}

// Journal represents journal information.
type Journal struct {
	Title               string `json:"title"`
	MedlineAbbreviation string `json:"medline_abbreviation,omitempty"`
	ISOAbbreviation     string `json:"iso_abbreviation,omitempty"`
	ISSN                string `json:"issn,omitempty"`
	ESSN                string `json:"essn,omitempty"`
	Volume              string `json:"volume,omitempty"`
	Issue               string `json:"issue,omitempty"`
	IssueID             int    `json:"issue_id,omitempty"`
	DateOfPublication   string `json:"date_of_publication,omitempty"`
	MonthOfPublication  int    `json:"month_of_publication,omitempty"`
	YearOfPublication   int    `json:"year_of_publication,omitempty"`
	NLMID               string `json:"nlm_id,omitempty"`
}

// MeshHeading represents MeSH heading information.
type MeshHeading struct {
	MajorTopic     bool            `json:"major_topic"`
	DescriptorName string          `json:"descriptor_name"`
	MeshQualifiers []MeshQualifier `json:"mesh_qualifiers,omitempty"`
}

// MeshQualifier represents MeSH qualifier information.
type MeshQualifier struct {
	QualifierName string `json:"qualifier_name"`
	MajorTopic    bool   `json:"major_topic"`
}

// Chemical represents chemical information.
type Chemical struct {
	Name        string `json:"name"`
	RegistryNum string `json:"registry_num"`
}

// Grant represents grant information.
type Grant struct {
	GrantID string `json:"grant_id"`
	Agency  string `json:"agency"`
	OrderIn int    `json:"order_in"`
}

// LiteratureError represents errors from literature API operations.
type LiteratureError struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	Code    string    `json:"code,omitempty"`
}

// ErrorType represents different types of literature API errors.
type ErrorType string

const (
	ErrorTypeInvalidInput    ErrorType = "invalid_input"
	ErrorTypeArticleNotFound ErrorType = "article_not_found"
	ErrorTypeNetworkError    ErrorType = "network_error"
	ErrorTypeAPIError        ErrorType = "api_error"
)

// Error implements the error interface.
func (e *LiteratureError) Error() string {
	return e.Message
}
