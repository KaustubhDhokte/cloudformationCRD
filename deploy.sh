echo "Creating namespace ... "

# kubectl create -f config/manager/namespace.yaml

echo '###################################################################################################'

echo "Creating iamserviceaccount ... "

# eksctl create iamserviceaccount \
# --name awsobjects \
# --namespace cloudformationcrd-system \
# --cluster $1 \
# --attach-policy-arn arn:aws:iam::aws:policy/AWSCloudFormationFullAccess \
# --attach-policy-arn arn:aws:iam::aws:policy/AmazonEC2FullAccess \
# --attach-policy-arn arn:aws:iam::aws:policy/AmazonS3FullAccess \
# --approve

echo '###################################################################################################'

echo "Creating resources ... "

make install

echo '###################################################################################################'

echo "Manager image Build and push ... "

make docker-build docker-push IMG=$2

echo '###################################################################################################'

echo "Deploying the controller image ... "

make deploy IMG=$2

echo '###################################################################################################'

echo "Creating sample resources in AWS ... "

#kubectl create -f config/samples/cloudf_v1_ec2instance.yaml

echo '###################################################################################################'

echo "Done !!!"

echo "Verify created resources through AWS management console."
