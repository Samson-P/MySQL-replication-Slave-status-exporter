# MySQL-replication-Slave-status-exporter
Golang MySQL replication Slave status exporter for prometheus

Постановка задачи:
  Настроить email уведомления на ошибки репликации базы данных MySQL, prometheus мониторинг репликации (ссылка на дашборд Grafana с текущим статусом).

Проблема, с которой прошлось столкнуться перед принятием решения писать свой exporter:
  mysqld_exporter (github.com/prometheus/mysqld_exporter), к большому удивлению, не собирает метрик по состоянию репликации


v2.1 не прошла код-ревью!!
