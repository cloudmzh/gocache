# gocache

gocache: Distributed cache framework based on consistent hashing algorithm, only applicable to Go language novice learning projects

#### 描述
gocache是一个基于Go语言实现的分布式缓存框架，参考了groupcache的设计，适合go语言、分布式初学者的开源项目。


#### 特点
- gocache采用了2Q算法的方式读写缓存
[2Q: A Low Overhead High Performance Buffer Management Replacement Algorithm](http://www.vldb.org/conf/1994/P439.PDF)
- gocache采用一致性哈希算法(虚拟节点)实现分布式储存
- gocache采用singleflight的方法预防缓存击穿
- 每个gocahe既具有客户端功能又具有服务端功能
- gocache具有Get/Delete操作，适合作为持久化数据源的分布式缓存

#### 使用Demo


```go
type dbSource map[string]float64
//模拟数据源
var dbsource = dbSource{
	"apple":      6.90,
	"pear":       18.20,
	"grape":      15.40,
	"banana":     11.0,
	"watermelon": 30.0,
	"melon":      35.0,
	"strawberry": 66.40, //....
}

func main() {
	group := gocache.NewGroup("price", 2<<5, 2<<4 , gocache.GetterFunction(
		// 调用远程数据源
		func(key string) ([]byte, error) {
			log.Println("access dbsource", key)
			if v, ok := dbsource[key]; ok {
				return []byte(fmt.Sprintf("%.2f", v)), nil
			}
			return nil, fmt.Errorf("[dbsource] %s is not exist ", key)
		}))
    
    // 集群
	addrs := []string{"http://localhost:8000", "http://localhost:8001", "http://localhost:8002", "http://localhost:8003"}
    
    // 启动gocache
	addr := localhost:8000
	peers := gocache.NewHTTPPool("http://" + addr)
	peers.Set(addrs...)
	group.RegisterPeers(peers)
	log.Println("gocache is runnning...", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}

```

> 分别启动4台主机，addrs分别修改成 *http://localhost:8000, http://localhost:8001, http://localhost:8002, http://localhost:8003*。我们可以选取任意一台主机当作API主机，
> 因为每个gocahe既具有客户端功能又具有服务端功能。
> 
现在访问8000端口主机寻找price命名组里的缓存pear：
```shell
curl http://localhost:8000/_gocache/price/pear
```

结果：

```
# http://localhost:8000 结果：
2020/08/13 21:55:24 gocache is runnning... localhost:8000
2020/08/13 21:57:09 [Server http://localhost:8000] GET /_gocache/price/pear
2020/08/13 21:57:09 [Server http://localhost:8000] Pick peer http://localhost:8003

# http://localhost:8003 结果：
2020/08/13 21:55:24 gocache is runnning... localhost:8003
2020/08/13 21:57:09 [Server http://localhost:8003] GET /_gocache/price/pear
2020/08/13 21:57:09 access dbsource pear
2020/08/13 21:57:09 [Group price] Add cache,key is pear
```

> 可以明显的看出8000端口主机通过一致性哈希算法定位到8003端口主机，8003端口主机的该缓存为空，于是访问本地数据源，将缓存添加至8003端口主机中。


再次访问8002端口主机寻找price命名组里的缓存pear：
```shell
curl http://localhost:8002/_gocache/price/pear
```

结果：

```
# http://localhost:8002 结果：
2020/08/13 22:01:52 [Server http://localhost:8002] GET /_gocache/price/pear
2020/08/13 22:01:52 [Server http://localhost:8002] Pick peer http://localhost:8003

# http://localhost:8003 结果
2020/08/13 22:01:52 [Server http://localhost:8003] GET /_gocache/price/pear
2020/08/13 22:01:52 [Group price] Hit cache,key is pear
```

当8003端口主机宕机，再次访问8002端口主机寻找price命名组里的缓存pear

结果：


```
# http://localhost:8002 结果：
2020/08/13 22:04:24 [Server http://localhost:8002] GET /_gocache/price/pear
2020/08/13 22:04:24 [Server http://localhost:8002] Pick peer http://localhost:8003
2020/08/13 22:04:27 [Group] Failed to get from peer price
2020/08/13 22:04:27 access dbsource pear
2020/08/13 22:04:27 [Group price] Add cache,key is pear
```
