package types

import (
	actor ".."
	stimulus "../../../stimulus/contracts"
)

type ProbeActorContract struct {
	Purpose      actor.ActorPurpose      `json:"purpose"`
	Commitments  []string                `json:"commitments,omitempty"`
	Capabilities actor.ActorCapabilities `json:"capabilities"`
	Stimulus     actor.ActorStimulus     `json:"stimulus"`
}

type RegistrationActorContract struct {
	Purpose      actor.ActorPurpose      `json:"purpose"`
	Commitments  []string                `json:"commitments,omitempty"`
	Capabilities actor.ActorCapabilities `json:"capabilities"`
	Stimulus     actor.ActorStimulus     `json:"stimulus"`
}

var ProbeContract = ProbeActorContract{
	Purpose: actor.ActorPurpose{
		Summary: "Discover candidate capabilities on a system subject, verify them through bounded invocation, and shape initial capability contracts from observed behavior.",
	},
	Capabilities: actor.ActorCapabilities{
		CapabilityIDs: []string{
			"act",
			"recall",
		},
	},
	Stimulus: actor.ActorStimulus{
		AcceptedTypes: []string{stimulus.StimulusTypeDeviceProbeRequest},
		EmittedTypes:  []string{stimulus.StimulusTypeDeviceProbeResult},
	},
}

var RegistrationContract = RegistrationActorContract{
	Purpose: actor.ActorPurpose{
		Summary: "Assemble the typed device registration envelope for a system subject from probe output and persist that envelope through the world-facing registration write path.",
	},
	Capabilities: actor.ActorCapabilities{
		CapabilityIDs: []string{
			"act",
			"recall",
		},
	},
	Stimulus: actor.ActorStimulus{
		AcceptedTypes: []string{stimulus.StimulusTypeDeviceProbeResult},
		EmittedTypes:  []string{stimulus.StimulusTypeRegistrationWriteResult},
	},
}
