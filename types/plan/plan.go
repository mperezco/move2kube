/*
Copyright IBM Corporation 2020

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

package plan

import (
	"os"
	"path/filepath"

	"github.com/konveyor/move2kube/internal/common"
	"github.com/konveyor/move2kube/types"
	log "github.com/sirupsen/logrus"
)

// SourceTypeValue defines the type of source
type SourceTypeValue string

// ContainerBuildTypeValue defines the containerization type
type ContainerBuildTypeValue string

// TranslationTypeValue defines the translation type
type TranslationTypeValue string

// TargetInfoArtifactTypeValue defines the target info type
type TargetInfoArtifactTypeValue string

// BuildArtifactTypeValue defines the build artifact type
type BuildArtifactTypeValue string

// SourceArtifactTypeValue defines the source artifact type
type SourceArtifactTypeValue string

// TargetArtifactTypeValue defines the target artifact type
type TargetArtifactTypeValue string

// PlanKind is kind of plan file
const PlanKind types.Kind = "Plan"

const (
	// Compose2KubeTranslation translation type is used when source is docker compose
	Compose2KubeTranslation TranslationTypeValue = "DockerCompose"
	// CfManifest2KubeTranslation translation type is used when source is cloud foundry manifest
	CfManifest2KubeTranslation TranslationTypeValue = "CloudFoundry"
	// Any2KubeTranslation translation type is used when source is of an unknown platform
	Any2KubeTranslation TranslationTypeValue = "Containerize"
	// Kube2KubeTranslation translation type is used when source is Kubernetes
	Kube2KubeTranslation TranslationTypeValue = "Kubernetes"
	// Dockerfile2KubeTranslation translation type is used when source is Knative
	Dockerfile2KubeTranslation TranslationTypeValue = "Dockerfile"
)

const (
	// ComposeSourceTypeValue defines the source as docker compose
	ComposeSourceTypeValue SourceTypeValue = "DockerCompose"
	// DirectorySourceTypeValue defines the source as a simple directory
	DirectorySourceTypeValue SourceTypeValue = "Directory"
	// CfManifestSourceTypeValue defines the source as cf manifest
	CfManifestSourceTypeValue SourceTypeValue = "CfManifest"
	// KNativeSourceTypeValue defines the source as KNative
	KNativeSourceTypeValue SourceTypeValue = "Knative"
	// K8sSourceTypeValue defines the source as Kubernetes
	K8sSourceTypeValue SourceTypeValue = "Kubernetes"
)

const (
	// DockerFileContainerBuildTypeValue defines the containerization type as docker file
	DockerFileContainerBuildTypeValue ContainerBuildTypeValue = "NewDockerfile"
	// ReuseDockerFileContainerBuildTypeValue defines the containerization type as reuse of dockerfile
	ReuseDockerFileContainerBuildTypeValue ContainerBuildTypeValue = "ReuseDockerfile"
	// ReuseContainerBuildTypeValue defines the containerization type as reuse of an existing container
	ReuseContainerBuildTypeValue ContainerBuildTypeValue = "Reuse"
	// CNBContainerBuildTypeValue defines the containerization type of cloud native buildpack
	CNBContainerBuildTypeValue ContainerBuildTypeValue = "CNB"
	// ManualContainerBuildTypeValue defines that the tool assumes that the image will be created manually
	ManualContainerBuildTypeValue ContainerBuildTypeValue = "Manual"
	// S2IContainerBuildTypeValue defines the containerization type of S2I
	S2IContainerBuildTypeValue ContainerBuildTypeValue = "S2I"
)

const (
	// K8sFileArtifactType defines the source artifact type of K8s
	K8sFileArtifactType SourceArtifactTypeValue = "Kubernetes"
	// KnativeFileArtifactType defines the source artifact type of KNative
	KnativeFileArtifactType SourceArtifactTypeValue = "Knative"
	// ComposeFileArtifactType defines the source artifact type of Docker compose
	ComposeFileArtifactType SourceArtifactTypeValue = "DockerCompose"
	// ImageInfoArtifactType defines the source artifact type of image info
	ImageInfoArtifactType SourceArtifactTypeValue = "ImageInfo"
	// CfManifestArtifactType defines the source artifact type of cf manifest
	CfManifestArtifactType SourceArtifactTypeValue = "CfManifest"
	// CfRunningManifestArtifactType defines the source artifact type of a manifest of a running instance
	CfRunningManifestArtifactType SourceArtifactTypeValue = "CfRunningManifest"
	// SourceDirectoryArtifactType defines the source artifact type of normal source code directory
	SourceDirectoryArtifactType SourceArtifactTypeValue = "SourceCode"
	// DockerfileArtifactType defines the source artifact type of dockerfile
	DockerfileArtifactType SourceArtifactTypeValue = "Dockerfile"
)

const (
	// SourceDirectoryBuildArtifactType defines source data artifact type
	SourceDirectoryBuildArtifactType BuildArtifactTypeValue = "SourceCode"
)

const (
	// K8sClusterArtifactType defines target info
	K8sClusterArtifactType TargetInfoArtifactTypeValue = "KubernetesCluster"
)

// Plan defines the format of plan
type Plan struct {
	types.TypeMeta   `yaml:",inline"`
	types.ObjectMeta `yaml:"metadata,omitempty"`
	Spec             PlanSpec `yaml:"spec,omitempty"`
}

// PlanSpec stores the data about the plan
type PlanSpec struct {
	Inputs  Inputs  `yaml:"inputs"`
	Outputs Outputs `yaml:"outputs"`
}

// Outputs defines the output section of plan
type Outputs struct {
	Kubernetes KubernetesOutput `yaml:"kubernetes"`
}

// KubernetesOutput defines the output format for kubernetes deployable artifacts
type KubernetesOutput struct {
	RegistryURL            string            `yaml:"registryURL,omitempty"`
	RegistryNamespace      string            `yaml:"registryNamespace,omitempty"`
	TargetCluster          TargetClusterType `yaml:"targetCluster,omitempty"`
	IgnoreUnsupportedKinds bool              `yaml:"ignoreUnsupportedKinds,omitempty"`
}

// TargetClusterType contains either the type of the target cluster or path to a file containing the target cluster metadata.
// Specify one or the other, not both.
type TargetClusterType struct {
	Type string `yaml:"type,omitempty"`
	Path string `yaml:"path,omitempty" m2kpath:"normal"`
}

// Merge allows merge of two Kubernetes Outputs
func (output *KubernetesOutput) Merge(newoutput KubernetesOutput) {
	if newoutput != (KubernetesOutput{}) {
		if newoutput.RegistryURL != "" {
			output.RegistryURL = newoutput.RegistryURL
		}
		if newoutput.RegistryNamespace != "" {
			output.RegistryNamespace = newoutput.RegistryNamespace
		}
		output.IgnoreUnsupportedKinds = newoutput.IgnoreUnsupportedKinds
		if newoutput.TargetCluster.Type != "" {
			output.TargetCluster = newoutput.TargetCluster
		}
	}
}

// Inputs defines the input section of plan
type Inputs struct {
	RootDir             string                                   `yaml:"rootDir"`
	K8sFiles            []string                                 `yaml:"kubernetesYamls,omitempty" m2kpath:"normal"`
	Services            map[string][]Service                     `yaml:"services"`                                       // [serviceName][Services]
	TargetInfoArtifacts map[TargetInfoArtifactTypeValue][]string `yaml:"targetInfoArtifacts,omitempty" m2kpath:"normal"` //[targetinfoartifacttype][List of artifacts]
}

// RepoInfo contains information specific to creating the CI/CD pipeline.
type RepoInfo struct {
	GitRepoDir    string `yaml:"gitRepoDir" m2kpath:"normal"`
	GitRepoURL    string `yaml:"gitRepoURL"`
	GitRepoBranch string `yaml:"gitRepoBranch"`
	TargetPath    string `yaml:"targetPath" m2kpath:"normal"`
}

// Service defines a plan service
type Service struct {
	ServiceName                   string                               `yaml:"serviceName"`
	ServiceRelPath                string                               `yaml:"serviceRelPath,omitempty"`
	Image                         string                               `yaml:"image"`
	TranslationType               TranslationTypeValue                 `yaml:"translationType"`
	ContainerBuildType            ContainerBuildTypeValue              `yaml:"containerBuildType"`
	SourceTypes                   []SourceTypeValue                    `yaml:"sourceType"`
	ContainerizationTargetOptions []string                             `yaml:"targetOptions,omitempty" m2kpath:"if:ContainerBuildType:in:NewDockerfile,ReuseDockerfile,S2I"`
	SourceArtifacts               map[SourceArtifactTypeValue][]string `yaml:"sourceArtifacts" m2kpath:"keys:Kubernetes,Knative,DockerCompose,CfManifest,CfRunningManifest,SourceCode,Dockerfile"` //[translationartifacttype][List of artifacts]
	BuildArtifacts                map[BuildArtifactTypeValue][]string  `yaml:"buildArtifacts,omitempty" m2kpath:"normal"`                                                                          //[buildartifacttype][List of artifacts]
	UpdateContainerBuildPipeline  bool                                 `yaml:"updateContainerBuildPipeline"`
	UpdateDeployPipeline          bool                                 `yaml:"updateDeployPipeline"`
	RepoInfo                      RepoInfo                             `yaml:"repoInfo,omitempty"`
}

// NewService creates a new service
func NewService(serviceName string, translationtype TranslationTypeValue) Service {
	return Service{
		ServiceName:                   serviceName,
		ServiceRelPath:                "/" + serviceName,
		Image:                         serviceName + ":latest",
		TranslationType:               translationtype,
		ContainerBuildType:            ReuseContainerBuildTypeValue,
		SourceTypes:                   []SourceTypeValue{},
		ContainerizationTargetOptions: []string{},
		SourceArtifacts:               map[SourceArtifactTypeValue][]string{},
		BuildArtifacts:                map[BuildArtifactTypeValue][]string{},
	}
}

// GatherGitInfo tries to find the git repo for the path if one exists.
// It returns true of it found a git repo.
func (service *Service) GatherGitInfo(path string, plan Plan) (bool, error) {
	if finfo, err := os.Stat(path); err != nil {
		log.Errorf("Failed to stat the path %q Error %q", path, err)
		return false, err
	} else if !finfo.IsDir() {
		pathDir := filepath.Dir(path)
		log.Debugf("The path %q is not a directory. Using %q instead.", path, pathDir)
		path = pathDir
	}

	preferredRemote := "upstream"
	remoteNames, err := common.GetGitRemoteNames(path)
	if err != nil || len(remoteNames) == 0 {
		log.Debugf("No remotes found at path %q Error: %q", path, err)
	} else {
		if !common.IsStringPresent(remoteNames, preferredRemote) {
			preferredRemote = "origin"
			if !common.IsStringPresent(remoteNames, preferredRemote) {
				preferredRemote = remoteNames[0]
			}
		}
	}

	remoteURLs, branch, repoDir, err := common.GetGitRepoDetails(path, preferredRemote)
	if err != nil {
		log.Debugf("Failed to get the git repo at path %q Error: %q", path, err)
		return false, err
	}

	service.RepoInfo.GitRepoBranch = branch
	if len(remoteURLs) == 0 {
		log.Debugf("The git repo at path %q has no remotes set.", path)
	} else {
		service.RepoInfo.GitRepoURL = remoteURLs[0]
	}

	service.RepoInfo.GitRepoDir = repoDir
	return true, nil
}

func (service *Service) merge(newservice Service) bool {
	if service.ServiceName != newservice.ServiceName || service.Image != newservice.Image || service.TranslationType != newservice.TranslationType || service.ContainerBuildType != newservice.ContainerBuildType {
		return false
	}
	if len(service.BuildArtifacts[SourceDirectoryBuildArtifactType]) > 0 && len(newservice.BuildArtifacts[SourceDirectoryBuildArtifactType]) > 0 && service.BuildArtifacts[SourceDirectoryBuildArtifactType][0] != newservice.BuildArtifacts[SourceDirectoryBuildArtifactType][0] {
		return false
	}
	service.UpdateContainerBuildPipeline = service.UpdateContainerBuildPipeline || newservice.UpdateContainerBuildPipeline
	service.UpdateDeployPipeline = service.UpdateDeployPipeline || newservice.UpdateDeployPipeline
	service.addSourceTypes(newservice.SourceTypes)
	service.addTargetOptions(newservice.ContainerizationTargetOptions)
	service.addSourceArtifacts(newservice.SourceArtifacts)
	service.addBuildArtifacts(newservice.BuildArtifacts)
	return true
}

// AddSourceArtifact adds a source artifact to a plan service
func (service *Service) AddSourceArtifact(sat SourceArtifactTypeValue, value string) {
	if val, ok := service.SourceArtifacts[sat]; ok {
		service.SourceArtifacts[sat] = append(val, value)
	} else {
		service.SourceArtifacts[sat] = []string{value}
	}
}

func (service *Service) addSourceArtifactArray(sat SourceArtifactTypeValue, values []string) {
	if val, ok := service.SourceArtifacts[sat]; ok {
		service.SourceArtifacts[sat] = common.MergeStringSlices(val, values)
	} else {
		service.SourceArtifacts[sat] = values
	}
}

func (service *Service) addSourceArtifacts(sats map[SourceArtifactTypeValue][]string) {
	for key2, value2 := range sats {
		service.addSourceArtifactArray(key2, value2)
	}
}

// AddBuildArtifact adds a build artifact to a plan service
func (service *Service) AddBuildArtifact(sat BuildArtifactTypeValue, value string) {
	if val, ok := service.BuildArtifacts[sat]; ok {
		service.BuildArtifacts[sat] = append(val, value)
	} else {
		service.BuildArtifacts[sat] = []string{value}
	}
}

func (service *Service) addBuildArtifactArray(sat BuildArtifactTypeValue, values []string) {
	if val, ok := service.BuildArtifacts[sat]; ok {
		service.BuildArtifacts[sat] = common.MergeStringSlices(val, values)
	} else {
		service.BuildArtifacts[sat] = values
	}
}

func (service *Service) addBuildArtifacts(sats map[BuildArtifactTypeValue][]string) {
	for key2, value2 := range sats {
		service.addBuildArtifactArray(key2, value2)
	}
}

// AddSourceType adds source type to a plan service
func (service *Service) AddSourceType(st SourceTypeValue) bool {
	found := false
	for _, est := range service.SourceTypes {
		if est == st {
			found = true
			break
		}
	}
	if !found {
		service.SourceTypes = append(service.SourceTypes, st)
	}
	return true
}

// addSourceTypes adds source types to a plan service
func (service *Service) addSourceTypes(sts []SourceTypeValue) {
	for _, st := range sts {
		service.AddSourceType(st)
	}
}

// addTargetOption adds target option to a plan service
func (service *Service) addTargetOption(st string) bool {
	found := false
	for _, est := range service.ContainerizationTargetOptions {
		if est == st {
			found = true
			break
		}
	}
	if !found {
		service.ContainerizationTargetOptions = append(service.ContainerizationTargetOptions, st)
	}
	return true
}

// addTargetOptions adds target options to a plan service
func (service *Service) addTargetOptions(sts []string) {
	for _, st := range sts {
		service.addTargetOption(st)
	}
}

// AddServicesToPlan adds a list of services to a plan
func (plan *Plan) AddServicesToPlan(services []Service) {
	for _, service := range services {
		if _, ok := plan.Spec.Inputs.Services[service.ServiceName]; !ok {
			plan.Spec.Inputs.Services[service.ServiceName] = []Service{}
			log.Debugf("Added new service to plan : %s", service.ServiceName)
		}
		merged := false
		existingServices := plan.Spec.Inputs.Services[service.ServiceName]
		for i := range existingServices {
			if existingServices[i].merge(service) {
				merged = true
			}
		}
		if !merged {
			plan.Spec.Inputs.Services[service.ServiceName] = append(plan.Spec.Inputs.Services[service.ServiceName], service)
		}
	}
}

// NewPlan creates a new plan
// Sets the version and optionally fills in some default values
func NewPlan() Plan {
	plan := Plan{
		TypeMeta: types.TypeMeta{
			Kind:       string(PlanKind),
			APIVersion: types.SchemeGroupVersion.String(),
		},
		ObjectMeta: types.ObjectMeta{
			Name: common.DefaultProjectName,
		},
		Spec: PlanSpec{
			Inputs: Inputs{
				Services:            map[string][]Service{},
				TargetInfoArtifacts: map[TargetInfoArtifactTypeValue][]string{},
			},
			Outputs: Outputs{
				Kubernetes: KubernetesOutput{
					TargetCluster:          TargetClusterType{Type: common.DefaultClusterType},
					IgnoreUnsupportedKinds: false,
				},
			},
		},
	}
	return plan
}
