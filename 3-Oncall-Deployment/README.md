# №3 "Развёртывание Oncall"
## 1 задание
Для закрепления пройденного материала Вам необходимо выполнить развёртывание OnCall в локальной инсталляции minikube.

Изучите материалы семинара. Выполните шаги для развёртывания, описанные в инструкции README, приложенной к материалам и ниже.

В качестве решения приложите видеозапись демонстрации экрана. На видео должно быть продемонстрировано выполнение всех пунктов. Обязательно продемонстрируйте доступность OnCall в браузере в конце видео.

```sh
export http_proxy=''
export https_proxy=''
export no_proxy='localhost,127.0.0.1,.1,.2,.3::1,.local'
```

1. Включаем ingress в minikube
```sh
minikube addons enable ingress
```

2. Проверяем запустился ли ingress controller
```sh
kubectl get pods -n ingress-nginx
```


```sh
git clone https://github.com/linkedin/oncall.git
cd oncall

# Настроить переменные для использования демона Docker в Minikube
eval $(minikube docker-env)

docker build -t oncall:latest .

# Если выполняете задания на ВМ, необходимо заменить Dockerfile и указать прокси-сервер для корректной сборки:
cp ../docker/Dockerfile Dockerfile

docker build -t oncall:latest --build-arg http_proxy=http://10.128.0.4:3128 --build-arg HTTP_PROXY=http://10.128.0.4:3128 --build-arg https_proxy=http://10.128.0.4:3128 --build-arg HTTPS_PROXY=http://10.128.0.4:3128 .
```

добавить адрес будущей локальной инсталляции OnCall в /etc/hosts
```sh
127.0.0.1 oncall.local
```

### 1. Запуск MySQL
- statefulset.yaml
- service.yaml
- secret.yaml

Создаём Secret с паролем к БД и Headless Service:
```sh
cd mysql
kubectl apply -f secret.yaml
kubectl apply -f service.yaml
```
Запускаем нагрузку MySQL
```sh
kubectl apply -f statefulset.yaml
```
Проверяем, запущен ли под:
```sh
kubectl get pod
```
```
NAME READY STATUS RESTARTS AGE
mysql-0 1/1 Running 0 10m
```

# 2. Запуск OnCall
Создаём ConfigMap с конфигом приложения и Service
```sh
cd oncall
kubectl apply -f config.yaml
kubectl apply -f service.yaml
```
Запускаем нагрузку OnCall
```sh
kubectl apply -f deployment.yaml
```

Проверяем, запущены ли поды
```sh
kubectl get pod
```
```
NAME READY STATUS RESTARTS AGE
mysql-0 1/1 Running 0 10m
oncall-94855d89d-6qggk 1/1 Running 0 10m
oncall-94855d89d-k78xm 1/1 Running 0 10m
```
Создаём Ingress для доступа к сервису:
```sh
kubectl apply -f ingress.yaml
```
Проверяем, создан ли ресурс Ingress
```sh
kubectl get ingress
```
```
NAME CLASS HOSTS ADDRESS PORTS AGE
oncall-ingress nginx oncall.local 192.168.49.2 80 5d2h
```
Делаем ingress доступным на хостовой машине, это связано с особенностями minikube. Не закрываем данное окно!!!
```
minikube tunnel
```
Проверяем доступность OnCall по адресу `oncall.local`
```sh
curl http://oncall.local
```

## 2 задание
Изучите манифесты, приложенные к материалам семинара.
Какие проблемы (нарушение best practices, потенциальные уязвимости и т.д.) Вы заметили?
Опишите минимум две обнаруженные проблемы и способы их исправления.

Грехи 12 factor architecture
Очень много хардкода, и паплайны с триажами AppSec'ов точно не пройдет.

📁 my-sql
Захардкожен заэнкоженный секрет `MTIzNA==`
Если это доступ к stg/prod бд, злоумышленник может сделать вот так;
```sh
echo MTIzNA== | base64 --decode
```
По хорошему в мы должны забирать его из Волта, или аналога:
```yml
---
apiVersion: v1
kind: Secret
metadata:
  name: vault-mysql-secrets
  namespace: vault-operator
type: Opaque
data:
  MYSQL_ROOT_PASSWORD: $(echo -n s.W4ndAJbuoDMsoLDZyVBG18F2 | base64)
```


📁 oncall
Конфигурация
- зашита прямо в образ — что-то меняяем придется каждый раз пересобирать образ.
- по хорошему должана быть защита в переменные окружения — а у нас порты захаркожены, вместо ваилдкардов.
