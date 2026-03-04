package util

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultCVEBaseURL = "https://cveawg.mitre.org/api"

var defaultHTTPClient = &http.Client{Timeout: 10 * time.Second}

// Client calls the CVE.org API for metadata about a given CVE identifier.
type CVEClient struct {
	baseURL    string
	httpClient *http.Client
}

// Option configures a CVE client before it is used.
type Option func(*CVEClient)

// WithBaseURL overrides the endpoint that will be queried.
func WithBaseURL(raw string) Option {
	return func(c *CVEClient) {
		if strings.TrimSpace(raw) != "" {
			c.baseURL = strings.TrimRight(raw, "/")
		}
	}
}

// WithHTTPClient overrides the HTTP transport used for requests.
func WithHTTPClient(client *http.Client) Option {
	return func(c *CVEClient) {
		if client != nil {
			c.httpClient = client
		}
	}
}

// NewClient builds a Client with zero or more options applied.
func NewClient(opts ...Option) *CVEClient {
	c := &CVEClient{baseURL: defaultCVEBaseURL, httpClient: defaultHTTPClient}
	for _, opt := range opts {
		if opt != nil {
			opt(c)
		}
	}
	return c
}

// CVE models the top-level response coming from cve.org.
type CVE struct {
	DataType    string         `json:"dataType"`
	DataVersion string         `json:"dataVersion"`
	Metadata    *CVEMetadata   `json:"cveMetadata"`
	Containers  *CVEContainers `json:"containers"`
}

// CVEMetadata contains the subset of fields we care about.
type CVEMetadata struct {
	AssignerOrgId     string `json:"assignerOrgId"`
	AssignerShortName string `json:"assignerShortName"`
	CVEID             string `json:"cveId"`
	State             string `json:"state"`
	DateReserved      string `json:"dateReserved"`
	DatePublished     string `json:"datePublished"`
	DateUpdated       string `json:"dateUpdated"`
}

// Published parses the date the CVE was published.
func (m *CVEMetadata) Published(layouts ...string) (*time.Time, error) {
	if m == nil || m.DatePublished == "" {
		return nil, fmt.Errorf("missing datePublished")
	}
	candidates := append(layouts, time.RFC3339, "2006-01-02", "2006-01-02T15:04:05Z")
	for _, layout := range candidates {
		if t, err := time.Parse(layout, m.DatePublished); err == nil {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("invalid date format %q", m.DatePublished)
}

// Fetch retrieves the CVE metadata for the supplied identifier.
func (c *CVEClient) Fetch(ctx context.Context, id string) (*CVE, error) {
	if c == nil {
		c = NewClient()
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("cve id must be provided")
	}
	endpoint := fmt.Sprintf("%s/cve/%s", c.baseURL, url.PathEscape(id))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cve request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("cve request failed with status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	var out CVE
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode cve response: %w", err)
	}
	return &out, nil
}

// Metadata returns the parsed metadata for id or an error when the API response is missing it.
func (c *CVEClient) Metadata(ctx context.Context, id string) (*CVEMetadata, error) {
	resp, err := c.Fetch(ctx, id)
	if err != nil {
		return nil, err
	}
	if resp.Metadata == nil {
		return nil, fmt.Errorf("cve metadata missing for %s", id)
	}
	return resp.Metadata, nil
}

// CVEContainers captures the payload under the containers key in the CVE response.
type CVEContainers struct {
	CNA *CNAContainer  `json:"cna"`
	ADP []ADPContainer `json:"adp"`
}

// CNAContainer represents the CNA-specific payload.
type CNAContainer struct {
	Title            string            `json:"title"`
	Descriptions     []LocalizedText   `json:"descriptions"`
	Affected         []AffectedProduct `json:"affected"`
	ProblemTypes     []ProblemType     `json:"problemTypes"`
	References       []Reference       `json:"references"`
	Metrics          []Metric          `json:"metrics"`
	Solutions        []LocalizedText   `json:"solutions"`
	Credits          []Credit          `json:"credits"`
	ProviderMetadata *ProviderMetadata `json:"providerMetadata"`
}

// ADPContainer represents an alternate data provider payload.
type ADPContainer struct {
	Metrics          []Metric          `json:"metrics"`
	Title            string            `json:"title"`
	ProviderMetadata *ProviderMetadata `json:"providerMetadata"`
}

// LocalizedText is a reusable structure for localized fields.
type LocalizedText struct {
	Lang  string `json:"lang"`
	Value string `json:"value"`
}

// AffectedProduct describes impacted vendors/products.
type AffectedProduct struct {
	Vendor        string            `json:"vendor"`
	Product       string            `json:"product"`
	Versions      []AffectedVersion `json:"versions"`
	DefaultStatus string            `json:"defaultStatus"`
}

// AffectedVersion represents a single version entry from the affected list.
type AffectedVersion struct {
	Version     string `json:"version"`
	Status      string `json:"status"`
	LessThan    string `json:"lessThan"`
	VersionType string `json:"versionType"`
}

// ProblemType groups CWE descriptions.
type ProblemType struct {
	Descriptions []ProblemDescription `json:"descriptions"`
}

// ProblemDescription adds CWE-specific details.
type ProblemDescription struct {
	Lang        string `json:"lang"`
	Description string `json:"description"`
	CWEID       string `json:"cweId"`
	Type        string `json:"type"`
}

// Reference points to external information.
type Reference struct {
	URL string `json:"url"`
}

// Metric contains CVSS or other scoring information.
type Metric struct {
	Format    string           `json:"format"`
	Scenarios []MetricScenario `json:"scenarios"`
	CVSSv31   *CVSSv31         `json:"cvssV3_1"`
	Other     *OtherMetric     `json:"other"`
}

// MetricScenario captures user-facing scenario labels.
type MetricScenario struct {
	Lang  string `json:"lang"`
	Value string `json:"value"`
}

// CVSSv31 mirrors the CVSSv3.1 vector structure.
type CVSSv31 struct {
	Version               string  `json:"version"`
	VectorString          string  `json:"vectorString"`
	AttackVector          string  `json:"attackVector"`
	AttackComplexity      string  `json:"attackComplexity"`
	PrivilegesRequired    string  `json:"privilegesRequired"`
	UserInteraction       string  `json:"userInteraction"`
	Scope                 string  `json:"scope"`
	ConfidentialityImpact string  `json:"confidentialityImpact"`
	IntegrityImpact       string  `json:"integrityImpact"`
	AvailabilityImpact    string  `json:"availabilityImpact"`
	BaseScore             float32 `json:"baseScore"`
	BaseSeverity          string  `json:"baseSeverity"`
}

// OtherMetric represents vendor-specific scores (e.g., SSVC).
type OtherMetric struct {
	Type    string       `json:"type"`
	Content OtherContent `json:"content"`
}

// OtherContent contains SSVC payload details.
type OtherContent struct {
	Timestamp string              `json:"timestamp"`
	ID        string              `json:"id"`
	Options   []map[string]string `json:"options"`
	Role      string              `json:"role"`
	Version   string              `json:"version"`
}

// Credit records contributors to the vulnerability discovery.
type Credit struct {
	Lang  string `json:"lang"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

// ProviderMetadata describes the entity supplying the container data.
type ProviderMetadata struct {
	OrgID       string `json:"orgId"`
	ShortName   string `json:"shortName"`
	DateUpdated string `json:"dateUpdated"`
}
