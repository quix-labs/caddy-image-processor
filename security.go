package CADDY_FILE_SERVER

import (
	"cmp"
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"net/url"
	"slices"
)

type OnSecurityFail string

const (
	// OnSecurityFailIgnore Deletes invalid parameters from the request but continues processing.
	OnSecurityFailIgnore OnSecurityFail = "ignore"

	// OnSecurityFailAbort Returns a 400 Bad Request error to the client.
	OnSecurityFailAbort OnSecurityFail = "abort"

	// OnSecurityFailBypass Forces the response to return the initial (unprocessed) image.
	OnSecurityFailBypass OnSecurityFail = "bypass"
)

type SecurityOptions struct {
	OnSecurityFail   OnSecurityFail `json:"on_security_fail,omitempty"`
	AllowedParams    *[]string      `json:"allowed_params,omitempty"`
	DisallowedParams *[]string      `json:"disallowed_params,omitempty"`
	Constraints      *Constraints   `json:"constraints,omitempty"`
}

// ProcessRequestForm
// Ensures that all security constraints are applied.
// May also remove specific parameters if they are not allowed.
func (s *SecurityOptions) ProcessRequestForm(form *url.Values) error {

	// If 'allowed' is specified, retain only the specified elements.
	if s.AllowedParams != nil {
		for param, _ := range *form {
			if !slices.Contains(*s.AllowedParams, param) {
				if s.OnSecurityFail == OnSecurityFailIgnore {
					form.Del(param)
				} else if s.OnSecurityFail == OnSecurityFailBypass {
					return BypassRequestError
				} else if s.OnSecurityFail == OnSecurityFailAbort {
					return &AbortRequestError{
						fmt.Sprintf("parameter '%s' is not allowed", param),
					}
				}
			}
		}
	}

	// If 'allowed' is not specified, remove the elements.
	if s.DisallowedParams != nil {
		for _, param := range *s.DisallowedParams {
			if form.Has(param) {
				if s.OnSecurityFail == OnSecurityFailIgnore {
					form.Del(param)
				} else if s.OnSecurityFail == OnSecurityFailBypass {
					return BypassRequestError
				} else if s.OnSecurityFail == OnSecurityFailAbort {
					return &AbortRequestError{
						fmt.Sprintf("parameter '%s' has been flagged as disallowed", param),
					}
				}
			}
		}
	}

	if s.Constraints != nil {
		if err := s.Constraints.ProcessRequestForm(form, s.OnSecurityFail); err != nil {
			return err
		}
	}

	return nil
}

// Provision Set default values if not defined
func (s *SecurityOptions) Provision(ctx caddy.Context) error {
	s.OnSecurityFail = cmp.Or(s.OnSecurityFail, OnSecurityFailIgnore)
	return nil
}

// Validate ensure security parameters are correctly defined
func (s *SecurityOptions) Validate() error {
	switch s.OnSecurityFail {
	case OnSecurityFailIgnore, OnSecurityFailAbort, OnSecurityFailBypass:
		// Valid values
	default:
		return fmt.Errorf("invalid value for 'on_security_fail': '%s' (expected 'ignore', 'abort', or 'bypass')", s.OnSecurityFail)
	}

	// Validate constraints if exists
	if s.Constraints != nil {
		if err := s.Constraints.Validate(); err != nil {
			return err
		}
	}

	// Check that AllowedParams and DisallowedParams are not both specified
	if s.AllowedParams != nil && s.DisallowedParams != nil {
		return fmt.Errorf("'allowed_params' and 'disallowed_params' cannot be specified together")
	}

	// Ensure that at least one of AllowedParams or DisallowedParams or 'Constraints' is specified
	if (s.AllowedParams == nil || len(*s.AllowedParams) == 0) &&
		(s.DisallowedParams == nil || len(*s.DisallowedParams) == 0) &&
		(s.Constraints == nil || len(*s.Constraints) == 0) {
		return fmt.Errorf("either 'allowed_params', 'disallowed_params', or 'constraints' must be specified")
	}

	// Validate that all elements in AllowedParams are in availableParams
	if s.AllowedParams != nil {
		for _, param := range *s.AllowedParams {
			if !slices.Contains(availableParams, param) {
				return fmt.Errorf("unknown parameter '%s' in 'allowed_params'", param)
			}
		}
	}

	// Validate that all elements in DisallowedParams are in availableParams
	if s.DisallowedParams != nil {
		for _, param := range *s.DisallowedParams {
			if !slices.Contains(availableParams, param) {
				return fmt.Errorf("unknown parameter '%s' in 'disallowed_params'", param)
			}
		}
	}

	// Verify OnSecurityFail is defined with a valid value
	switch s.OnSecurityFail {
	case "ignore", "abort", "bypass":
		// Valid values
	default:
		return fmt.Errorf("invalid value for 'on_security_fail': '%s'. Valid values are 'ignore', 'abort', or 'bypass'", s.OnSecurityFail)
	}

	return nil
}

func (s *SecurityOptions) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for nesting := d.Nesting(); d.NextBlock(nesting); {
		switch d.Val() {
		case "on_security_fail":
			// Check if argument provided
			if !d.NextArg() {
				return d.ArgErr()
			}
			s.OnSecurityFail = OnSecurityFail(d.Val())

			// Ensure there are no more arguments
			if d.NextArg() {
				return d.ArgErr() // More than one argument provided
			}
			break
		case "allowed_params":
			allowedParams := d.RemainingArgs()
			if len(allowedParams) == 0 {
				return d.Err("allowed_params requires at least one parameter")
			}
			s.AllowedParams = &allowedParams
			break
		case "disallowed_params":
			disallowedParams := d.RemainingArgs()
			if len(disallowedParams) == 0 {
				return d.Err("disallowed_params requires at least one parameter")
			}
			s.DisallowedParams = &disallowedParams
			break
		case "constraints":
			// If it's a nested block, process it
			s.Constraints = &Constraints{}
			if err := s.Constraints.UnmarshalCaddyfile(d); err != nil {
				return err
			}
			break
		default:
			return d.Errf("unexpected directive '%s' in security block", d.Val())
		}
	}

	return nil
}
