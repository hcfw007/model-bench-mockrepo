export type MaterialType = 'Text' | 'Image' | 'File' | 'Voice' | 'Emoticon';

export interface OfficialMaterialDto {
  materialId: string;
  /** 前端声明的类型。服务端不应盲信，需按真实素材类型解析。 */
  materialType: MaterialType;
}

export interface CreateOfficialTaskDto {
  name: string;
  notifyMode: 0 | 1;
  /** 精确选择时传；全部客户不传 */
  externalUserIds?: string[];
  materials: OfficialMaterialDto[];
  markTag?: {
    tagType: number;
    tagIds: string[];
  };
}

export interface ConfirmItemDto {
  jobId: string;
  name: string;
  jobStatus: number;
  jobDeleted: boolean;
  canConfirm: boolean;
  denyReason?: string;
  targetCount: number;
  estimatedCount?: number | null;
}

export const MATERIALS_MAX = 10;
export const ATTACHMENTS_MAX = 9;
