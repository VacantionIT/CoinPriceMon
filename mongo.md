# Подсказка для работы с Mongo через Docker

```
# скачать образ
sudo docker pull mongo

# список локальных образов
sudo docker images

# запустить монго с доступом по localhost
sudo docker run --name mongodb -d -p 27017:27017 mongo

# зупустить монго, который уже запускали с пред. параметрами
sudo docker start mongodb

# список работающих контейнеров 
sudo docker ps

# остановить контейнер
sudo docker stop mongodb 

```
