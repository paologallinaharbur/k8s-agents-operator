/*
Copyright 2024.

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

package instrumentation

import (
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	volumeName        = "newrelic-instrumentation"
	initContainerName = "newrelic-instrumentation"
	sideCarName       = "opentelemetry-auto-instrumentation"

	// indicates whether newrelic agents should be injected or not.
	// Possible values are "true", "false" or "<Instrumentation>" name.
	annotationInjectJava                 = "instrumentation.newrelic.com/inject-java"
	annotationInjectJavaContainersName   = "instrumentation.newrelic.com/java-container-names"
	annotationInjectNodeJS               = "instrumentation.newrelic.com/inject-nodejs"
	annotationInjectNodeJSContainersName = "instrumentation.newrelic.com/nodejs-container-names"
	annotationInjectPython               = "instrumentation.newrelic.com/inject-python"
	annotationInjectPythonContainersName = "instrumentation.newrelic.com/python-container-names"
	annotationInjectDotNet               = "instrumentation.newrelic.com/inject-dotnet"
	annotationInjectDotnetContainersName = "instrumentation.newrelic.com/dotnet-container-names"
	annotationInjectPhp                  = "instrumentation.newrelic.com/inject-php"
	annotationInjectPhpContainersName    = "instrumentation.newrelic.com/php-container-names"
	annotationInjectRuby                 = "instrumentation.newrelic.com/inject-ruby"
	annotationInjectRubyContainersName   = "instrumentation.newrelic.com/ruby-container-names"
	annotationInjectContainerName        = "instrumentation.newrelic.com/container-name"
	annotationInjectGo                   = "instrumentation.opentelemetry.io/inject-go"
	annotationGoExecPath                 = "instrumentation.opentelemetry.io/otel-go-auto-target-exe"
	annotationInjectGoContainerName      = "instrumentation.opentelemetry.io/go-container-name"
	annotationPhpVersion                 = "instrumentation.newrelic.com/php-version"
)

// annotationValue returns the effective annotation value, based on the annotations from the pod and namespace.
func annotationValue(ns metav1.ObjectMeta, pod metav1.ObjectMeta, annotation string) string {
	// is the pod annotated with instructions to inject sidecars? is the namespace annotated?
	// if any of those is true, a sidecar might be desired.
	podAnnValue := pod.Annotations[annotation]
	nsAnnValue := ns.Annotations[annotation]

	// if the namespace value is empty, the pod annotation should be used, whatever it is
	if len(nsAnnValue) == 0 {
		return podAnnValue
	}

	// if the pod value is empty, the annotation should be used (true, false, instance)
	if len(podAnnValue) == 0 {
		return nsAnnValue
	}

	// the pod annotation isn't empty -- if it's an instance name, or false, that's the decision
	if !strings.EqualFold(podAnnValue, "true") {
		return podAnnValue
	}

	// pod annotation is 'true', and if the namespace annotation is false, we just return 'true'
	if strings.EqualFold(nsAnnValue, "false") {
		return podAnnValue
	}

	// by now, the pod annotation is 'true', and the namespace annotation is either true or an instance name
	// so, the namespace annotation can be used
	return nsAnnValue
}
