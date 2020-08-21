kubectl create -f config/manager/namespace.yaml

eksctl create iamserviceaccount --name cloudfservice --namespace cloudformationcrd-system --cluster $1 --attach-policy-arn arn:aws:iam::aws:policy/AWSCloudFormationFullAccess --approve

make install

make docker-build docker-push IMG=$2

make deploy IMG=$2

kubectl create -f config/samples/cloudf_v1_ec2instance.yaml
