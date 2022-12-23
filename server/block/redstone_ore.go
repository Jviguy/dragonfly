package block

import "github.com/df-mc/dragonfly/server/item"

type RedstoneOre struct {
	solid
	Type OreType
}

// BreakInfo ...
func (r RedstoneOre) BreakInfo() BreakInfo {
	i := newBreakInfo(r.Type.Hardness(), func(t item.Tool) bool {
		return t.ToolType() == item.TypePickaxe && t.HarvestLevel() >= item.ToolTierIron.HarvestLevel
	}, pickaxeEffective, silkTouchOneOf(RedstoneDust{}, r))
	if r.Type == DeepslateOre() {
		i = i.withBlastResistance(9)
	}
	return i
}

// EncodeItem ...
func (r RedstoneOre) EncodeItem() (name string, meta int16) {
	return "minecraft:" + r.Type.Prefix() + "redstone_ore", 0
}

// EncodeBlock ...
func (r RedstoneOre) EncodeBlock() (string, map[string]any) {
	return "minecraft:" + r.Type.Prefix() + "redstone_ore", nil
}
