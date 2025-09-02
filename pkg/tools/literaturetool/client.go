package literaturetool

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dictybase/literature"
)

const (
	IDTypePMID = "pmid"
	IDTypeDOI  = "doi"
)

// LiteratureClient wraps the dictyBase literature clients.
type LiteratureClient struct {
	pubmedClient    *literature.Client
	europePMCClient *literature.EuropePMCClient
	logger          *log.Logger
}

// Option represents a configuration option for LiteratureClient.
type Option func(*Config)

// Config holds the configuration for the literature client.
type Config struct {
	timeout time.Duration
	logger  *log.Logger
}

// WithTimeout sets the HTTP timeout for requests.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.timeout = timeout
	}
}

// WithLogger sets the logger for the client.
func WithLogger(logger *log.Logger) Option {
	return func(c *Config) {
		c.logger = logger
	}
}

// NewLiteratureClient creates a new literature client with both PubMed and EuropePMC support.
func NewLiteratureClient(opts ...Option) (*LiteratureClient, error) {
	cfg := &Config{
		timeout: 30 * time.Second,
		logger:  log.Default(),
	}

	for _, opt := range opts {
		opt(cfg)
	}

	// Create PubMed client
	pubmedClient, err := literature.New(
		literature.WithTimeout(cfg.timeout),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create PubMed client: %w", err)
	}

	// Create EuropePMC client
	europePMCClient, err := literature.NewEuropePMCClient(
		literature.WithEuropePMCTimeout(cfg.timeout),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create EuropePMC client: %w", err)
	}

	return &LiteratureClient{
		pubmedClient:    pubmedClient,
		europePMCClient: europePMCClient,
		logger:          cfg.logger,
	}, nil
}

// GetArticleFromPubMed fetches article information from PubMed.
func (c *LiteratureClient) GetArticleFromPubMed(ctx context.Context, identifier, idType string) (*Article, error) {
	var article interface{}
	var err error

	switch idType {
	case IDTypePMID:
		article, err = c.pubmedClient.GetArticle(identifier)
	case IDTypeDOI:
		// PubMed doesn't directly support DOI lookup, so we'll use EuropePMC as fallback
		return c.GetArticleFromEuropePMC(ctx, identifier, idType)
	default:
		return nil, fmt.Errorf("unsupported ID type for PubMed: %s", idType)
	}

	if err != nil {
		// Convert to our standard error format
		if isNotFoundError(err) {
			return nil, &LiteratureError{
				Type:    ErrorTypeArticleNotFound,
				Message: fmt.Sprintf("article not found in PubMed for %s: %s", idType, identifier),
				Code:    "PUBMED_NOT_FOUND",
			}
		}
		return nil, &LiteratureError{
			Type:    ErrorTypeAPIError,
			Message: fmt.Sprintf("PubMed API error: %v", err),
			Code:    "PUBMED_API_ERROR",
		}
	}

	return c.convertToStandardArticle(article, "pubmed")
}

// GetArticleFromEuropePMC fetches article information from EuropePMC.
func (c *LiteratureClient) GetArticleFromEuropePMC(ctx context.Context, identifier, idType string) (*Article, error) {
	var article interface{}
	var err error

	switch idType {
	case IDTypePMID:
		article, err = c.europePMCClient.GetArticle(identifier)
	case IDTypeDOI:
		// For DOI, we need to search first to get the article
		searchResult, searchErr := c.europePMCClient.Search(
			fmt.Sprintf("DOI:%s", identifier),
			literature.WithEuropePMCLimit(1),
		)
		if searchErr != nil {
			return nil, fmt.Errorf("EuropePMC search error: %w", searchErr)
		}

		if len(searchResult.Articles) == 0 {
			return nil, &LiteratureError{
				Type:    ErrorTypeArticleNotFound,
				Message: fmt.Sprintf("no article found for DOI: %s", identifier),
				Code:    "DOI_NOT_FOUND",
			}
		}

		article = searchResult.Articles[0]
	default:
		return nil, fmt.Errorf("unsupported ID type for EuropePMC: %s", idType)
	}

	if err != nil {
		// Convert to our standard error format
		if isNotFoundError(err) {
			return nil, &LiteratureError{
				Type:    ErrorTypeArticleNotFound,
				Message: fmt.Sprintf("article not found in EuropePMC for %s: %s", idType, identifier),
				Code:    "EUROPEPMC_NOT_FOUND",
			}
		}
		return nil, &LiteratureError{
			Type:    ErrorTypeAPIError,
			Message: fmt.Sprintf("EuropePMC API error: %v", err),
			Code:    "EUROPEPMC_API_ERROR",
		}
	}

	return c.convertToStandardArticle(article, "europepmc")
}

