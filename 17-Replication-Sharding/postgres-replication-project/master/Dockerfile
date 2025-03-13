FROM postgres:13-alpine

# Копируем скрипт инициализации
COPY init.sql /docker-entrypoint-initdb.d/

# Копируем конфигурационные файлы
COPY postgresql.conf /etc/postgresql/postgresql.conf
COPY pg_hba.conf /etc/postgresql/pg_hba.conf

# Указываем PostgreSQL использовать наши конфигурационные файлы
CMD ["postgres", "-c", "config_file=/etc/postgresql/postgresql.conf", "-c", "hba_file=/etc/postgresql/pg_hba.conf"]
