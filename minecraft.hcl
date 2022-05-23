minecraft "Burpcraft" {
  version  = "1.18.2"
  snapshot = false
  path     = "./burpcraft/vanilla"

  fabric_loader    = "0.13.3"
  fabric_installer = "0.10.2"

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

// allows for remote control over RPC
control_address = "127.0.0.1:2311"
control_token   = "s3cr3t"