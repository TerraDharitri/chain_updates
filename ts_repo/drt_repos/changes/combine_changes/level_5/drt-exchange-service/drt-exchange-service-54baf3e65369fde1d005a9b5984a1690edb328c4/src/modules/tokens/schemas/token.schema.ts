import { Prop, Schema, SchemaFactory } from '@nestjs/mongoose';
import { Document } from 'mongoose';

export type TokenDocument = DcdtTokenDbModel & Document;

@Schema()
export class DcdtTokenDbModel {
    @Prop()
    tokenID: string;
    @Prop()
    type: string;
}

export const DcdtTokenSchema = SchemaFactory.createForClass(DcdtTokenDbModel);
