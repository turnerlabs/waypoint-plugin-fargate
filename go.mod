module github.com/turnerlabs/waypoint-plugin-fargate

go 1.14

require (
	github.com/aws/aws-sdk-go v1.35.15
	github.com/golang/protobuf v1.4.3
	github.com/hashicorp/waypoint v0.1.4
	github.com/hashicorp/waypoint-plugin-sdk v0.0.0-20201016002013-59421183d54f
	golang.org/x/sys v0.0.0-20201018230417-eeed37f84f13 // indirect
	google.golang.org/protobuf v1.25.0
)

replace golang.org/x/sys => golang.org/x/sys v0.0.0-20200826173525-f9321e4c35a6
