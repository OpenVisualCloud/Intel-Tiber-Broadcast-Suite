/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

package bcs

import "bcs.pod.launcher.intel/resources_library/resources/general"

type BcsApp struct {
	Name       string
	Namespace  string
	Containers general.Containers
}
