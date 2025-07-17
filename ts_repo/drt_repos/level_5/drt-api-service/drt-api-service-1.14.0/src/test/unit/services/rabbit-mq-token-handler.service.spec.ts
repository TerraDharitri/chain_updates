import { CacheService } from "@terradharitri/sdk-nestjs-cache";
import { TestingModule, Test } from "@nestjs/testing";
import { NotifierEvent } from "src/common/rabbitmq/entities/notifier.event";
import { RabbitMqTokenHandlerService } from "src/common/rabbitmq/rabbitmq.token.handler.service";
import { DcdtService } from "src/endpoints/dcdt/dcdt.service";

const cacheServiceMock = {
  set: jest.fn(),
};

const dcdtServiceMock = {
  getDcdtTokenPropertiesRaw: jest.fn(),
};

const clientProxyMock = {
  emit: jest.fn(),
};

describe('RabbitMqTokenHandlerService', () => {
  let service: RabbitMqTokenHandlerService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [
        RabbitMqTokenHandlerService,
        { provide: CacheService, useValue: cacheServiceMock },
        { provide: DcdtService, useValue: dcdtServiceMock },
        { provide: 'PUBSUB_SERVICE', useValue: clientProxyMock },
      ],
    }).compile();

    service = module.get<RabbitMqTokenHandlerService>(RabbitMqTokenHandlerService);
  });

  afterEach(() => {
    jest.resetAllMocks();
  });

  it('should be defined', () => {
    expect(service).toBeDefined();
  });

  it('should handle transfer ownership event with valid token properties', async () => {
    const tokenIdentifier = 'WREWA-bd4d79';
    const event: NotifierEvent = {
      topics: [Buffer.from(tokenIdentifier).toString('base64')],
      address: "drt1",
      identifier: "WREWA-bd4d79",
    };

    const dcdtProperties = { someProperty: 'value' };

    dcdtServiceMock.getDcdtTokenPropertiesRaw.mockResolvedValue(dcdtProperties);

    const result = await service.handleTransferOwnershipEvent(event);

    expect(result).toBe(true);
    expect(dcdtServiceMock.getDcdtTokenPropertiesRaw).toHaveBeenCalledWith(tokenIdentifier);
    expect(cacheServiceMock.set).toHaveBeenCalled();
    expect(clientProxyMock.emit).toHaveBeenCalled();
  });

  it('should handle transfer ownership event with no token properties', async () => {
    const tokenIdentifier = 'WREWA-bd4d79';
    const event: NotifierEvent = {
      topics: [Buffer.from(tokenIdentifier).toString('base64')],
      address: "drt1",
      identifier: "WREWA-bd4d79",
    };

    dcdtServiceMock.getDcdtTokenPropertiesRaw.mockResolvedValue(null);

    const result = await service.handleTransferOwnershipEvent(event);

    expect(result).toBe(false);
    expect(dcdtServiceMock.getDcdtTokenPropertiesRaw).toHaveBeenCalledWith(tokenIdentifier);
    expect(cacheServiceMock.set).not.toHaveBeenCalled();
    expect(clientProxyMock.emit).not.toHaveBeenCalled();
  });

  it('should handle transfer ownership event with exception', async () => {
    const tokenIdentifier = 'WREWA-bd4d79';
    const event: NotifierEvent = {
      topics: [Buffer.from(tokenIdentifier).toString('base64')],
      address: "drt1",
      identifier: "WREWA-bd4d79",
    };

    dcdtServiceMock.getDcdtTokenPropertiesRaw.mockRejectedValue(new Error('Test error'));

    const loggerSpy = jest.spyOn(service['logger'], 'error').mockImplementation(() =>
      "An unhandled error occurred when processing transferOwnership event for token with identifier 'WREWA-bd4d79'");

    const result = await service.handleTransferOwnershipEvent(event);

    expect(result).toBe(false);
    expect(dcdtServiceMock.getDcdtTokenPropertiesRaw).toHaveBeenCalledWith(tokenIdentifier);
    expect(cacheServiceMock.set).not.toHaveBeenCalled();
    expect(clientProxyMock.emit).not.toHaveBeenCalled();
    expect(loggerSpy).toHaveBeenCalledTimes(2);
  });

  it('should invalidate and refresh cache key', async () => {
    const key = 'test-key';
    const data = { someData: 'value' };
    const ttl = 1000;

    await (service as any).invalidateKey(key, data, ttl);

    expect(cacheServiceMock.set).toHaveBeenCalledWith(key, data, ttl);
    expect(clientProxyMock.emit).toHaveBeenCalledWith('refreshCacheKey', { key, ttl });
  });

  it('should refresh cache key', () => {
    const key = 'test-key';
    const ttl = 1000;

    (service as any).refreshCacheKey(key, ttl);

    expect(clientProxyMock.emit).toHaveBeenCalledWith('refreshCacheKey', { key, ttl });
  });
});
