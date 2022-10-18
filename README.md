# deep-thought

## regular ol' app

![deep thought diagram](static//0-deep-thought-diagram.png)

`cd 0-deep-thought-uninstrumented/`

`docker-compose up --build`

```bash
$ curl localhost:4242
what is the answer to the ultimate question of life, the universe, and everything?
42
```

## instrumented app, with instrumentation for http requests

![deep thought diagram instrumented](static/1-deep-thought-diagram-instrumented.png)

`cd 1-deep-thought-instrumented/`

`export HONEYCOMB_API_KEY=<api-key>`

`docker-compose up --build`

```bash
$ curl localhost:4242
what is the answer to the ultimate question of life, the universe, and everything?
42
```

![deep thought instrumented http requests](static/1-deep-thought-instrumented.png)

## instrumented app with auto and manual instrumentation

`cd 2-deep-thought-instrumented-manual/`

`export HONEYCOMB_API_KEY=<api-key>`

`docker-compose up --build`

```bash
$ curl localhost:4242
what is the answer to the ultimate question of life, the universe, and everything?
42
```

![deep thought instrumented http requests and manual instrumentation](static/2-deep-thought-instrumented-manual.png)
