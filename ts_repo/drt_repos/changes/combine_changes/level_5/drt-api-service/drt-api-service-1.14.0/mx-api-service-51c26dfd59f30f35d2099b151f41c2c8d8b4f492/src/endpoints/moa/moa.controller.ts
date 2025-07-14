import { MoaEconomics } from './entities/moa.economics';
import { MoaToken } from './entities/moa.token';
import { Controller, DefaultValuePipe, Get, NotFoundException, Param, Query } from "@nestjs/common";
import { ApiExcludeEndpoint, ApiNotFoundResponse, ApiOkResponse, ApiOperation, ApiQuery, ApiResponse, ApiTags } from "@nestjs/swagger";
import { MoaPair } from "./entities/moa.pair";
import { MoaSettings } from "./entities/moa.settings";
import { MoaEconomicsService } from "./moa.economics.service";
import { MoaPairService } from "./moa.pair.service";
import { MoaSettingsService } from "./moa.settings.service";
import { MoaTokenService } from "./moa.token.service";
import { MoaFarmService } from './moa.farm.service';
import { MoaFarm } from './entities/moa.farm';
import { QueryPagination } from 'src/common/entities/query.pagination';
import { ParseIntPipe, ParseTokenPipe, ParseEnumPipe, ParseBoolPipe } from '@terradharitri/sdk-nestjs-common';
import { MoaPairExchange } from './entities/moa.pair.exchange';
import { MoaPairsFilter } from './entities/moa.pairs..filter';
import { MoaTokenChartsService } from './moa.token.charts.service';
import { MoaTokenChart } from './entities/moa.token.chart';

@Controller()
@ApiTags('dharitrix')
export class MoaController {
  constructor(
    private readonly moaEconomicsService: MoaEconomicsService,
    private readonly moaSettingsService: MoaSettingsService,
    private readonly drtPairsService: MoaPairService,
    private readonly moaTokensService: MoaTokenService,
    private readonly moaFarmsService: MoaFarmService,
    private readonly moaTokenChartsService: MoaTokenChartsService
  ) { }

  @Get("/moa/settings")
  @ApiExcludeEndpoint()
  @ApiResponse({ status: 200, description: 'The settings of the xExchange' })
  @ApiNotFoundResponse({ description: 'MOA settings not found' })
  async getMoaSettings(): Promise<MoaSettings> {
    const settings = await this.moaSettingsService.getSettings();
    if (!settings) {
      throw new NotFoundException('MOA settings not found');
    }
    return settings;
  }

  @Get("/moa/economics")
  @ApiOperation({ summary: 'xExchange economics', description: 'Returns economics details of xExchange' })
  @ApiOkResponse({ type: MoaEconomics })
  async getMoaEconomics(): Promise<MoaEconomics> {
    return await this.moaEconomicsService.getMoaEconomics();
  }

  @Get("/moa/pairs")
  @ApiOperation({ summary: 'xExchange pairs', description: 'Returns active liquidity pools available on xExchange' })
  @ApiOkResponse({ type: [MoaPair] })
  @ApiQuery({ name: 'from', description: 'Number of items to skip for the result set', required: false })
  @ApiQuery({ name: 'size', description: 'Number of items to retrieve', required: false })
  @ApiQuery({ name: 'exchange', description: 'Filter by exchange', required: false, enum: MoaPairExchange })
  @ApiQuery({ name: 'includeFarms', description: 'Include farms information in response', required: false, type: Boolean })
  async getMoaPairs(
    @Query('from', new DefaultValuePipe(0), ParseIntPipe) from: number,
    @Query("size", new DefaultValuePipe(25), ParseIntPipe) size: number,
    @Query('exchange', new ParseEnumPipe(MoaPairExchange)) exchange?: MoaPairExchange,
    @Query('includeFarms', new DefaultValuePipe(false), ParseBoolPipe) includeFarms?: boolean,
  ): Promise<MoaPair[]> {
    const filter = new MoaPairsFilter({ exchange, includeFarms });
    return await this.drtPairsService.getMoaPairs(from, size, filter);
  }

  @Get("/moa-pairs")
  @ApiOperation({ summary: 'xExchange pairs', description: 'Returns active liquidity pools available on xExchange', deprecated: true })
  @ApiOkResponse({ type: [MoaPair] })
  @ApiQuery({ name: 'from', description: 'Number of items to skip for the result set', required: false })
  @ApiQuery({ name: 'size', description: 'Number of items to retrieve', required: false })
  @ApiQuery({ name: 'exchange', description: 'Filter by exchange', required: false, enum: MoaPairExchange })
  async getMoaPairsTemp(
    @Query('from', new DefaultValuePipe(0), ParseIntPipe) from: number,
    @Query("size", new DefaultValuePipe(25), ParseIntPipe) size: number,
    @Query('exchange', new ParseEnumPipe(MoaPairExchange)) exchange?: MoaPairExchange,
  ): Promise<MoaPair[]> {
    const filter = new MoaPairsFilter({ exchange });
    return await this.drtPairsService.getMoaPairs(from, size, filter);
  }

