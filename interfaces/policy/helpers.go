// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2016 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package policy

import (
	"fmt"

	"github.com/snapcore/snapd/asserts"
	"github.com/snapcore/snapd/snap"
)

// check helpers

func checkSnapType(snapType snap.Type, types []string) error {
	if len(types) == 0 {
		return nil
	}
	s := string(snapType)
	if s == "os" { // we use "core" in the assertions
		s = "core"
	}
	for _, t := range types {
		if t == s {
			return nil
		}
	}
	return fmt.Errorf("snap type does not match")
}

func checkPlugConnectionConstraints(connc *ConnectCandidate, cstrs *asserts.PlugConnectionConstraints) error {
	if err := cstrs.PlugAttributes.Check(connc.plugAttrs()); err != nil {
		return err
	}
	if err := cstrs.SlotAttributes.Check(connc.slotAttrs()); err != nil {
		return err
	}
	if err := checkSnapType(connc.slotSnapType(), cstrs.SlotSnapTypes); err != nil {
		return err
	}
	// TODO: IDs constraints
	return nil
}

func checkSlotConnectionConstraints(connc *ConnectCandidate, cstrs *asserts.SlotConnectionConstraints) error {
	if err := cstrs.PlugAttributes.Check(connc.plugAttrs()); err != nil {
		return err
	}
	if err := cstrs.SlotAttributes.Check(connc.slotAttrs()); err != nil {
		return err
	}
	if err := checkSnapType(connc.plugSnapType(), cstrs.PlugSnapTypes); err != nil {
		return err
	}
	// TODO: IDs constraints
	return nil
}
