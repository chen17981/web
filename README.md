# web version

A golang web App


Docker Image Build
===

```
git clone https://github.com/chen17981/web

cd web 

docker build -t my-web .
```


Run app in Docker
===

```
docker run --rm -p 8080:8080 my-web

```

1. Launch a web browser, go to http://localhost:8080/

2. Input the price list like "CH1,3.11,AP1,6.00,CF1,11.23,MK1,4.75,OM1,3.69" in the form "Product Prices" or leave it empty (the default value will be used). The default value is "CH1,3.11,AP1,6.00,CF1,11.23,MK1,4.75,OM1,3.69".

3. Input the shopping items like "AP1,AP1,AP1,OM1" in the form "Shopping Items", leaving it empty will cause an error of "invalid input".

4. Click "submit" button, it will display the result.
