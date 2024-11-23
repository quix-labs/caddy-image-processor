package CADDY_FILE_SERVER

import (
	"fmt"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"slices"
	"strconv"
)

func init() {
	RegisterConstraintType(func() Constraint {
		return new(RangeConstraint)
	})
}

type RangeConstraint struct {
	From int `json:"from,omitempty"`
	To   int `json:"to,omitempty"`
}

func (r *RangeConstraint) ID() string {
	return "range"
}

func (r *RangeConstraint) Validate(param string) error {
	if !slices.Contains([]string{"w", "h", "q", "ah", "aw", "t", "l", "r", "b"}, param) {
		return fmt.Errorf("range constraint cannot be applied on param: '%s'", param)
	}
	if r.From < 0 {
		return fmt.Errorf("range constraint must have minimum value less than 0")
	}
	if r.From >= r.To {
		return fmt.Errorf("range constraint must have minimum value less than max")
	}
	return nil
}

func (r *RangeConstraint) ValidateParam(param string, value string) error {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("invalid integer value for %s: %s", param, value)
	}

	if intValue < r.From || intValue > r.To {
		return fmt.Errorf("%s must be in range %d to %d", param, r.From, r.To)
	}

	return nil
}

func (r *RangeConstraint) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	var nested bool

	// Try to load nested block if present
	for nesting := d.Nesting(); d.NextBlock(nesting); {
		nested = true
		param := d.Val()

		switch param {
		case "from":
			if !d.NextArg() {
				return d.Err("missing value for from")
			}
			var err error
			r.From, err = strconv.Atoi(d.Val())
			if err != nil {
				return d.Errf("invalid from value for range: %v", err)
			}
		case "to":
			if !d.NextArg() {
				return d.Err("missing value for to")
			}
			var err error
			r.To, err = strconv.Atoi(d.Val())
			if err != nil {
				return d.Errf("invalid to value for range: %v", err)
			}
		default:
			return d.Errf("unexpected parameter '%s' in range constraint", param)
		}
	}

	// If not a nested block, process inline arguments
	if !nested {
		if !d.NextArg() {
			return d.Err("missing from value for range constraint")
		}
		var err error
		r.From, err = strconv.Atoi(d.Val())
		if err != nil {
			return d.Errf("invalid from value for range: %v", err)
		}

		if !d.NextArg() {
			return d.Err("missing to value for range constraint")
		}
		r.To, err = strconv.Atoi(d.Val())
		if err != nil {
			return d.Errf("invalid to value for range: %v", err)
		}

		if d.NextArg() {
			return d.ArgErr()
		}
	}
	return nil
}
