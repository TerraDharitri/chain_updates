export type RuleType<TValue> = {
  name: string;
  message: string;
  test: (value: TValue) => boolean;
};

export type RulesTypes<TValue> = Array<RuleType<TValue>>;
