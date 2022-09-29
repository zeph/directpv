// This file is part of MinIO DirectPV
// Copyright (c) 2021, 2022 MinIO, Inc.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package installer

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/minio/directpv/pkg/consts"
	"github.com/minio/directpv/pkg/credential"
	"github.com/minio/directpv/pkg/k8s"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func installSecretsDefault(ctx context.Context, c *Config) error {
	if err := installCredSecret(ctx, c); err != nil {
		return err
	}
	// Add more secrets here..
	return nil
}

func uninstallSecretsDefault(ctx context.Context, c *Config) error {
	if err := uninstallCredSecret(ctx, c); err != nil {
		return err
	}
	// Delete more secrets here
	return nil
}

func installCredSecret(ctx context.Context, c *Config) error {
	cred, err := credential.Load(c.ConfigFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = fmt.Errorf("credentials not found: please configure access and secret keys via ENV or config file (%s)", c.ConfigFile)
		}
		return err
	}
	return createOrUpdateSecret(ctx, consts.CredentialsSecretName, map[string][]byte{
		consts.AccessKeyDataKey: []byte(cred.AccessKey),
		consts.SecretKeyDataKey: []byte(cred.SecretKey),
	}, c)
}

func uninstallCredSecret(ctx context.Context, c *Config) error {
	if err := k8s.KubeClient().CoreV1().Secrets(c.namespace()).Delete(ctx, consts.CredentialsSecretName, metav1.DeleteOptions{}); err != nil && !apierrors.IsNotFound(err) {
		return err
	}
	return nil
}