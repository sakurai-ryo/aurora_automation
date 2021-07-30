# !/bin/bash

set -eu

stage=""
while getopts s: OPT; do
    case $OPT in
    s)
        stage="$OPTARG"
        ;;
    esac
done

if [ -z ${stage} ] || [ -z ${AWS_PROFILE} ] || [ -z ${CLUSTER_NAME} ]; then
    exit 1
fi

cluster_arn=$(aws \
                rds describe-db-clusters \
                --db-cluster-identifier ${CLUSTER_NAME} \
                --profile ${AWS_PROFILE} \
                --region ap-northeast-1 \
                --query 'DBClusters[].DBClusterArn' \
                --output text
            )

if [ -z ${cluster_arn} ]; then
    echo "Cluster Not Found Error"
    exit 1
fi
echo ${cluster_arn}

GOOS=linux go build -o bin/main

sls deploy --aws-profile ${AWS_PROFILE} -s ${stage} --clusterName ${CLUSTER_NAME} --clusterArn ${cluster_arn}
