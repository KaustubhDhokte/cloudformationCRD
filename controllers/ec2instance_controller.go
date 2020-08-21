/*


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

package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	cloudfv1 "cf/api/v1"

	"github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/ec2"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	//"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	cf "github.com/awslabs/goformation/v4/cloudformation"
	ectwo "github.com/awslabs/goformation/v4/cloudformation/ec2"
)

// EC2InstanceReconciler reconciles a EC2Instance object
type EC2InstanceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=cloudf.cloudfdomain,resources=ec2instances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cloudf.cloudfdomain,resources=ec2instances/status,verbs=get;update;patch

func (r *EC2InstanceReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("ec2instance", req.NamespacedName)

	// your logic here

	// sess, err := session.NewSession(&aws.Config{
	// 	Region: aws.String("us-east-2")},
	// )

	//svc := ec2.New(sess)

	var ec2I cloudfv1.EC2Instance

	if err := r.Get(ctx, req.NamespacedName, &ec2I); err != nil {
		//log.Error(err, "unable to fetch EC2Instance")
		fmt.Println("unable to fetch EC2Instance")
	} else {

		stackName := ec2I.Spec.StackName

		template := cf.NewTemplate()

		template.Resources["EC2Instance"] = &ectwo.Instance{
			ImageId:      ec2I.Spec.ImageId,
			InstanceType: ec2I.Spec.InstanceType,
		}

		j, err := template.JSON()
		if err != nil {
			fmt.Printf("Failed to generate JSON: %s\n", err)
		} else {
			fmt.Printf("%s\n", string(j))
		}

		templateBody := string(j)
		fmt.Println(templateBody)

		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
			Config: aws.Config{
				Region: aws.String(ec2I.Spec.Region),
			},
		}))

		svc := cloudformation.New(sess)

		_, errr := svc.CreateStack(&cloudformation.CreateStackInput{
			TemplateBody: &templateBody,
			StackName:    &stackName,
		})

		if errr != nil {
			fmt.Println("Got an error creating stack " + stackName)
			fmt.Println(errr)
		}

		err = svc.WaitUntilStackCreateComplete(&cloudformation.DescribeStacksInput{
			StackName: &stackName,
		})
		if err != nil {
			fmt.Println("Got an error waiting for stack to be created")
			fmt.Println(err)
		} else {
			fmt.Println("Created stack " + stackName)
		}
	}

	// runR, err := svc.RunInstances(&ec2.RunInstancesInput{

	// 	ImageId:      aws.String(ec2I.Spec.ImageId),
	// 	InstanceType: aws.String(ec2I.Spec.InstanceType),
	// 	MinCount:     aws.Int64(1),
	// 	MaxCount:     aws.Int64(1),
	// })

	// if err != nil {
	// 	fmt.Println("\nCould not create instance")
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println("Created instance")
	// 	fmt.Println(runR)
	// }

	return ctrl.Result{}, nil
}

func (r *EC2InstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cloudfv1.EC2Instance{}).
		Complete(r)
}
