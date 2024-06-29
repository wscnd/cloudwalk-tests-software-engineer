package logparser

type meansOfDeath string

const (
	MOD_UNKNOWN        meansOfDeath = "MOD_UNKNOWN"
	MOD_SHOTGUN        meansOfDeath = "MOD_SHOTGUN"
	MOD_GAUNTLET       meansOfDeath = "MOD_GAUNTLET"
	MOD_MACHINEGUN     meansOfDeath = "MOD_MACHINEGUN"
	MOD_GRENADE        meansOfDeath = "MOD_GRENADE"
	MOD_GRENADE_SPLASH meansOfDeath = "MOD_GRENADE_SPLASH"
	MOD_ROCKET         meansOfDeath = "MOD_ROCKET"
	MOD_ROCKET_SPLASH  meansOfDeath = "MOD_ROCKET_SPLASH"
	MOD_PLASMA         meansOfDeath = "MOD_PLASMA"
	MOD_PLASMA_SPLASH  meansOfDeath = "MOD_PLASMA_SPLASH"
	MOD_RAILGUN        meansOfDeath = "MOD_RAILGUN"
	MOD_LIGHTNING      meansOfDeath = "MOD_LIGHTNING"
	MOD_BFG            meansOfDeath = "MOD_BFG"
	MOD_BFG_SPLASH     meansOfDeath = "MOD_BFG_SPLASH"
	MOD_WATER          meansOfDeath = "MOD_WATER"
	MOD_SLIME          meansOfDeath = "MOD_SLIME"
	MOD_LAVA           meansOfDeath = "MOD_LAVA"
	MOD_CRUSH          meansOfDeath = "MOD_CRUSH"
	MOD_TELEFRAG       meansOfDeath = "MOD_TELEFRAG"
	MOD_FALLING        meansOfDeath = "MOD_FALLING"
	MOD_SUICIDE        meansOfDeath = "MOD_SUICIDE"
	MOD_TARGET_LASER   meansOfDeath = "MOD_TARGET_LASER"
	MOD_TRIGGER_HURT   meansOfDeath = "MOD_TRIGGER_HURT"
	MOD_NAIL           meansOfDeath = "MOD_NAIL"
	MOD_CHAINGUN       meansOfDeath = "MOD_CHAINGUN"
	MOD_PROXIMITY_MINE meansOfDeath = "MOD_PROXIMITY_MINE"
	MOD_KAMIKAZE       meansOfDeath = "MOD_KAMIKAZE"
	MOD_JUICED         meansOfDeath = "MOD_JUICED"
	MOD_GRAPPLE        meansOfDeath = "MOD_GRAPPLE"
)

func (d meansOfDeath) String() string {
	return string(d)
}
