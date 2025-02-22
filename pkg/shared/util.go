// Copyright 2019 Hewlett Packard Enterprise Development LP

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shared

import (
	"fmt"
	"os"

	kdv1 "github.com/bluek8s/kubedirector/pkg/apis/kubedirector/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// StrPtr convert a string to a pointer
func StrPtr(s string) *string {
	return &s
}

// StringInList is a utility function that checks if a given string is
// present at least once in the given slice of strings.
func StringInList(
	test string,
	list []string,
) bool {

	for _, s := range list {
		if s == test {
			return true
		}
	}
	return false
}

// ListIsUnique is a utility function that checks if a given slice of strings
// is free of duplicates.
func ListIsUnique(
	list []string,
) bool {

	seen := make(map[string]bool)
	for _, s := range list {
		if _, ok := seen[s]; ok {
			return false
		}
		seen[s] = true
	}
	return true
}

// GetKubeDirectorNamespace is a utility function to fetch the namespace
// where kubedirector is running
func GetKubeDirectorNamespace() (string, error) {

	ns, found := os.LookupEnv(KubeDirectorNamespaceEnvVar)
	if !found {
		return "", fmt.Errorf("%s must be set", KubeDirectorNamespaceEnvVar)
	}
	return ns, nil
}

// OwnerReferences creates an owner reference spec that identifies the
// custom resource as the owner.
func OwnerReferences(
	cr KubeDirectorObject,
) []metav1.OwnerReference {

	// IF THIS IS EVER CHANGED TO RETURN MORE THAN ONE REFERENCE for some
	// reason, then ownerReferencesPresent below will also need to be
	// changed.
	return []metav1.OwnerReference{
		*metav1.NewControllerRef(cr, schema.GroupVersionKind{
			Group:   kdv1.SchemeGroupVersion.Group,
			Version: kdv1.SchemeGroupVersion.Version,
			Kind:    cr.GetObjectKind().GroupVersionKind().Kind,
		}),
	}
}

// OwnerReferencesPresent determines whether the desired references (from
// the ownerReferences func) are present in the CR.
func OwnerReferencesPresent(
	cr KubeDirectorObject,
	currentRefs []metav1.OwnerReference,
) bool {

	// As mentioned above, for simplicity we leverage the fact that
	// we only require one owner reference. Also we probably don't need/want
	// to do an entire struct compare; only the fields we really care about.
	desiredRef := &(OwnerReferences(cr)[0])
	for _, ref := range currentRefs {
		if (ref.APIVersion == desiredRef.APIVersion) &&
			(ref.Kind == desiredRef.Kind) &&
			(ref.Name == desiredRef.Name) &&
			(ref.UID == desiredRef.UID) &&
			(ref.Controller != nil) &&
			(*ref.Controller == true) {
			return true
		}
	}
	return false
}
