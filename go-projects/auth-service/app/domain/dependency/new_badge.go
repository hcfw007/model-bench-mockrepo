package dependency

import "time"

// NewBadgeDefaultWindow 未显式指定时，NEW 角标的默认有效期。
const NewBadgeDefaultWindow = 30 * 24 * time.Hour

// NewBadgeRegionLaunchFeatures 按区域登记功能的上线时间。
// 只有登记在此处的功能才会计算并展示 NEW 角标。
//
// 新增功能时，key 必须与前端菜单使用的 feature key 完全一致，
// 且需要与前端同批 release 上线。
var NewBadgeRegionLaunchFeatures = map[string]map[string]time.Time{
	"cn": {
		"AI_CHAT":    time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
		"BROADCAST":  time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC),
		"MATERIAL":   time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC),
		"AUTOMATION": time.Date(2026, 6, 10, 0, 0, 0, 0, time.UTC),
	},
	"intl": {
		"AI_CHAT":   time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
		"BROADCAST": time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
	},
}

// BadgeWindow 返回某区域某功能的 NEW 角标有效期；未登记则返回 0 表示不展示。
func BadgeWindow(region, feature string) time.Duration {
	features, ok := NewBadgeRegionLaunchFeatures[region]
	if !ok {
		return 0
	}
	if _, ok := features[feature]; !ok {
		return 0
	}
	return NewBadgeDefaultWindow
}
