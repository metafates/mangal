package main

import "time"

const AppName = "Mangal"
const CachePrefix = AppName + "Cache"
const TempPrefix = AppName + "Temp"
const Parallelism = 100
const TestQuery = "Death Note"
const Forever = time.Duration(1<<63 - 1) // 292 years
