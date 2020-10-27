package platform

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/hashicorp/waypoint/builtin/docker"
)

const defaultRegion = "us-east-1"

type fargateDeployment struct {
	region    string
	cluster   string
	service   string
	container string
	image     *docker.Image
	client    *ecs.ECS
}

//initializes a new instance of a fargateDeployment
func newFargateDeployment(region, cluster, service, container string, image *docker.Image) (*fargateDeployment, error) {

	r := region
	if r == "" {
		r = defaultRegion
	}

	result := fargateDeployment{
		region:    r,
		cluster:   cluster,
		service:   service,
		container: container,
		image:     image,
	}

	//create a new aws session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(result.region)},
	)
	if err != nil {
		return nil, err
	}

	//create client
	result.client = ecs.New(sess)

	return &result, nil
}

//deploy deploys a fargateDeployment to fargate
func (f *fargateDeployment) deploy() (*int64, error) {

	//get the running task definition attached to the service
	td, err := f.getRunningTaskDefinition()
	if err != nil {
		return nil, fmt.Errorf("fetching service's existing task definition: %w", err)
	}
	dtd, err := f.describeTaskDefinition(td)
	if err != nil {
		return nil, fmt.Errorf("describing task definition: %w", err)
	}

	//look for specified container (or default) in the task definition
	if len(dtd.TaskDefinition.ContainerDefinitions) == 0 {
		return nil, fmt.Errorf("no containers found in the task definition")
	}
	targetContainer := dtd.TaskDefinition.ContainerDefinitions[0]
	if f.container != "" {
		for _, c := range dtd.TaskDefinition.ContainerDefinitions {
			if *c.Name == f.container {
				targetContainer = c
			}
		}
	}

	//register a new task definition with the updated image
	targetContainer.Image = aws.String(f.image.Name())
	newTD, err := f.registerTaskDefinition(dtd)
	if err != nil {
		return nil, fmt.Errorf("registering task definition: %w", err)
	}

	//update the service to use the newly registered task definition
	err = f.updateServiceTaskDefinition(*newTD.TaskDefinitionArn)
	if err != nil {
		return nil, fmt.Errorf("updating service to use new task definition: %w", err)
	}

	return newTD.Revision, nil
}

func (f *fargateDeployment) getRunningTaskDefinition() (string, error) {

	result, err := f.client.DescribeServices(&ecs.DescribeServicesInput{
		Cluster:  aws.String(f.cluster),
		Services: aws.StringSlice([]string{f.service}),
	})
	if err != nil {
		return "", err
	}
	if len(result.Services) == 0 {
		return "", fmt.Errorf("no services found for cluster :%s, service: %s", f.cluster, f.service)
	}

	return *result.Services[0].TaskDefinition, nil
}

func (f *fargateDeployment) describeTaskDefinition(arn string) (*ecs.DescribeTaskDefinitionOutput, error) {

	result, err := f.client.DescribeTaskDefinition(
		&ecs.DescribeTaskDefinitionInput{
			TaskDefinition: aws.String(arn),
			Include:        aws.StringSlice([]string{"TAGS"}),
		},
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (f *fargateDeployment) registerTaskDefinition(dtd *ecs.DescribeTaskDefinitionOutput) (*ecs.TaskDefinition, error) {

	input := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions:    dtd.TaskDefinition.ContainerDefinitions,
		Cpu:                     dtd.TaskDefinition.Cpu,
		ExecutionRoleArn:        dtd.TaskDefinition.ExecutionRoleArn,
		Family:                  dtd.TaskDefinition.Family,
		Memory:                  dtd.TaskDefinition.Memory,
		NetworkMode:             dtd.TaskDefinition.NetworkMode,
		RequiresCompatibilities: dtd.TaskDefinition.RequiresCompatibilities,
		TaskRoleArn:             dtd.TaskDefinition.TaskRoleArn,
		Volumes:                 dtd.TaskDefinition.Volumes,
	}

	//it's unfortunate that the tags aren't included in the task definition itself :(
	if len(dtd.Tags) > 0 {
		input.Tags = dtd.Tags
	}

	//register a new task definition
	result, err := f.client.RegisterTaskDefinition(input)
	if err != nil {
		return nil, err
	}

	return result.TaskDefinition, nil
}

func (f *fargateDeployment) updateServiceTaskDefinition(arn string) error {

	_, err := f.client.UpdateService(
		&ecs.UpdateServiceInput{
			Cluster:        aws.String(f.cluster),
			Service:        aws.String(f.service),
			TaskDefinition: aws.String(arn),
		},
	)
	if err != nil {
		return err
	}

	return nil
}