  @Get("/moa/pairs/count")
  @ApiOperation({ summary: 'Maiar Exchange pairs count', description: 'Returns active liquidity pools count available on Maiar Exchange' })
  @ApiQuery({ name: 'exchange', description: 'Filter by exchange', required: false, enum: MoaPairExchange })
  @ApiQuery({ name: 'includeFarms', description: 'Include farms information in response', required: false, type: Boolean })
  async getMoaPairsCount(
    @Query('exchange', new ParseEnumPipe(MoaPairExchange)) exchange?: MoaPairExchange,
    @Query('includeFarms', new DefaultValuePipe(false), ParseBoolPipe) includeFarms?: boolean,
  ): Promise<number> {
    const filter = new MoaPairsFilter({ exchange, includeFarms });
    return await this.drtPairsService.getMoaPairsCount(filter);
  }

  @Get("/moa/tokens")
  @ApiOperation({ summary: 'xExchange tokens details', description: 'Returns a list of tokens listed on xExchange' })
  @ApiOkResponse({ type: [MoaToken] })
  @ApiQuery({ name: 'from', description: 'Number of items to skip for the result set', required: false })
  @ApiQuery({ name: 'size', description: 'Number of items to retrieve', required: false })
  async getMoaTokens(
    @Query('from', new DefaultValuePipe(0), ParseIntPipe) from: number,
    @Query("size", new DefaultValuePipe(25), ParseIntPipe) size: number,
  ): Promise<MoaToken[]> {
    return await this.moaTokensService.getMoaTokens(new QueryPagination({ from, size }));
  }

  @Get("/moa/tokens/count")
  @ApiOperation({ summary: 'Maiar Exchange tokens count', description: 'Returns tokens count available on Maiar Exchange' })
  async getMoaTokensCount(
  ): Promise<number> {
    return await this.moaTokensService.getMoaTokensCount();
  }

  @Get("/moa/tokens/:identifier")
  @ApiOperation({ summary: 'xExchange token details', description: 'Returns a specific token listed on xExchange' })
  @ApiOkResponse({ type: MoaToken })
  @ApiNotFoundResponse({ description: 'Token not found' })
  async getMoaTokenIdentifier(
    @Param('identifier', ParseTokenPipe) identifier: string
  ): Promise<MoaToken> {
    const moaToken = await this.moaTokensService.getMoaTokenByIdentifier(identifier);
    if (!moaToken) {
      throw new NotFoundException('Token not found');
    }

    return moaToken;
  }

  @Get("/moa/farms")
  @ApiOperation({ summary: 'xExchange farms details', description: 'Returns a list of farms listed on xExchange' })
  @ApiOkResponse({ type: [MoaFarm] })
  @ApiQuery({ name: 'from', description: 'Number of items to skip for the result set', required: false })
  @ApiQuery({ name: 'size', description: 'Number of items to retrieve', required: false })
  async getMoaFarms(
    @Query('from', new DefaultValuePipe(0), ParseIntPipe) from: number,
    @Query("size", new DefaultValuePipe(25), ParseIntPipe) size: number
  ): Promise<MoaFarm[]> {
    return await this.moaFarmsService.getMoaFarms(new QueryPagination({ from, size }));
  }

  @Get("/moa/farms/count")
  @ApiOperation({ summary: 'Maiar Exchange farms count', description: 'Returns farms count available on Maiar Exchange' })
  async getMoaFarmsCount(
  ): Promise<number> {
    return await this.moaFarmsService.getMoaFarmsCount();
  }

  @Get("/moa/pairs/:baseId/:quoteId")
  @ApiOperation({ summary: 'xExchange pairs details', description: 'Returns liquidity pool details by providing a combination of two tokens' })
  @ApiOkResponse({ type: MoaPair })
  @ApiNotFoundResponse({ description: 'Pair not found' })
  @ApiQuery({ name: 'includeFarms', description: 'Include farms information in response', required: false, type: Boolean })
  async getMoaPair(
    @Param('baseId') baseId: string,
    @Param('quoteId') quoteId: string,
    @Query('includeFarms', new DefaultValuePipe(false), ParseBoolPipe) includeFarms?: boolean,
  ): Promise<MoaPair> {
    const pair = await this.drtPairsService.getMoaPair(baseId, quoteId, includeFarms);
    if (!pair) {
      throw new NotFoundException('Pair not found');
    }

    return pair;
  }

  @Get('moa/tokens/prices/hourly/:identifier')
  @ApiOperation({ summary: 'xExchange token prices hourly', description: 'Returns token prices hourly' })
  @ApiOkResponse({ type: [MoaTokenChart] })
  @ApiNotFoundResponse({ description: 'Price not available for given token identifier' })
  async getTokenPricesHourResolution(
    @Param('identifier', ParseTokenPipe) identifier: string): Promise<MoaTokenChart[] | undefined> {
    const charts = await this.moaTokenChartsService.getTokenPricesHourResolution(identifier);
    if (!charts) {
      throw new NotFoundException('Price not available for given token identifier');
    }

    return charts;
  }

  @Get('moa/tokens/prices/daily/:identifier')
  @ApiOperation({
    summary: 'xExchange token prices daily',
    description: 'Returns token prices daily, ordered by timestamp in ascending order. The entries represent the latest complete daily values for the given token series.',
  })
  @ApiOkResponse({ type: [MoaTokenChart] })
  @ApiNotFoundResponse({ description: 'Price not available for given token identifier' })
  async getTokenPricesDayResolution(
    @Param('identifier', ParseTokenPipe) identifier: string): Promise<MoaTokenChart[] | undefined> {
    const charts = await this.moaTokenChartsService.getTokenPricesDayResolution(identifier);
    if (!charts) {
      throw new NotFoundException('Price not available for given token identifier');
    }

    return charts;
  }
}
