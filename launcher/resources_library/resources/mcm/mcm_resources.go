/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

package mcm

import "bcs.pod.launcher.intel/resources_library/resources/general"

type McmApp struct {
	Name       string
	Namespace  string
	Containers general.Containers
}
