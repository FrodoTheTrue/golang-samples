// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package snippets

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

func TestComputeSnippets(t *testing.T) {
	tc := testutil.SystemTest(t)
	zone := "europe-central2-b"
	instanceName := "test-instance-name"
	instanceName2 := "test-instance-name-2"
	machineType := "n1-standard-1"
	sourceImage := "projects/debian-cloud/global/images/family/debian-10"
	networkName := "global/networks/default"

	buf := &bytes.Buffer{}

	if err := createInstance(buf, tc.ProjectID, zone, instanceName, machineType, sourceImage, networkName); err != nil {
		t.Errorf("createInstance got err: %v", err)
	}

	expectedResult := "Instance created"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("createInstance got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, expectedResult)
	}

	buf.Reset()

	if err := listInstances(buf, tc.ProjectID, zone); err != nil {
		t.Errorf("listInstances got err: %v", err)
	}

	expectedResult = "Instances found in zone"
	expectedResult2 := fmt.Sprintf("- %s", instanceName)
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("listInstances got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, expectedResult)
	}
	if got := buf.String(); !strings.Contains(got, expectedResult2) {
		t.Errorf("listInstances got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, expectedResult2)
	}

	buf.Reset()

	if err := listAllInstances(buf, tc.ProjectID); err != nil {
		t.Errorf("listAllInstances got err: %v", err)
	}

	expectedResult = "Instances found:"
	expectedResult2 = fmt.Sprintf("zones/%s\n", zone)
	expectedResult3 := fmt.Sprintf("- %s", instanceName)
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("listAllInstances got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, expectedResult)
	}
	if got := buf.String(); !strings.Contains(got, expectedResult2) {
		t.Errorf("listAllInstances got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, expectedResult2)
	}
	if got := buf.String(); !strings.Contains(got, expectedResult2) {
		t.Errorf("listAllInstances got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, expectedResult3)
	}

	buf.Reset()

	if err := deleteInstance(buf, tc.ProjectID, zone, instanceName); err != nil {
		t.Errorf("deleteInstance got err: %v", err)
	}

	expectedResult = "Instance deleted"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("deleteInstance got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, expectedResult)
	}

	buf.Reset()

	if err := createInstance(buf, tc.ProjectID, zone, instanceName2, machineType, sourceImage, networkName); err != nil {
		t.Errorf("createInstance got err: %v", err)
	}

	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		t.Errorf("NewInstancesRESTClient: %v", err)
	}
	defer instancesClient.Close()

	req := &computepb.DeleteInstanceRequest{
		Project:  tc.ProjectID,
		Zone:     zone,
		Instance: instanceName2,
	}

	op, err := instancesClient.Delete(ctx, req)
	if err != nil {
		t.Errorf("Delete instance request: %v", err)
	}

	waitForOperation(buf, op, tc.ProjectID)

	expectedResult = "Operation finished"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("waitForOperation got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, expectedResult)
	}
}
