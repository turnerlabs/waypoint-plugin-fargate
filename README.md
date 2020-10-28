# waypoint-plugin-fargate

A [waypoint](https://www.waypointproject.io/) plugin that deploys [AWS ECS/Fargate](https://aws.amazon.com/fargate/) applications.

This plugin is similar to the built-in [aws-ecs](https://www.waypointproject.io/plugins/aws-ecs) plugin, however it does not create or provision any cloud infrastructure resources.  The plugin assumes that you already have an existing ECS service and simply deploys your application container on top of the infrastructure.

The plugin is optimized to work well with the [fargate-create](https://github.com/turnerlabs/fargate-create/) tool which uses [Terraform](https://github.com/turnerlabs/terraform-ecs-fargate) to provision the AWS cloud infrastructure.

The great thing about `waypoint` is that it enables simple, declarative, and portable build/push/deploy flows that run the same on your laptop as they do in any of your CI/CD pipelines.


### install

To install the plugin, download the asset for your platform from [releases](https://github.com/turnerlabs/waypoint-plugin-fargate/releases), unzip it and put it in your `${HOME}/.config/waypoint/plugins/` directory (or use one of the [other options](https://www.waypointproject.io/docs/extending-waypoint/creating-plugins/compiling#installing-the-plugin)).


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

If you have more than one container in your app, you can specify which container you'd like to deploy using the `container` field.

```hcl
  deploy {
    use "fargate" {
      region    = "us-east-1"
      cluster   = "my-cluster"
      service   = "my-service"
      container = "my-container"
    }
  }
```


Later, you can list your deployments.

```
waypoint deployment list

     | ID | PLATFORM |   DETAILS   |    STARTED    |   COMPLETED    
-----+----+----------+-------------+---------------+----------------
  🚀 | 74 | fargate  | artifact:32 | 8 minutes ago | 8 minutes ago  
  ✔  | 73 | fargate  | artifact:30 | 21 hours ago  | 21 hours ago
```


### build from source

```
make build
make install
```