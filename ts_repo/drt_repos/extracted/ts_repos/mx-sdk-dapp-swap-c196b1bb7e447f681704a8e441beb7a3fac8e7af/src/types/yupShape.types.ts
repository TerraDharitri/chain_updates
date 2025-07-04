import { AnySchema } from 'yup/lib/schema';

export type Shape<Fields extends object> = {
  [Key in keyof Fields]: AnySchema<Fields[Key]>;
};
