/**
 * 群发1：群聊待确认列表面板。
 *
 * 与私聊版 ConfirmListPanel 结构几乎一致（历史遗留的重复实现，
 * 两处逻辑需同步维护）。请勿与私聊版混淆。
 */
import React from 'react';
import type { ConfirmItem } from './ConfirmListPanel';

export interface RoomConfirmListPanelProps {
  items: ConfirmItem[];
  onConfirm: (jobId: string) => void;
  onCancel: (jobId: string) => void;
}

export function RoomConfirmListPanel({ items, onConfirm, onCancel }: RoomConfirmListPanelProps) {
  return (
    <table className="room-confirm-list">
      <thead>
        <tr>
          <th>群聊任务</th>
          <th>状态</th>
          <th>目标数</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        {items.map((item) => (
          <tr key={item.jobId} data-testid={`room-confirm-row-${item.jobId}`}>
            <td>{item.name}</td>
            <td>{item.jobStatus}</td>
            <td>{item.targetCount}</td>
            <td>
              <button
                data-testid={`bulk-msg-official-room-list-settings-edit-btn-${item.jobId}`}
                onClick={() => onConfirm(item.jobId)}
              >
                通过
              </button>
              <button
                data-testid={`bulk-msg-official-room-list-settings-cancel-btn-${item.jobId}`}
                onClick={() => onCancel(item.jobId)}
              >
                取消
              </button>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
