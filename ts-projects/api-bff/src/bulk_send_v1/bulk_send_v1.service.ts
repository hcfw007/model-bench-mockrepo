import {
  ATTACHMENTS_MAX,
  CreateOfficialTaskDto,
  MATERIALS_MAX,
  OfficialMaterialDto,
} from './bulk_send_v1.dto';

export interface MaterialRecord {
  materialId: string;
  /** 素材库里记录的真实类型 */
  realType: string;
}

/** 素材库读取。真实实现走 material 服务。 */
export interface MaterialRepo {
  get(materialId: string): Promise<MaterialRecord | undefined>;
}

export class BadRequestError extends Error {
  constructor(public readonly code: string, message: string) {
    super(message);
  }
}

export class BroadcastOfficialService {
  constructor(
    private readonly materials: MaterialRepo,
    private readonly automation: {
      create(input: unknown): Promise<{ jobId: string }>;
      confirm(jobId: string): Promise<void>;
    },
  ) {}

  /**
   * 素材规则校验。
   *
   * 注意：materialType 是前端声明值，不可信。真实实现按素材库的真实类型解析，
   * 因此「传 Voice 但素材实际是 Text」不是缺陷，而是设计如此。
   */
  async validateMaterials(items: OfficialMaterialDto[]): Promise<void> {
    if (items.length > MATERIALS_MAX) {
      throw new BadRequestError(
        'MATERIALS_TOO_MANY',
        'materials must contain no more than 10 elements',
      );
    }
    const attachments = items.filter((m) => m.materialType !== 'Text');
    const texts = items.filter((m) => m.materialType === 'Text');
    if (attachments.length > ATTACHMENTS_MAX) {
      throw new BadRequestError('ATTACHMENTS_TOO_MANY', 'attachments must not exceed 9');
    }
    if (texts.length > 1) {
      throw new BadRequestError('TEXT_TOO_MANY', 'text material must not exceed 1');
    }
    for (const item of items) {
      const record = await this.materials.get(item.materialId);
      if (!record) {
        throw new BadRequestError('MATERIAL_NOT_FOUND', `material ${item.materialId} not found`);
      }
      if (record.realType === 'Emoticon') {
        throw new BadRequestError('MATERIAL_TYPE_UNSUPPORTED', 'Emoticon is not supported');
      }
    }
  }

  async create(dto: CreateOfficialTaskDto): Promise<{ jobId: string }> {
    await this.validateMaterials(dto.materials);
    return this.automation.create(dto);
  }

  async confirm(jobId: string): Promise<void> {
    return this.automation.confirm(jobId);
  }
}
