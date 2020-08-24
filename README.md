# cloudformationCRD
AWS objects in kubernetes using respective cloudformation specs

Prerequisites:
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

Run the deploy script:

`./deploy.sh cluster_name container_image_repo_path`

cluster_name: Name of the EKS cluster created above
container_image_repo_path: Full valid path to an image into authorized container registry. e.g. GCR or Docker

example command:

`./deploy.sh Cluster2 docker.io/kaustud/cloudformationcrd:latest`