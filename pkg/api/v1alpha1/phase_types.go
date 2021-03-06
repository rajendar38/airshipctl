/*
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true

// Phase object to control deployment steps
type Phase struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Config            PhaseConfig `json:"config,omitempty"`
}

// PhaseConfig represents configuration for a particular phase. It contains a reference to
// phase runner object which should contain runner configuration
type PhaseConfig struct {
	ExecutorRef        *corev1.ObjectReference `json:"executorRef"`
	DocumentEntryPoint string                  `json:"documentEntryPoint"`
}

// DefaultPhase can be used to safely unmarshal phase object without nil pointers
func DefaultPhase() *Phase {
	return &Phase{
		Config: PhaseConfig{},
	}
}
