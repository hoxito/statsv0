


docker build -t dev-stats-go .
docker run -it --name dev-stats-go -p 3010:3010 -v ${PWD}:/go/src/statsv0 dev-stats-go