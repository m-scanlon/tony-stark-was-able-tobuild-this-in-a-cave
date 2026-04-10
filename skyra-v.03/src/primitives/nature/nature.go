package nature

import (
	"fmt"
	"strings"

	"skyra-v03/src/primitives/identity"
	"skyra-v03/src/primitives/purpose"
)

type Nature struct {
	Identity identity.Identity
	Purpose  purpose.Purpose
}

func (n Nature) Validate() error {
	if strings.TrimSpace(n.Identity.Value) == "" {
		return fmt.Errorf("nature: identity is required")
	}
	if strings.TrimSpace(n.Purpose.Value) == "" {
		return fmt.Errorf("nature: purpose is required")
	}
	return nil
}

func CreateNature(expression string) (Nature, error) {
	id, err := identity.CreateIdentity(expression)
	if err != nil {
		return Nature{}, err
	}

	p, err := purpose.CreatePurpose(expression)
	if err != nil {
		return Nature{}, err
	}

	return Nature{
		Identity: id,
		Purpose:  p,
	}, nil
}
