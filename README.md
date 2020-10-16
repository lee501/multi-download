#### 多线程现在文件

#### 使用方式：
- 说明：文件需要支持多点续传
----
    1. go build
```
    生成可执行文件multi-download
```

    2. ./multi-download -url="file site" -name="filename" -total=10 -path="filepath" -sha256=""
        或./multi-download -url="file site"
```
    参数说明：
        url: 文件下载路径(必填)
        name: 文件名字(可选, 不传使用下载文件的名字)
        total: 线程数量(可选, 默认为10)
        path: 文件保存路径(可选, 不传使用当前项目路径)
        sha256: 文件256校验字符串(可选, 需下载文件提供)
```
