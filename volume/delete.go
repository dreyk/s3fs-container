/*
Copyright 2016 The Kubernetes Authors.

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

package volume

import (
	"github.com/golang/glog"
	"fmt"
	"strconv"
	"os/exec"
	"github.com/childsb/s3fs-container/controller"
	"k8s.io/client-go/pkg/api/v1"
)

// Delete removes the directory that was created by Provision backing the given
// PV and removes its export from the NFS server.
func (p *s3fsProvisioner) Delete(volume *v1.PersistentVolume) error {
	// Ignore the call if this provisioner was not the one to provision the
	// volume. It doesn't even attempt to delete it, so it's neither a success
	// (nil error) nor failure (any other error)


	glog.Infof("Delete called for volume:" , volume.Name)

	provisioned, err := p.provisioned(volume)
	if err != nil {
		return fmt.Errorf("error determining if this provisioner was the one to provision volume %q: %v", volume.Name, err)
	}
	if !provisioned {
		strerr := fmt.Sprintf("this provisioner id %s didn't provision volume %q and so can't delete it; id %s did & can", p.identity, volume.Name, volume.Annotations[annProvisionerId])
		return &controller.IgnoredError{strerr}
	}

	cmd := exec.Command(p.execCommand, "delete", volume.Spec.FlexVolume.Options[annAwss3bucket], volume.Spec.FlexVolume.Options[annAwsAccessKeyId], volume.Spec.FlexVolume.Options[annAwsSecretAccessKey] )
	output, err := cmd.CombinedOutput()
	if err != nil {
		glog.Errorf("Failed to delete volume %s, output: %s, error: %s",  volume.Spec.FlexVolume.Options[annAwss3bucket], output, err.Error())
		//_, err := handleCmdResponse(mountCmd, output)
		return err
	}
	return nil
}

func (p *s3fsProvisioner) provisioned(volume *v1.PersistentVolume) (bool, error) {
	provisionerId, ok := volume.Annotations[annProvisionerId]
	if !ok {
		return false, fmt.Errorf("PV doesn't have an annotation %s", annProvisionerId)
	}

	return provisionerId == string(p.identity), nil
}

func getBlockAndId(volume *v1.PersistentVolume, annBlock, annId string) (string, uint16, error) {
	block, ok := volume.Annotations[annBlock]
	if !ok {
		return "", 0, fmt.Errorf("PV doesn't have an annotation with key %s", annBlock)
	}

	idStr, ok := volume.Annotations[annId]
	if !ok {
		return "", 0, fmt.Errorf("PV doesn't have an annotation %s", annId)
	}
	id, _ := strconv.ParseUint(idStr, 10, 16)

	return block, uint16(id), nil
}
