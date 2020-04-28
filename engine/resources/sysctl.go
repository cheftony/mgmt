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

// TODOs:
// - Find a cleaner way to store and use the path from parameter
// - Use obj.Name() if parameter is not defined
// - Have a better understading of makeComposite()

package resources

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/purpleidea/mgmt/engine"
	"github.com/purpleidea/mgmt/engine/traits"
	"github.com/purpleidea/mgmt/recwatch"
	"github.com/purpleidea/mgmt/util/errwrap"
)

func init() {
	engine.RegisterResource("sysctl", func() engine.Res { return &SysctlRes{} })
}

const (
	// sysctlDir is the location of the kernel parameters.
	sysctlDir  = "/proc/sys/"
	persistDir = "/etc/sysctl.d/"
	fileName   = "10-%s.conf"
)

// SysctlRes is a resource that manage kernel parameters
type SysctlRes struct {
	traits.Base // add the base methods without re-implementation

	init *engine.Init

	Parameter  string `lang:"parameter" yaml:"parameter"`
	Value      string `lang:"value" yaml:"value"`
	Persistent bool   `lang:"persistent" yaml:"persistent"`

	recWatcher *recwatch.RecWatcher
	file       *FileRes
}

// getPath returns the actual path to use for this parameter.
func (obj *SysctlRes) getPath() string {
	p := obj.Parameter
	if obj.Parameter == "" { // use the name as the parameter default if missing
		p = obj.Name()
	}

	return filepath.Join(sysctlDir, strings.Replace(p, ".", "/", -1))
}

// Default returns some sensible defaults for this resource.
func (obj *SysctlRes) Default() engine.Res {
	return &SysctlRes{
		Persistent: true,
	}
}

// makeComposite creates a pointer to a FileRes. The pointer is used to validate
// and initialize the nested file resource and to apply the file state in
// CheckApply.
// func (obj *SysctlRes) makeComposite() (*FileRes, error) {
// 	res, err := engine.NewNamedResource("file", fmt.Sprintf(persistDir+fileName, obj.Name()))
// 	if err != nil {
// 		return nil, errwrap.Wrapf(err, "error creating nested file resource")
// 	}
// 	file, ok := res.(*FileRes)
// 	if !ok {
// 		return nil, fmt.Errorf("error casting file resource")
// 	}

// 	if obj.Persistent {
// 		file.State = "exists"
// 	} else {
// 		file.State = "absent"
// 	}

// 	if obj.Persistent {
// 		// TODO: Build the content of the file
// 		// file.Content = &s
// 	}
// 	return file, nil
// }

// Validate reports any problems with the struct definition.
func (obj *SysctlRes) Validate() error {
	p := obj.getPath()

	if p == "" {
		return fmt.Errorf("parameter is empty")
	}
	if obj.Value == "" {
		return fmt.Errorf("value is empty")
	}
	if _, err := os.Stat(p); err != nil { // check if path for parameter exists
		return fmt.Errorf("kernel parameter does not exist: %v", p)
	}

	fmt.Println(obj.Name())
	return nil
}

// Init runs some startup code for this resource.
func (obj *SysctlRes) Init(init *engine.Init) error {
	obj.init = init // save for later

	return nil
}

// Close is run by the engine to clean up after the resource is done.
func (obj *SysctlRes) Close() error {
	return nil
}

// CheckApply does the idempotent work of checking and applying resource state.
func (obj *SysctlRes) CheckApply(apply bool) (bool, error) {
	p := obj.getPath()

	value, err := ioutil.ReadFile(p)
	if err != nil {
		return false, err
	}

	if string(value) == obj.Value {
		return true, nil
	}

	if !apply {
		return false, nil
	}

	err = ioutil.WriteFile(p, []byte(obj.Value), 0644)
	if err != nil {
		return false, err
	}

	return false, nil
}

// Watch is the listener and main loop for this resource.
func (obj *SysctlRes) Watch() error {
	var err error
	obj.recWatcher, err = recwatch.NewRecWatcher(obj.getPath(), false)
	if err != nil {
		return err
	}
	defer obj.recWatcher.Close()

	// To Delete
	obj.init.Logf("test: name: %v parameter: %v value: %v", obj.Name(), obj.Parameter, obj.Value)

	obj.init.Logf("XXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
	fmt.Printf("YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY")
	obj.init.Running() // when started, notify engine that we're running

	var send = false // send event?

	for {
		if obj.init.Debug {
			obj.init.Logf("watching: %s", obj.getPath()) // attempting to watch...
		}

		select {
		case event, ok := <-obj.recWatcher.Events():
			if !ok {
				return fmt.Errorf("unexpected close")
			}
			if err := event.Error; err != nil {
				return errwrap.Wrapf(err, "unknown %s watcher error", obj)
			}
			if obj.init.Debug { // don't access event.Body if event.Error isn't nil
				obj.init.Logf("event(%s): %v", event.Body.Name, event.Body.Op)
			}
			send = true
		case <-obj.init.Done: // signal for shutdown request
			return nil
		}

		// do all our event sending all together to avoid duplicate msgs
		if send {
			send = false
			obj.init.Event()
		}
	}

	return nil
}

// Cmp compares two resources and returns an error if they are not equivalent.
func (obj *SysctlRes) Cmp(r engine.Res) error {
	res, ok := r.(*SysctlRes)
	if !ok {
		return fmt.Errorf("not a %s", obj.Kind())
	}
	if res.Parameter != obj.Parameter {
		return fmt.Errorf("the Paramater differs")
	}
	if res.Value != obj.Value {
		return fmt.Errorf("the Value differs")
	}
	return nil
}

// UnmarshalYAML is the custom unmarshal handler for this struct. It is
// primarily useful for setting the defaults.
func (obj *SysctlRes) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawRes SysctlRes // indirection to avoid infinite recursion

	def := obj.Default()        // get the default
	res, ok := def.(*SysctlRes) // put in the right format
	if !ok {
		return fmt.Errorf("could not convert to SysctlRes")
	}
	raw := rawRes(*res) // convert; the defaults go here

	if err := unmarshal(&raw); err != nil {
		return err
	}

	*obj = SysctlRes(raw) // restore from indirection with type conversion!
	return nil
}
