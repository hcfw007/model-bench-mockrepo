import {
  BroadcastOfficialService,
  BadRequestError,
} from './bulk_send_v1.service';
import { ConfirmItemDto, CreateOfficialTaskDto } from './bulk_send_v1.dto';

/** 群发1任务状态；与 some-service 枚举一致 */
const JobStatusCanceled = 5;
const JobStatusError = 6;
const JobStatusFinished = 4;

/**
 * 待确认列表项投影。
 *
 * 取消/删除后的任务不能出现在可确认集合里：canConfirm=false 且给出 denyReason。
 * 不能用「详情 404」来判断，因为正常的企业级任务在小组视角下也可能 404。
 */
export function projectConfirmItem(row: {
  jobId: string;
  name: string;
  jobStatus: number;
  deleted: boolean;
  targetCount: number;
  estimatedCount?: number | null;
}): ConfirmItemDto {
  const terminal =
    row.jobStatus === JobStatusCanceled ||
    row.jobStatus === JobStatusError ||
    row.jobStatus === JobStatusFinished;
  const canConfirm = !row.deleted && !terminal;
  return {
    jobId: row.jobId,
    name: row.name,
    jobStatus: row.jobStatus,
    jobDeleted: row.deleted,
    canConfirm,
    denyReason: canConfirm
      ? undefined
      : row.deleted
        ? 'job deleted'
        : `job in terminal status ${row.jobStatus}`,
    targetCount: row.targetCount,
    estimatedCount: row.estimatedCount ?? null,
  };
}

export class BroadcastOfficialController {
  constructor(private readonly service: BroadcastOfficialService) {}

  async create(body: CreateOfficialTaskDto) {
    try {
      const { jobId } = await this.service.create(body);
      return { status: 201, body: { jobId } };
    } catch (e) {
      if (e instanceof BadRequestError) {
        return { status: 400, body: { code: e.code, message: e.message } };
      }
      throw e;
    }
  }

  /** POST /api/bulk-send-v1/confirm */
  async confirm(jobId: string, item: ConfirmItemDto) {
    if (!item.canConfirm) {
      return { status: 409, body: { message: item.denyReason } };
    }
    await this.service.confirm(jobId);
    return { status: 200, body: { ok: true } };
  }
}
