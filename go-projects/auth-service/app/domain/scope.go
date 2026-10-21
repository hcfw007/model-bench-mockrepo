package domain

import "sort"

// ScopeResult 双重范围解析结果。
type ScopeResult struct {
	// Allowed 同时满足应用可见范围与客户联系使用范围的账号。
	Allowed []string
	// Filtered 被过滤掉的账号，需要在前端弹窗中展示。
	Filtered []string
}

// ResolveScope 解析托管账号的双重使用范围：
// 应用可见范围 ∩ 客户联系使用范围。
//
// 超管继承其可见范围内的全部权限。
func ResolveScope(appVisible, contactUsable []string, isSuperAdmin bool) ScopeResult {
	inContact := make(map[string]struct{}, len(contactUsable))
	for _, id := range contactUsable {
		inContact[id] = struct{}{}
	}

	var allowed, filtered []string
	for _, id := range appVisible {
		if _, ok := inContact[id]; ok {
			allowed = append(allowed, id)
		} else {
			filtered = append(filtered, id)
		}
	}
	return ScopeResult{Allowed: allowed, Filtered: filtered}
}

// StableScope 返回排序去重后的范围，用于前端稳定展示与测试断言。
func StableScope(ids []string) []string {
	seen := make(map[string]struct{}, len(ids))
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}
