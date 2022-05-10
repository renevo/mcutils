server "Burpcraft" {
  version  = "1.18.2"
  snapshot = false
  path     = "./burpcraft/vanilla"

  properties = {
    motd             = "Hello world"
    allow-flying     = false
    spawn-protection = 16
  }
}
