package literaturetool

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLiteratureTool(t *testing.T) {
	t.Parallel()

	logger := log.New(os.Stderr, "[test] ", log.LstdFlags)
	tool, err := NewLiteratureTool(logger)

	require.NoError(t, err)
	assert.NotNil(t, tool)
	assert.Equal(t, "literature-fetch", tool.GetName())
	assert.Contains(t, tool.GetDescription(), "literature information")
	assert.NotNil(t, tool.GetTool())
	assert.NotNil(t, tool.GetSchema())
}

func TestLiteratureTool_GetMethods(t *testing.T) {
	t.Parallel()

	logger := log.New(os.Stderr, "[test] ", log.LstdFlags)
	tool, err := NewLiteratureTool(logger)
	require.NoError(t, err)

	assert.Equal(t, "literature-fetch", tool.GetName())
	expectedDescription := "Fetches scientific literature information using PubMed or DOI IDs via the dictyBase literature API"
	assert.Equal(t, expectedDescription, tool.GetDescription())

	schema := tool.GetSchema()
	assert.NotNil(t, schema)

	mcpTool := tool.GetTool()
	assert.Equal(t, "literature-fetch", mcpTool.Name)
}

func TestNormalizePMID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "valid PMID",
			input:   "12345678",
			want:    "12345678",
			wantErr: false,
		},
		{
			name:    "PMID with whitespace",
			input:   "  12345678  ",
			want:    "12345678",
			wantErr: false,
		},
		{
			name:    "empty PMID",
			input:   "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "non-numeric PMID",
			input:   "abc123",
			want:    "",
			wantErr: true,
		},
		{
			name:    "PMID with letters",
			input:   "1234abc",
			want:    "",
			wantErr: true,
		},
	}

	logger := log.New(os.Stderr, "[test] ", log.LstdFlags)
	tool, err := NewLiteratureTool(logger)
	require.NoError(t, err)

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			got, err := tool.normalizePMID(testCase.input)

			if testCase.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, testCase.want, got)
		})
	}
}

func TestNormalizeDOI(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "valid DOI",
			input:   "10.1038/nature12373",
			want:    "10.1038/nature12373",
			wantErr: false,
		},
		{
			name:    "DOI with prefix",
			input:   "DOI:10.1038/nature12373",
			want:    "10.1038/nature12373",
			wantErr: false,
		},
		{
			name:    "DOI with lowercase prefix",
			input:   "doi:10.1038/nature12373",
			want:    "10.1038/nature12373",
			wantErr: false,
		},
		{
			name:    "DOI with URL prefix",
			input:   "https://doi.org/10.1038/nature12373",
			want:    "10.1038/nature12373",
			wantErr: false,
		},
		{
			name:    "DOI with HTTP URL prefix",
			input:   "http://doi.org/10.1038/nature12373",
			want:    "10.1038/nature12373",
			wantErr: false,
		},
		{
			name:    "DOI with whitespace",
			input:   "  10.1038/nature12373  ",
			want:    "10.1038/nature12373",
			wantErr: false,
		},
		{
			name:    "empty DOI",
			input:   "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid DOI format - no 10. prefix",
			input:   "1038/nature12373",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid DOI format - no slash",
			input:   "10.1038nature12373",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid DOI format - empty suffix",
			input:   "10.1038/",
			want:    "",
			wantErr: true,
		},
	}

	logger := log.New(os.Stderr, "[test] ", log.LstdFlags)
	tool, err := NewLiteratureTool(logger)
	require.NoError(t, err)

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			got, err := tool.normalizeDOI(testCase.input)

			if testCase.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, testCase.want, got)
		})
	}
}

