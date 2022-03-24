Summary for profiling.

During investigation profile my application I didn't find any big troubles with CPU/Memory. 
For memory optimization I can use in application library [fastjson](https://github.com/valyala/fastjson)
for validation and parsing JSON data instead of `json.Decode` and `go-playground/validator`. 

I used standard router because I not really need any complex features for routes.
But I can used more optimized router such as [httprouter](https://github.com/julienschmidt/httprouter)
or [fasthttp](https://github.com/valyala/fasthttp). Its conclusion based on 
[benchmarks](https://medium.com/@smallnest/go-web-framework-benchmark-93a34403ef0a).

Logger. I used logrus for logging. It can store logs asynchronous with implementing custom hook. 
For easier support this feature can use next solutions:
* [zap](https://github.com/uber-go/zap) - have an [issue](https://github.com/uber-go/zap/issues/988) where 
discussed that this package haven't fully asynchronously functionality
* [logr](https://github.com/mattermost/logr) - haven't any benchmarks.
* [go-logging](https://github.com/ccding/go-logging)

I used concurrency for reading messages because it's 
can be faster to get income messages and save them in several streams that in one stream.