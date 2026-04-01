package registration

import "time"

type SubjectKind string

const (
	SubjectKindSelfHosted    SubjectKind = "self_hosted"
	SubjectKindNetworkDevice SubjectKind = "network_device"
	SubjectKindPeripheral    SubjectKind = "peripheral"
	SubjectKindMobile        SubjectKind = "mobile"
	SubjectKindUnknown       SubjectKind = "unknown"
)

type TransportAttachment string

const (
	TransportAttachmentLocal   TransportAttachment = "local"
	TransportAttachmentNetwork TransportAttachment = "network"
	TransportAttachmentProxied TransportAttachment = "proxied"
)

type ProbeStrategyConfidence string

const (
	ProbeStrategyConfidenceHigh   ProbeStrategyConfidence = "high"
	ProbeStrategyConfidenceMedium ProbeStrategyConfidence = "medium"
	ProbeStrategyConfidenceLow    ProbeStrategyConfidence = "low"
)

type VerifiedCapabilityKind string

const (
	VerifiedCapabilityInput   VerifiedCapabilityKind = "input"
	VerifiedCapabilityOutput  VerifiedCapabilityKind = "output"
	VerifiedCapabilityCompute VerifiedCapabilityKind = "compute"
	VerifiedCapabilityStorage VerifiedCapabilityKind = "storage"
	VerifiedCapabilitySensor  VerifiedCapabilityKind = "sensor"
	VerifiedCapabilityNetwork VerifiedCapabilityKind = "network"
	VerifiedCapabilityOther   VerifiedCapabilityKind = "other"
)

type VerifiedCapabilityStatus string

const (
	VerifiedCapabilityVerified VerifiedCapabilityStatus = "verified"
	VerifiedCapabilityPartial  VerifiedCapabilityStatus = "partial"
	VerifiedCapabilityRevoked  VerifiedCapabilityStatus = "revoked"
)

type RegistrationState string

const (
	RegistrationStateActive  RegistrationState = "active"
	RegistrationStatePartial RegistrationState = "partial"
	RegistrationStateOffline RegistrationState = "offline"
	RegistrationStateFailed  RegistrationState = "failed"
)

type DeviceRegistration struct {
	Subject              SubjectRegistration       `json:"subject"`
	Transport            TransportRegistration     `json:"transport"`
	ProbeStrategy        ProbeStrategyRegistration `json:"probe_strategy"`
	VerifiedCapabilities []VerifiedCapability      `json:"verified_capabilities"`
	RegistrationState    RegistrationState         `json:"registration_state"`
	LastVerifiedAt       time.Time                 `json:"last_verified_at"`
}

type SubjectRegistration struct {
	SubjectID   string         `json:"subject_id"`
	SubjectKind SubjectKind    `json:"subject_kind"`
	DisplayName string         `json:"display_name,omitempty"`
	Identity    map[string]any `json:"identity,omitempty"`
}

type TransportRegistration struct {
	Kind       string              `json:"kind"`
	Attachment TransportAttachment `json:"attachment,omitempty"`
	Details    map[string]any      `json:"details,omitempty"`
}

type ProbeStrategyRegistration struct {
	StrategyID string                  `json:"strategy_id"`
	Version    string                  `json:"version,omitempty"`
	Confidence ProbeStrategyConfidence `json:"confidence,omitempty"`
}

type VerifiedCapability struct {
	CapabilityID    string                   `json:"capability_id"`
	Name            string                   `json:"name"`
	Kind            VerifiedCapabilityKind   `json:"kind,omitempty"`
	Status          VerifiedCapabilityStatus `json:"status"`
	Interface       string                   `json:"interface,omitempty"`
	Constraints     []string                 `json:"constraints,omitempty"`
	EvidenceSummary string                   `json:"evidence_summary,omitempty"`
}
