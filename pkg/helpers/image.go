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
	"encoding/json"
	log "k8s.io/klog"
	dockerv10 "github.com/openshift/api/image/docker10"
	imagev1 "github.com/openshift/api/image/v1"
	imagev1client "github.com/openshift/client-go/image/clientset/versioned/typed/image/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
)

func GetImageStreamTags(c *imagev1client.ImageV1Client, ns string, n string) *imagev1.ImageStreamTag {
	ist, err := c.ImageStreamTags(ns).Get(n, metav1.GetOptions{})
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	if len(ist.Image.DockerImageMetadata.Raw) != 0 {
		di := &dockerv10.DockerImage{}
		err = json.Unmarshal(ist.Image.DockerImageMetadata.Raw, di)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		ist.Image.DockerImageMetadata.Object = di
	}


	return ist
}

// Unmarshals json data into a DockerImage type
func GetDockerImageFromIST(i *imagev1.ImageStreamTag) *dockerv10.DockerImage {
	di := dockerv10.DockerImage{}
	if len(i.Image.DockerImageMetadata.Raw) != 0 {
		err := json.Unmarshal(i.Image.DockerImageMetadata.Raw, &di)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		i.Image.DockerImageMetadata.Object = &di
	}
	return &di
}

// Unmarshals json data into a DockerImage type
func GetDockerImageFromImage(i *imagev1.Image) *dockerv10.DockerImage {
	di := dockerv10.DockerImage{}
	if len(i.DockerImageMetadata.Raw) != 0 {
		err := json.Unmarshal(i.DockerImageMetadata.Raw, &di)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		i.DockerImageMetadata.Object = &di
	}
	return &di
}

func GetImage(c *imagev1client.ImageV1Client, i string) *imagev1.Image {
	img, err := c.Images().Get(i, metav1.GetOptions{})
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	return img
}

