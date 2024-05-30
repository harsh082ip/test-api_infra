package main

import (
	// "github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	"fmt"
	"log"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// CREATE SECURITY GROUP
		// sg- args-ingress
		sgArgs := &ec2.SecurityGroupArgs{
			// ingress
			Ingress: ec2.SecurityGroupIngressArray{
				ec2.SecurityGroupIngressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(8002),
					ToPort:     pulumi.Int(8002),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
				ec2.SecurityGroupIngressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(22),
					ToPort:     pulumi.Int(22),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},

			//egress
			Egress: ec2.SecurityGroupEgressArray{
				ec2.SecurityGroupEgressArgs{
					Protocol:   pulumi.String("-1"),
					FromPort:   pulumi.Int(0),
					ToPort:     pulumi.Int(0),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
		}

		// create sg
		sg, err := ec2.NewSecurityGroup(ctx, "jenkins-sg", sgArgs)
		if err != nil {
			log.Println("Security Group Creation Failed, :", err.Error())
			return err
		}

		// key-pair
		kp, err := ec2.NewKeyPair(ctx, "go-infra", &ec2.KeyPairArgs{
			PublicKey: pulumi.String("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCTRXP2AfM1mNR0XSM7KhDTDpAwAbG4Hxlhfx4ahV86h4JI4Kn7nVOvscNGgfqaZr2lIio6y+A30IKASrIQttI9pm/4KIPzh2uGSuedO4HDWaiwWMkXLXeymQm+qYlWbGBvz0CZyoHeHpvbXl0TZAEaW6YSIiAzaicI9vUn+lOVjOst+/gyconivI93XggRHSQXkdr+Lx6LU4hgRTS0/FbD7bsZIjaSOSF63tLX6MD9nck0Amk2sPydOfBYVuOhaW9Lqn2Nb0VkL/q2eyYyycoWyexSvVW7idJ/7sRMFzU0eVrDO2ZkoE3wUa/EsLjen1Qyjm+5lB3knolUcSfllFdUOftuGsf/JU6YtzeZ6ptSiGPngSg9ix2gma3L0kgsPll0owoMgZy+nNgiaN7vrEm1Pf3Uc+3tVE0XZtVAhchYrsLdHgdb2c07dvm736sl9143Dhjt3DUzUluN1bao+UXnBLfvlJ173R0NmbV5VgVMV7u8qJ1WbY4IFOu9OFbJdkE= nightshade@cloud1"),
		})
		if err != nil {
			log.Println("Key-pair creation failed: ", err.Error())
			return err
		}

		// create jenkins-server
		jenkinsServer, err := ec2.NewInstance(ctx, "jenkins-server", &ec2.InstanceArgs{
			InstanceType:        pulumi.String("t2.micro"),
			VpcSecurityGroupIds: pulumi.StringArray{sg.ID()},
			Ami:                 pulumi.String("ami-0f58b397bc5c1f2e8"),
			KeyName:             kp.KeyName,
		})
		if err != nil {
			log.Println("Error Creating the Server:", err.Error())
			return err
		}

		// log details
		fmt.Println("Public IP:", jenkinsServer.PublicIp)
		fmt.Println("Public DNS:", jenkinsServer.PublicDns)

		ctx.Export("publicIP", jenkinsServer.PublicIp)
		ctx.Export("publicDNS", jenkinsServer.PublicDns)
		return nil
	})
}
