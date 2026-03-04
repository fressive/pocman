package vuln

import (
	"time"

	"gorm.io/gorm"
)

// Vuln represents a vulnerability report with contextual metadata.
type Vuln struct {
	gorm.Model

	// Identifier of the vulnerability, typically a CVE or vendor-specific ID.
	Identifier string `gorm:"size:32;not null;unique"`

	// User-friendly title for quick scans.
	Title *string `gorm:"size:256"`

	// Longer description or summary of the issue.
	Description *string `gorm:"type:text"`

	// Severity from 0-10 or other scoring scheme.
	Severity *float32 `gorm:"type:decimal(3,1)"`

	// Normalized severity label such as low/medium/high/critical.
	SeverityLabel *string `gorm:"size:32"`

	// CVSS score if available.
	CvssScore *float32 `gorm:"type:decimal(4,2)"`

	// Vector string from CVSS v3+.
	CvssVector *string `gorm:"size:64"`

	// Date the vendor publicly disclosed the issue.
	DisclosureDate *time.Time

	// Date the vulnerability was first published or observed.
	PublishedDate *time.Time

	// Suggested remediation steps or patches.
	Remediation *string `gorm:"type:text"`

	// Free-form notes for operators.
	Notes *string `gorm:"type:text"`

	// Relationship to a single vendor/product record.
	VendorID *uint
	Vendor   *Vendor `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	// References for this vulnerability.
	References []*Reference `gorm:"foreignKey:VulnID;constraint:OnDelete:CASCADE;"`
}
