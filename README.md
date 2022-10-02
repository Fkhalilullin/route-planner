# route-planner
Построение маршрутов 

# Пример запроса:
```
curl http://localhost:8000/route\?from_lat\=5\&from_lon\=10\&to_lat\=25\&to_lon\=40 | jq
```

## Todo backend
+ Получить запрос от map в виде xml
+ Запарсить xml и получить координаты препятсвий
+ Рассчитать маршрут и вернуть в json
