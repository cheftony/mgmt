// Mgmt
// Copyright (C) 2013-2020+ James Shubin and the project contributors
// Written by James Shubin <james@shubin.ca> and the project contributors
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package resources

import (
	"reflect"
	"testing"

	"github.com/purpleidea/mgmt/engine"
	"github.com/purpleidea/mgmt/engine/traits"
	"github.com/purpleidea/mgmt/recwatch"
)

func TestSysctlRes_makeComposite(t *testing.T) {
	type fields struct {
		Base       traits.Base
		init       *engine.Init
		Parameter  string
		Value      string
		Persistent bool
		recWatcher *recwatch.RecWatcher
		file       *FileRes
	}
	tests := []struct {
		name    string
		fields  fields
		want    *FileRes
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := &SysctlRes{
				Base:       tt.fields.Base,
				init:       tt.fields.init,
				Parameter:  tt.fields.Parameter,
				Value:      tt.fields.Value,
				Persistent: tt.fields.Persistent,
				recWatcher: tt.fields.recWatcher,
				file:       tt.fields.file,
			}
			got, err := obj.makeComposite()
			if (err != nil) != tt.wantErr {
				t.Errorf("SysctlRes.makeComposite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SysctlRes.makeComposite() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSysctlRes_Validate(t *testing.T) {
	type fields struct {
		Base       traits.Base
		init       *engine.Init
		Parameter  string
		Value      string
		Persistent bool
		recWatcher *recwatch.RecWatcher
		file       *FileRes
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := &SysctlRes{
				Base:       tt.fields.Base,
				init:       tt.fields.init,
				Parameter:  tt.fields.Parameter,
				Value:      tt.fields.Value,
				Persistent: tt.fields.Persistent,
				recWatcher: tt.fields.recWatcher,
				file:       tt.fields.file,
			}
			if err := obj.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("SysctlRes.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
