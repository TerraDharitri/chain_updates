import { TestingModule, Test } from "@nestjs/testing";
import { NotifierEvent } from "src/common/rabbitmq/entities/notifier.event";
import { NotifierEventIdentifier } from "src/common/rabbitmq/entities/notifier.event.identifier";
import { RabbitMqConsumer } from "src/common/rabbitmq/rabbitmq.consumer";
import { RabbitMqNftHandlerService } from "src/common/rabbitmq/rabbitmq.nft.handler.service";
import { RabbitMqTokenHandlerService } from "src/common/rabbitmq/rabbitmq.token.handler.service";

describe('RabbitMqConsumer', () => {
  let service: RabbitMqConsumer;
  let nftHandlerService: RabbitMqNftHandlerService;
  let tokenHandlerService: RabbitMqTokenHandlerService;

  beforeEach(async () => {
    const nftHandlerServiceMock = {
      handleNftCreateEvent: jest.fn(),
      handleNftUpdateAttributesEvent: jest.fn(),
    };

    const tokenHandlerServiceMock = {
      handleTransferOwnershipEvent: jest.fn(),
    };

    const module: TestingModule = await Test.createTestingModule({
      providers: [
        RabbitMqConsumer,
        { provide: RabbitMqNftHandlerService, useValue: nftHandlerServiceMock },
        { provide: RabbitMqTokenHandlerService, useValue: tokenHandlerServiceMock },
      ],
    }).compile();

    service = module.get<RabbitMqConsumer>(RabbitMqConsumer);
    nftHandlerService = module.get<RabbitMqNftHandlerService>(RabbitMqNftHandlerService);
    tokenHandlerService = module.get<RabbitMqTokenHandlerService>(RabbitMqTokenHandlerService);
  });

  it('should be defined', () => {
    expect(service).toBeDefined();
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  it('should call the appropriate handler based on the event identifier', async () => {
    const event1: NotifierEvent = {
      identifier: NotifierEventIdentifier.DCDTNFTCreate,
      address: "drt1",
      topics: [''],
    };

    const event2: NotifierEvent = {
      identifier: NotifierEventIdentifier.transferOwnership,
      address: "drt1",
      topics: [''],
    };

    await service.consumeEvents({ events: [event1, event2] });

    expect(nftHandlerService.handleNftCreateEvent).toHaveBeenCalledWith(event1);
    expect(tokenHandlerService.handleTransferOwnershipEvent).toHaveBeenCalledWith(event2);
  });

  it('should log the error when an unhandled error occurs', async () => {
    const event: NotifierEvent = {
      identifier: NotifierEventIdentifier.DCDTNFTCreate,
      address: "drt1qga7ze0l03chfgru0a32wxqf2226nzrxnyhzer9lmudqhjgy7ycq0wn4su",
      topics: [''],
    };

    const error = new Error('Test error');

    jest.spyOn(nftHandlerService, 'handleNftCreateEvent').mockImplementationOnce(() => {
      throw error;
    });

    const loggerSpy = jest.spyOn(service['logger'], 'error').mockImplementation(() =>
      `An unhandled error occurred when consuming events: {"events":[{"identifier":"DCDTNFTCreate","address":"drt1qga7ze0l03chfgru0a32wxqf2226nzrxnyhzer9lmudqhjgy7ycq0wn4su","topics":[""]}]}`);

    await service.consumeEvents({ events: [event] });

    expect(loggerSpy).toHaveBeenCalledWith(`An unhandled error occurred when consuming events: ${JSON.stringify({ events: [event] })}`);
    expect(loggerSpy).toHaveBeenCalledWith(error);
  });
});
