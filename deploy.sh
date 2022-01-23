env=$1 #prod or dev
hasura_endpoint=$2
admin_secret=$3

git clone https://github.com/CudoVentures/cudos-bdjuno.git CudosBDJuno
git clone https://github.com/CudoVentures/big-dipper-2.0-cosmos.git BigDipper2
cp -R /bdjuno-config CudosBDJuno/
cp .env.bdjuno CudosBDJuno/.env
cp .env.big-dipper-2 BigDipper2/.env

cd CudosBDJuno
echo "Starting BDJuno docker-compose"
if [ $env = "prod" ]; then
    docker-compose up -d --file=docker-compose-prod.yml
elif [ $env = "dev" ]; then
    docker-compose up -d --file=docker-compose-dev.yml
else
    echo "Wrong env passed: can be either dev or prod"
fi
cd hasura
curl -L https://github.com/hasura/graphql-engine/raw/stable/cli/get.sh | bash
hasura metadata apply --endpoint $hasura_endpoint --admin-secret $admin_secret

cd ../../BigDipper2
echo "Starting BigDipper2 docker-compose"
if [ $env = "prod" ]; then
    docker-compose up -d --file=docker-compose-prod.yml
elif [ $env = "dev" ]; then
    docker-compose up -d --file=docker-compose-dev.yml
else
    echo "Wrong env passed: can be either dev or prod"
fi
