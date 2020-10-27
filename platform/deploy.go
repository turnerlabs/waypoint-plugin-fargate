package platform

import (
	"context"
	"fmt"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/docker"
)

// Config represents the fargate plugin config
type Config struct {
	Cluster   string `hcl:"cluster"`
	Service   string `hcl:"service"`
	Region    string `hcl:"region,optional"`
	Container string `hcl:"container,optional"`
}

// Platform is the platform component
type Platform struct {
	config Config
}

// Config implements Configurable
func (p *Platform) Config() (interface{}, error) {
	return &p.config, nil
}

// DeployFunc implements component.Platform
func (p *Platform) DeployFunc() interface{} {
	return p.deploy
}

func (p *Platform) deploy(ctx context.Context, ui terminal.UI, src *component.Source, dockerImage *docker.Image, deployConfig *component.DeploymentConfig) (*Deployment, error) {

	deployment, err := newFargateDeployment(
		p.config.Region,
		p.config.Cluster,
		p.config.Service,
		p.config.Container,
		dockerImage,
	)
	if err != nil {
		return nil, err
	}

	u := ui.Status()
	defer u.Close()
	sg := ui.StepGroup()
	step := sg.Add(fmt.Sprintf("Deploying %v to %v", deployment.image.Name(), deployment.service))
	defer step.Abort()

	revision, err := deployment.deploy()
	if err != nil {
		e := fmt.Errorf("deployment failed: %w", err)
		return nil, e
	}

	step.Update(fmt.Sprintf("Deployed %v to %v as revision %v", deployment.image.Name(), deployment.service, *revision))
	step.Done()

	return &Deployment{
		Id: string(*revision),
	}, nil
}
