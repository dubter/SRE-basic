# Домашнее задание №18 "Grafana, Prometheus и Range Query"

У нас есть сервис, на который жалуются пользователи, но мы не видим проблем по метрикам в панелях.

**Твоя задача:**
1) развернуть инфраструктуру по ниже стоящей инструкции - 2 балла
2) изучить дешборд и найти проблему - 2 балла
3) прислать исправленный вариант(export Dashboard) и расписать проблемы, которые нашли - 6 баллов

**Инструкция:**
1) Запуск инфраструктуры
2) Загрузить тестовые данные
```
# docker compose exec -it prometheus sh
# promtool tsdb create-blocks-from openmetrics test-series.txt /tmp/backfill_data
# cp -r /tmp/backfill_data/* /prometheus
```
3) Добавить Data Source
4) Загрузить дешборд
