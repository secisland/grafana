# grafana
基于grafana-6.3.5改进报警机制


一、查询告警API

1、查询所有告警

curl http://api_key:eyJrIjoiUnFpS3k3MzNlcG52ZWJmbjFrZlNpMElaS3BFZTFVR3oiLCJuIjoieXNmIiwiaWQiOjF9@http://192.168.56.102:3000/api/alerts

or 

curl -H "Authorization: Bearer eyJrIjoiUnFpS3k3MzNlcG52ZWJmbjFrZlNpMElaS3BFZTFVR3oiLCJuIjoieXNmIiwiaWQiOjF9" http://192.168.56.102:3000/api/alerts


2、查具体某个仪表盘的告警

curl http://api_key:eyJrIjoiUnFpS3k3MzNlcG52ZWJmbjFrZlNpMElaS3BFZTFVR3oiLCJuIjoieXNmIiwiaWQiOjF9@http://192.168.56.102:3000/api/alerts?dashboardId=93


3、查具体某个图形的告警

curl http://api_key:eyJrIjoiUnFpS3k3MzNlcG52ZWJmbjFrZlNpMElaS3BFZTFVR3oiLCJuIjoieXNmIiwiaWQiOjF9@http://192.168.56.102:3000/api/alerts?dashboardId=93&panelId=7



二、编译

编译环境要求：
go 1.12
v10.22.0


cd $GOPATH/src/github.com/
mv secisland grafana
cd grafana

go run build.go setup
go run build.go build 


编译前端源码步骤
npm install -g yarn --registry=http://registry.npm.taobao.org
yarn install --pure-lockfile --registry=http://registry.npm.taobao.org
npm install -g grunt-cli --registry=http://registry.npm.taobao.org
// grunt
yarn start #, yarn start:hot, or yarn build


如果grunt中使用sass，还需要安装sass
mac自带ruby，所以直接在命令行输入: 
gem install -n /usr/loca/bin sass
