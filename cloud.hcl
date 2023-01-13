minecraft "Cloud" {
  version  = "1.19.3"
  snapshot = false
  path     = "./burpcraft/cloud"

  fabric_loader    = "0.14.12"
  fabric_installer = "0.11.1"

  memory_min = 4
  memory_max = 8

  java_extra_args = [
    "-XX:+UseG1GC",
    "-XX:ParallelGCThreads=2",
    "-XX:MinHeapFreeRatio=5",
    "-XX:MaxHeapFreeRatio=10",
  ]

  // only need to put stuff we are overwriting here
  properties = {
    level-name        = "coleslaw"
    motd              = "Welcome to Minecraft!"
    difficulty        = "hard"
    allow-flight      = true
    spawn-protection  = 0
    enforce-whitelist = true
    white-list        = true
    pvp               = false
    level-seed        = "7274395869746300621"
    gamemode          = "survival"
    max-players       = 10
    force-gamemode    = true
    enable-query      = true
    "query.port"      = 25565
    server-port       = 25565
  }

  purge_datapacks = false

  datapack "minecraft-datapack" {
    url = "https://github.com/renevo/minecraft-datapack/releases/download/v1.19-1.1.0/minecraft-datapack.zip"
  }
  datapack "afk-display" {
    url = "https://github.com/renevo/minecraft-datapack/releases/download/v1.19-1.1.0/afk-display-v1.1.3.zip"
  }
  datapack "anti-enderman-grief" {
    url = "https://github.com/renevo/minecraft-datapack/releases/download/v1.19-1.1.0/anti-enderman-grief-v1.1.3.zip"
  }
  datapack "cauldron-concret" {
    url = "https://github.com/renevo/minecraft-datapack/releases/download/v1.19-1.1.0/cauldron-concrete-v2.0.6.zip"
  }
  datapack "classic-fishing-loot" {
    url = "https://github.com/renevo/minecraft-datapack/releases/download/v1.19-1.1.0/classic-fishing-loot-v1.1.3.zip"
  }
  datapack "double-shulker-shells" {
    url = "https://github.com/renevo/minecraft-datapack/releases/download/v1.19-1.1.0/double-shulker-shells-v1.3.3.zip"
  }
  datapack "dragon-drops" {
    url = "https://github.com/renevo/minecraft-datapack/releases/download/v1.19-1.1.0/dragon-drops-v1.3.3.zip"
  }
  datapack "mob-heads" {
    url = "https://github.com/renevo/minecraft-datapack/releases/download/v1.19-1.1.0/mob-heads-v2.10.0.zip"
  }
  datapack "player-heads" {
    url = "https://github.com/renevo/minecraft-datapack/releases/download/v1.19-1.1.0/player-heads-v1.1.3.zip"
  }
  datapack "silence-mobs" {
    url = "https://github.com/renevo/minecraft-datapack/releases/download/v1.19-1.1.0/silence-mobs-v1.1.3.zip"
  }
  datapack "recipe-unlock" {
    url = "https://github.com/renevo/minecraft-datapack/releases/download/v1.19-1.1.0/unlock-all-recipes-v2.0.4.zip"
  }
  datapack "custom-recipes" {
    url = "https://github.com/renevo/minecraft-datapack/releases/download/v1.19-1.1.0/VanillaTweaks-crafting.zip"
  }

  purge_mods = true

  mod "carpet" {
    url = "https://github.com/gnembon/fabric-carpet/releases/download/1.4.93/fabric-carpet-1.19.3-1.4.93+v221230.jar"

    config "coleslaw/carpet.conf" {
      content = <<EOC
locked

antiCheatDisabled true
commandLog true
commandScript ops
defaultLoggers mobcaps,tps
lagFreeSpawning true
leadFix true
lightningKillsDropsFix true
persistentParrots true
scriptsAutoload true
stackableShulkerBoxes true
xpNoCooldown true
      EOC
    }
  }

  mod "fabric-api" {
    url = "https://github.com/FabricMC/fabric/releases/download/0.72.0%2B1.19.3/fabric-api-0.72.0+1.19.3.jar"
  }

  mod "malilib" {
    url = "https://pena2.dy.fi/tmp/minecraft/mods/malilib/malilib-fabric-1.19.3-0.14.0.jar"
  }
  
  mod "item-scroller" {
    url = "https://pena2.dy.fi/tmp/minecraft/mods/itemscroller/itemscroller-fabric-1.19.3-0.18.0.jar"
  }

  mod "inventory-sorter" {
    url = "https://mediafilez.forgecdn.net/files/4168/828/InventorySorter-1.8.10-1.19.3.jar"
  }

  mod "servux" {
    url = "https://kosma.pl/masamods/archive/servux-fabric-1.19-0.1.0.jar"
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

startup_commands = [
  "script download survival/silk_budding_amethyst.sc",
  "script download survival/silk_spawners.sc",
  "script download survival/simply_harvest.sc",
]