// isNotFoundError checks if an error indicates that an article was not found.
func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := strings.ToLower(err.Error())
	notFoundIndicators := []string{
		"not found",
		"404",
		"no results",
		"no articles found",
		"article not found",
	}

	for _, indicator := range notFoundIndicators {
		if strings.Contains(errMsg, indicator) {
			return true
		}
	}

	return false
}

// GetArticleWithFallback implements the recommended logic: EuropePMC first, then PubMed fallback.
func (c *LiteratureClient) GetArticleWithFallback(ctx context.Context, identifier, idType string) (*Article, error) {
	// Try EuropePMC first
	article, err := c.GetArticleFromEuropePMC(ctx, identifier, idType)
	if err == nil {
		return article, nil
	}

	c.logger.Printf("EuropePMC failed for %s %s: %v, trying PubMed fallback", idType, identifier, err)

	// Only try PubMed fallback for PMIDs (since PubMed doesn't handle DOIs directly)
	if idType == IDTypePMID {
		fallbackArticle, fallbackErr := c.GetArticleFromPubMed(ctx, identifier, idType)
		if fallbackErr == nil {
			return fallbackArticle, nil
		}
		c.logger.Printf("PubMed fallback also failed for PMID %s: %v", identifier, fallbackErr)
	}

	// Return the original EuropePMC error
	return nil, err
}

