/*
Copyright 2024 The Kubernetes Authors.

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

package client

import (
	"fmt"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

// DynamicClientImpersonator is an interface that provides functionality to impersonate
// a Kubernetes user and obtain a dynamic client with the impersonated credentials.
// This is useful for scenarios where you need to perform operations on behalf of
// different users while maintaining a single client instance.
type DynamicClientImpersonator interface {
	Impersonate(impersonateConfig rest.ImpersonationConfig) (dynamic.Interface, error)
}

func (g *clientset) Impersonate(impersonateConfig rest.ImpersonationConfig) (dynamic.Interface, error) {
	if cli, ok := g.impersonationCache[impersonateConfig.UserName]; ok {
		return cli, nil
	}

	restConfig, err := g.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	newRestConfig := rest.CopyConfig(restConfig)
	newRestConfig.Impersonate = impersonateConfig
	cli, err := dynamic.NewForConfig(newRestConfig)
	if err != nil {
		return nil, fmt.Errorf("could not get Kubernetes dynamicClient: %w", err)
	}

	g.impersonationCache[impersonateConfig.UserName] = cli
	return cli, nil
}
