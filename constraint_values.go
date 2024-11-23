package CADDY_FILE_SERVER

import (
	"errors"
	"fmt"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"slices"
	"strconv"
)

func init() {
	RegisterConstraintType(func() Constraint {
		return new(ValuesConstraint)
	})
}

type ValuesConstraint struct {
	Values []int `json:"values"`
}

func (r *ValuesConstraint) ID() string {
	return "values"
}

func (r *ValuesConstraint) Validate(param string) error {
	if !slices.Contains([]string{"w", "h", "q", "ah", "aw", "t", "l", "r", "b"}, param) {
		return fmt.Errorf("values constraint cannot be applied on param: '%s'", param)
	}
	if len(r.Values) == 0 {
		return errors.New("you need to provide at least one value for values constraint")
	}
	return nil
}

func (r *ValuesConstraint) ValidateParam(param string, value string) error {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("invalid integer value for %s: %s", param, value)
	}

	if !slices.Contains(r.Values, intValue) {
		return fmt.Errorf("parameter %s has an invalid value: %d", param, intValue)
	}

	return nil
}

func (r *ValuesConstraint) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	values := d.RemainingArgs()
	r.Values = make([]int, len(values))
	for idx, v := range values {
		var err error
		if r.Values[idx], err = strconv.Atoi(v); err != nil {
			return err
		}
	}

	return nil
}
