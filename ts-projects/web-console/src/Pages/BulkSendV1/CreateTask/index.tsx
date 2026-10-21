/**
 * 群发1：创建任务页。
 *
 * 已知缺陷（对应真实测试用例：定时任务待执行期间人数显示为 0）：
 * 目标人数只在创建后由服务端解析，待执行阶段 targetCount=0，
 * 页面直接把 0 当成「预计触达人数」展示，用户会误以为没选中对象。
 */
import React, { useEffect, useState } from 'react';

export interface PreviewMetrics {
  /** 服务端返回的预估触达数；尚未解析时为 0 */
  count: number;
  /** 是否为创建时预估（区别于下发时实际目标数） */
  estimated?: boolean;
}

export function CreateTask() {
  const [metrics, setMetrics] = useState<PreviewMetrics | null>(null);
  const [selected, setSelected] = useState<string[]>([]);

  useEffect(() => {
    let alive = true;
    fetch('/api/queryCondition/metricsUnique', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ type: 1, externalUserIds: selected }),
    })
      .then((r) => r.json())
      .then((d) => {
        if (alive) setMetrics({ count: d.data?.count ?? 0 });
      });
    return () => {
      alive = false;
    };
  }, [selected]);

  return (
    <section className="bulk-send-v1-create">
      <h2>新建群发1任务</h2>
      <input
        data-testid="bulk-msg-official-contact-create-name-input"
        placeholder="任务名称"
      />
      <div className="preview">
        预计触达人数：<span>{metrics?.count ?? 0}</span>
      </div>
      <button
        data-testid="bulk-msg-official-contact-create-submit-btn"
        onClick={() => setSelected([])}
      >
        确定
      </button>
    </section>
  );
}
