# waypoint-plugin-fargate

A [waypoint](https://www.waypointproject.io/) plugin that deploys [AWS ECS/Fargate](https://aws.amazon.com/fargate/) applications.

This plugin is similar to the built-in [aws-ecs](https://www.waypointproject.io/plugins/aws-ecs) plugin, however it does not create or provision any cloud infrastructure resources.  The plugin assumes that you already have an existing ECS service and simply deploys your application container on top of the infrastructure.

The plugin is optimized to work well with the [fargate-create](https://github.com/turnerlabs/fargate-create/) tool which uses [Terraform](https://github.com/turnerlabs/terraform-ecs-fargate) to provision the AWS cloud infrastructure.

The great thing about `waypoint` is that it enables simple, declarative, and portable build/push/deploy flows that run the same on your laptop as they do in any of your CI/CD pipelines.


### usage example

The following example will build your application container (assuming you have a `Dockerfile`), push it to AWS ECR, register a new task definition which references the newly build image, and finally update the service to run the new task definition.

`waypoint.hcl`

```hcl
project = "waypoint-test"

app "waypoint-test" {

  build {
    use "docker" {}

    registry {
      use "aws-ecr" {
        region     = "us-east-1"
        repository = "waypoint-test"
        tag        = gitrefpretty()
      }
    }
  }

  deploy {
    use "fargate" {
      cluster = "waypoint-test-dev"
      service = "waypoint-test-dev"
    }
  }
}
```

```
waypoint up
```


### build from source

```
make build
make install
```