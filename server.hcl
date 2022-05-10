server "Burpcraft" {
  version  = "1.18.2"
  snapshot = false
  path     = "./burpcraft/vanilla"

  memory_min = 1
  memory_max = 4

  java_extra_args = [
    "-XX:+UseG1GC",
    "-XX:ParallelGCThreads=2",
    "-XX:MinHeapFreeRatio=5",
    "-XX:MaxHeapFreeRatio=10",
  ]

  // only need to put stuff we are overwriting here
  properties = {
    level-name        = "burpcraft"
    motd              = "Hello world"
    difficulty        = "hard"
    allow-flight      = true
    spawn-protection  = 0
    enforce-whitelist = true
    white-list        = true
    pvp               = false
    level-seed        = "r0flc0pt3r"
    gamemode          = "survival"
    max-players       = 20
    force-gamemode    = true
  }
}
