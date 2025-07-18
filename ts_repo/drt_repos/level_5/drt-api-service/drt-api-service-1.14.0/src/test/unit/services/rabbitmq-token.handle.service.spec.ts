import { CacheService } from "@terradharitri/sdk-nestjs-cache";
import { BinaryUtils } from "@terradharitri/sdk-nestjs-common";
import { TestingModule, Test } from "@nestjs/testing";
import { NotifierEvent } from "src/common/rabbitmq/entities/notifier.event";
import { RabbitMqTokenHandlerService } from "src/common/rabbitmq/rabbitmq.token.handler.service";
import { DcdtService } from "src/endpoints/dcdt/dcdt.service";

describe('RabbitMqTokenHandlerService', () => {
  let service: RabbitMqTokenHandlerService;
  let dcdtService: DcdtService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [
        RabbitMqTokenHandlerService,
        {
          provide: CacheService,
          useValue: {
            set: jest.fn(),
          },
        },
        {
          provide: DcdtService,
          useValue: {
            getDcdtTokenPropertiesRaw: jest.fn(),
          },
        },
        {
          provide: 'PUBSUB_SERVICE',
          useValue: {
            emit: jest.fn(),
          },
        },
      ],
    }).compile();

    service = module.get<RabbitMqTokenHandlerService>(RabbitMqTokenHandlerService);
    dcdtService = module.get<DcdtService>(DcdtService);
  });

  it('should be defined', () => {
    expect(service).toBeDefined();
  });

  describe('handleTransferOwnershipEvent', () => {
    const event: NotifierEvent = {
      identifier: 'DCDTNFTCreate',
      address: 'drt1',
      topics: [BinaryUtils.base64Encode('test-topic')],
    };

    const tokenIdentifier = 'test-topic';

    it('should return false if no properties are found for the token', async () => {
      jest.spyOn(dcdtService, 'getDcdtTokenPropertiesRaw').mockResolvedValue(null);

      expect(await service.handleTransferOwnershipEvent(event)).toBe(false);
      expect(dcdtService.getDcdtTokenPropertiesRaw).toHaveBeenCalledWith(tokenIdentifier);
    });
  });
});
