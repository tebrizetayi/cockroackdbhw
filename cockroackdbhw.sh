docker network create -d bridge roachnet
docker container stop roach1
docker container rm roach1
docker run -d \
--name=roach1 \
--hostname=roach1 \
--net=roachnet \
-p 26257:26257 -p 3333:8080  \
-v "${PWD}/cockroach-data/roach1:/cockroach/cockroach-data"  \
cockroachdb/cockroach:v20.1.6 start \
--insecure \
--join=roach1,roach2,roach3

docker container stop roach2
docker container rm roach2
docker run -d \
--name=roach2 \
--hostname=roach2 \
--net=roachnet \
-v "${PWD}/cockroach-data/roach2:/cockroach/cockroach-data" \
cockroachdb/cockroach:v20.1.6 start \
--insecure \
--join=roach1,roach2,roach3

docker container stop roach3
docker container rm roach3
docker run -d \
--name=roach3 \
--hostname=roach3 \
--net=roachnet \
-v "${PWD}/cockroach-data/roach3:/cockroach/cockroach-data" \
cockroachdb/cockroach:v20.1.6 start \
--insecure \
--join=roach1,roach2,roach3



#docker exec -it roach1 ./cockroach init --insecure
#docker exec -it roach1 ./cockroach sql --insecure
#CREATE DATABASE ACCOUNT

