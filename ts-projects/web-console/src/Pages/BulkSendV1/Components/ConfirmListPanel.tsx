/**
 * 群发1：待确认列表面板。
 *
 * 已知缺陷（对应真实测试用例 #30）：
 * 任务已被取消（status=5）后，列表里仍然渲染「通过」按钮，
 * 用户点击会走 confirm 接口并拿到 409。正确行为是取消/删除后不展示按钮。
 */
import React from 'react';

export type JobStatus = 1 | 2 | 3 | 4 | 5 | 6 | 7;

export interface ConfirmItem {
  jobId: string;
  name: string;
  jobStatus: JobStatus;
  jobDeleted?: boolean;
  notifyMode: 0 | 1;
  targetCount: number;
}

export interface ConfirmListPanelProps {
  items: ConfirmItem[];
  onConfirm: (jobId: string) => void;
  onCancel: (jobId: string) => void;
}

export function ConfirmListPanel({ items, onConfirm, onCancel }: ConfirmListPanelProps) {
  return (
    <table className="confirm-list">
      <thead>
        <tr>
          <th>任务名称</th>
          <th>状态</th>
          <th>目标数</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        {items.map((item) => (
          <tr key={item.jobId} data-testid={`confirm-row-${item.jobId}`}>
            <td>{item.name}</td>
            <td>{item.jobStatus}</td>
            <td>{item.targetCount}</td>
            <td>
              <button
                data-testid={`bulk-msg-official-contact-list-settings-edit-btn-${item.jobId}`}
                onClick={() => onConfirm(item.jobId)}
              >
                通过
              </button>
              <button
                data-testid={`bulk-msg-official-contact-list-settings-cancel-btn-${item.jobId}`}
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
