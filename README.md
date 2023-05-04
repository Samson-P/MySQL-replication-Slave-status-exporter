# MySQL Slave replication status exporter
Golang MySQL Slave replication status exporter for prometheus. Gives database replication metrics to the scraper.

<p align="center"><img src="metrics.png" width=70%/></p>

## TODO
* HTTP сервер хостящий страничку с метриками;
* страница с метриками репликации;
* email уведомления на ошибки репликации базы данных MySQL;
* prometheus мониторинг репликации (ссылка на дашборд Grafana с текущим статусом).


## Актуальность
  Проблема, с которой пришлось столкнуться перед принятием решения писать свой exporter: [mysqld_exporter](https://github.com/prometheus/mysqld_exporter), к большому удивлению, не собирает метрик по состоянию репликации


## Stack
<p>
    <img src="https://img.icons8.com/color/48/000000/golang.png" alt="Go" width="30" height="30" />
    <img src="https://img.icons8.com/color/48/000000/mysql.png" alt="MySQL" width="30" height="30" />
    <img src="https://github.com/Samson-P/Samson-P/blob/main/img/prometheus.png" alt="prometheus" width=30px height=30px />
    <img src="https://img.icons8.com/color/48/000000/grafana.png" alt="grafana" width="30" height="30" />
    <img src="https://img.icons8.com/color/48/000000/centos.png" alt="docker" width="30" height="30" />
    <img src="https://img.icons8.com/color/48/000000/ubuntu.png" alt="docker" width="30" height="30" />
</p>

**Настройки IPTABLES можно посмотреть в [man](https://github.com/Samson-P/MySQL-replication-Slave-status-exporter/blob/main/man)**
