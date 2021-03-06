/*
Copyright 2019 The MetaGraph Authors

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

package modules

import (
	"fmt"
	log "k8s.io/klog"
	"k8s.io/apimachinery/pkg/util/intstr"
	"metagraf/pkg/helpers"
	"os"
	"reflect"
	"sort"
	"strings"

	"metagraf/mg/ocpclient"
	"metagraf/pkg/imageurl"
	"metagraf/pkg/metagraf"

	routev1 "github.com/openshift/api/route/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GenRoute(mg *metagraf.MetaGraf) {
	var weight int32 = 100

	objname := Name(mg)

	var DockerImage string
	if len(mg.Spec.BaseRunImage) > 0 {
		DockerImage = mg.Spec.BaseRunImage
	} else if len(mg.Spec.BuildImage) > 0 {
		DockerImage = mg.Spec.BuildImage
	} else {
		DockerImage = ""
	}

	client := ocpclient.GetImageClient()
	var imgurl imageurl.ImageURL
	err := imgurl.Parse(DockerImage)
	ist := helpers.GetImageStreamTags(
		client,
		imgurl.Namespace,
		imgurl.Image+":"+imgurl.Tag)
	if err != nil {
		log.Errorf("%v", err)
	}

	ImageInfo := helpers.GetDockerImageFromIST(ist)
	log.V(2).Infof("Docker image ports: %v", ImageInfo.Config.ExposedPorts)

	var ports []string

	for _,v := range reflect.ValueOf(ImageInfo.Config.ExposedPorts).MapKeys() {
		ports = append(ports, v.String())
	}
	sort.Strings(ports)

	log.V(2).Infof("First port: %v, %t", ports[0], ports[0])


	l := make(map[string]string)
	l["app"] = objname
	l["deploymentconfig"] = objname

	obj := routev1.Route{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Route",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   objname,
			Labels: l,
		},
		Spec: routev1.RouteSpec{
			To: routev1.RouteTargetReference{
				Kind: "Service",
				Name: objname,
				Weight: &weight,
			},
			Path: Context,
			Port: &routev1.RoutePort{
				TargetPort: intstr.IntOrString{
					Type: 1,
					StrVal: strings.Replace(ports[0], "/", "-", -1),

				},
			},
		},
	}

	if !Dryrun {
		StoreRoute(obj)
	}
	if Output {
		MarshalObject(obj.DeepCopyObject())
	}
}

func StoreRoute(obj routev1.Route) {
	client := ocpclient.GetRouteClient().Routes(NameSpace)
    route, _ := client.Get(obj.Name, metav1.GetOptions{} )
	if len(route.ResourceVersion) > 0 {
		obj.ResourceVersion = route.ResourceVersion
		_, err := client.Update(&obj)
		if err != nil {
			log.Error(err)
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Updated Route: ", obj.Name, " in Namespace: ", NameSpace)
	} else {
		_, err := client.Create(&obj)
		if err != nil {
			log.Error(err)
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Created Route: ", obj.Name, " in Namespace: ", NameSpace)
	}
}

func DeleteRoute(name string) {
	client := ocpclient.GetRouteClient().Routes(NameSpace)

	_, err := client.Get(name, metav1.GetOptions{})
	if err != nil {
		fmt.Println("Route: ", name, "does not exist in namespace: ", NameSpace,", skipping...")
		return
	}

	err = client.Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		fmt.Println( "Unable to delete Route: ", name, " in namespace: ", NameSpace)
		return
	}
	fmt.Println("Deleted Route: ", name, ", in namespace: ", NameSpace)
}