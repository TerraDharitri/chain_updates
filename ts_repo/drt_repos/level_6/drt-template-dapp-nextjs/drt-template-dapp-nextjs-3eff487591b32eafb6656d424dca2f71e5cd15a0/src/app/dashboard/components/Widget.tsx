'use client';
import { Card } from '@/components';
import { WidgetType } from '@/types/widget.types';
import { TransactionsPropsType } from '../widgets/Transactions/types';

export const Widget = ({
  title,
  description,
  reference,
  anchor,
  widget: DrtWidget,
  props = {}
}: WidgetType<TransactionsPropsType>) => {
  return (
    <Card
      title={title}
      description={description}
      reference={reference}
      anchor={anchor}
    >
      <DrtWidget {...props} />
    </Card>
  );
};
