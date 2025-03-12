// Code generated by go-swagger; DO NOT EDIT.

//
// Copyright NetFoundry Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// __          __              _
// \ \        / /             (_)
//  \ \  /\  / /_ _ _ __ _ __  _ _ __   __ _
//   \ \/  \/ / _` | '__| '_ \| | '_ \ / _` |
//    \  /\  / (_| | |  | | | | | | | | (_| | : This file is generated, do not edit it.
//     \/  \/ \__,_|_|  |_| |_|_|_| |_|\__, |
//                                      __/ |
//                                     |___/

package rest_model

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// PostureChecksFailureMfaCriteria posture checks failure mfa criteria
//
// swagger:model postureChecksFailureMfaCriteria
type PostureChecksFailureMfaCriteria struct {

	// passed mfa at
	// Required: true
	// Format: date-time
	PassedMfaAt *strfmt.DateTime `json:"passedMfaAt"`

	// timeout remaining seconds
	// Required: true
	TimeoutRemainingSeconds *int64 `json:"timeoutRemainingSeconds"`

	// timeout seconds
	// Required: true
	TimeoutSeconds *int64 `json:"timeoutSeconds"`

	// unlocked at
	// Required: true
	// Format: date-time
	UnlockedAt *strfmt.DateTime `json:"unlockedAt"`

	// woken at
	// Required: true
	// Format: date-time
	WokenAt *strfmt.DateTime `json:"wokenAt"`
}

// Validate validates this posture checks failure mfa criteria
func (m *PostureChecksFailureMfaCriteria) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validatePassedMfaAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTimeoutRemainingSeconds(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTimeoutSeconds(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUnlockedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateWokenAt(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PostureChecksFailureMfaCriteria) validatePassedMfaAt(formats strfmt.Registry) error {

	if err := validate.Required("passedMfaAt", "body", m.PassedMfaAt); err != nil {
		return err
	}

	if err := validate.FormatOf("passedMfaAt", "body", "date-time", m.PassedMfaAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *PostureChecksFailureMfaCriteria) validateTimeoutRemainingSeconds(formats strfmt.Registry) error {

	if err := validate.Required("timeoutRemainingSeconds", "body", m.TimeoutRemainingSeconds); err != nil {
		return err
	}

	return nil
}

func (m *PostureChecksFailureMfaCriteria) validateTimeoutSeconds(formats strfmt.Registry) error {

	if err := validate.Required("timeoutSeconds", "body", m.TimeoutSeconds); err != nil {
		return err
	}

	return nil
}

func (m *PostureChecksFailureMfaCriteria) validateUnlockedAt(formats strfmt.Registry) error {

	if err := validate.Required("unlockedAt", "body", m.UnlockedAt); err != nil {
		return err
	}

	if err := validate.FormatOf("unlockedAt", "body", "date-time", m.UnlockedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *PostureChecksFailureMfaCriteria) validateWokenAt(formats strfmt.Registry) error {

	if err := validate.Required("wokenAt", "body", m.WokenAt); err != nil {
		return err
	}

	if err := validate.FormatOf("wokenAt", "body", "date-time", m.WokenAt.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this posture checks failure mfa criteria based on context it is used
func (m *PostureChecksFailureMfaCriteria) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *PostureChecksFailureMfaCriteria) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PostureChecksFailureMfaCriteria) UnmarshalBinary(b []byte) error {
	var res PostureChecksFailureMfaCriteria
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
