# go-practice
goの勉強用

### やっていること
* DBからのSELECT結果をグローバル変数に設定
* SQLのログ出力
* Prometheusの導入

### Reference
* [labstack/echo](https://github.com/labstack/echo) -> webフレームワーク
* [go-sql-tracer](https://github.com/walf443/go-sql-tracer) -> SQLをログ出力できるライブラリ
* [Realizeを使ってGoでホットリロードを実現する](https://qiita.com/godgarden/items/f73e4a717f1a27b9a3b0)
* [Prometheusを使ってISUCON9の監視をやってみる](https://qiita.com/K-jun/items/17a66c3d691e94bd8c45)
* [Prometheus Middleware](https://echo.labstack.com/middleware/prometheus)
* [felixge/fgprof](https://github.com/felixge/fgprof)

#### Prometheus
```
cp -i ./prometheus/prometheus.yml /tmp/
docker run -d -p 9090:9090 -v /tmp/prometheus.yml:/etc/prometheus/prometheus.yml --name prometheus prom/prometheus
```

http://localhost:9090

