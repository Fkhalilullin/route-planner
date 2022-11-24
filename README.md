# route-planner
Построение маршрутов 

# Пример запроса:
```
		//50.6874,63.9391,50.7074,63.9454
curl http://localhost:8000/route\?from_lat\=5\&from_lon\=10\&to_lat\=25\&to_lon\=40 | jq
curl http://localhost:8000/points\?min_lat\=63.9391\&min_lon\=50.6874\&max_lat\=63.9454\&max_lon\=50.7074 | jq
```

