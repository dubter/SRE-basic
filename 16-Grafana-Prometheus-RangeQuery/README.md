# Домашнее задание №18 "Grafana, Prometheus и Range Query"

У нас есть сервис, на который жалуются пользователи, но мы не видим проблем по метрикам в панелях.

**Твоя задача:**
1) развернуть инфраструктуру по ниже стоящей инструкции - 2 балла
2) изучить дешборд и найти проблему - 2 балла
3) прислать исправленный вариант(export Dashboard) и расписать проблемы, которые нашли - 6 баллов

**Инструкция:**
1) Запуск инфраструктуры
```
docker-compose.yml
test-series.txt
```
2) Загрузить тестовые данные
```
# docker compose exec -it prometheus sh
# promtool tsdb create-blocks-from openmetrics test-series.txt /tmp/backfill_data
# cp -r /tmp/backfill_data/* /prometheus
```
3) Добавить Data Source
4) Загрузить дешборд
```
http-requests.json
```

# Решение

[Дашборд](http-requests-corrected.json)

1) на исходном дашборде не было данных на половине экрана. Соответственно, я увеличил масштаб. 
2) С нормализацией масштаба, `_interval` стал `20m`. Как мы знаем `_interval == step == timeWindow`, чтобы не пропускать пики. Поэтому я выставил timeWindow == step == 20m. Увидели пики в 100 requests per scrape interval 
3) Также выставил `min_interval = 1m = 2 * 30s = 2 * scrape_interval`
4) Поменял название для http_requests gauge единиц измерения. Были `ms`, выставил `req`.	


