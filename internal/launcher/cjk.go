package launcher

import (
	"slices"

	"github.com/erik-adelbert/firework/internal/helper"
	"github.com/erik-adelbert/firework/internal/particle"
)

type runset struct {
	thresh float64
	runes  []rune
}

var palettes = [][]runset{
	particle.Spawning: {
		{0.3, []rune("。，”“』 『￥")},
		{0.5, []rune("一二三二三五十十已于上下义天")},
		{0.7, []rune("时中自字木月日目火田左右点以")},
		{1.0, []rune("龖龠龜")},
	},

	particle.Active: {
		{0.2, []rune("？。， 『』 ||")},
		{0.6, []rune("（）【】*￥|十一二三六")},
		{0.85, []rune("人中亿入上下火土")},
		{1.0, []rune("繁荣昌盛国泰民安龍龖龠龜耋")},
	},

	particle.Fading: {
		{0.6, []rune("。 『 』 、： |。，— ……")},
		{1.0, []rune("|￥人 上十入乙小 下")},
	},

	particle.Dead: {
		{1.0, []rune(" ")},
	},
}

func getrune(p *particle.Particle, d float64) rune {
	lut := palettes[p.State()]

	runes := lut[slices.IndexFunc(lut, func(rs runset) bool {
		return d <= rs.thresh
	})].runes

	return randomRune(runes)
}

func randomRune(runes []rune) rune {
	if len(runes) == 0 {
		return ' '
	}

	r, _ := helper.Pick(runes)

	return r
}