func TestNormalizeID(t *testing.T) {
	t.Parallel()

	logger := log.New(os.Stderr, "[test] ", log.LstdFlags)
	tool, err := NewLiteratureTool(logger)
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      string
		idType  string
		want    string
		wantErr bool
	}{
		{
			name:    "normalize PMID",
			id:      "12345678",
			idType:  "pmid",
			want:    "12345678",
			wantErr: false,
		},
		{
			name:    "normalize DOI",
			id:      "DOI:10.1038/nature12373",
			idType:  "doi",
			want:    "10.1038/nature12373",
			wantErr: false,
		},
		{
			name:    "unsupported ID type",
			id:      "12345",
			idType:  "isbn",
			want:    "",
			wantErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			got, err := tool.normalizeID(testCase.id, testCase.idType)

			if testCase.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, testCase.want, got)
		})
	}
}

func TestHandler_ValidationErrors(t *testing.T) {
	t.Parallel()
	logger := log.New(os.Stderr, "[test] ", log.LstdFlags)
	tool, err := NewLiteratureTool(logger)
	require.NoError(t, err)

	tests := []struct {
		name            string
		args            map[string]any
		wantErrContains string
	}{
		{
			name: "missing id",
			args: map[string]any{
				"id_type": "pmid",
			},
			wantErrContains: "missing required parameters",
		},
		{
			name: "missing id_type",
			args: map[string]any{
				"id": "12345678",
			},
			wantErrContains: "missing required parameters",
		},
		{
			name: "invalid id_type",
			args: map[string]any{
				"id":      "12345678",
				"id_type": "invalid",
			},
			wantErrContains: "validation error",
		},
		{
			name: "invalid provider",
			args: map[string]any{
				"id":       "12345678",
				"id_type":  "pmid",
				"provider": "invalid",
			},
			wantErrContains: "validation error",
		},
		{
			name: "invalid PMID format",
			args: map[string]any{
				"id":      "abc123",
				"id_type": "pmid",
			},
			wantErrContains: "invalid pmid format",
		},
		{
			name: "invalid DOI format",
			args: map[string]any{
				"id":      "invalid-doi",
				"id_type": "doi",
			},
			wantErrContains: "invalid doi format",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      "literature-fetch",
					Arguments: testCase.args,
				},
			}

			result, err := tool.Handler(context.Background(), request)

			require.Error(t, err)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), testCase.wantErrContains)
		})
	}
}

func TestFormatArticleResult(t *testing.T) {
	t.Parallel()

	logger := log.New(os.Stderr, "[test] ", log.LstdFlags)
	tool, err := NewLiteratureTool(logger)
	require.NoError(t, err)

	t.Run("nil article", func(t *testing.T) {
		t.Parallel()
		result, err := tool.formatArticleResult(nil)
		require.NoError(t, err)
		assert.Equal(t, "No article found", result)
	})

	t.Run("full article", func(t *testing.T) {
		t.Parallel()
		article := &Article{
			ID:       "PMC123456",
			Source:   "europepmc",
			PMID:     "12345678",
			DOI:      "10.1038/nature12373",
			Title:    "Test Article Title",
			Abstract: "This is a test abstract.",
			Authors: []Author{
				{
					FullName:  "John Doe",
					FirstName: "John",
					LastName:  "Doe",
				},
				{
					FullName:  "Jane Smith",
					FirstName: "Jane",
					LastName:  "Smith",
				},
			},
			Journal: Journal{
				Title: "Nature",
			},
			PubYear:      "2023",
			CitedByCount: 42,
		}

		result, err := tool.formatArticleResult(article)
		require.NoError(t, err)
		assert.Contains(t, result, "Test Article Title")
		assert.Contains(t, result, "John Doe")
		assert.Contains(t, result, "Jane Smith")
		assert.Contains(t, result, "Nature")
		assert.Contains(t, result, "2023")
		assert.Contains(t, result, "12345678")
		assert.Contains(t, result, "10.1038/nature12373")
		assert.Contains(t, result, "42")
		assert.Contains(t, result, "This is a test abstract")
		assert.Contains(t, result, "Raw JSON Data")
	})
}
