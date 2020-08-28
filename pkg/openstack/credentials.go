// Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package openstack

import (
	"context"
	"fmt"

	extensionscontroller "github.com/gardener/gardener/extensions/pkg/controller"
	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Credentials contains the necessary OpenStack credential information.
type Credentials struct {
	DomainName string
	TenantName string
	TenantID   string
	Username   string
	Password   string
	AuthURL    string
}

// GetCredentialsBySecretBinding computes for a given context and infrastructure the corresponding credentials object.
func GetCredentialsBySecretBinding(ctx context.Context, c client.Client, secretBindingKey client.ObjectKey) (*Credentials, error) {
	binding := &gardencorev1beta1.SecretBinding{}
	if err := c.Get(ctx, secretBindingKey, binding); err != nil {
		return nil, err
	}

	return GetCredentials(ctx, c, binding.SecretRef)
}

// GetCredentials computes for a given context and infrastructure the corresponding credentials object.
func GetCredentials(ctx context.Context, c client.Client, secretRef corev1.SecretReference) (*Credentials, error) {
	secret, err := extensionscontroller.GetSecretByReference(ctx, c, &secretRef)
	if err != nil {
		return nil, err
	}
	return ExtractCredentials(secret)
}

// ExtractCredentials generates a credentials object for a given provider secret.
func ExtractCredentials(secret *corev1.Secret) (*Credentials, error) {
	if secret.Data == nil {
		return nil, fmt.Errorf("secret does not contain any data")
	}
	domainName, err := getRequired(secret.Data, DomainName)
	if err != nil {
		return nil, err
	}
	tenantName, err := getRequired(secret.Data, TenantName)
	if err != nil {
		return nil, err
	}
	tenantID, err := getRequired(secret.Data, TenantID)
	if err != nil {
		return nil, err
	}
	userName, err := getRequired(secret.Data, UserName)
	if err != nil {
		return nil, err
	}
	password, err := getRequired(secret.Data, Password)
	if err != nil {
		return nil, err
	}
	authURL := secret.Data[AuthURL]

	return &Credentials{
		DomainName: domainName,
		TenantName: tenantName,
		TenantID:   tenantID,
		Username:   userName,
		Password:   password,
		AuthURL:    string(authURL),
	}, nil
}

// getRequired checks if the provided map has a valid value for a corresponding key.
func getRequired(data map[string][]byte, key string) (string, error) {
	value, ok := data[key]
	if !ok {
		return "", fmt.Errorf("map %v does not contain key %s", data, key)
	}
	if len(value) == 0 {
		return "", fmt.Errorf("key %s may not be empty", key)
	}
	return string(value), nil
}
