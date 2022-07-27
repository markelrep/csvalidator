package schema

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-multierror"

	"github.com/stretchr/testify/assert"
)

func TestContains_Contain(t *testing.T) {
	// TODO: this test can suddenly fail, because of map is using under `exactContain` therefore error order is not constant
	cases := []struct {
		values   map[string]struct{}
		contains exactContain
		expected func() error
	}{
		{
			values: map[string]struct{}{
				"value1": {},
				"value2": {},
				"value3": {},
				"value4": {},
			},
			contains: exactContain{"value1", "value2", "value3", "value4"},
			expected: func() error {
				return nil
			},
		},
		{
			values: map[string]struct{}{
				"value1": {},
				"value2": {},
				"value3": {},
			},
			contains: exactContain{"value1", "value2", "value3", "value4"},
			expected: func() (err error) {
				err = multierror.Append(err, fmt.Errorf("value4 is defined in schema, but absent in column"))
				return err
			},
		},
		{
			values: map[string]struct{}{
				"value1": {},
				"value2": {},
				"value3": {},
				"value4": {},
			},
			contains: exactContain{"value1", "value2", "value5", "value3", "value4"},
			expected: func() (err error) {
				err = multierror.Append(err, fmt.Errorf("value5 is defined in schema, but absent in column"))
				return err
			},
		},
		{
			values: map[string]struct{}{
				"value1": {},
				"value2": {},
				"value3": {},
			},
			contains: exactContain{"value2", "value3", "value4"},
			expected: func() (err error) {
				err = multierror.Append(err, fmt.Errorf("value4 is defined in schema, but absent in column"))
				err = multierror.Append(err, fmt.Errorf("value1 is not defined in schema, but exist in column"))
				return err
			},
		},
	}

	for _, tc := range cases {
		assert.Equal(t, tc.expected(), tc.contains.Contain(tc.values))
	}
}

func TestContains_IsNoOp(t *testing.T) {
	cases := []struct {
		contains exactContain
		expected bool
	}{
		{
			contains: nil,
			expected: true,
		},
		{
			contains: exactContain{},
			expected: true,
		},
		{
			contains: exactContain{"somevalue"},
			expected: false,
		},
	}

	for _, tc := range cases {
		assert.Equal(t, tc.expected, tc.contains.IsNoOp())
	}
}
