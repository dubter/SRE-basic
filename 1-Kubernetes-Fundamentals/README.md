# №1 "Основы Kubernetes"

## 1 задание
Выстройте последовательность действий в процессе создания Pod.

1. Клиент вызывает kube-API, сервер записывает данные в etcd
2. Controller-manager видит изменение в состоянии и инициирует процесс создания pod
3. Controller-manager вызывает kube-API, сервер записывает состояние в etcd

4. Scheduler видит новые pod’ы и планирует их размещение на worker-нодах
5. Scheduler передает информацию о размещении pod’а в kube-API
6. Сервер записывает данные в etcd
7. Kubelet запускает pod, назначенный на worker-ноду, на которой расположен
8. Kubelet передает статус pod’а в kube-API, сервер записывает статус pod’а в etcd

## 2 задание
Вам как инженеру требуется запустить приложение в Kubernetes.
При этом нужно уметь автоматически обновлять приложение, управлять сетевым трафиком(запросы могут попадать на разные pod’ы).
Также приложение на старте пишет в директорию /my-dir, при этом данные в этой директории в последующем приложению не нужны и могут быть удалены.
У приложения есть конфигурация. Какие манифесты и директивы примените для развертывания данного приложения?

```sh
# При этом нужно уметь автоматически обновлять приложение,
Deployment
# управлять сетевым трафиком(запросы могут попадать на разные pod’ы).
Service
# У приложения есть конфигурация.
ConfigMap
# Также приложение на старте пишет в директорию /my-dir, при этом данные в этой директории в последующем приложению не нужны и могут быть удалены.
emptyDir
```

## 3 задание
В рамках задачи Вам необходимо:
Установить minikube для выполнения последующих практических заданий
```sh
brew install minikube
```

Ознакомиться с информацией в README в приложенных материалах к семинару (manifests.zip) и выполнить развёртывание нагрузок-примеров по указанным инструкциям.

В качестве решения приложите видеозапись экрана, на которой демонстрируются: