minecraft "Burpcraft" {
  version  = "1.19"
  snapshot = false
  path     = "./burpcraft/1.19.fabric"

  fabric_loader    = "0.14.7"
  fabric_installer = "0.11.0"

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
    level-seed        = "-156227665"
    gamemode          = "survival"
    max-players       = 10
    force-gamemode    = true
    enable-query      = true
    "query.port"      = 25565
  }
}

// allows for remote control over RPC
control_address = "127.0.0.1:2311"
control_token   = "s3cr3t"

// sets the game rules for the server, these are set as soon as the world has finished loading
// these are not part of the server block as they are executed as commands through the console
game_rules = {
  disableElytraMovementCheck = "true"
  doFireTick                 = "false"
  doLimitedCrafting          = "false"
  forgiveDeadPlayers         = "true"
  playersSleepingPercentage  = "1"
  showDeathMessages          = "true"
  spawnRadius                = "0"
  universalAnger             = "true"
}

discord_server_id      = "152083503767486464"
discord_server_channel = "984537987117383740"