// convertToStandardArticle converts provider-specific article structs to our standard Article struct.
func (c *LiteratureClient) convertToStandardArticle(article interface{}, provider string) (*Article, error) {
	switch provider {
	case "pubmed":
		return c.convertPubMedArticle(article)
	case "europepmc":
		return c.convertEuropePMCArticle(article)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

// convertPubMedArticle converts a PubMed article to our standard format.
func (c *LiteratureClient) convertPubMedArticle(article interface{}) (*Article, error) {
	// Type assertion for PubMed article
	pubmedArticle, ok := article.(*literature.Article)
	if !ok {
		return nil, fmt.Errorf("invalid PubMed article type")
	}

	// Convert authors
	authors := make([]Author, len(pubmedArticle.Authors))
	for i, author := range pubmedArticle.Authors {
		authors[i] = Author{
			FullName:  author.FullName,
			FirstName: author.FirstName,
			LastName:  author.LastName,
		}
	}

	// Extract year from publish date
	pubYear := ""
	if !pubmedArticle.PublishDate.IsZero() {
		pubYear = fmt.Sprintf("%d", pubmedArticle.PublishDate.Year())
	}

	return &Article{
		ID:           pubmedArticle.PMID,
		Source:       "pubmed",
		PMID:         pubmedArticle.PMID,
		DOI:          pubmedArticle.DOI,
		Title:        pubmedArticle.Title,
		AuthorString: "", // Will be constructed from authors
		Authors:      authors,
		Abstract:     pubmedArticle.Abstract,
		Journal: Journal{
			Title:  pubmedArticle.Journal,
			Volume: pubmedArticle.Volume,
			Issue:  pubmedArticle.Issue,
		},
		PubYear:      pubYear,
		PageInfo:     pubmedArticle.Pages,
		Keywords:     pubmedArticle.Keywords,
		IsOpenAccess: false,
		HasPDF:       false,
		CitedByCount: 0,
		PublishDate:  &pubmedArticle.PublishDate,
	}, nil
}

// convertEuropePMCArticle converts a EuropePMC article to our standard format.
func (c *LiteratureClient) convertEuropePMCArticle(article interface{}) (*Article, error) {
	europePMCArticle, ok := article.(*literature.EuropePMCArticle)
	if !ok {
		return nil, fmt.Errorf("invalid EuropePMC article type")
	}

	authors := c.convertAuthors(europePMCArticle.Authors)
	meshHeadings := c.convertMeshHeadings(europePMCArticle.MeshHeadings)
	chemicals := c.convertChemicals(europePMCArticle.Chemicals)
	grants := c.convertGrants(europePMCArticle.Grants)
	journal := c.convertJournal(europePMCArticle.Journal)

	return &Article{
		ID:           europePMCArticle.ID,
		Source:       "europepmc",
		PMID:         europePMCArticle.PMID,
		PMCID:        europePMCArticle.PMCID,
		DOI:          europePMCArticle.DOI,
		Title:        europePMCArticle.Title,
		AuthorString: europePMCArticle.AuthorString,
		Authors:      authors,
		Abstract:     europePMCArticle.Abstract,
		Journal:      journal,
		PubYear:      europePMCArticle.PubYear,
		PageInfo:     europePMCArticle.PageInfo,
		Keywords:     europePMCArticle.Keywords,
		IsOpenAccess: europePMCArticle.IsOpenAccess,
		HasPDF:       europePMCArticle.HasPDF,
		License:      europePMCArticle.License,
		CitedByCount: europePMCArticle.CitedByCount,
		Language:     europePMCArticle.Language,
		PubTypes:     europePMCArticle.PubTypes,
		MeshHeadings: meshHeadings,
		Chemicals:    chemicals,
		Grants:       grants,
		PublishDate:  europePMCArticle.PublishDate,
		CreationDate: europePMCArticle.CreationDate,
		RevisionDate: europePMCArticle.RevisionDate,
	}, nil
}

// convertAuthors converts EuropePMC authors to standard format.
func (c *LiteratureClient) convertAuthors(europePMCAuthors []literature.EuropePMCAuthor) []Author {
	authors := make([]Author, len(europePMCAuthors))
	for authorIndex, author := range europePMCAuthors {
		affiliations := make([]Affiliation, len(author.Affiliations))
		for affiliationIndex, affil := range author.Affiliations {
			affiliations[affiliationIndex] = Affiliation{
				Affiliation: affil.Affiliation,
			}
		}

		authors[authorIndex] = Author{
			FullName:     author.FullName,
			FirstName:    author.FirstName,
			LastName:     author.LastName,
			Initials:     author.Initials,
			ORCID:        author.ORCID,
			Affiliations: affiliations,
		}
	}
	return authors
}

// convertMeshHeadings converts EuropePMC MeSH headings to standard format.
func (c *LiteratureClient) convertMeshHeadings(europePMCMeshHeadings []literature.EuropePMCMeshHeading) []MeshHeading {
	meshHeadings := make([]MeshHeading, len(europePMCMeshHeadings))
	for meshIndex, mesh := range europePMCMeshHeadings {
		qualifiers := make([]MeshQualifier, len(mesh.MeshQualifiers))
		for qualifierIndex, qual := range mesh.MeshQualifiers {
			qualifiers[qualifierIndex] = MeshQualifier{
				QualifierName: qual.QualifierName,
				MajorTopic:    qual.MajorTopic,
			}
		}

		meshHeadings[meshIndex] = MeshHeading{
			MajorTopic:     mesh.MajorTopic,
			DescriptorName: mesh.DescriptorName,
			MeshQualifiers: qualifiers,
		}
	}
	return meshHeadings
}

// convertChemicals converts EuropePMC chemicals to standard format.
func (c *LiteratureClient) convertChemicals(europePMCChemicals []literature.EuropePMCChemical) []Chemical {
	chemicals := make([]Chemical, len(europePMCChemicals))
	for chemicalIndex, chem := range europePMCChemicals {
		chemicals[chemicalIndex] = Chemical{
			Name:        chem.Name,
			RegistryNum: chem.RegistryNumber,
		}
	}
	return chemicals
}

// convertGrants converts EuropePMC grants to standard format.
func (c *LiteratureClient) convertGrants(europePMCGrants []literature.EuropePMCGrant) []Grant {
	grants := make([]Grant, len(europePMCGrants))
	for grantIndex, grant := range europePMCGrants {
		grants[grantIndex] = Grant{
			GrantID: grant.GrantID,
			Agency:  grant.Agency,
			OrderIn: grant.OrderIn,
		}
	}
	return grants
}

// convertJournal converts EuropePMC journal to standard format.
func (c *LiteratureClient) convertJournal(europePMCJournal literature.EuropePMCJournal) Journal {
	return Journal{
		Title:               europePMCJournal.Title,
		MedlineAbbreviation: europePMCJournal.MedlineAbbreviation,
		ISOAbbreviation:     europePMCJournal.ISOAbbreviation,
		ISSN:                europePMCJournal.ISSN,
		ESSN:                europePMCJournal.ESSN,
		Volume:              europePMCJournal.Volume,
		Issue:               europePMCJournal.Issue,
		IssueID:             europePMCJournal.IssueID,
		DateOfPublication:   europePMCJournal.DateOfPublication,
		MonthOfPublication:  europePMCJournal.MonthOfPublication,
		YearOfPublication:   europePMCJournal.YearOfPublication,
		NLMID:               europePMCJournal.NLMID,
	}
}
