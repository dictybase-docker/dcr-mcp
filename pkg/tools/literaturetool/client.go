package literaturetool

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dictybase/literature"
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
func (c *LiteratureClient) GetArticleFromPubMed(ctx context.Context, id, idType string) (*Article, error) {
	var article interface{}
	var err error

	switch idType {
	case "pmid":
		article, err = c.pubmedClient.GetArticle(id)
	case "doi":
		// PubMed doesn't directly support DOI lookup, so we'll use EuropePMC as fallback
		return c.GetArticleFromEuropePMC(ctx, id, idType)
	default:
		return nil, fmt.Errorf("unsupported ID type for PubMed: %s", idType)
	}

	if err != nil {
		// Convert to our standard error format
		if isNotFoundError(err) {
			return nil, &LiteratureError{
				Type:    ErrorTypeArticleNotFound,
				Message: fmt.Sprintf("article not found in PubMed for %s: %s", idType, id),
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
func (c *LiteratureClient) GetArticleFromEuropePMC(ctx context.Context, id, idType string) (*Article, error) {
	var article interface{}
	var err error

	switch idType {
	case "pmid":
		article, err = c.europePMCClient.GetArticle(id)
	case "doi":
		// For DOI, we need to search first to get the article
		searchResult, searchErr := c.europePMCClient.Search(
			fmt.Sprintf("DOI:%s", id),
			literature.WithEuropePMCLimit(1),
		)
		if searchErr != nil {
			return nil, fmt.Errorf("EuropePMC search error: %w", searchErr)
		}

		if len(searchResult.Articles) == 0 {
			return nil, &LiteratureError{
				Type:    ErrorTypeArticleNotFound,
				Message: fmt.Sprintf("no article found for DOI: %s", id),
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
				Message: fmt.Sprintf("article not found in EuropePMC for %s: %s", idType, id),
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
func (c *LiteratureClient) GetArticleWithFallback(ctx context.Context, id, idType string) (*Article, error) {
	// Try EuropePMC first
	article, err := c.GetArticleFromEuropePMC(ctx, id, idType)
	if err == nil {
		return article, nil
	}

	c.logger.Printf("EuropePMC failed for %s %s: %v, trying PubMed fallback", idType, id, err)

	// Only try PubMed fallback for PMIDs (since PubMed doesn't handle DOIs directly)
	if idType == "pmid" {
		fallbackArticle, fallbackErr := c.GetArticleFromPubMed(ctx, id, idType)
		if fallbackErr == nil {
			return fallbackArticle, nil
		}
		c.logger.Printf("PubMed fallback also failed for PMID %s: %v", id, fallbackErr)
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
	// Type assertion for EuropePMC article
	europePMCArticle, ok := article.(*literature.EuropePMCArticle)
	if !ok {
		return nil, fmt.Errorf("invalid EuropePMC article type")
	}

	// Convert authors
	authors := make([]Author, len(europePMCArticle.Authors))
	for i, author := range europePMCArticle.Authors {
		affiliations := make([]Affiliation, len(author.Affiliations))
		for j, affil := range author.Affiliations {
			affiliations[j] = Affiliation{
				Affiliation: affil.Affiliation,
			}
		}

		authors[i] = Author{
			FullName:     author.FullName,
			FirstName:    author.FirstName,
			LastName:     author.LastName,
			Initials:     author.Initials,
			ORCID:        author.ORCID,
			Affiliations: affiliations,
		}
	}

	// Convert MeSH headings
	meshHeadings := make([]MeshHeading, len(europePMCArticle.MeshHeadings))
	for i, mesh := range europePMCArticle.MeshHeadings {
		qualifiers := make([]MeshQualifier, len(mesh.MeshQualifiers))
		for j, qual := range mesh.MeshQualifiers {
			qualifiers[j] = MeshQualifier{
				QualifierName: qual.QualifierName,
				MajorTopic:    qual.MajorTopic,
			}
		}

		meshHeadings[i] = MeshHeading{
			MajorTopic:     mesh.MajorTopic,
			DescriptorName: mesh.DescriptorName,
			MeshQualifiers: qualifiers,
		}
	}

	// Convert chemicals
	chemicals := make([]Chemical, len(europePMCArticle.Chemicals))
	for i, chem := range europePMCArticle.Chemicals {
		chemicals[i] = Chemical{
			Name:        chem.Name,
			RegistryNum: chem.RegistryNumber,
		}
	}

	// Convert grants
	grants := make([]Grant, len(europePMCArticle.Grants))
	for i, grant := range europePMCArticle.Grants {
		grants[i] = Grant{
			GrantID: grant.GrantID,
			Agency:  grant.Agency,
			OrderIn: grant.OrderIn,
		}
	}

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
		Journal: Journal{
			Title:               europePMCArticle.Journal.Title,
			MedlineAbbreviation: europePMCArticle.Journal.MedlineAbbreviation,
			ISOAbbreviation:     europePMCArticle.Journal.ISOAbbreviation,
			ISSN:                europePMCArticle.Journal.ISSN,
			ESSN:                europePMCArticle.Journal.ESSN,
			Volume:              europePMCArticle.Journal.Volume,
			Issue:               europePMCArticle.Journal.Issue,
			IssueID:             europePMCArticle.Journal.IssueID,
			DateOfPublication:   europePMCArticle.Journal.DateOfPublication,
			MonthOfPublication:  europePMCArticle.Journal.MonthOfPublication,
			YearOfPublication:   europePMCArticle.Journal.YearOfPublication,
			NLMID:               europePMCArticle.Journal.NLMID,
		},
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
