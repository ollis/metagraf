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

package helpers

import (
		"github.com/fsouza/go-dockerclient"
)

// todo: verify image string
func DockerInspectImage(image string, tag string, auth docker.AuthConfiguration) *docker.Image {

	endpoint := "unix://var/run/docker.sock"

	cli, err := docker.NewClient(endpoint)
	if err != nil {
		panic(err)
	}

	pull := docker.PullImageOptions{
		Tag: tag,
		Repository: image,
	}

	err = cli.PullImage( pull, auth )
	if err != nil {
		panic(err)
	}

	i, err := cli.InspectImage(image)
	if err != nil {
		panic(err)
	}
	return i
}