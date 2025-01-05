# Домашнее задание №5 "Метрики. Prometheus"

## 1 задание
Домашнее задание по разворачиванию Prometheus:
- можно выполнить на локальном компьютере используя docker-compose, просто развернув prometheus из образа официальной документации https://prometheus.io/docs/prometheus/latest/installation/
- можно выполнив команды представленные ниже на виртуальной машине
```
# Обновим зависимости
sudo apt update && sudo apt install wget systemctl -y

# Скачиваем Prometheus
wget "https://github.com/prometheus/prometheus/releases/download/v2.37.0/prometheus-2.37.0.linux-amd64.tar.gz"

# Распаковываем
tar xvf prometheus-2.37.0.linux-amd64.tar.gz -C /tmp

# Копируем бинарники
sudo cp /tmp/prometheus-2.37.0.linux-amd64/prometheus /usr/local/bin/
sudo cp /tmp/prometheus-2.37.0.linux-amd64/promtool /usr/local/bin/

# Создаем папки
sudo mkdir -p /etc/prometheus
sudo mkdir /var/lib/prometheus

# Копируем файлы
sudo cp -r /tmp/prometheus-2.37.0.linux-amd64/consoles /etc/prometheus
sudo cp -r /tmp/prometheus-2.37.0.linux-amd64/console_libraries /etc/prometheus

# Создаем пользователя в системе без возможности логина (для безопасности)
sudo useradd --no-create-home --shell /bin/false prometheus

# Ставим правильного пользователя владельцем
sudo chown prometheus:prometheus -R /etc/prometheus
sudo chown prometheus:prometheus -R /var/lib/prometheus/

# Создаем базовый конфиг
sudo cat << EOF > /etc/prometheus/prometheus.yml
global:
  scrape_interval: 15s
scrape_configs:
  - job_name: 'prometheus'
    scrape_interval: 5s
    static_configs:
      - targets: ['localhost:9090']
EOF

# Создаем инитскрипт для system.d
sudo cat << EOF > /lib/systemd/system/prometheus.service

[Unit]
Description=Prometheus service

[Service]
User=prometheus
Group=prometheus

ExecStart=/usr/local/bin/prometheus --config.file /etc/prometheus/prometheus.yml --storage.tsdb.path /var/lib/prometheus/ --web.console.templates=/etc/prometheus/consoles --web.console.libraries=/etc/prometheus/console_libraries

[Install]
WantedBy=multi-user.target

EOF

# Релоадим демона, ставим в автозапуск и запускаем
sudo systemctl daemon-reload
sudo systemctl enable prometheus.service
sudo systemctl start prometheus.service

asdas
```

Теперь по адресу <ip>:9090 будет доступен Prometheus
<ip> - либо адресс виртуальной машины, либо адресс вашей локальной машины

> [!NOTE]
> **Критерии оценки**:
> * 3 балла - скриншот/видеозапись развернутого Prometheus 
> * 2 балла - скриншот/видеозапись выполнить запрос базовых метрик, которые есть Prometheus

> [!INFO]
> Написал базовый [prometheus.yml](build/prometheus/prometheus.yml)

```sh
make up
```

```sh
make down
```
## 2 задание
Включить метрики Oncall и добавить их targetом в Prometheus.

### Пререквизиты
1. Включаем ingress в minikube
```sh
minikube addons enable ingress
```

2. Проверяем запустился ли ingress controller
```sh
kubectl get pods -n ingress-nginx
```

### Флоу задания

