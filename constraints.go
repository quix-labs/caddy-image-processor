package CADDY_FILE_SERVER

import (
	"encoding/json"
	"fmt"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"net/url"
)

var constraintsRegistry = make(map[string]func() Constraint)

func RegisterConstraintType(factory func() Constraint) {
	typ := factory().ID()
	constraintsRegistry[typ] = factory
}

// Constraints represent constraints as {params:[{type:..., customConfig..}]}
type Constraints map[string][]Constraint

type Constraint interface {
	Validate(param string) error
	ValidateParam(param string, value string) error
	UnmarshalCaddyfile(d *caddyfile.Dispenser) error
	ID() string
}

// Temporary map to hold serialized constraints with type as a key

// MarshalJSON serializes Constraints to JSON, adding a `type` field to each entry.
func (cs *Constraints) MarshalJSON() ([]byte, error) {
	type jsonConstraintsWrapper map[string][]map[string]Constraint

	// Initialize the map for wrapped constraints
	wrappedConstraints := make(jsonConstraintsWrapper, len(*cs))

	for param, constraintList := range *cs {
		var wrappedList []map[string]Constraint
		for _, constraint := range constraintList {
			wrappedList = append(wrappedList, map[string]Constraint{
				constraint.ID(): constraint,
			})
		}
		wrappedConstraints[param] = wrappedList
	}

	return json.Marshal(wrappedConstraints)
}

// UnmarshalJSON deserializes JSON into Constraints, dynamically instantiating types using the registry.
func (cs *Constraints) UnmarshalJSON(data []byte) error {
	type jsonConstraintsUnwrapper map[string][]map[string]json.RawMessage

	var wrappedConstraints jsonConstraintsUnwrapper
	if err := json.Unmarshal(data, &wrappedConstraints); err != nil {
		return err
	}

	*cs = make(Constraints, len(wrappedConstraints))

	for param, wrappedList := range wrappedConstraints {
		(*cs)[param] = make([]Constraint, len(wrappedList))
		for idx, wrappedConstraint := range wrappedList {
			for constraintType, constraintData := range wrappedConstraint {
				// Look up the factory function for the given constraint type
				factory, found := constraintsRegistry[constraintType]
				if !found {
					return fmt.Errorf("unknown constraint type: %s for param: %s", constraintType, param)
				}

				// Instantiate the correct constraint type
				constraint := factory()

				// Unmarshal the constraint data into the instantiated constraint
				if err := json.Unmarshal(constraintData, constraint); err != nil {
					return fmt.Errorf("error unmarshaling constraint for param %s: %v", param, err)
				}
				// Add the deserialized constraint to the slice for this parameter
				(*cs)[param][idx] = constraint
			}
		}
	}
	return nil
}

func (cs *Constraints) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for nesting := d.Nesting(); d.NextBlock(nesting); {
		param := d.Val()
		var constraintsForParam []Constraint

		var nested bool

		// Try to detect and process a nested block if any (like `w { range ... , range ... }`)
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			nested = true
			constraintName := d.Val()

			factory, found := constraintsRegistry[constraintName]
			if !found {
				return d.Errf("unknown constraint type: %s", constraintName)
			}

			constraint := factory()

			if err := constraint.UnmarshalCaddyfile(d); err != nil {
				return d.Errf("error unmarshaling parameters for %s constraint: %v", constraintName, err)
			}

			constraintsForParam = append(constraintsForParam, constraint)
		}

		// If no nested block was found, process inline arguments (like `range 10 20`)
		if !nested {
			if !d.NextArg() {
				return d.Errf("missing constraint name for parameter %s", param)
			}

			constraintName := d.Val()

			factory, found := constraintsRegistry[constraintName]
			if !found {
				return d.Errf("unknown constraint type: %s", constraintName)
			}

			constraint := factory()

			if err := constraint.UnmarshalCaddyfile(d); err != nil {
				return d.Errf("error unmarshaling parameters for %s constraint: %v", constraintName, err)
			}

			if d.NextArg() {
				return d.ArgErr()
			}

			constraintsForParam = append(constraintsForParam, constraint)

		}

		(*cs)[param] = constraintsForParam
	}
	return nil
}

func (cs *Constraints) Validate() error {
	for param, constraints := range *cs {
		for _, constraint := range constraints {
			if err := constraint.Validate(param); err != nil {
				return err
			}
		}
	}
	return nil
}

func (cs *Constraints) ProcessRequestForm(form *url.Values, onSecurityFail OnSecurityFail) error {
	for param, constraints := range *cs {
		if !form.Has(param) {
			continue
		}

		for _, constraint := range constraints {
			if err := constraint.ValidateParam(param, form.Get(param)); err != nil {
				if onSecurityFail == OnSecurityFailIgnore {
					form.Del(param)
				} else if onSecurityFail == OnSecurityFailBypass {
					return BypassRequestError
				} else if onSecurityFail == OnSecurityFailAbort {
					return &AbortRequestError{
						err.Error(),
					}
				}

				return err
			}
		}
	}
	return nil
}
