// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2015 Canonical Ltd
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

package partition

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

const (
	bootBase        = "/boot"
	ubootDir        = bootBase + "/uboot"
	grubDir         = bootBase + "/grub"
	ubootConfigFile = ubootDir + "/snappy-system.txt"
	grubConfigFile  = grubDir + "/grubenv"
)

// BootSystem returns the name of the boot system, grub or uboot.
func BootSystem() (string, error) {
	matches, err := filepath.Glob(bootBase + "/grub")
	if err != nil {
		return "", err
	}
	if len(matches) == 1 {
		return "grub", nil
	}
	return "uboot", nil
}

// BootDir returns the directory used by the boot system.
func BootDir(bootSystem string) string {
	if bootSystem == "grub" {
		return grubDir
	}
	return ubootDir
}

// CurrentPartition returns the current partition, a or b.
func CurrentPartition() (partition string, err error) {
	bootConfigFile, err := bootConf()
	if err != nil {
		return
	}
	file, err := os.Open(bootConfigFile)
	if err != nil {
		return
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "snappy_ab") {
			fields := strings.Split(scanner.Text(), "=")
			if len(fields) > 1 {
				var system string
				system, err = BootSystem()
				if err != nil {
					return
				}
				if system == "grub" {
					partition = fields[1]
				} else {
					partition = OtherPartition(fields[1])
				}
			}
			return
		}
	}
	return
}

func bootConf() (string, error) {
	bootSystem, err := BootSystem()
	if err != nil {
		return "", err
	}
	if bootSystem == "grub" {
		return grubConfigFile, nil
	}
	return ubootConfigFile, nil
}

// OtherPartition returns the backup partition, a or b.
func OtherPartition(current string) string {
	if current == "a" {
		return "b"
	}
	return "a"
}