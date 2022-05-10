server "Burpcraft" {
  version  = "1.12.2"
  snapshot = false
  path     = "./burpcraft/legacy"

  memory_min = 1
  memory_max = 4

  java_extra_args = [
    "-XX:+UseG1GC",
    "-XX:ParallelGCThreads=2",
    "-XX:MinHeapFreeRatio=5",
    "-XX:MaxHeapFreeRatio=10",
  ]

  properties = {
    motd             = "Hello world"
    allow-flying     = false
    spawn-protection = 16
  }
}
