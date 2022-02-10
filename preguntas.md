1-> controllers: utilizacion de validator. En este caso, cuando declaro en cada archivo un var validate = validator.New()
salta la advertencia de redeclared on this block. Se declara en un archivo aparte como variable global o se redeclara en cada
controller con un nombre diferente?


docker build -t dev-stats-go .
docker run -it --name dev-stats-go -p 3010:3010 -v ${PWD}:/go/src/statsv0 dev-stats-go