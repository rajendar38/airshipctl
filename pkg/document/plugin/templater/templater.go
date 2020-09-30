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

package templater

import (
	"io"
	"text/template"

	"github.com/Masterminds/sprig"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"

	airshipv1 "opendev.org/airship/airshipctl/pkg/api/v1alpha1"
	plugtypes "opendev.org/airship/airshipctl/pkg/document/plugin/types"
)

type plugin struct {
	*airshipv1.Templater
}

// New creates new instance of the plugin
func New(obj map[string]interface{}) (plugtypes.Plugin, error) {
	cfg := &airshipv1.Templater{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj, cfg)
	if err != nil {
		return nil, err
	}
	return &plugin{
		Templater: cfg,
	}, nil
}

// Run templater plugin
func (t *plugin) Run(_ io.Reader, out io.Writer) error {
	funcMap := sprig.TxtFuncMap()
	funcMap["toYaml"] = toYaml
	tmpl, err := template.New("tmpl").Funcs(funcMap).Parse(t.Template)
	if err != nil {
		return err
	}
	return tmpl.Execute(out, t.Values)
}

// Render input yaml as output yaml
// This function is from the Helm project:
// https://github.com/helm/helm
// Copyright The Helm Authors
func toYaml(v interface{}) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return string(data)
}