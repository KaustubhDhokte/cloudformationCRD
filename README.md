## cloudformationCRD
AWS objects in kubernetes using respective cloudformation specs

# Prerequisites:
Listed below are the dependencies for the use of this package
1. golang version >= 1.14
2. kubectl server version >= v1.16.8-eks-fd1ea7
3. kubectl client version >= v1.18.6
4. eksctl >= 0.24.0
5. kubebuilder >= 2.3.1

Refer step 1 (only) from https://aws.amazon.com/blogs/opensource/introducing-fine-grained-iam-roles-service-accounts/
to create EKS cluster and create OIDC ID provider using eksctl

Configure kubectl to work with AWS EKS cluster. i.e. Valid AWS credentials are configured in ~/.aws/ directory
https://stackoverflow.com/questions/53266960/how-do-you-get-kubectl-to-log-in-to-an-aws-eks-cluster

Once EKS cluster is set up double check that it is set as correct context to be used to work with kubectl

`#> kubectl config get-contexts`

`#> kubectl config current-context`

Once you have a working kubernetes setup clone this repository and change present working directory to ./cloudformationCRD

Make sure all your GO binaries can be accessed from this location. (set $GOPATH and $GOROOT environment variables appropriately)

# Usage:

Run the deploy script:

`./deploy.sh cluster_name container_image_repo_path`

cluster_name: Name of the EKS cluster created above
container_image_repo_path: Full valid path to an image into authorized container registry. e.g. GCR or Docker

example command:

`./deploy.sh Cluster2 docker.io/kaustud/cloudformationcrd:latest`

OR

If you do not want to provide container image path, run the commands in the script one by one in the same sequence without providing image path to commands `make docker-build docker-push` and `make deploy` (untested)

What happens behind the scenes?
1. The script creates a new namespace 'cloudformationcrd-system'.
2. Using eksctl utility, a new service account is created in this namespace annotated by the newly created IAM role with
   required policies attached such as AWSCloudFormationFullAccess and AmazonEC2FullAccess and AmazonS3FullAccess.
3. New kinds defined in the api/v1/ directory are created. Currently only EC2 Instance is defined.
4. A container image for the manager is built and pushed to the repository provided.
5. The manager is deployed as a deployment object in the namespace created above with necessary permissions to run.
6. Create a new object of the each kind defined. The specs are defined in config/samples directory.

# Verify the installation:

1. If deployed correctly the IAM-IRSA-WEB-IDENTITY-HOOK populates the controller pod with two enviroment variables.

`#> kubectl -n cloudformationcrd-system get pods`

`#> kubectl -n cloudformationcrd-system get pod <pod_name> -o yaml`

exmple output:

> env:
> - name: AWS_ROLE_ARN
>   value: arn:aws:iam::<account_id>:role/rolename
> - name: AWS_WEB_IDENTITY_TOKEN_FILE
>   value: /var/run/secrets/eks.amazonaws.com/serviceaccount/token

