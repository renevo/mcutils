server "Burpcraft" {
    version = "latest"
    snapshot = true
    path = "./burpcraft"

    world "burpcraft" {
        seed = "Burpcraft5"
    }

    game {
        difficulty = "hard"
        mode = "survival"
        mode_forced = true

        hardcore = false
        pvp = false
        flying = false

        spawning {
            animals = true
            monsters = true
            npcs = true
        }
    }
}
