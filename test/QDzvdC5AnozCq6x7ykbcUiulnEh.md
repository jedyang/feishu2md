# https

建议通过 cert-manager 创建和管理证书

# 配置 gateway

gateway 配置增加 https 的 listener

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: ys-nginx-gateway
  namespace: yunsheng
spec:
  gatewayClassName: nginx
  
  listeners:
  - name: https
    port: 443
    protocol: HTTPS
    hostname: ys.test.com
    tls:
      mode: Terminate
      certificateRefs:
      - name: test-tls-secret
        namespace: yunsheng
        kind: Secret
  - name: http
    protocol: HTTP
    port: 80
```

![](https://miaoji360.oss-cn-qingdao.aliyuncs.com/feishu2md/NehSb6c5UoMlruxRivCccsKRnqd.png)

# 配置 http 自动跳转 https

调整 httproute

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: ys-route-http
  namespace: yunsheng
spec:
  parentRefs:
    - name: ys-nginx-gateway
      sectionName: http
  hostnames:
    - "ys.test.com"
  rules:
    - matches:
        - path:
            type: PathPrefix
            value: /
      filters:
      - type: RequestRedirect
        requestRedirect:
          scheme: https  # 跳转至 HTTPS 协议
          port: 443      # 跳转至 443 端口
          statusCode: 308  # 308 = 永久重定向（推荐，比 301 更规范）
```

创建一个 httpsroute

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: ys-route-https
  namespace: yunsheng
spec:
  parentRefs:
    - name: ys-nginx-gateway
      sectionName: https
  hostnames:
    - "ys.test.com"
  rules:
    - matches:
        - path:
            type: PathPrefix
            value: /
      backendRefs:
      - name: ys-svc
        port: 80
```

也可以将两个 httproute 写在一个 yaml 中，但是分成两个更清晰，生产环境更推荐分开的做法。

测试一下

```
$ curl -v http://ys.test.com
```

![](https://miaoji360.oss-cn-qingdao.aliyuncs.com/feishu2md/QP6ubmg3QoeqmaxLhYKcXKhTnug.png)

发生了 301 跳转
