import { BaseSchema } from 'yup';
import { RulesTypes } from 'validation/types/rules.types';

export const applyValidationSchemaRules = <TValue>(
  initialSchema: BaseSchema,
  rules: RulesTypes<TValue>
) => {
  rules.forEach(({ name, message, test }) => {
    initialSchema = initialSchema.test(name, message, test);
  });

  return initialSchema;
};
