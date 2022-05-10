server "Burpcraft" {
  version  = "1.12.2"
  snapshot = false
  path     = "./burpcraft/legacy"

  properties = {
    motd             = "Hello world"
    allow-flying     = false
    spawn-protection = 16
  }
}
