/*
Copyright 2018 The MetaGraph Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"metagraf/pkg/metagraf"
	"os"
	"strings"
)

// Returns a slice of strings of potential parameterized variables in a
// metaGraf specification that can be found in the execution environment.
func VarsFromEnv(mgv MGVars) EnvVars {
	envs := EnvVars{}

	for _,v := range os.Environ() {
		if
	}
	return envs
}

func keyFromEnv(s string) string {
	return strings.Split(s,"=")[0]
}

// Returns a slice of strings of alle parameterized fields in a metaGraf
// specification.
func VarsFromMetaGraf(mg *metagraf.MetaGraf) MGVars {
	vars := MGVars{}

	for _,env := range mg.Spec.Environment.Local {
		vars[env.Name] = ""
	}
	for _,env := range mg.Spec.Environment.External.Introduces {
		vars[env.Name] = ""
	}
	for _,env := range mg.Spec.Environment.External.Consumes {
		vars[env.Name] = ""
	}

	return vars
}