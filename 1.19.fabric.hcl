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

  purge_datapacks = true

  datapack "minecraft-datapack" {
    url = "https://github.com/renevo/minecraft-datapack/archive/refs/tags/v1.19-1.0.0.zip"
  }

  purge_mods = true

  mod "fabric-api" {
    url = "https://github.com/FabricMC/fabric/releases/download/0.55.3%2B1.19/fabric-api-0.55.3+1.19.jar"
  }

  mod "bluemap" {
    url = "https://github.com/BlueMap-Minecraft/BlueMap/releases/download/v1.7.3/BlueMap-1.7.3-fabric-1.18.jar"

    config "config/bluemap/core.conf" {
      content = <<EOC
accept-download: true
renderThreadCount: 1
metrics: false
data: "bluemap"
EOC
    }

    config "config/bluemap/plugin.conf" {
      content = <<EOC
liveUpdates: true
skinDownload: true
hiddenGameModes: [
        "spectator"
]
hideInvisible: true
hideSneaking: false
fullUpdateInterval: 1440
EOC
    }

    config "config/bluemap/render.conf" {
      content = <<EOC
webroot: "bluemap/web"
useCookies: true
enableFreeFlight: true
maps: [
        {
                id: "world"
                name: "Burpcraft"
                world: "burpcraft"
                skyColor: "#7dabff"
                ambientLight: 0
                renderCaves: true
                renderEdges: true
                useCompression: true
                ignoreMissingLightData: false
        }
        {
                id: "end"
                name: "End"
                world: "burpcraft/DIM1"
                skyColor: "#080010"
                renderCaves: true
                ambientLight: 0.6
        }
        {
                id: "nether"
                name: "Nether"
                world: "burpcraft/DIM-1"
                skyColor: "#290000"
                renderCaves: true
                ambientLight: 0.6
                renderEdges: true
        }
]
EOC
    }

    config "config/bluemap/webserver.conf" {
      content = <<EOC
enabled: true
webroot: "bluemap/web"
port: 8100
maxConnectionCount: 100
EOC
    }
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
