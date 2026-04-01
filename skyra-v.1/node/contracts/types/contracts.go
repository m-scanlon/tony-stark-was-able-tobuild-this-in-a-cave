package types

import (
	node ".."
	primitives "../../../protocol/primitives/contracts"
	stimulus "../../../stimulus/contracts"
)

type ProbeNodeContract struct {
	NodeType        string             `json:"node_type"`
	Purpose         node.NodePurpose   `json:"purpose"`
	Stimulus        node.NodeStimulus  `json:"stimulus"`
	Cognition       node.NodeCognition `json:"cognition"`
	Commands        node.NodeCommands  `json:"commands"`
	LearningEnabled bool               `json:"learning_enabled,omitempty"`
}

type RegistrationNodeContract struct {
	NodeType        string             `json:"node_type"`
	Purpose         node.NodePurpose   `json:"purpose"`
	Stimulus        node.NodeStimulus  `json:"stimulus"`
	Cognition       node.NodeCognition `json:"cognition"`
	Commands        node.NodeCommands  `json:"commands"`
	LearningEnabled bool               `json:"learning_enabled,omitempty"`
}

var ProbeContract = ProbeNodeContract{
	NodeType: "probe",
	Purpose: node.NodePurpose{
		Summary: "Discover candidate capabilities on a system subject, verify them through bounded invocation, and shape initial capability contracts from observed behavior.",
		Limits: []string{
			"Does not persist registration truth",
			"Does not birth nodes",
			"Does not perform unconstrained exploration",
		},
	},
	Stimulus: node.NodeStimulus{
		AcceptedTypes: []string{stimulus.StimulusTypeDeviceProbeRequest},
		EmittedTypes:  []string{stimulus.StimulusTypeDeviceProbeResult},
	},
	Cognition: node.NodeCognition{
		Mode:     "bounded_probe",
		MaxSteps: 1,
	},
	Commands: node.NodeCommands{
		AllowedCommands: []primitives.PrimitiveName{
			primitives.PrimitiveInteract,
			primitives.PrimitiveRecall,
		},
	},
	LearningEnabled: true,
}

var RegistrationContract = RegistrationNodeContract{
	NodeType: "registration",
	Purpose: node.NodePurpose{
		Summary: "Assemble the typed device registration envelope for a system subject from probe output and persist that envelope through the world-facing registration write path.",
		Limits: []string{
			"Does not rediscover capabilities on its own",
			"Does not shape initial capability contracts from scratch",
			"Does not birth nodes",
		},
	},
	Stimulus: node.NodeStimulus{
		AcceptedTypes: []string{stimulus.StimulusTypeDeviceProbeResult},
		EmittedTypes:  []string{stimulus.StimulusTypeRegistrationWriteResult},
	},
	Cognition: node.NodeCognition{
		Mode:     "bounded_registration",
		MaxSteps: 1,
	},
	Commands: node.NodeCommands{
		AllowedCommands: []primitives.PrimitiveName{
			primitives.PrimitiveInteract,
			primitives.PrimitiveRecall,
		},
	},
	LearningEnabled: true,
}
