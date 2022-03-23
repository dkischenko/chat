Summary for profiling.
During investigation profile my application I didn't find any big troubles with CPU/Memory. 
For memory optimization I can use in application library [fastjson](https://github.com/valyala/fastjson) for validation and parsing JSON data instead of `json.Decode` and `go-playground/validator`. 