Включить метрики в формате prometheus в приложении Oncall
В склонированном ранее
- для этого придется пересобрать Dockerfile с установкой prometheus_client
  ```dockerfile
  # Включить метрики в формате prometheus в приложении Oncall,
  # пересобираем Dockerfile с установкой prometheus_client
  RUN chown -R oncall:oncall /home/oncall/source /var/log/nginx /var/lib/nginx \
      && sudo -Hu oncall mkdir -p /home/oncall/var/log/uwsgi /home/oncall/var/log/nginx /home/oncall/var/run /home/oncall/var/relay \
      && sudo -Hu oncall python3 -m venv /home/oncall/env \
      && sudo -Hu oncall /bin/bash -c 'source /home/oncall/env/bin/activate && cd /home/oncall/source && pip install wheel && pip install prometheus_client && pip install .'
      #                                                                                                                     !-------------------------------!
  ```
  В миникубе билдим новый образ
  ```sh
  eval $(minikube docker-env)
  # Переходим в склониную репу
  # cd ./oncall
  docker build --no-cache -t oncall:latest .
  # Check that image of oncall now in minikube
  docker images
  ```
  ```md
  REPOSITORY                                TAG        IMAGE ID       CREATED         SIZE
  oncall                                    latest     7e72879b6da6   2 minutes ago   717MB
  registry.k8s.io/kube-apiserver            v1.31.0    cd0f0ae0ec9e   3 months ago    91.5MB
  registry.k8s.io/kube-controller-manager   v1.31.0    fcb0683e6bdb   3 months ago    85.9MB
  registry.k8s.io/kube-scheduler            v1.31.0    fbbbd428abb4   3 months ago    66MB
  registry.k8s.io/kube-proxy                v1.31.0    71d55d66fd4e   3 months ago    94.7MB
  registry.k8s.io/etcd                      3.5.15-0   27e3830e1402   3 months ago    139MB
  registry.k8s.io/pause                     3.10       afb61768ce38   5 months ago    514kB
  registry.k8s.io/coredns/coredns           v1.11.1    2437cf762177   15 months ago   57.4MB
  gcr.io/k8s-minikube/storage-provisioner   v5         ba04bb24b957   3 years ago     29MB
  ```

  ```yml
  # config.yml
  oncall.conf: |
    ---
    metrics: prometheus
    # ...
    prometheus:
        oncall-notifier:
        server_port: 9091
  ```
  - вытащить дополнительные порты в Deployment и Service
  ```yml
  # deployment.yml
  ports:
    - containerPort: 8080
    - containerPort: 9091 # Вытаскиваем дополнительные порты в Deployment
  ```
  ```yml
  # service.yml
  ports:
    - protocol: TCP
      port: 8000
      targetPort: 8000
  ```
  - прописать новый path в Ingress
  ```yml
  # ingress.yml
  pathType: Prefix
  metrics:
    service:
      name: oncall
      port:
        number: 8080
  ```

- Добавить адрес с метриками в target Prometheus prometheus.yml
- Сделать запрос к метрикам Oncall

0 добавить адрес будущей локальной инсталляции OnCall в /etc/hosts
```sh
127.0.0.1 oncall.local
127.0.0.1 oncall.metrics.local
```

1. Включаем ingress в minikube
```sh
minikube addons enable ingress
```

2. Проверяем запустился ли ingress controller
```sh
kubectl get pods -n ingress-nginx
```

3. Применяем манифесты
```sh
kubectl apply -f <folder>
```


Делаем ingress доступным на хостовой машине, это связано с особенностями minikube. Не закрываем данное окно!!!
```
minikube tunnel
```

NB.
Обязательно нужно чекнуть что поднялось:
```sh
kubectl exec -it имя_пода_с_oncall -- bash
```
Для отвратительнейший шел, в котором нужно чекнуть типа:
```sh
# По этому пути лежат гавно логи
# ls /home/oncall/var/log/uwsgi
# в файлике error.log


# Все заебумба если он не падает
cat /home/oncall/var/log/uwsgi/error.log | grep notifier
```

```
2024-11-17 13:36:57,683 INFO root Setting metrics gauge message_blackhole_cnt to 0
2024-11-17 13:36:57,684 INFO root Setting metrics gauge message_sent_cnt to 0
2024-11-17 13:36:57,684 INFO root Setting metrics gauge message_fail_cnt to 0
```


## 2 задание

В предыдущих семинарах, мы разворачивали Oncall, в текущем задании нам надо разобраться как включить метрики Oncall и добавить их targetом в Prometheus.

* Включить метрики в формате prometheus в приложении Oncall
  * для этого придется пересобрать Dockerfile с установкой prometheus_client
  * вытащить дополнительные порты в Deployment и Service
  * прописать новый path в Ingress
* Добавить адрес с метриками в target Prometheus prometheus.yml
* Сделать запрос к метрикам Oncall

> [!NOTE]
> **Критерии оценки**:
> * 1 балл - Показать/прислать конфигурацию, как включить метрики в Oncall
> * 1 балл - Показать/прислать prometheus.yml
> * 3 балла - Скриншот/видеозапись из Prometheus UI с запросом любой метрики из Oncall

## Решение

Включить метрики Oncall и добавить их targetом в Prometheus.

Включить метрики в формате prometheus в приложении Oncall

Меняем Dockerfile

Смотрим как метрики используются, для модификаций кофигмапа
```py
# linkedin src/oncall/metrics
config['prometheus'][appname]['server_port']
```

```yml
prometheus

    server_port
```





добавил CLUSTER-IP сервиса в `/etc/hosts`
10.106.81.214

Проверим что ингресс применился
```
kubectl get ingress
kubectl describe ingress
```


Увеличиваем скрейп интервалы и
--add-host=host.docker.internal:host-gateway
