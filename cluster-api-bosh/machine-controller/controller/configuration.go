/*
Copyright 2017 The Kubernetes Authors.

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

package controller

import (
	"github.com/spf13/pflag"
)

type Configuration struct {
	Kubeconfig      string
	Cloud           string
	KubeadmToken    string
	UaaURL          string
	DirectorURL     string
	UaaClient       string
	UaaClientSecret string
}

func NewConfiguration() *Configuration {
	return &Configuration{}
}

func (c *Configuration) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.Kubeconfig, "kubeconfig", c.Kubeconfig, "Path to kubeconfig file with authorization and master location information.")
	fs.StringVar(&c.Cloud, "cloud", c.Cloud, "Cloud provider (google/azure).")
	fs.StringVar(&c.KubeadmToken, "token", c.KubeadmToken, "Kubeadm token to use to join new machines.")
	fs.StringVar(&c.UaaURL, "UaaURL", c.UaaURL, "Kubeadm token to use to join new machines.")
	fs.StringVar(&c.DirectorURL, "DirectorURL", c.DirectorURL, "Kubeadm token to use to join new machines.")
	fs.StringVar(&c.UaaClient, "UaaClient", c.UaaClient, "Kubeadm token to use to join new machines.")
	fs.StringVar(&c.UaaClientSecret, "UaaClientSecret", c.UaaClientSecret, "Kubeadm token to use to join new machines.")
}
