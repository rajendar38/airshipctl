/*
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

package config

import (
	"os"

	"opendev.org/airship/airshipctl/pkg/errors"
)

// ContextOptions holds all configurable options for context
type ContextOptions struct {
	Name                    string
	CurrentContext          bool
	Manifest                string
	Current                 bool
	ManagementConfiguration string
	EncryptionConfig        string
}

// ManifestOptions holds all configurable options for manifest configuration
type ManifestOptions struct {
	Name         string
	RepoName     string
	URL          string
	Branch       string
	CommitHash   string
	Tag          string
	RemoteRef    string
	Force        bool
	IsPhase      bool
	SubPath      string
	TargetPath   string
	MetadataPath string
}

// EncryptionConfigOptions holds all configurable options for encryption configuration
type EncryptionConfigOptions struct {
	Name               string
	EncryptionKeyPath  string
	DecryptionKeyPath  string
	KeySecretName      string
	KeySecretNamespace string
}

// TODO(howell): The following functions are tightly coupled with flags passed
// on the command line. We should find a way to remove this coupling, since it
// is possible to create (and validate) these objects without using the command
// line.

// Validate checks for the possible context option values and returns
// Error when invalid value or incompatible choice of values given
func (o *ContextOptions) Validate() error {
	if !o.Current && o.Name == "" {
		return ErrEmptyContextName{}
	}

	if o.Current && o.Name != "" {
		return ErrConflictingContextOptions{}
	}

	// If the user simply wants to change the current context, no further validation is needed
	if o.CurrentContext {
		return nil
	}

	// TODO Manifest, Cluster could be validated against the existing config maps
	return nil
}

// Validate checks for the possible manifest option values and returns
// Error when invalid value or incompatible choice of values given
func (o *ManifestOptions) Validate() error {
	if o.Name == "" {
		return ErrMissingManifestName{}
	}
	if o.RemoteRef != "" {
		return errors.ErrNotImplemented{What: "repository checkout by RemoteRef"}
	}
	if o.IsPhase && o.RepoName == "" {
		return ErrMissingRepositoryName{}
	}
	possibleValues := [3]string{o.CommitHash, o.Branch, o.Tag}
	var count int
	for _, val := range possibleValues {
		if val != "" {
			count++
		}
	}
	if count > 1 {
		return ErrMutuallyExclusiveCheckout{}
	}
	return nil
}

// Validate checks for the possible errors with encryption config
// Error when invalid value, incompatible choice of values given or if the
// key file paths do not exist in the file system
func (o *EncryptionConfigOptions) Validate() error {
	switch {
	case o.Name == "":
		return ErrMissingEncryptionConfigName{}

	case o.backedByFileSystem() == o.backedByAPIServer():
		return ErrMutuallyExclusiveEncryptionConfigType{}

	case o.backedByFileSystem():
		if o.EncryptionKeyPath == "" || o.DecryptionKeyPath == "" {
			return ErrInvalidEncryptionKeyPath{}
		}

	case o.backedByAPIServer():
		if o.KeySecretName == "" || o.KeySecretNamespace == "" {
			return ErrInvalidEncryptionKey{}
		}
	}

	if o.backedByFileSystem() {
		if err := checkExists("encryption-key-path", o.EncryptionKeyPath); err != nil {
			return err
		}
		if err := checkExists("decryption-key-path", o.EncryptionKeyPath); err != nil {
			return err
		}
	}

	return nil
}

func (o EncryptionConfigOptions) backedByFileSystem() bool {
	return o.EncryptionKeyPath != "" || o.DecryptionKeyPath != ""
}

func (o EncryptionConfigOptions) backedByAPIServer() bool {
	return o.KeySecretName != "" || o.KeySecretNamespace != ""
}

func checkExists(flagName, path string) error {
	if path == "" {
		return ErrMissingFlag{FlagName: flagName}
	}
	if _, err := os.Stat(path); err != nil {
		return ErrCheckFile{
			FlagName:    flagName,
			Path:        path,
			InternalErr: err,
		}
	}
	return nil
}
