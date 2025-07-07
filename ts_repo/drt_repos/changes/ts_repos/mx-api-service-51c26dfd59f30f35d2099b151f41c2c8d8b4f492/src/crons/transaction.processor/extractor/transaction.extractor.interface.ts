import { ShardTransaction } from "@terradharitri/sdk-transaction-processor";
import { TransactionDetailed } from "src/endpoints/transactions/entities/transaction.detailed";

export interface TransactionExtractorInterface<T> {
  extract(transaction: ShardTransaction, transactionDetailed?: TransactionDetailed): T | undefined;
}
