# 1 Billion Row Challenge

WIP experimentation to do the [1BRC](https://www.morling.dev/blog/one-billion-row-challenge/).

The original one was for Java but I decided to go with Golang

## Dataset

https://huggingface.co/datasets/nietras/1brc.data/tree/main

## Iterative Process

* `v1` : First try (~2 min 50s)
* `v2` : Optimizations like goroutine and manual int parsing (~35s)
* `v3` : Hashcode (~14s)