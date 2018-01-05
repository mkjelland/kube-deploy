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
	"os"

	"github.com/spf13/pflag"
)

type Configuration struct {
	Kubeconfig      string
	BOSHDirectorURL string
	UaaURL          string
	UaaClient       string
	UaaClientSecret string
	UaaCACert       string
}

func NewConfiguration() *Configuration {
	return &Configuration{}
}

func (c *Configuration) AddFlags(fs *pflag.FlagSet) {
	c.Kubeconfig = os.Getenv("Kubeconfig")
	c.BOSHDirectorURL = os.Getenv("BOSHDirectorURL")
	c.UaaURL = os.Getenv("UaaURL")
	c.UaaClient = os.Getenv("UaaClient")
	c.UaaClientSecret = os.Getenv("UaaClientSecret")
	c.UaaCACert = os.Getenv("UaaCACert")
}
