# Домашнее задание №11 "Балансировка"

## 1 задание

Что использует под собой сущность Service для балансировки 

- [x] iptables 
- [ ] fwgw 
- [ ] Внешний тип балансировки

## 2 задание

Что такое мультиплексированное соединение? 

- [x] Сетевое соединение в рамках которого можно посылать и принимать разные сообщения от разных потоков данных 
- [ ] Способ балансировки в кластерах kubernetes 
- [ ] Сетевое соединение, которое гарантирует доставку пакетов 
- [ ] Сетевое соединение, которое не гарантирует доставку пакетов

## 3 задание

Какой тип балансировки отвечает за круговое распределение запросов с использованием весов? 

- [ ] Leastconn 
- [x] Round-robin weight 
- [ ] Round-robin 
- [ ] Sticky Session 
- [ ] Source Hash 
- [ ] Destination Hash

## 4 задание

В текущей инсталяции необходимо реализовать Sticky Session на основе след. параметров(любой на выбор):
* cookie
* headers

В скринкасте или видео показать, что разные запросы, гарантировано прилипают к своим подам, приложить манифест Ingress файла

## Решение

1) Разворачиваем манифест:
```sh
kubectl apply -f manifests
```
2) Первый запрос (сохраните cookie)
```sh
curl -v http://sticky-test.info
```
3) Последующие запросы с сохраненной cookie
```sh
curl -v --cookie "STICKYSESSIONID=<STICKYSESSIONID>" http://sticky-test.info
```