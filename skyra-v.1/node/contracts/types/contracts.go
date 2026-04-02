package types

import (
	node ".."
	stimulus "../../../stimulus/contracts"
)

type ProbeNodeContract struct {
	Purpose      node.NodePurpose      `json:"purpose"`
	Capabilities node.NodeCapabilities `json:"capabilities"`
	Stimulus     node.NodeStimulus     `json:"stimulus"`
}

type RegistrationNodeContract struct {
	Purpose      node.NodePurpose      `json:"purpose"`
	Capabilities node.NodeCapabilities `json:"capabilities"`
	Stimulus     node.NodeStimulus     `json:"stimulus"`
}

var ProbeContract = ProbeNodeContract{
	Purpose: node.NodePurpose{
		Summary: "Discover candidate capabilities on a system subject, verify them through bounded invocation, and shape initial capability contracts from observed behavior.",
	},
	Capabilities: node.NodeCapabilities{
		CapabilityIDs: []string{
			"act",
			"recall",
		},
	},
	Stimulus: node.NodeStimulus{
		AcceptedTypes: []string{stimulus.StimulusTypeDeviceProbeRequest},
		EmittedTypes:  []string{stimulus.StimulusTypeDeviceProbeResult},
	},
}

var RegistrationContract = RegistrationNodeContract{
	Purpose: node.NodePurpose{
		Summary: "Assemble the typed device registration envelope for a system subject from probe output and persist that envelope through the world-facing registration write path.",
	},
	Capabilities: node.NodeCapabilities{
		CapabilityIDs: []string{
			"act",
			"recall",
		},
	},
	Stimulus: node.NodeStimulus{
		AcceptedTypes: []string{stimulus.StimulusTypeDeviceProbeResult},
		EmittedTypes:  []string{stimulus.StimulusTypeRegistrationWriteResult},
	},
}
