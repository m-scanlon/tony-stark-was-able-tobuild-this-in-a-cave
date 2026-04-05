package stimulus

import (
	capability "../../capability/contracts"
	registration "../../registration/contracts"
)

const (
	StimulusTypeBootstrapFingerprint     = "bootstrap_fingerprint"
	StimulusTypeDeviceProbeRequest       = "device_probe_request"
	StimulusTypeDeviceProbeResult        = "device_probe_result"
	StimulusTypeRegistrationWriteRequest = "registration_write_request"
	StimulusTypeRegistrationWriteResult  = "registration_write_result"
)

type BootstrapFingerprintPayload struct {
	SubjectID          string   `json:"subject_id"`
	PlatformFamily     string   `json:"platform_family,omitempty"`
	PlatformVersion    string   `json:"platform_version,omitempty"`
	Architecture       string   `json:"architecture,omitempty"`
	HostKind           string   `json:"host_kind,omitempty"`
	ObservedTransports []string `json:"observed_transports,omitempty"`
}

type DeviceProbeRequestPayload struct {
	SubjectID  string `json:"subject_id"`
	ProbeClass string `json:"probe_class,omitempty"`
}

type DeviceProbeCapability struct {
	VerifiedCapability registration.VerifiedCapability `json:"verified_capability"`
	CapabilityContract *capability.CapabilityContract  `json:"capability_contract,omitempty"`
}

type DeviceProbeResultPayload struct {
	SubjectID       string                  `json:"subject_id"`
	ProbeStrategyID string                  `json:"probe_strategy_id,omitempty"`
	Capabilities    []DeviceProbeCapability `json:"capabilities,omitempty"`
}

type RegistrationWriteRequestPayload struct {
	Registration registration.DeviceRegistration `json:"registration"`
}

type RegistrationWriteResultPayload struct {
	SubjectID         string                         `json:"subject_id"`
	RegistrationState registration.RegistrationState `json:"registration_state"`
}
